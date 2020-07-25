#pragma once

#include <QObject>

class PPDatabase : public QObject
{
    Q_OBJECT

private:
    PPDatabase(QObject *parent);
    class Private;
    Private *d_ptr;

public:
    static PPDatabase* instance();
};

