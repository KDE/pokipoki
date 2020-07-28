package parser

import (
	"log"
	"strings"
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
