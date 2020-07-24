#include "Database.h"
#include "Database+Private.h"

const QString DRIVER("QSQLITE");

FBDatabase::FBDatabase(QObject *parent) : QObject(parent)
{
    d_ptr = new Private;

    Q_ASSERT(QSqlDatabase::isDriverAvailable(DRIVER));
    Q_ASSERT(QDir().mkpath(QDir::cleanPath(QStandardPaths::writableLocation(QStandardPaths::DataLocation))));

    d_ptr->db = QSqlDatabase::addDatabase(DRIVER);
    d_ptr->db.setDatabaseName(QDir::cleanPath(QStandardPaths::writableLocation(QStandardPaths::DataLocation) + "/" + qAppName()));

    QDir::cleanPath(QStandardPaths::writableLocation(QStandardPaths::DataLocation) + "/" + qAppName());

    Q_ASSERT(d_ptr->db.open());
}

void FBDatabase::save(FBObject* object)
{
    d_ptr->initializeForObject(object);
    d_ptr->insert(object);
}

FBDatabase* FBDatabase::instance()
{
    static QMutex mutex;
    mutex.lock();
    static QPointer<FBDatabase> db;
    if (db.isNull()) {
        db = new FBDatabase(qApp);
    }
    mutex.unlock();
    return db;
};
