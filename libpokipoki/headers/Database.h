#pragma once

#include <QObject>
#include <QSqlQuery>
#include <QList>
#include <QVariant>

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

struct Predicate {
    virtual QString toWhere() = 0;
    virtual void bindToQuery(QSqlQuery *query) = 0;
    virtual ~Predicate() {}
};

#define eq(key, val) new Equals(QStringLiteral(#key), val)
#define neq(key, val) new NotEquals(QStringLiteral(#key), val)
#define lt(key, val) new LessThan(QStringLiteral(#key), val)
#define gt(key, val) new GreaterThan(QStringLiteral(#key), val)
#define lte(key, val) new LessThanOrEqualTo(QStringLiteral(#key), val)
#define gte(key, val) new GreaterThanOrEqualTo(QStringLiteral(#key), val)
#define like(key, val) new Like(QStringLiteral(#key), val)
#define between(key, first, second) new Between(QStringLiteral(#key), first, second)

#define operatorPredicate(name, operator) struct name : Predicate {\
    QString column;\
    QVariant value;\
\
    name(QString col, QVariant val) : column(col), value(val) {}\
    QString toWhere() override { return QStringLiteral("%1 " #operator " :" #name "_%2").arg(this->column).arg(this->column); }\
    void bindToQuery(QSqlQuery *query) override { query->bindValue(QStringLiteral(":" #name "_%1").arg(this->column), this->value); }\
};

operatorPredicate(Equals, =)
operatorPredicate(NotEquals, !=)
operatorPredicate(LessThan, <)
operatorPredicate(GreaterThan, >)
operatorPredicate(LessThanOrEqualTo, <=)
operatorPredicate(GreaterThanOrEqualTo, >=)
operatorPredicate(Like, LIKE)

struct Between : Predicate {
    QString column;
    QVariant first;
    QVariant second;
    Between(QString col, QVariant first, QVariant second) : column(col), first(first), second(second) {}
    QString toWhere() override { return QStringLiteral("%1 BETWEEN :between_first_%2, :between_second_%2").arg(this->column).arg(this->column).arg(this->column); }
    void bindToQuery(QSqlQuery *query) override {
        query->bindValue(QStringLiteral(":between_first_%1").arg(this->column), this->first);
        query->bindValue(QStringLiteral(":between_second_%1").arg(this->column), this->second);
    }
};

class PredicateList : public QList<Predicate*>
{
public:
    PredicateList(Predicate* item...) {
        va_list args;
        va_start(args, item);

        *this << item;

        va_end(args);
    }
    ~PredicateList() {
        for (auto item : *this) {
            delete item;
        }
    }
    QStringList allPredicatesToWhere() {
        QStringList ret;
        for (auto predicate : *this) {
            ret << predicate->toWhere();
        }
        return ret;
    };
    void bindAllPredicates(QSqlQuery *query) {
        for (auto predicate : *this) {
            predicate->bindToQuery(query);
        }
    }
};