#include <QMetaProperty>
#include <QSqlQuery>
#include <QSqlError>
#include <QUuid>

#include "Database.h"
#include "Database+Private.h"

QList<FBTableField> FBDatabase::Private::extractProperties(FBObject* object)
{
    auto metaObject = object->metaObject();

    QList<FBTableField> fields;

    for (int i = metaObject->propertyOffset(); i < metaObject->propertyCount(); ++i) {
        auto prop = metaObject->property(i);
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
            bool isNotNull = object->notNullKeys().contains(propName);
            fields << FBTableField{i, propName, type, prop.type(), isNotNull};
        }
    }

    return fields;
}

bool FBDatabase::Private::initializeForObject(FBObject* object)
{
    QList<FBTableField> fields = extractProperties(object);

    auto tableQuery = QStringLiteral("CREATE TABLE IF NOT EXISTS %1 (\n").arg(cleanClassName(object->metaObject()->className()));
    for (const auto& field : fields) {
        tableQuery.append(QStringLiteral("\t%1 %2%3,").arg(field.name).arg(field.sqlType).arg(field.notNull ? QStringLiteral(" NOT NULL") : QStringLiteral("")));
        tableQuery.append("\n");
    }
    tableQuery.append(QStringLiteral("\tID BLOB NOT NULL,\n"));
    tableQuery.append(QStringLiteral("\tPRIMARY KEY (ID)\n"));
    tableQuery.append(QStringLiteral(");\n"));

    QSqlQuery query;
    auto ret = query.exec(tableQuery);
    if (!ret) {
        qCritical() << qUtf8Printable(tableQuery);
        qCritical() << query.lastError();
    }
    return ret;
}

bool FBDatabase::Private::insert(FBObject* object)
{
    QList<FBTableField> fields = extractProperties(object);
    QStringList fieldNames = [=]() {
        QStringList ret;
        for (const auto& field : fields) {
            ret << field.name;
        }
        return ret;
    }();
    fieldNames << "ID";
    QStringList placeHolders = [=]() {
        QStringList ret;
        for (const auto& field : fields) {
            Q_UNUSED(field)
            ret << QStringLiteral("?");
        }
        return ret;
    }();
    placeHolders << "?";

    QVariantList properties;
    for (const auto& field : fields) {
        QMetaProperty prop = object->metaObject()->property(field.metaIndex);
        properties << object->property(prop.name());
    }

    QSqlQuery query;
    auto tq = QStringLiteral("INSERT OR REPLACE INTO %1(%2) VALUES (%3)").
        arg(cleanClassName(object->metaObject()->className())).
        arg(fieldNames.join(",")).
        arg(placeHolders.join(","));
    query.prepare(tq);
    for (auto val : properties) {
        query.addBindValue(val);
    }
    if (object->p_id.length() == 0) {
        object->p_id = QUuid::createUuid().toByteArray();
    }
    query.addBindValue(object->p_id);

    auto ret = query.exec();
    if (!ret) {
        qCritical() << qUtf8Printable(tq);
        qCritical() << query.lastError();
    }
    return ret;
}

QList<FBObject*> FBDatabase::Private::whereQuery(QList<QPair<QString,QVariant>> clauses, FBObject *like)
{
    QStringList c;

    for (auto idx : clauses) {
        c << QStringLiteral("%1 = ?").arg(idx.first);
    }

    auto tq = QStringLiteral("SELECT * FROM %1 WHERE %2").arg(cleanClassName(like->metaObject()->className())).arg(c.join(" AND "));
    QList<FBTableField> fields = extractProperties(like);

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

    QList<FBObject*> ret;

    while (query.next()) {
        auto add = qobject_cast<FBObject*>(like->metaObject()->newInstance());
        for (const auto& field : fields) {
            QMetaProperty prop = add->metaObject()->property(field.metaIndex);
            add->setProperty(field.name.toStdString().c_str(), query.value(field.name));
        }
        ret << add;
    }

    return ret;
}