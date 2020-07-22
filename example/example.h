#include "Database.h"

struct Note : public FBObject
{
    primary_keys("id")
    not_null_keys("id")

    Q_PROPERTY(QString id MEMBER id)
    QString id;

    Q_PROPERTY(QString text MEMBER text)
    QString text;

private:
    Q_GADGET
};