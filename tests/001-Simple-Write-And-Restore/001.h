

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

#include "Database.h"

enum ModelTypes {
	ItemKind,
	};
class Item;
class ItemModel;


class Item : public QObject, PPUndoRedoable {
	Q_OBJECT
	Q_INTERFACES(PPUndoRedoable)

	struct Change {
		
		
		
		Optional<QString> previouspropValue;
		
	};

	Item(QUuid ID) : QObject(nullptr), m_ID(ID) {
		static bool db_initialized = false;
		if (!db_initialized) {
			prepareDatabase();
			db_initialized = true;
		}
	}

	~Item() {
		if (m_DELETE_PENDING) {
			QSqlQuery query(PPDatabase::instance()->connection());
			query.prepare(QStringLiteral("DELETE FROM Item WHERE ID = :ID"));
			query.bindValue(":ID", QVariant::fromValue(m_ID));
			query.exec();
		}
	}

	static QSharedPointer<Item> withID(QUuid ID) {
		static QMap<QUuid,QPointer<Item>> s_instances;
		static QMutex s_mutex;

		QMutexLocker locker(&s_mutex);
		auto val = s_instances.value(ID, nullptr);
		if (val.isNull()) {
			s_instances[ID] = new Item(ID);
		}
		return QSharedPointer<Item>(s_instances[ID].data(), &QObject::deleteLater);
	}

	
	friend class ItemModel;

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

	
	
	
	Q_PROPERTY(QString prop READ prop WRITE set_prop NOTIFY propChanged)
	QString m_prop;
	QString m_prop_prev;
	bool m_prop_dirty;
	

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
			
			if (m_prop_dirty) {
				new_dirty = true;
			}
			
			if (!new_dirty) {
				m_DIRTY = false;
			}
			Q_EMIT dirtyChanged();
		} else {
			
			if (m_prop_dirty) {
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
			
			if (last.previouspropValue.has_value()) {
				last.previouspropValue.swap(m_prop);
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
			
			if (last.previouspropValue.has_value()) {
				last.previouspropValue.swap(m_prop);
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

	
	
	
	Q_SIGNAL void propChanged();
	QString prop() const { return m_prop; };
	void set_prop(const QString& val) {
		if (val == m_prop) {
			return;
		}
		m_prop_prev = m_prop;
		m_prop_dirty = true;
		m_prop = val;
		Q_EMIT void propChanged();
		clear_redo();
		evaluate_dirty_changed();
		evaluate_can_undo_changed();
	}
	void discard_prop_changes() {
		if (m_prop_dirty) {
			m_prop_dirty = false;
			m_prop = m_prop_prev;
			Q_EMIT void propChanged();
			evaluate_dirty_changed();
		}
	}
	

	void discard_all_changes() {
		
		if (m_prop_dirty) {
			m_prop_dirty = false;
			m_prop = m_prop_prev;
			Q_EMIT void propChanged();
		}
		
		evaluate_dirty_changed();
	}

	Q_INVOKABLE void save() {
		if (m_NEW || m_DELETE_PENDING) {
			auto tq = QStringLiteral(R"RJIENRLWEY(
INSERT INTO Item
(ID,prop)
VALUES
(:ID,   :prop );
			)RJIENRLWEY");
			QSqlQuery query(PPDatabase::instance()->connection());
			query.prepare(tq);
			query.bindValue(":ID", QVariant::fromValue(m_ID));
			query.bindValue(":prop", QVariant::fromValue(m_prop));
			auto res = query.exec();
			if (!res) {
				qCritical() << query.lastError() << "when creating a new item of Item";
			}
			m_NEW = false;
			if (m_DELETE_PENDING) {
				m_DELETE_PENDING = false;
				pendingDeleteChanged();
			}
		} else {
		Change changes;
		
		if (m_prop_dirty) {
			changes.previouspropValue.copy(m_prop_prev);
			QSqlQuery query(PPDatabase::instance()->connection());
			auto tq = QStringLiteral(R"RJIENRLWEY( UPDATE Item SET prop = :val WHERE ID = :id )RJIENRLWEY");
			query.prepare(tq);
			query.bindValue(":val", QVariant::fromValue(m_prop));
			query.bindValue(":id", QVariant::fromValue(m_ID));
			auto res = query.exec();
			if (!res) {
				qCritical() << query.lastError() << "when updating an item of type Item at row prop";
			}
			m_prop_dirty = false;
		}
		
		m_UNDO_STACK << changes;
		pUR->undoItemAdded(this);
		evaluate_can_undo_changed();
		}
	}

	

	static QSharedPointer<Item> newItem() {
		auto ret = Item::withID(QUuid::createUuid());
		ret->m_NEW = true;
		return ret;
	}

	static QSharedPointer<Item> load(const QUuid& ID) {
		auto tq = QStringLiteral("SELECT * FROM Item WHERE ID = :id");
		QSqlQuery query(PPDatabase::instance()->connection());
		query.prepare(tq);
		query.bindValue(":id", ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when loading an item of type Item";
		}
		auto ret = Item::withID(ID);
		while (query.next()) {
			ret->setProperty("prop", query.value("prop"));
			
		}
		return ret;
	}

	static QList<QSharedPointer<Item>> where(PredicateList predicates) {
		auto tq = QStringLiteral("SELECT * FROM Item WHERE %1").arg(predicates.allPredicatesToWhere().join(","));
		QSqlQuery query(PPDatabase::instance()->connection());
		query.prepare(tq);
		predicates.bindAllPredicates(&query);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when running a where query on items of type Item";
		}
		QList<QSharedPointer<Item>> ret;
		while (query.next()) {
			auto add = Item::withID(query.value("ID").value<QUuid>());
			
			add->setProperty("prop", query.value("prop"));
			
			ret << add;
		}
		return ret;
	}

	static void prepareDatabase() {
		volatile auto db = PPDatabase::instance();
		Q_UNUSED(db)

		auto tq = QStringLiteral(R"RJIENRLWEY(
		CREATE TABLE IF NOT EXISTS Item(
			ID BLOB NOT NULL,
			
			prop TEXT NOT NULL,
			PRIMARY KEY (ID))
		)RJIENRLWEY");
		QSqlQuery query(PPDatabase::instance()->connection());
		query.prepare(tq);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError();
		}
	}

};

class ItemModel : public QAbstractListModel {
	Q_OBJECT

	mutable QSqlQuery m_query = QSqlQuery(PPDatabase::instance()->connection());
	static const int fetch_size = 255;
	int m_rowCount = 0;
	int m_bottom = 0;
	bool m_atEnd = false;
	QUuid m_parentID;
	mutable QMap<int,QSharedPointer<Item>> m_items;
	QSharedPointer<Item> m_staging;
	ModelTypes m_parentedKind;

	Q_PROPERTY(Item* staging READ staging NOTIFY stagingItemChanged)

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
				do {
					i++;
				} while (m_query.next());
				newBottom = i;
			} else {
				newBottom = -1;
			}
			m_atEnd = true;
		}
		if (newBottom >= 0 && newBottom >= oldBottom) {
			beginInsertRows(QModelIndex(), oldBottom, newBottom - 1);
			m_bottom = newBottom;
			endInsertRows();
		}
	}

public:

	Q_SIGNAL void stagingItemChanged();

	enum ItemData {
		prop = Qt::UserRole,
		
		
		object
	};

	Item* staging() const {
		return m_staging.data();
	}

	Q_INVOKABLE void createStaging() {
		m_staging = Item::newItem();
		Q_EMIT stagingItemChanged();
	}

	Q_INVOKABLE void commitStaging() {
		m_staging->save();
		if (!m_parentID.isNull()) {
			
		}
		m_query.exec();
		m_atEnd = false;
		prefetch(fetch_size);
		m_staging = nullptr;
		Q_EMIT stagingItemChanged();
	}

	ItemModel(QObject *parent = nullptr) : QAbstractListModel(parent)
	{
		m_query.prepare("SELECT * FROM Item");
		m_query.exec();
		prefetch(fetch_size);
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

	QHash<int, QByteArray> roleNames() const override {
		auto rn = QAbstractItemModel::roleNames();
		rn[ItemData::prop] = QByteArray("prop");
		
		
		rn[ItemData::object] = QByteArray("Item-object");
		return rn;
	}

	QVariant data(const QModelIndex &item, int role) const override {
		if (!item.isValid()) return QVariant();

		if (!m_items.contains(item.row())) {
			if (!m_query.seek(item.row())) {
				qCritical() << m_query.lastError() << "when seeking data for Item";
				return QVariant();
			}

			auto add = Item::withID(m_query.value("ID").value<QUuid>());
			
			add->setProperty("prop", m_query.value("prop"));
			
			m_items.insert(item.row(), add);
		}

		switch (role) {
		case ItemData::prop:
			return QVariant::fromValue(m_items[item.row()]->prop());
		
		
		case ItemData::object:
			return QVariant::fromValue(m_items[item.row()].data());
		}

		return QVariant();
	}

	Qt::ItemFlags flags(const QModelIndex &index) const override {
		return Qt::ItemIsEditable | Qt::ItemIsSelectable | Qt::ItemIsEnabled;
	}

	bool setData(const QModelIndex &item, const QVariant &value, int role = Qt::EditRole) override {
		if (!m_items.contains(item.row())) {
			if (!m_query.seek(item.row())) {
				qCritical() << m_query.lastError() << "when seeking data for Item";
				return false;
			}

			auto add = Item::withID(m_query.value("ID").value<QUuid>());
			
			add->setProperty("prop", m_query.value("prop"));
			
			m_items.insert(item.row(), add);
		}

		switch (role) {
			
			
			case ItemData::prop:
				m_items[item.row()]->set_prop(value.value<QString>());
				Q_EMIT dataChanged(item, item, {role});
				return true;
			
		}

		return false;
	}
};

