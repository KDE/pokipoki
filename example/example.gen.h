

#pragma once

#include <QAbstractListModel>
#include <QDebug>
#include <QMutex>
#include <QMutexLocker>
#include <QObject>
#include <QPointer>
#include <QSharedPointer>
#include <QSqlError>
#include <QSqlQuery>
#include <QUuid>
#include <QVariant>
#include <QString>
#include <QMap>

#include "Database.h"

class NoteModel;

class Note : public QObject, PPUndoRedoable {
	Q_OBJECT
	Q_INTERFACES(PPUndoRedoable)

	struct Change {
		
		
		
		Optional<QString> previoustitleValue;
		
		
		
		Optional<QMap<QString,QString>> previousmetadataValue;
		
	};

	Note(QUuid ID) : QObject(nullptr), m_ID(ID) {
		static bool db_initialized = false;
		if (!db_initialized) {
			volatile auto db = PPDatabase::instance();
			Q_UNUSED(db)
			prepareDatabase();
			db_initialized = true;
		}
	}

	~Note() {
		if (m_DELETE_PENDING) {
			QSqlQuery query;
			query.prepare(QStringLiteral("DELETE FROM Note WHERE ID = :ID"));
			query.bindValue(":ID", QVariant::fromValue(m_ID));
			query.exec();
		}
	}

	static QSharedPointer<Note> withID(QUuid ID) {
		static QMap<QUuid,QPointer<Note>> s_instances;
		static QMutex s_mutex;

		QMutexLocker locker(&s_mutex);
		auto val = s_instances.value(ID, nullptr);
		if (val.isNull()) {
			s_instances[ID] = new Note(ID);
		}
		return QSharedPointer<Note>(s_instances[ID].data(), &QObject::deleteLater);
	}

	
	QUuid m_parent_Note_ID;
	
	
	friend class NoteModel;

	QUuid m_ID;
	bool m_NEW = false;
	bool m_DELETE_PENDING = false;
	bool m_DIRTY = false;
	bool m_CAN_UNDO = false;
	bool m_CAN_REDO = false;

	Q_PROPERTY(bool pendingDelete READ pendingDelete NOTIFY pendingDeleteChanged)
	Q_PROPERTY(bool dirty READ dirty NOTIFY dirtyChanged)
	Q_PROPERTY(bool canUndo READ canUndo NOTIFY canUndoChanged)
	Q_PROPERTY(bool canRedo READ canRedo NOTIFY canRedoChanged)

	QList<Change> m_UNDO_STACK;
	QList<Change> m_REDO_STACK;

	
	
	
	Q_PROPERTY(QString title READ title WRITE set_title NOTIFY titleChanged)
	QString m_title;
	QString m_title_prev;
	bool m_title_dirty;
	
	
	
	Q_PROPERTY(QMap<QString,QString> metadata READ metadata WRITE set_metadata NOTIFY metadataChanged)
	QMap<QString,QString> m_metadata;
	QMap<QString,QString> m_metadata_prev;
	bool m_metadata_dirty;
	

	void evaluate_can_undo_changed() {
		auto setUndo = [this](bool newUndo){
			if (newUndo != m_CAN_UNDO) {
				m_CAN_UNDO = newUndo;
				Q_EMIT canUndoChanged();
			}
		};
		if (!m_UNDO_STACK.empty()) {
			setUndo(true);
			return;
		}
		setUndo(false);
	}

	void evaluate_can_redo_changed() {
		auto setRedo = [this](bool newRedo){
			if (newRedo != m_CAN_REDO) {
				m_CAN_REDO = newRedo;
				Q_EMIT canRedoChanged();
			}
		};
		if (!m_REDO_STACK.empty()) {
			setRedo(true);
			return;
		}
		setRedo(false);
	}

	void clear_redo() {
		m_REDO_STACK.clear();
		evaluate_can_redo_changed();
	}

	void evaluate_dirty_changed() {
		if (m_DIRTY) {
			auto new_dirty = false;
			
			if (m_title_dirty) {
				new_dirty = true;
			}
			
			if (m_metadata_dirty) {
				new_dirty = true;
			}
			
			if (!new_dirty) {
				m_DIRTY = false;
			}
			Q_EMIT dirtyChanged();
		} else {
			
			if (m_title_dirty) {
				m_DIRTY = true;
				Q_EMIT dirtyChanged();
				return;
			}
			
			if (m_metadata_dirty) {
				m_DIRTY = true;
				Q_EMIT dirtyChanged();
				return;
			}
			
		}
		evaluate_can_undo_changed();
	}

public:

	Q_SIGNAL void pendingDeleteChanged();
	Q_SIGNAL void dirtyChanged();
	Q_SIGNAL void canUndoChanged();
	Q_SIGNAL void canRedoChanged();
	bool pendingDelete() const { return m_DELETE_PENDING; }
	bool dirty() const { return m_DIRTY; }
	bool canUndo() const { return m_CAN_UNDO; }
	bool canRedo() const { return m_CAN_REDO; }

	Q_INVOKABLE void undo() override {
		if (!m_UNDO_STACK.empty()) {
			pUR->undoItemRemoved(this);
			auto last = m_UNDO_STACK.takeLast();
			
			if (last.previoustitleValue.has_value()) {
				last.previoustitleValue.swap(m_title);
			}
			
			if (last.previousmetadataValue.has_value()) {
				last.previousmetadataValue.swap(m_metadata);
			}
			
			m_REDO_STACK << last;
			pUR->redoItemAdded(this);
			evaluate_can_undo_changed();
			evaluate_can_redo_changed();
		}
	}

	Q_INVOKABLE void redo() override {
		if (!m_REDO_STACK.empty()) {
			pUR->redoItemRemoved(this);
			auto last = m_REDO_STACK.takeLast();
			
			if (last.previoustitleValue.has_value()) {
				last.previoustitleValue.swap(m_title);
			}
			
			if (last.previousmetadataValue.has_value()) {
				last.previousmetadataValue.swap(m_metadata);
			}
			
			m_UNDO_STACK << last;
			evaluate_can_undo_changed();
			evaluate_can_redo_changed();
		}
	}

	Q_INVOKABLE void stageDelete() {
		if (!m_DELETE_PENDING) {
			m_DELETE_PENDING = true;
			pendingDeleteChanged();
		}
	}

	
	
	
	Q_SIGNAL void titleChanged();
	QString title() const { return m_title; };
	void set_title(const QString& val) {
		if (val == m_title) {
			return;
		}
		m_title_prev = m_title;
		m_title_dirty = true;
		m_title = val;
		Q_EMIT void titleChanged();
		clear_redo();
		evaluate_dirty_changed();
		evaluate_can_undo_changed();
	}
	void discard_title_changes() {
		if (m_title_dirty) {
			m_title_dirty = false;
			m_title = m_title_prev;
			Q_EMIT void titleChanged();
			evaluate_dirty_changed();
		}
	}
	
	
	
	Q_SIGNAL void metadataChanged();
	QMap<QString,QString> metadata() const { return m_metadata; };
	void set_metadata(const QMap<QString,QString>& val) {
		if (val == m_metadata) {
			return;
		}
		m_metadata_prev = m_metadata;
		m_metadata_dirty = true;
		m_metadata = val;
		Q_EMIT void metadataChanged();
		clear_redo();
		evaluate_dirty_changed();
		evaluate_can_undo_changed();
	}
	void discard_metadata_changes() {
		if (m_metadata_dirty) {
			m_metadata_dirty = false;
			m_metadata = m_metadata_prev;
			Q_EMIT void metadataChanged();
			evaluate_dirty_changed();
		}
	}
	

	void discard_all_changes() {
		
		if (m_title_dirty) {
			m_title_dirty = false;
			m_title = m_title_prev;
			Q_EMIT void titleChanged();
		}
		
		if (m_metadata_dirty) {
			m_metadata_dirty = false;
			m_metadata = m_metadata_prev;
			Q_EMIT void metadataChanged();
		}
		
		evaluate_dirty_changed();
	}

	Q_INVOKABLE void save() {
		if (m_NEW || m_DELETE_PENDING) {
			auto tq = QStringLiteral(R"RJIENRLWEY(
INSERT INTO Note
(ID,title,metadata)
VALUES
(:ID,   :title  , :metadata );
			)RJIENRLWEY");
			QSqlQuery query;
			query.prepare(tq);
			query.bindValue(":ID", QVariant::fromValue(m_ID));
			query.bindValue(":title", QVariant::fromValue(m_title));
			
			query.bindValue(":metadata", QVariant::fromValue(m_metadata));
			auto res = query.exec();
			if (!res) {
				qCritical() << query.lastError() << "when creating a new item of Note";
			}
			m_NEW = false;
			if (m_DELETE_PENDING) {
				m_DELETE_PENDING = false;
				pendingDeleteChanged();
			}
		} else {
		Change changes;
		
		if (m_title_dirty) {
			changes.previoustitleValue.copy(m_title_prev);
			QSqlQuery query;
			auto tq = QStringLiteral(R"RJIENRLWEY( UPDATE Note SET title = :val WHERE ID = :id )RJIENRLWEY");
			query.prepare(tq);
			query.bindValue(":val", QVariant::fromValue(m_title));
			query.bindValue(":id", QVariant::fromValue(m_ID));
			auto res = query.exec();
			if (!res) {
				qCritical() << query.lastError() << "when updating an item of type Note at row title";
			}
			m_title_dirty = false;
		}
		
		if (m_metadata_dirty) {
			changes.previousmetadataValue.copy(m_metadata_prev);
			QSqlQuery query;
			auto tq = QStringLiteral(R"RJIENRLWEY( UPDATE Note SET metadata = :val WHERE ID = :id )RJIENRLWEY");
			query.prepare(tq);
			query.bindValue(":val", QVariant::fromValue(m_metadata));
			query.bindValue(":id", QVariant::fromValue(m_ID));
			auto res = query.exec();
			if (!res) {
				qCritical() << query.lastError() << "when updating an item of type Note at row metadata";
			}
			m_metadata_dirty = false;
		}
		
		m_UNDO_STACK << changes;
		pUR->undoItemAdded(this);
		evaluate_can_undo_changed();
		}
	}

	
	Q_INVOKABLE QList<QSharedPointer<Note>> childNotes() {
		auto tq = QStringLiteral("SELECT * FROM Note WHERE PARENT_Note_ID = :parent_id");
		QSqlQuery query;
		query.prepare(tq);
		query.bindValue(":parent_id", m_ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when loading an Note children of a Note";
		}
		QList<QSharedPointer<Note>> ret;
		
		while (query.next()) {
			auto add = Note::withID(query.value("ID").value<QUuid>());
		
			add->setProperty("title", query.value("title"));
		
			add->setProperty("metadata", query.value("metadata"));
		
			ret << add;
		}
		return ret;
	}
	Q_INVOKABLE void addChildNote(QSharedPointer<Note> child) {
		auto tq = QStringLiteral("UPDATE Note SET PARENT_Note_ID = :new_parent_id WHERE ID = :child_id ");
		QSqlQuery query;
		query.prepare(tq);
		query.bindValue(":new_parent_id", m_ID);
		query.bindValue(":child_id", child->m_ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when adding a new Note to a parent Note";
		}
		child->m_parent_Note_ID = m_ID;
	}
	Q_INVOKABLE void removeChildNote(QSharedPointer<Note> child) {
		auto tq = QStringLiteral("UPDATE Note SET PARENT_Note_ID = NULL WHERE ID = :child_id ");
		QSqlQuery query;
		query.prepare(tq);
		query.bindValue(":child_id", child->m_ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when removing a Note from a parent Note";
		}
		child->m_parent_Note_ID = QUuid();
	}
	

	static QSharedPointer<Note> newNote() {
		auto ret = Note::withID(QUuid::createUuid());
		ret->m_NEW = true;
		return ret;
	}

	static QSharedPointer<Note> load(const QUuid& ID) {
		auto tq = QStringLiteral("SELECT * FROM Note WHERE ID = :id");
		QSqlQuery query;
		query.prepare(tq);
		query.bindValue(":id", ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when loading an item of type Note";
		}
		auto ret = Note::withID(ID);
		while (query.next()) {
			ret->setProperty("title", query.value("title"));
			ret->setProperty("metadata", query.value("metadata"));
			
		}
		return ret;
	}

	static QList<QSharedPointer<Note>> where(PredicateList predicates) {
		auto tq = QStringLiteral("SELECT * FROM Note WHERE %1").arg(predicates.allPredicatesToWhere().join(","));
		QSqlQuery query;
		query.prepare(tq);
		predicates.bindAllPredicates(&query);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when running a where query on items of type Note";
		}
		QList<QSharedPointer<Note>> ret;
		while (query.next()) {
			auto add = Note::withID(query.value("ID").value<QUuid>());
			
			add->setProperty("title", query.value("title"));
			
			add->setProperty("metadata", query.value("metadata"));
			
			ret << add;
		}
		return ret;
	}

	static void prepareDatabase() {
		auto tq = QStringLiteral(R"RJIENRLWEY(
		CREATE TABLE IF NOT EXISTS Note(
			ID BLOB NOT NULL,
			
			PARENT_Note_ID BLOB,
			
			title TEXT NOT NULL,
			metadata BLOB NOT NULL,
			PRIMARY KEY (ID))
		)RJIENRLWEY");
		QSqlQuery query;
		query.prepare(tq);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError();
		}
	}

};

class NoteModel : public QAbstractListModel {
	Q_OBJECT

	mutable QSqlQuery m_query;
	static const int fetch_size = 255;
	int m_rowCount = 0;
	int m_bottom = 0;
	bool m_atEnd = false;
	QUuid m_parentID;
	mutable QMap<int,QSharedPointer<Note>> m_items;
	QSharedPointer<Note> m_staging;

	Q_PROPERTY(Note* staging READ staging NOTIFY stagingItemChanged)

	void prefetch(int toRow) {
		if (m_atEnd || toRow <= m_bottom)
			return;

		int oldBottom = m_bottom;
		int newBottom = 0;
		
		if (m_query.seek(toRow)) {
			newBottom = toRow;
		} else {
			int i =	oldBottom;
			if (m_query.seek(i)) {
				while (m_query.next()) i++;
				newBottom = i;
			} else {
				newBottom = -1;
			}
			m_atEnd = true;
		}
		if (newBottom >= 0 && newBottom >= oldBottom) {
			beginInsertRows(QModelIndex(), oldBottom + 1, newBottom);
			m_bottom = newBottom;
			endInsertRows();
		}
	}

public:

	Q_SIGNAL void stagingItemChanged();

	enum NoteData {
		title = Qt::UserRole,
		metadata ,
		
		childrenNote,
		
		object
	};

	Note* staging() const {
		return m_staging.data();
	}

	Q_INVOKABLE void createStaging() {
		m_staging = Note::newNote();
		Q_EMIT stagingItemChanged();
	}

	Q_INVOKABLE void commitStaging() {
		m_staging->save();
		prefetch(fetch_size);
		m_staging = nullptr;
		Q_EMIT stagingItemChanged();
	}

	NoteModel(QObject *parent = nullptr) : QAbstractListModel(parent)
	{
		m_query.prepare("SELECT * FROM Note");
		m_query.exec();
		prefetch(fetch_size);
	}

	
	static NoteModel* withNoteParent(const QUuid& id) {
		static QMap<QUuid,QPointer<NoteModel>> s_models;
		if (s_models.value(id).isNull()) {
			auto childModel = new NoteModel();
			childModel->m_query.prepare("SELECT * FROM Note WHERE PARENT_Note_ID = :parent_id");
			childModel->m_query.bindValue(":parent_id", id);
			childModel->m_bottom = 0;
			childModel->m_parentID = id;
			childModel->m_query.exec();
			childModel->prefetch(fetch_size);
			s_models[id] = childModel;
		}
		return s_models[id].data();
	}
	

	void fetchMore(const QModelIndex &parent) override {
		prefetch(m_bottom + fetch_size);
	}

	bool canFetchMore(const QModelIndex &parent) const override {
		return !m_atEnd;
	}

	int rowCount(const QModelIndex &parent = QModelIndex()) const override {
		Q_UNUSED(parent)
		return m_bottom;
	}

	QHash<int, QByteArray> roleNames() const {
		auto rn = QAbstractItemModel::roleNames();
		rn[NoteData::title] = QByteArray("title");
		rn[NoteData::metadata] = QByteArray("metadata");
		
		rn[NoteData::childrenNote] = QByteArray("children-Note");
		
		rn[NoteData::object] = QByteArray("Note-object");
		return rn;
	}

	QVariant data(const QModelIndex &item, int role) const override {
		if (!item.isValid()) return QVariant();

		if (!m_items.contains(item.row())) {
			if (!m_query.seek(item.row())) {
				qCritical() << m_query.lastError() << "when seeking data for Note";
				return QVariant();
			}

			auto add = Note::withID(m_query.value("ID").value<QUuid>());
			
			add->setProperty("title", m_query.value("title"));
			
			add->setProperty("metadata", m_query.value("metadata"));
			
			m_items.insert(item.row(), add);
		}

		switch (role) {
		case NoteData::title:
			return QVariant::fromValue(m_items[item.row()]->title());
		case NoteData::metadata:
			return QVariant::fromValue(m_items[item.row()]->metadata());
		
		case NoteData::childrenNote:
			return QVariant::fromValue(NoteModel::withNoteParent(m_items[item.row()]->m_ID));
		
		case NoteData::object:
			return QVariant::fromValue(m_items[item.row()]);
		}

		return QVariant();
	}

	Qt::ItemFlags flags(const QModelIndex &index) const {
		return Qt::ItemIsEditable | Qt::ItemIsSelectable | Qt::ItemIsEnabled;
	}

	bool setData(const QModelIndex &item, const QVariant &value, int role = Qt::EditRole) override {
		if (!m_items.contains(item.row())) {
			if (!m_query.seek(item.row())) {
				qCritical() << m_query.lastError() << "when seeking data for Note";
				return false;
			}

			auto add = Note::withID(m_query.value("ID").value<QUuid>());
			
			add->setProperty("title", m_query.value("title"));
			
			add->setProperty("metadata", m_query.value("metadata"));
			
			m_items.insert(item.row(), add);
		}

		switch (role) {
			
			
			case NoteData::title:
				m_items[item.row()]->set_title(value.value<QString>());
				Q_EMIT dataChanged(item, item, {role});
				return true;
			
			
			case NoteData::metadata:
				m_items[item.row()]->set_metadata(value.value<QMap<QString,QString>>());
				Q_EMIT dataChanged(item, item, {role});
				return true;
			
		}

		return false;
	}
};

