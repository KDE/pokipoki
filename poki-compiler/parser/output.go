package parser

import (
	"log"
	"strings"
	"text/template"
)

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"StringJoin": strings.Join,
	"TypeDef":    SqlType,
}).Parse(`
{{ $root := . }}
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
{{ StringJoin $root.LocateImports "\n" }}

#include "Database.h"

{{- range $item := .Objects }}

class {{ .Name }}Model;

class {{ .Name }} : public QObject, PPUndoRedoable {
	Q_OBJECT
	Q_INTERFACES(PPUndoRedoable)

	struct Change {
		{{ range $prop := .Properties }}
		{{ $propType := $root.AlwaysType $prop.Type }}
		{{ $propTypeName := StringJoin $propType "" }}
		Optional<{{ $propTypeName }}> previous{{ $prop.Name }}Value;
		{{ end }}
	};

	{{ .Name }}(QUuid ID) : QObject(nullptr), m_ID(ID) {
		static bool db_initialized = false;
		if (!db_initialized) {
			volatile auto db = PPDatabase::instance();
			Q_UNUSED(db)
			prepareDatabase();
			db_initialized = true;
		}
	}

	~{{ .Name }}() {
		if (m_DELETE_PENDING) {
			QSqlQuery query;
			query.prepare(QStringLiteral("DELETE FROM {{ .Name }} WHERE ID = :ID"));
			query.bindValue(":ID", QVariant::fromValue(m_ID));
			query.exec();
		}
	}

	static QSharedPointer<{{ .Name }}> withID(QUuid ID) {
		static QMap<QUuid,QPointer<{{.Name}}>> s_instances;
		static QMutex s_mutex;

		QMutexLocker locker(&s_mutex);
		auto val = s_instances.value(ID, nullptr);
		if (val.isNull()) {
			s_instances[ID] = new {{.Name}}(ID);
		}
		return QSharedPointer<{{ .Name }}>(s_instances[ID].data(), &QObject::deleteLater);
	}

	{{ range $parent := $root.ParentedBy .Name }}
	QUuid m_parent_{{ $parent }}_ID;
	{{ if ne $parent $item.Name }}
	friend class {{ $parent }};
	{{ end }}
	{{ end }}
	friend class {{ .Name }}Model;

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

	{{ range $prop := .Properties }}
	{{ $propType := $root.AlwaysType $prop.Type }}
	{{ $propTypeName := StringJoin $propType "" }}
	Q_PROPERTY({{ $propTypeName }} {{ $prop.Name }} READ {{ $prop.Name }} WRITE set_{{ $prop.Name }} NOTIFY {{$prop.Name}}Changed)
	{{ $propTypeName }} m_{{$prop.Name}};
	{{ $propTypeName }} m_{{$prop.Name}}_prev;
	bool m_{{$prop.Name}}_dirty;
	{{ end }}

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
			{{ range $prop := .Properties }}
			if (m_{{$prop.Name}}_dirty) {
				new_dirty = true;
			}
			{{ end }}
			if (!new_dirty) {
				m_DIRTY = false;
			}
			Q_EMIT dirtyChanged();
		} else {
			{{ range $prop := .Properties }}
			if (m_{{$prop.Name}}_dirty) {
				m_DIRTY = true;
				Q_EMIT dirtyChanged();
				return;
			}
			{{ end }}
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
			{{ range $prop := .Properties }}
			if (last.previous{{ $prop.Name }}Value.has_value()) {
				last.previous{{ $prop.Name }}Value.swap(m_{{$prop.Name}});
			}
			{{ end }}
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
			{{ range $prop := .Properties }}
			if (last.previous{{ $prop.Name }}Value.has_value()) {
				last.previous{{ $prop.Name }}Value.swap(m_{{$prop.Name}});
			}
			{{ end }}
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

	{{ range $prop := .Properties }}
	{{ $propType := $root.AlwaysType $prop.Type }}
	{{ $propTypeName := StringJoin $propType "" }}
	Q_SIGNAL void {{$prop.Name}}Changed();
	{{ $propTypeName }} {{$prop.Name}}() const { return m_{{$prop.Name}}; };
	void set_{{$prop.Name}}(const {{ $propTypeName }}& val) {
		if (val == m_{{$prop.Name}}) {
			return;
		}
		m_{{$prop.Name}}_prev = m_{{$prop.Name}};
		m_{{$prop.Name}}_dirty = true;
		m_{{$prop.Name}} = val;
		Q_EMIT void {{$prop.Name}}Changed();
		clear_redo();
		evaluate_dirty_changed();
		evaluate_can_undo_changed();
	}
	void discard_{{$prop.Name}}_changes() {
		if (m_{{$prop.Name}}_dirty) {
			m_{{$prop.Name}}_dirty = false;
			m_{{$prop.Name}} = m_{{$prop.Name}}_prev;
			Q_EMIT void {{$prop.Name}}Changed();
			evaluate_dirty_changed();
		}
	}
	{{ end }}

	void discard_all_changes() {
		{{ range $prop := .Properties }}
		if (m_{{$prop.Name}}_dirty) {
			m_{{$prop.Name}}_dirty = false;
			m_{{$prop.Name}} = m_{{$prop.Name}}_prev;
			Q_EMIT void {{$prop.Name}}Changed();
		}
		{{ end }}
		evaluate_dirty_changed();
	}

	Q_INVOKABLE void save() {
		if (m_NEW || m_DELETE_PENDING) {
			auto tq = QStringLiteral(R"RJIENRLWEY(
INSERT INTO {{ $item.Name }}
(ID,
{{- range $index, $prop := .Properties -}}
{{- if $index -}},{{- end -}}
{{- $prop.Name -}}
{{- end -}}
)
VALUES
(:ID, {{ range $index, $prop := .Properties }} {{ if $index }},{{ end }} :{{- $prop.Name }} {{ end }});
			)RJIENRLWEY");
			QSqlQuery query;
			query.prepare(tq);
			query.bindValue(":ID", QVariant::fromValue(m_ID));
			{{- range $prop := .Properties }}
			query.bindValue(":{{- $prop.Name -}}", QVariant::fromValue(m_{{$prop.Name}}));
			{{ end -}}
			auto res = query.exec();
			if (!res) {
				qCritical() << query.lastError() << "when creating a new item of {{ $item.Name }}";
			}
			m_NEW = false;
			if (m_DELETE_PENDING) {
				m_DELETE_PENDING = false;
				pendingDeleteChanged();
			}
		} else {
		Change changes;
		{{ range $prop := .Properties }}
		if (m_{{$prop.Name}}_dirty) {
			changes.previous{{$prop.Name}}Value.copy(m_{{ $prop.Name }}_prev);
			QSqlQuery query;
			auto tq = QStringLiteral(R"RJIENRLWEY( UPDATE {{ $item.Name}} SET {{$prop.Name}} = :val WHERE ID = :id )RJIENRLWEY");
			query.prepare(tq);
			query.bindValue(":val", QVariant::fromValue(m_{{$prop.Name}}));
			query.bindValue(":id", QVariant::fromValue(m_ID));
			auto res = query.exec();
			if (!res) {
				qCritical() << query.lastError() << "when updating an item of type {{ $item.Name }} at row {{ $prop.Name }}";
			}
			m_{{$prop.Name}}_dirty = false;
		}
		{{ end }}
		m_UNDO_STACK << changes;
		pUR->undoItemAdded(this);
		evaluate_can_undo_changed();
		}
	}

	{{ range $child := .Children }}
	Q_INVOKABLE QList<QSharedPointer<{{ $child }}>> child{{ $child }}s() {
		auto tq = QStringLiteral("SELECT * FROM {{ $child }} WHERE PARENT_{{ $item.Name }}_ID = :parent_id");
		QSqlQuery query;
		query.prepare(tq);
		query.bindValue(":parent_id", m_ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when loading an {{ $child }} children of a {{ $item.Name }}";
		}
		QList<QSharedPointer<{{ $child }}>> ret;
		{{ $childKind := index $root.Objects $child }}
		while (query.next()) {
			auto add = {{ $child }}::withID(query.value("ID").value<QUuid>());
		{{ range $prop := $childKind.Properties }}
			add->setProperty("{{ $prop.Name }}", query.value("{{ $prop.Name }}"));
		{{ end }}
			ret << add;
		}
		return ret;
	}
	Q_INVOKABLE void addChild{{ $child }}(QSharedPointer<{{ $child }}> child) {
		auto tq = QStringLiteral("UPDATE {{ $child }} SET PARENT_{{ $item.Name }}_ID = :new_parent_id WHERE ID = :child_id ");
		QSqlQuery query;
		query.prepare(tq);
		query.bindValue(":new_parent_id", m_ID);
		query.bindValue(":child_id", child->m_ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when adding a new {{ $child }} to a parent {{ $item.Name }}";
		}
		child->m_parent_{{ $item.Name }}_ID = m_ID;
	}
	Q_INVOKABLE void removeChild{{ $child }}(QSharedPointer<{{ $child }}> child) {
		auto tq = QStringLiteral("UPDATE {{ $child }} SET PARENT_{{ $item.Name }}_ID = NULL WHERE ID = :child_id ");
		QSqlQuery query;
		query.prepare(tq);
		query.bindValue(":child_id", child->m_ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when removing a {{ $child }} from a parent {{ $item.Name }}";
		}
		child->m_parent_{{ $item.Name }}_ID = QUuid();
	}
	{{ end }}

	static QSharedPointer<{{ .Name }}> new{{ .Name }}() {
		auto ret = {{.Name}}::withID(QUuid::createUuid());
		ret->m_NEW = true;
		return ret;
	}

	static QSharedPointer<{{ .Name }}> load(const QUuid& ID) {
		auto tq = QStringLiteral("SELECT * FROM {{ $item.Name }} WHERE ID = :id");
		QSqlQuery query;
		query.prepare(tq);
		query.bindValue(":id", ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when loading an item of type {{ $item.Name }}";
		}
		auto ret = {{.Name}}::withID(ID);
		while (query.next()) {
			{{ range $prop := .Properties -}}
			ret->setProperty("{{ $prop.Name }}", query.value("{{ $prop.Name }}"));
			{{ end }}
		}
		return ret;
	}

	static QList<QSharedPointer<{{ .Name }}>> where(PredicateList predicates) {
		auto tq = QStringLiteral("SELECT * FROM {{ $item.Name }} WHERE %1").arg(predicates.allPredicatesToWhere().join(","));
		QSqlQuery query;
		query.prepare(tq);
		predicates.bindAllPredicates(&query);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when running a where query on items of type {{ $item.Name }}";
		}
		QList<QSharedPointer<{{ .Name }}>> ret;
		while (query.next()) {
			auto add = {{ .Name }}::withID(query.value("ID").value<QUuid>());
			{{ range $prop := .Properties }}
			add->setProperty("{{ $prop.Name }}", query.value("{{ $prop.Name }}"));
			{{ end }}
			ret << add;
		}
		return ret;
	}

	static void prepareDatabase() {
		auto tq = QStringLiteral(R"RJIENRLWEY(
		CREATE TABLE IF NOT EXISTS {{ $item.Name }}(
			ID BLOB NOT NULL,
			{{ range $parent := $root.ParentedBy .Name }}
			PARENT_{{ $parent }}_ID BLOB,
			{{ end }}
			{{ range $prop := .Properties -}}
			{{ $prop.Name }} {{ TypeDef $prop.Type }} NOT NULL,
			{{ end -}}
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

class {{ .Name }}Model : public QAbstractListModel {
	Q_OBJECT

	mutable QSqlQuery m_query;
	static const int fetch_size = 255;
	int m_rowCount = 0;
	int m_bottom = 0;
	bool m_atEnd = false;
	QUuid m_parentID;
	mutable QMap<int,QSharedPointer<{{.Name}}>> m_items;
	QSharedPointer<{{ .Name }}> m_staging;

	Q_PROPERTY({{ .Name }}* staging READ staging NOTIFY stagingItemChanged)

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

	enum {{ .Name }}Data {
		{{ range $index, $prop := .Properties -}}
		{{ $prop.Name }} {{ if eq $index 0 }}= Qt::UserRole{{ end }},
		{{ end }}
		{{ range $child := .Children -}}
		children{{ $child }},
		{{ end }}
		object
	};

	{{ .Name }}* staging() const {
		return m_staging.data();
	}

	Q_INVOKABLE void createStaging() {
		m_staging = {{ .Name }}::new{{ .Name }}();
		Q_EMIT stagingItemChanged();
	}

	Q_INVOKABLE void commitStaging() {
		m_staging->save();
		prefetch(fetch_size);
		m_staging = nullptr;
		Q_EMIT stagingItemChanged();
	}

	{{.Name}}Model(QObject *parent = nullptr) : QAbstractListModel(parent)
	{
		m_query.prepare("SELECT * FROM {{ .Name }}");
		m_query.exec();
		prefetch(fetch_size);
	}

	{{ range $parent := $root.ParentedBy .Name }}
	static {{ $item.Name }}Model* with{{ $parent }}Parent(const QUuid& id) {
		static QMap<QUuid,QPointer<{{ $item.Name }}Model>> s_models;
		if (s_models.value(id).isNull()) {
			auto childModel = new {{ $item.Name }}Model();
			childModel->m_query.prepare("SELECT * FROM {{ $item.Name }} WHERE PARENT_{{ $parent}}_ID = :parent_id");
			childModel->m_query.bindValue(":parent_id", id);
			childModel->m_bottom = 0;
			childModel->m_parentID = id;
			childModel->m_query.exec();
			childModel->prefetch(fetch_size);
			s_models[id] = childModel;
		}
		return s_models[id].data();
	}
	{{ end }}

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
		{{ range $index, $prop := .Properties -}}
		rn[{{ $item.Name }}Data::{{ $prop.Name }}] = QByteArray("{{ $prop.Name }}");
		{{ end }}
		{{ range $child := .Children -}}
		rn[{{ $item.Name }}Data::children{{ $child }}] = QByteArray("children-{{ $child }}");
		{{ end }}
		rn[{{ $item.Name }}Data::object] = QByteArray("{{ $item.Name }}-object");
		return rn;
	}

	QVariant data(const QModelIndex &item, int role) const override {
		if (!item.isValid()) return QVariant();

		if (!m_items.contains(item.row())) {
			if (!m_query.seek(item.row())) {
				qCritical() << m_query.lastError() << "when seeking data for {{ $item.Name }}";
				return QVariant();
			}

			auto add = {{ .Name }}::withID(m_query.value("ID").value<QUuid>());
			{{ range $prop := .Properties }}
			add->setProperty("{{ $prop.Name }}", m_query.value("{{ $prop.Name }}"));
			{{ end }}
			m_items.insert(item.row(), add);
		}

		switch (role) {
		{{ range $index, $prop := .Properties -}}
		case {{ $item.Name }}Data::{{ $prop.Name }}:
			return QVariant::fromValue(m_items[item.row()]->{{ $prop.Name }}());
		{{ end }}
		{{ range $child := .Children -}}
		case {{ $item.Name }}Data::children{{ $child }}:
			return QVariant::fromValue({{ $child }}Model::with{{ $item.Name }}Parent(m_items[item.row()]->m_ID));
		{{ end }}
		case {{ $item.Name }}Data::object:
			return QVariant::fromValue(m_items[item.row()]);
		}

		return QVariant();
	}

	Qt::ItemFlags flags(const QModelIndex &index) const override {
		return Qt::ItemIsEditable | Qt::ItemIsSelectable | Qt::ItemIsEnabled;
	}

	bool setData(const QModelIndex &item, const QVariant &value, int role = Qt::EditRole) override {
		if (!m_items.contains(item.row())) {
			if (!m_query.seek(item.row())) {
				qCritical() << m_query.lastError() << "when seeking data for {{ $item.Name }}";
				return false;
			}

			auto add = {{ .Name }}::withID(m_query.value("ID").value<QUuid>());
			{{ range $prop := .Properties }}
			add->setProperty("{{ $prop.Name }}", m_query.value("{{ $prop.Name }}"));
			{{ end }}
			m_items.insert(item.row(), add);
		}

		switch (role) {
			{{ range $index, $prop := .Properties -}}
			{{ $propType := $root.AlwaysType $prop.Type }}
			{{ $propTypeName := StringJoin $propType "" }}
			case {{ $item.Name }}Data::{{ $prop.Name }}:
				m_items[item.row()]->set_{{ $prop.Name }}(value.value<{{ $propTypeName }}>());
				Q_EMIT dataChanged(item, item, {role});
				return true;
			{{ end }}
		}

		return false;
	}
};

{{ end -}}
`))

// Output formats a PokiPokiDocument into a C++ file
func (d PokiPokiDocument) Output() string {
	d.Verify()
	var sb strings.Builder
	err := tmpl.Execute(&sb, d)
	if err != nil {
		log.Fatal(err)
	}
	return sb.String()
}
