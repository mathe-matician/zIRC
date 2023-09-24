#include "socketmanager.h"

#include <QString>
#include <QDebug>

SocketManager::SocketManager()
{
    qDebug() << "SocketManager created";
}


void SocketManager::ServerConnect()
{
    qDebug() << "ServerConnect started";
    QString host = QString();
    quint16 port;
#ifdef QT_DEBUG
    host.append("localhost");
    port = 6667;
#endif
    this->connectToHost(host, port);
    if (this->waitForConnected(1000))
        qDebug("SocketManager::ServerConnect(): Connected to Chat Server!!");
}
