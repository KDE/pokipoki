#include "Database.h"

struct Note : public FBObject
{
    not_null_keys("id")

    Q_PROPERTY(QString text MEMBER text)
    QString text;

private:
    Q_OBJECT
};