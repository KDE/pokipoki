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

    assert(QSqlDatabase::isDriverAvailable(DRIVER));
    auto ok = QDir().mkpath(QDir::cleanPath(QStandardPaths::writableLocation(QStandardPaths::DataLocation)));
    assert(ok);

    d_ptr->db = QSqlDatabase::addDatabase(DRIVER);
    d_ptr->db.setDatabaseName(QDir::cleanPath(QStandardPaths::writableLocation(QStandardPaths::DataLocation) + "/" + qAppName()));

    QDir::cleanPath(QStandardPaths::writableLocation(QStandardPaths::DataLocation) + "/" + qAppName());

    auto result = d_ptr->db.open();
    assert(result);
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

QSqlDatabase& PPDatabase::connection()
{
    return d_ptr->db;
}

class PPUndoRedoStack::Private
{
    QList<PPUndoRedoable*> undoItems;
    QList<PPUndoRedoable*> redoItems;
    friend class PPUndoRedoStack;
};

PPUndoRedoStack::PPUndoRedoStack(QObject *parent) : QObject(parent)
{
    d_ptr = new Private;
}

PPUndoRedoStack* PPUndoRedoStack::instance()
{
    static QMutex mutex;
    mutex.lock();
    static QPointer<PPUndoRedoStack> stack;
    if (stack.isNull()) {
        stack = new PPUndoRedoStack(qApp);
    }
    mutex.unlock();
    return stack;
};

bool PPUndoRedoStack::canUndo() const {
    return d_ptr->undoItems.length() > 0;
}

bool PPUndoRedoStack::canRedo() const {
    return d_ptr->redoItems.length() > 0;
}


void PPUndoRedoStack::undoItemAdded(PPUndoRedoable* item) {
    d_ptr->undoItems << item;
    if (d_ptr->undoItems.length()-1 == 0) {
        Q_EMIT canUndoChanged();
    }
}

void PPUndoRedoStack::undoItemRemoved(PPUndoRedoable* item) {
    auto idx = d_ptr->undoItems.lastIndexOf(item);
    if (idx == -1) {
        return;
    }
    d_ptr->undoItems.removeAt(idx);
    if (d_ptr->undoItems.length() == 0) {
        Q_EMIT canUndoChanged();
    }
}

void PPUndoRedoStack::redoItemAdded(PPUndoRedoable* item) {
    d_ptr->redoItems << item;
    if (d_ptr->redoItems.length()-1 == 0) {
        Q_EMIT canRedoChanged();
    }
}

void PPUndoRedoStack::redoItemRemoved(PPUndoRedoable* item) {
    auto idx = d_ptr->redoItems.lastIndexOf(item);
    if (idx == -1) {
        return;
    }
    d_ptr->redoItems.removeAt(idx);
    if (d_ptr->redoItems.length() == 0) {
        Q_EMIT canRedoChanged();
    }
}


void PPUndoRedoStack::undo() {
    if (d_ptr->undoItems.empty()) return;
    d_ptr->undoItems.last()->undo();
}

void PPUndoRedoStack::redo() {
    if (d_ptr->redoItems.empty()) return;
    d_ptr->redoItems.last()->redo();
}
