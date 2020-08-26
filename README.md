# PokiPoki

formerly fim(t)beb(f)vew(f)seb(p)san(t)ap(t)vlir(t)sang(b)es(p)u(t)vom(b)ngag(t)vlim(p)kay(f)sna(f)kay(f)ga(f)bop(t)veg(p)daf(f)shof(b)*om(p)vlim(p)ga(f)vlim(p)ga(f).

(unspeakably epic storage infinite in number that matches the speaker's expectations that's not food which isn't gendered that will hopefully die by sheep with a market value less than 10 US dollars)

## Persistent Object Framework

PokiPoki is a persistent object framework which allows you to write persistent objects
and maintain a graph between them.

# Using PokiPoki

PokiPoki centralises around `.pokipoki` schema files, which look like this:

```go
schema PokiNotes 1

object Category {
    title       String
    description String
    Note
}

object Note {
    title    String
    metadata Map[String]String
}
```

Each file is given a schema name and version, which are defined as `schema $NAME $VERSION`. `$VERSION` must be an integer.

`$NAME` must be a valid Name, which is an alphabetic string that starts with a capital letter.

Objects are defined in a PokiPoki schema file with the `object $NAME {}` syntax, which defines a unique type. Like schemas, object names must be a valid Name.

Inside an `object` declaration, properties are given as `ident type`, where ident is a valid Identifier. An Identifier is an alphabetic string that starts with a lowercase letter. `type` must be a scalar type, or a compound type. `type` cannot be another object. Objects can have other objects as children, which is indicated with the name of an object on its own line without an identifier.

```go
object Parent {
    Child
}
object Child {
    name String
}
```

## Scalar Types

The following scalar types are recognised with the following names corresponding to the following C++ types:

- `Boolean`: `bool`
- `Int8`: `qint8`
- `Int16`: `qint16`
- `Int32`: `qint32`
- `Int64`: `qint64`
- `Uint8`: `quint8`
- `Uint16`: `quint16`
- `Uint32`: `quint32`
- `Uint64`: `quint64`
- `Float32`: `float`
- `Float64`: `double`
- `ConstChar`: `const char *`
- `BitArray`: `QBitArray`
- `Brush`: `QBrush`
- `ByteArray`: `QByteArray`
- `Color`: `QColor`
- `Cursor`: `QCursor`
- `Date`: `QDate`
- `DateTime`: `QDateTime`
- `EasingCurve`: `QEasingCurve`
- `Font`: `QFont`
- `GenericMatrix`: `QGenericMatrix`
- `Icon`: `QIcon`
- `Image`: `QImage`
- `KeySequence`: `QKeySequence`
- `Margins`: `QMargins`
- `Matrix4x4`: `QMatrix4x4`
- `Palette`: `QPalette`
- `Pen`: `QPen`
- `Picture`: `QPicture`
- `Pixmap`: `QPixmap`
- `Point`: `QPoint`
- `Quaternion`: `QQuaternion`
- `Rect`: `QRect`
- `RegExp`: `QRegExp`
- `RegularExpression`: `QRegularExpression`
- `Region`: `QRegion`
- `Size`: `QSize`
- `String`: `QString`
- `Time`: `QTime`
- `Transform`: `QTransform`
- `URL`: `QUrl`
- `Variant`: `QVariant`
- `Vector2D`: `QVector2D`
- `Vector3D`: `QVector3D`
- `Vector4D`: `QVector4D`

## Singular Generic Types

The following generic types that take one type argument are recognised:

- `LinkedList[T]`: `QLinkedList<T>`
- `List[T]`: `QList<T>`
- `Vector[T]`: `QVector<T>`

## Doubly Generic Types

The following generic types that take two type arguments are recognised:

- `Hash[K]V`: `QHash<K, V>`
- `Map[K]V`: `QMap<K, V>`
- `Pair[1]2`: `QPair<1, 2>`

# Generating Code

Generating code with pokic is fairly straightforward. pokic takes two flags: `-input file.pokipoki` and `-output file.gen.h`. Both flags are required.

The following code can be used in CMake:

```cmake
add_custom_command(
    OUTPUT File.gen.h
    COMMAND pokic -input ${CMAKE_CURRENT_SOURCE_DIR}/File.pokipoki -output ${CMAKE_CURRENT_BINARY_DIR}/File.gen.h
    DEPENDS File.pokipoki
)

cmake_policy(SET CMP0071 NEW)
```

# Formatting PokiPoki Files

For keeping PokiPoki files well-formatted, adhere to the following conventions:

- Names (Object and Schema) are in PascalCase
- Identifiers (Properties) are in camelCase
- The `schema` declaration is the first line of the file
- There is one empty line after the `schema` declaration
- There is one empty line in between every `object` declaration
- Objects list exactly one property per line
- Properties are alinged like this:
```
object Name {
    longPropertyNameYaddaYadda String
    small                      String
    mediumSizedName            String
}
```
- Object children go after object property declarations
- Generics have no spaces (`Hash[A]B`, not `Hash [A] B`)
