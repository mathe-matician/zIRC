#ifndef BACKEND_H
#define BACKEND_H

#include <QObject>
#include <QtQml>

class BackEnd : public QObject
{
    Q_OBJECT
    QML_ELEMENT
public:
    explicit BackEnd(QObject *parent = nullptr);

    Q_INVOKABLE void generateNumber(int min, int max);

signals:
    void numberEmitted(int num);
};

#endif // BACKEND_H
