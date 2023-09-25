#include "socketmanager.h"

#include <QString>
#include <QDebug>

SocketManager::SocketManager()
{
    qDebug() << "SocketManager created";

    connect(this, SIGNAL(hostFound()), this, SLOT(Success_HostLookup()));
    connect(this, SIGNAL(connected()), this, SLOT(Success_Connected()));
    connect(this, SIGNAL(disconnected()), this, SLOT(Success_Disconnected()));
    connect(this, SIGNAL(errorOccurred(QAbstractSocket::SocketError)), this, SLOT(Error_Occurred(QAbstractSocket::SocketError)));
    connect(this, SIGNAL(bytesWritten(qint64)), this, SLOT(Bytes_Written(qint64)));
    connect(this, SIGNAL(readyRead()), this, SLOT(Bytes_Written(qint64)));
}


void SocketManager::ServerConnect()
{
    qDebug() << "ServerConnect started";
    QString l_host = QString();
    quint16 l_port;
#ifdef QT_DEBUG
    l_host.append("localhost");
    l_port = 6667;
#endif
    this->connectToHost(l_host, l_port);
    if (this->waitForConnected(1000))
        qDebug("SocketManager::ServerConnect(): Waiting for connected success");

    //if (this->isValid())
}

void SocketManager::Success_HostLookup()
{
    qDebug() << "Host lookup successful!!";
}

void SocketManager::Success_Connected()
{
    qDebug() << "Successfully connected to server: " << this->peerName() << this->peerAddress().toString() << ":" << this->peerPort();
}

void SocketManager::Success_Disconnected()
{
    qDebug() << "Successfully disconnected to server";
}

void SocketManager::Error_Occurred(QAbstractSocket::SocketError socketError)
{
    qDebug() << "SocketError occurred: " << socketError;
}

void SocketManager::Bytes_Written(qint64 bytes)
{
    qDebug() << "Bytes_Written num of bytes: " << bytes;
}

void SocketManager::Read_Data()
{
    QByteArray res = this->readAll();
    qDebug() << "RES: " << res.toStdString();
}
