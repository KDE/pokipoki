package parser

import (
	"log"
	"strings"
	"text/template"
)

var sqlKind = map[string]string{
	"String":   "TEXT",
	"Date":     "DATE",
	"DateTime": "DATETIME",
}

func SqlType(typeDef []string) string {
	if val, ok := sqlKind[typeDef[0]]; ok {
		return val
	}
	return "BLOB"
}

var imports = map[string]string{
	"QBitArray":          "QBitArray",
	"QBrush":             "QBrush",
	"QByteArray":         "QByteArray",
	"QColor":             "QColor",
	"QCursor":            "QCursor",
	"QDate":              "QDate",
	"QDateTime":          "QDateTime",
	"QEasingCurve":       "QEasingCurve",
	"QFont":              "QFont",
	"QGenericMatrix":     "QGenericMatrix",
	"QIcon":              "QIcon",
	"QImage":             "QImage",
	"QKeySequence":       "QKeySequence",
	"QMargins":           "QMargins",
	"QMatrix4x4":         "QMatrix4x4",
	"QPalette":           "QPalette",
	"QPen":               "QPen",
	"QPicture":           "QPicture",
	"QPixmap":            "QPixmap",
	"QPoint":             "QPoint",
	"QQuaternion":        "QQuaternion",
	"QRect":              "QRect",
	"QRegExp":            "QRegExp",
	"QRegularExpression": "QRegularExpression",
	"QRegion":            "QRegion",
	"QSize":              "QSize",
	"QString":            "QString",
	"QTime":              "QTime",
	"QTransform":         "QTransform",
	"QUrl":               "QUrl",
	"QVariant":           "QVariant",
	"QVector2D":          "QVector2D",
	"QVector3D":          "QVector3D",
	"QVector4D":          "QVector4D",
	"QLinkedList":        "QLinkedList",
	"QList":              "QList",
	"QVector":            "QVector",
	"QHash":              "QHash",
	"QMap":               "QMap",
	"QPair":              "QPair",
}

func (d PokiPokiDocument) LocateImports() []string {
	ret := map[string]struct{}{}
	for _, obj := range d.Objects {
		for _, prop := range obj.Properties {
			propType := d.AlwaysType(prop.Type)
			for _, chunk := range propType {
				if val, ok := imports[chunk]; ok {
					ret[val] = struct{}{}
				}
			}
		}
	}
	var arrRet []string
	for key := range ret {
		arrRet = append(arrRet, "#include <"+key+">")
	}
	return arrRet
}

var scalarTypes = map[string]string{
	"Boolean":           "bool",
	"Int8":              "qint8",
	"Int16":             "qint16",
	"Int32":             "qint32",
	"Int64":             "qint64",
	"Uint8":             "quint8",
	"Uint16":            "quint16",
	"Uint32":            "quint32",
	"Uint64":            "quint64",
	"Float32":           "float",
	"Float64":           "double",
	"ConstChar":         "const char *",
	"BitArray":          "QBitArray",
	"Brush":             "QBrush",
	"ByteArray":         "QByteArray",
	"Color":             "QColor",
	"Cursor":            "QCursor",
	"Date":              "QDate",
	"DateTime":          "QDateTime",
	"EasingCurve":       "QEasingCurve",
	"Font":              "QFont",
	"GenericMatrix":     "QGenericMatrix",
	"Icon":              "QIcon",
	"Image":             "QImage",
	"KeySequence":       "QKeySequence",
	"Margins":           "QMargins",
	"Matrix4x4":         "QMatrix4x4",
	"Palette":           "QPalette",
	"Pen":               "QPen",
	"Picture":           "QPicture",
	"Pixmap":            "QPixmap",
	"Point":             "QPoint",
	"Quaternion":        "QQuaternion",
	"Rect":              "QRect",
	"RegExp":            "QRegExp",
	"RegularExpression": "QRegularExpression",
	"Region":            "QRegion",
	"Size":              "QSize",
	"String":            "QString",
	"Time":              "QTime",
	"Transform":         "QTransform",
	"URL":               "QUrl",
	"Variant":           "QVariant",
	"Vector2D":          "QVector2D",
	"Vector3D":          "QVector3D",
	"Vector4D":          "QVector4D",
}

// ScalarType returns the typedef for a scalar type
func (d PokiPokiDocument) ScalarType(typeDef []string) ([]string, bool) {
	if len(typeDef) > 1 {
		return []string{}, false
	}
	if val, ok := scalarTypes[typeDef[0]]; ok {
		return []string{val}, true
	}
	return []string{}, false
}

var singularGenericTypes = map[string]string{
	"LinkedList": "QLinkedList",
	"List":       "QList",
	"Vector":     "QVector",
}

// SingularGenericTypes returns the typedef for a singular generic type
func (d PokiPokiDocument) SingularGenericTypes(typeDef []string) ([]string, bool) {
	if _, ok := singularGenericTypes[typeDef[0]]; !ok {
		return []string{}, false
	}
	if typeDef[len(typeDef)-1] != "]" {
		return []string{}, false
	}
	inner, ok := d.Type(typeDef[1 : len(typeDef)-1])
	if !ok {
		return []string{}, false
	}
	ret := []string{singularGenericTypes[typeDef[0]], "<"}
	ret = append(ret, inner...)
	ret = append(ret, ">")
	return ret, true
}

var dualGenericTypes = map[string]string{
	"Hash": "QHash",
	"Map":  "QMap",
	"Pair": "QPair",
}

// DualGenericTypes returns the typedef for a dual generic type
func (d PokiPokiDocument) DualGenericTypes(typeDef []string) ([]string, bool) {
	if _, ok := dualGenericTypes[typeDef[0]]; !ok {
		return []string{}, false
	}
	var first []string
	var second []string

	opener := 1

	for _, str := range typeDef[2:] {
		if opener == 0 {
			second = append(second, str)
			continue
		}
		if str == "[" {
			first = append(first, "[")
			opener++
			continue
		}
		if str == "]" {
			opener--
			if opener != 0 {
				first = append(first, "]")
			}
			continue
		}
		first = append(first, str)
	}

	prefix := dualGenericTypes[typeDef[0]]

	firstInner, ok := d.Type(first)
	if !ok {
		return []string{}, false
	}
	secondInner, ok := d.Type(second)
	if !ok {
		return []string{}, false
	}

	ret := []string{prefix, "<"}
	ret = append(ret, firstInner...)
	ret = append(ret, ",")
	ret = append(ret, secondInner...)
	ret = append(ret, ">")

	return ret, true
}

// Type returns whether the type is valid and its C++ representation
func (d PokiPokiDocument) Type(typeDef []string) ([]string, bool) {
	if len(typeDef) == 0 {
		return []string{}, false
	}
	if val, ok := d.ScalarType(typeDef); ok {
		return val, ok
	}
	if val, ok := d.SingularGenericTypes(typeDef); ok {
		return val, ok
	}
	if val, ok := d.DualGenericTypes(typeDef); ok {
		return val, ok
	}
	return []string{}, false
}

// AlwaysType is Type but it only has one return value
func (d PokiPokiDocument) AlwaysType(typeDef []string) []string {
	ret, _ := d.Type(typeDef)
	return ret
}

// Verify verifies that a PokiPokiDocument is valid
func (d PokiPokiDocument) Verify() {
	for _, obj := range d.Objects {
	childrenLoop:
		for _, child := range obj.Children {
			for kind := range d.Objects {
				if kind == child {
					continue childrenLoop
				}
			}
			log.Fatalf("Unknown type '%s'", child)
		}
		for _, prop := range obj.Properties {
			if _, ok := d.Type(prop.Type); !ok {
				log.Fatalf("Unknown type '%s'", strings.Join(prop.Type, ""))
			}
		}
	}
}

// ParentedBy
func (d PokiPokiDocument) ParentedBy(typ string) []string {
	ret := map[string]struct{}{}

	for _, obj := range d.Objects {
		for _, child := range obj.Children {
			if child == typ {
				ret[obj.Name] = struct{}{}
			}
		}
	}

	arrRet := []string{}
	for key := range ret {
		arrRet = append(arrRet, key)
	}

	return arrRet
}

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

class {{ .Name }} : public QObject {
	Q_OBJECT

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
	bool m_DIRTY = false;
	bool dirty() const { return m_DIRTY; }
	bool m_CAN_UNDO = false;
	bool canUndo() const { return m_CAN_UNDO; }
	bool m_CAN_REDO = false;
	bool canRedo() const { return m_CAN_REDO; }

	Q_PROPERTY(bool dirty READ dirty NOTIFY dirtyChanged)
	Q_SIGNAL void dirtyChanged();

	Q_PROPERTY(bool canUndo READ canUndo NOTIFY canUndoChanged)
	Q_SIGNAL void canUndoChanged();

	Q_PROPERTY(bool canRedo READ canRedo NOTIFY canRedoChanged)
	Q_SIGNAL void canRedoChanged();

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
		if (m_DIRTY) {
			setUndo(true);
			return;
		}
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
		if (m_DIRTY) {
			setRedo(false);
			return;
		}
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

	Q_INVOKABLE void undo() {
		if (m_DIRTY) {
			discard_all_changes();
			evaluate_can_undo_changed();
			return;
		}
		if (!m_UNDO_STACK.empty()) {
			auto last = m_UNDO_STACK.takeLast();
			{{ range $prop := .Properties }}
			if (last.previous{{ $prop.Name }}Value.has_value()) {
				last.previous{{ $prop.Name }}Value.swap(m_{{$prop.Name}});
			}
			{{ end }}
			m_REDO_STACK << last;
			evaluate_can_undo_changed();
			evaluate_can_redo_changed();
		}
	}

	Q_INVOKABLE void redo() {
		if (!m_REDO_STACK.empty()) {
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

	void save() {
		if (m_NEW) {
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
		evaluate_can_undo_changed();
		}
	}

	{{ range $child := .Children }}
	QList<QSharedPointer<{{ $child }}>> child{{ $child }}s() {
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
	void addChild{{ $child }}(QSharedPointer<{{ $child }}> child) {
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
	void removeChild{{ $child }}(QSharedPointer<{{ $child }}> child) {
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
	mutable QSqlQuery m_query;
	static const int fetch_size = 255;
	int m_rowCount = 0;
	int m_bottom = 0;
	bool m_atEnd = false;
	mutable QMap<int,QSharedPointer<{{.Name}}>> m_items;

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

	enum {{ .Name }}Data {
		{{ range $index, $prop := .Properties -}}
		{{ $prop.Name }} {{ if eq $index 0 }}= Qt::UserRole{{ end }},
		{{ end }}
		{{ range $child := .Children -}}
		children{{ $child }},
		{{ end }}
	};

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

	QHash<int, QByteArray> roleNames() const {
		auto rn = QAbstractItemModel::roleNames();
		{{ range $index, $prop := .Properties -}}
		rn[{{ $item.Name }}Data::{{ $prop.Name }}] = QByteArray("{{ $prop.Name }}");
		{{ end }}
		{{ range $child := .Children -}}
		rn[{{ $item.Name }}Data::children{{ $child }}] = QByteArray("children-{{ $child }}");
		{{ end }}
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
		}

		return QVariant();
	}

	Qt::ItemFlags flags(const QModelIndex &index) const {
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
