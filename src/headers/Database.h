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
#define not_null_keys(keys) QStringList notNullKeys() const override { return { keys }; };

struct FBTableField {
    int metaIndex;
    QString name;
    QString sqlType;
    QVariant::Type qtType;
    bool notNull;
};

struct FBObject;

class FBDatabase : public QObject
{
    Q_OBJECT

private:
    FBDatabase(QObject *parent);
    class Private;
    friend class FBObject;
    Private *d_ptr;

public:
    static FBDatabase* instance();

    void save(FBObject* object);
};

struct FBObject : public QObject
{
    virtual QStringList notNullKeys() const { return {}; };
    virtual ~FBObject() {};

    Q_PROPERTY(QByteArray id READ id)
    QByteArray id() const { return p_id; }

private:
    friend class FBDatabase::Private;
    QByteArray p_id;

    Q_OBJECT
};

