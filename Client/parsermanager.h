#ifndef PARSERMANAGER_H
#define PARSERMANAGER_H

#include <QObject>

class ParserManager : public QObject
{
    Q_OBJECT
public:
    ParserManager();
    ParserManager(QByteArray data);
    //~ParserManager();
};

#endif // PARSERMANAGER_H
