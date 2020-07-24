#pragma once

#include <QObject>
#include <QSqlDatabase>
#include <QVariant>

#include "Database.h"

inline QString cleanClassName(QString in)
{
    return in.replace(QStringLiteral(":"), QStringLiteral("_")).toLower();
}

class FBDatabase::Private : public QObject
{
    Q_OBJECT

public:
    friend class FBDatabase;
    friend struct FBObject;

    QList<FBTableField> extractProperties(FBObject* object);
    bool initializeForObject(FBObject* object);
    bool insert(FBObject* object);
    QList<FBObject*> whereQuery(QList<QPair<QString,QVariant>> clauses, FBObject *like);

    QSqlDatabase db;
};