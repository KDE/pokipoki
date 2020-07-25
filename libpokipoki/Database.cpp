#include <QCoreApplication>
#include <QDebug>
#include <QDir>
#include <QMetaProperty>
#include <QMutex>
#include <QPointer>
#include <QSqlDatabase>
#include <QSqlError>
#include <QSqlQuery>
#include <QStandardPaths>
#include <QStringList>
#include <QVariant>

#include "Database.h"

const QString DRIVER("QSQLITE");

class PPDatabase::Private
{
    friend class PPDatabase;
    QSqlDatabase db;
};

PPDatabase::PPDatabase(QObject *parent) : QObject(parent)
{
    d_ptr = new Private;

    Q_ASSERT(QSqlDatabase::isDriverAvailable(DRIVER));
    Q_ASSERT(QDir().mkpath(QDir::cleanPath(QStandardPaths::writableLocation(QStandardPaths::DataLocation))));

    d_ptr->db = QSqlDatabase::addDatabase(DRIVER);
    d_ptr->db.setDatabaseName(QDir::cleanPath(QStandardPaths::writableLocation(QStandardPaths::DataLocation) + "/" + qAppName()));

    QDir::cleanPath(QStandardPaths::writableLocation(QStandardPaths::DataLocation) + "/" + qAppName());

    Q_ASSERT(d_ptr->db.open());
}

PPDatabase* PPDatabase::instance()
{
    static QMutex mutex;
    mutex.lock();
    static QPointer<PPDatabase> db;
    if (db.isNull()) {
        db = new PPDatabase(qApp);
    }
    mutex.unlock();
    return db;
};