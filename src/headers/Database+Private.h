#pragma once

#include <QObject>
#include <QSqlDatabase>
#include <QVariant>

#include "Database.h"

class FBDatabase::Private : public QObject
{
    Q_OBJECT

public:
    friend class FBDatabase;
    friend struct FBObject;

    QSqlDatabase db;
};