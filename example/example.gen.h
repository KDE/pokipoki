#include <QObject>
#include <QUuid>
#include <QSqlQuery>
#include <QSqlError>
#include <QVariant>
#include <QDebug>
#include <QString>
#include <QMap>

#include "Database.h"
class Note : public QObject {
        Q_OBJECT

        Note(QUuid ID) : QObject(nullptr), m_ID(ID) {
                static bool db_initialized = false;
                if (!db_initialized) {
                        volatile auto db = PPDatabase::instance();
                        prepareDatabase();
                        db_initialized = true;
                }
        }


        QUuid m_parent_Note_ID;



        QUuid m_ID;
        bool m_NEW = false;




        Q_PROPERTY(QString title READ title WRITE set_title)
        QString m_title;
        QString m_title_prev;
        bool m_title_dirty;



        Q_PROPERTY(QMap<QString,QString> metadata READ metadata WRITE set_metadata)
        QMap<QString,QString> m_metadata;
        QMap<QString,QString> m_metadata_prev;
        bool m_metadata_dirty;


public:



        Q_SIGNAL void titleChanged();
        QString title() const { return m_title; };
        void set_title(const QString& val) {
                if (val == m_title) {
                        return;
                }
                m_title_prev = m_title;
                m_title_dirty = true;
                m_title = val;
                Q_EMIT void titleChanged();
        }
        void discard_title_changes() {
                if (m_title_dirty) {
                        m_title_dirty = false;
                        m_title = m_title_prev;
                        Q_EMIT void titleChanged();
                }
        }



        Q_SIGNAL void metadataChanged();
        QMap<QString,QString> metadata() const { return m_metadata; };
        void set_metadata(const QMap<QString,QString>& val) {
                if (val == m_metadata) {
                        return;
                }
                m_metadata_prev = m_metadata;
                m_metadata_dirty = true;
                m_metadata = val;
                Q_EMIT void metadataChanged();
        }
        void discard_metadata_changes() {
                if (m_metadata_dirty) {
                        m_metadata_dirty = false;
                        m_metadata = m_metadata_prev;
                        Q_EMIT void metadataChanged();
                }
        }


        void discard_all_changes() {

                if (m_title_dirty) {
                        m_title_dirty = false;
                        m_title = m_title_prev;
                        Q_EMIT void titleChanged();
                }

                if (m_metadata_dirty) {
                        m_metadata_dirty = false;
                        m_metadata = m_metadata_prev;
                        Q_EMIT void metadataChanged();
                }

        }

        void commit() {
                if (m_NEW) {
                        auto tq = QStringLiteral(R"RJIENRLWEY(
INSERT INTO Note
(ID,title,metadata)
VALUES
(:ID,   :title  , :metadata );
                        )RJIENRLWEY");
                        QSqlQuery query;
                        query.prepare(tq);
                        query.bindValue(":ID", QVariant::fromValue(m_ID));
                        query.bindValue(":title", QVariant::fromValue(m_title));

                        query.bindValue(":metadata", QVariant::fromValue(m_metadata));
                        auto res = query.exec();
                        if (!res) {
                                qCritical() << query.lastError() << "when creating a new item of Note";
                        }
                        m_NEW = false;
                } else {

                if (m_title_dirty) {
                        QSqlQuery query;
                        auto tq = QStringLiteral(R"RJIENRLWEY( UPDATE Note SET title = :val WHERE ID = :id )RJIENRLWEY");
                        query.prepare(tq);
                        query.bindValue(":val", QVariant::fromValue(m_title));
                        query.bindValue(":id", QVariant::fromValue(m_ID));
                        auto res = query.exec();
                        if (!res) {
                                qCritical() << query.lastError() << "when updating an item of type Note at row title";
                        }
                        m_title_dirty = false;
                }

                if (m_metadata_dirty) {
                        QSqlQuery query;
                        auto tq = QStringLiteral(R"RJIENRLWEY( UPDATE Note SET metadata = :val WHERE ID = :id )RJIENRLWEY");
                        query.prepare(tq);
                        query.bindValue(":val", QVariant::fromValue(m_metadata));
                        query.bindValue(":id", QVariant::fromValue(m_ID));
                        auto res = query.exec();
                        if (!res) {
                                qCritical() << query.lastError() << "when updating an item of type Note at row metadata";
                        }
                        m_metadata_dirty = false;
                }


                }
        }


        QList<Note*> childNotes() {
                auto tq = QStringLiteral("SELECT * FROM Note WHERE PARENT_Note_ID = :parent_id");
                QSqlQuery query;
                query.prepare(tq);
                query.bindValue(":parent_id", m_ID);
                auto ok = query.exec();
                if (!ok) {
                        qCritical() << query.lastError() << "when loading an Note children of a Note";
                }
                QList<Note*> ret;

                while (query.next()) {
                        auto add = new Note(query.value("ID").value<QUuid>());

                        add->setProperty("title", query.value("title"));

                        add->setProperty("metadata", query.value("metadata"));

                        ret << add;
                }
                return ret;
        }
        void addChildNote(Note* child) {
                auto tq = QStringLiteral("UPDATE Note SET PARENT_Note_ID = :new_parent_id WHERE ID = :child_id ");
                QSqlQuery query;
                query.prepare(tq);
                query.bindValue(":new_parent_id", m_ID);
                query.bindValue(":child_id", child->m_ID);
                auto ok = query.exec();
                if (!ok) {
                        qCritical() << query.lastError() << "when adding a new Note to a parent Note";
                }
                child->m_parent_Note_ID = m_ID;
        }
        void removeChildNote(Note* child) {
                auto tq = QStringLiteral("UPDATE Note SET PARENT_Note_ID = NULL WHERE ID = :child_id ");
                QSqlQuery query;
                query.prepare(tq);
                query.bindValue(":child_id", child->m_ID);
                auto ok = query.exec();
                if (!ok) {
                        qCritical() << query.lastError() << "when removing a Note from a parent Note";
                }
                child->m_parent_Note_ID = QUuid();
        }


        static Note* newNote() {
                auto ret = new Note(QUuid::createUuid());
                ret->m_NEW = true;
                return ret;
        }

        static Note* load(const QUuid& ID) {
                auto tq = QStringLiteral("SELECT * FROM Note WHERE ID = :id");
                QSqlQuery query;
                query.prepare(tq);
                query.bindValue(":id", ID);
                auto ok = query.exec();
                if (!ok) {
                        qCritical() << query.lastError() << "when loading an item of type Note";
                }
                auto ret = new Note(ID);
                while (query.next()) {
                        ret->setProperty("title", query.value("title"));
                        ret->setProperty("metadata", query.value("metadata"));

                }
                return ret;
        }

        static void prepareDatabase() {
                auto tq = QStringLiteral(R"RJIENRLWEY(
                CREATE TABLE IF NOT EXISTS Note(
                        ID BLOB NOT NULL,

                        PARENT_Note_ID BLOB,

                        title BLOB NOT NULL,
                        metadata BLOB NOT NULL,
                        PRIMARY KEY (ID))
                )RJIENRLWEY");
                QSqlQuery query;
                query.prepare(tq);
                auto ok = query.exec();
                if (!ok) {
                        qCritical() << query.lastError();
                }
        }

};