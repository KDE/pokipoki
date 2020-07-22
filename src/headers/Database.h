#pragma once

#include <QCoreApplication>
#include <QDebug>
#include <QDir>
#include <QMetaProperty>
#include <QMutex>
#include <QObject>
#include <QPointer>
#include <QSqlDatabase>
#include <QSqlError>
#include <QSqlQuery>
#include <QStandardPaths>
#include <QStringList>
#include <QVariant>

#define en ,
#define primary_keys(keys) QStringList primaryKeys() const override { return { keys }; };
#define not_null_keys(keys) QStringList notNullKeys() const override { return { keys }; };

struct FBObject
{
    virtual QStringList primaryKeys() const { return {}; };
    virtual QStringList notNullKeys() const { return {}; };
    virtual ~FBObject() {};

private:
    Q_GADGET
};

struct FBTableField {
    int metaIndex;
    QString name;
    QString sqlType;
    QVariant::Type qtType;
    bool notNull;
};

class FBDatabase : public QObject
{
    Q_OBJECT

private:
    FBDatabase(QObject *parent);
    class Private;
    friend class FBObject;
    Private *d_ptr;

    template <class T>
    QList<FBTableField> extractProperties(T& object);

    template <class T>
    bool createTable();

    template<class T>
    QList<T> whereQuery(QList<QPair<QString,QVariant>> clauses);

    template <class T>
    bool insert(T value);

public:
    static FBDatabase* instance();

    template<class T>
    void save(const T &obj);
};

inline QString cleanClassName(QString in)
{
    return in.replace(QStringLiteral(":"), QStringLiteral("_")).toLower();
}

template <class T>
void FBDatabase::save(const T& obj)
{
    createTable<T>();
    insert(obj);
}

template <class T>
QList<FBTableField> FBDatabase::extractProperties(T& object)
{
    static_assert(std::is_base_of<FBObject, T>::value, "type parameter of this class must derive from FBObject");
    auto metaObject = object.staticMetaObject;

    QList<FBTableField> fields;

    for (int i = metaObject.propertyOffset(); i < metaObject.propertyCount(); ++i) {
        auto prop = metaObject.property(i);
        if (prop.isReadable() && prop.isWritable()) {
            auto propName = QString(prop.name());
            QString type;
            switch (prop.type()) {
            case QVariant::String:
                type = QStringLiteral("TEXT");
                break;
            case QVariant::Int:
                type = QStringLiteral("INTEGER");
                break;
            default:
                type = QStringLiteral("BLOB");
                break;
            }
            bool isNotNull = object.notNullKeys().contains(propName);
            fields << FBTableField{i, propName, type, prop.type(), isNotNull};
        }
    }

    return fields;
}

template <class T>
bool FBDatabase::createTable()
{
    static_assert(std::is_base_of<FBObject, T>::value, "type parameter of this class must derive from FBObject");

    T object;
    QList<FBTableField> fields = extractProperties(object);

    auto tableQuery = QStringLiteral("CREATE TABLE IF NOT EXISTS %1 (\n").arg(cleanClassName(object.staticMetaObject.className()));
    for (const auto& field : fields) {
        tableQuery.append(QStringLiteral("\t%1 %2%3,").arg(field.name).arg(field.sqlType).arg(field.notNull ? QStringLiteral(" NOT NULL") : QStringLiteral("")));
        tableQuery.append("\n");
    }
    tableQuery.append(QStringLiteral("\tPRIMARY KEY (%1)\n").arg(object.primaryKeys().join(", ")));
    tableQuery.append(QStringLiteral(");\n"));

    QSqlQuery query;
    auto ret = query.exec(tableQuery);
    if (!ret) {
        qCritical() << qUtf8Printable(tableQuery);
        qCritical() << query.lastError();
    }
    return ret;
}

template<class T>
QList<T> FBDatabase::whereQuery(QList<QPair<QString,QVariant>> clauses)
{
    static_assert(std::is_base_of<FBObject, T>::value, "type parameter of this class must derive from FBObject");

    T tmpl;
    QStringList c;

    for (auto idx : clauses) {
        c << QStringLiteral("%1 = ?").arg(idx.first);
    }

    auto tq = QStringLiteral("SELECT * FROM %1 WHERE %2").arg(cleanClassName(tmpl.staticMetaObject.className())).arg(c.join(" AND "));
    QList<FBTableField> fields = extractProperties(tmpl);

    QSqlQuery query;
    query.prepare(tq);
    for (auto val : clauses) {
        query.addBindValue(val.second);
    }
    auto ok = query.exec();
    if (!ok) {
        qCritical() << qUtf8Printable(tq);
        qCritical() << query.lastError();
    }

    QList<T> ret;

    while (query.next()) {
        T add;
        for (const auto& field : fields) {
            QMetaProperty prop = add.staticMetaObject.property(field.metaIndex);
            prop.writeOnGadget(&add, query.value(field.name));
        }
        ret << add;
    }

    return ret;
}

template <class T>
bool FBDatabase::insert(T value)
{
    static_assert(std::is_base_of<FBObject, T>::value, "type parameter of this class must derive from FBObject");

    QList<FBTableField> fields = extractProperties(value);
    QStringList fieldNames = [=]() {
        QStringList ret;
        for (const auto& field : fields) {
            ret << field.name;
        }
        return ret;
    }();
    QStringList placeHolders = [=]() {
        QStringList ret;
        for (const auto& field : fields) {
            Q_UNUSED(field)
            ret << QStringLiteral("?");
        }
        return ret;
    }();

    QVariantList properties;
    for (const auto& field : fields) {
        QMetaProperty prop = value.staticMetaObject.property(field.metaIndex);
        properties << prop.readOnGadget(&value);
    }

    QSqlQuery query;
    auto tq = QStringLiteral("INSERT OR REPLACE INTO %1(%2) VALUES (%3)").arg(cleanClassName(value.staticMetaObject.className())).arg(fieldNames.join(",")).arg(placeHolders.join(","));
    query.prepare(tq);
    for (auto val : properties) {
        query.addBindValue(val);
    }

    auto ret = query.exec();
    if (!ret) {
        qCritical() << qUtf8Printable(tq);
        qCritical() << query.lastError();
    }
    return ret;
}