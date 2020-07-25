package parser

import (
	"log"
	"strings"
	"text/template"
)

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
}).Parse(`
{{ $root := . }}
#include <QObject>
#include <QUuid>
#include <QSqlQuery>
#include <QSqlError>
#include <QVariant>
{{ StringJoin $root.LocateImports "\n" }}

#include "Database.h"

{{- range $item := .Objects }}
class {{ .Name }} : public QObject {
	Q_OBJECT

	{{ .Name }}(QUuid ID) : QObject(nullptr), m_ID(ID) {
		static bool db_initialized = false;
		if (!db_initialized) {
			volatile auto db = PPDatabase::instance();
			prepareDatabase();
			db_initialized = true;
		}
	}

	QUuid m_ID;
	bool m_NEW = false;

	{{ range $prop := .Properties }}
	{{ $propType := $root.AlwaysType $prop.Type }}
	{{ $propTypeName := StringJoin $propType "" }}
	Q_PROPERTY({{ $propTypeName }} {{ $prop.Name }} READ {{ $prop.Name }} WRITE set_{{ $prop.Name }})
	{{ $propTypeName }} m_{{$prop.Name}};
	{{ $propTypeName }} m_{{$prop.Name}}_prev;
	bool m_{{$prop.Name}}_dirty;
	{{ end }}

public:
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
	}
	void discard_{{$prop.Name}}_changes() {
		if (m_{{$prop.Name}}_dirty) {
			m_{{$prop.Name}}_dirty = false;
			m_{{$prop.Name}} = m_{{$prop.Name}}_prev;
			Q_EMIT void {{$prop.Name}}Changed();
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
	}

	void commit() {
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
		{{ range $prop := .Properties }}
		if (m_{{$prop.Name}}_dirty) {
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
			
		}
	}

	static {{ .Name }}* new{{ .Name }}() {
		auto ret = new {{.Name}}(QUuid::createUuid());
		ret->m_NEW = true;
		return ret;
	}

	static {{ .Name }}* load(const QUuid& ID) {
		auto tq = QStringLiteral("SELECT * FROM {{ $item.Name }} WHERE ID = :id");
		QSqlQuery query;
		query.prepare(tq);
		query.bindValue(":id", ID);
		auto ok = query.exec();
		if (!ok) {
			qCritical() << query.lastError() << "when loading an item of type {{ $item.Name }}";
		}
		auto ret = new {{.Name}}(ID);
		while (query.next()) {
			{{ range $prop := .Properties -}}
			ret->setProperty("{{ $prop.Name }}", query.value("{{ $prop.Name }}"));
			{{ end }}
		}
		return ret;
	}

	static void prepareDatabase() {
		auto tq = QStringLiteral(R"RJIENRLWEY(
		CREATE TABLE IF NOT EXISTS {{ $item.Name }}(
			ID BLOB NOT NULL,
			{{ range $prop := .Properties -}}
			{{ $prop.Name }} BLOB NOT NULL,
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
