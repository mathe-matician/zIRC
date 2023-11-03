#include "socketmanager.h"

#include <QString>
#include <QDebug>

/*
#ifdef __EMSCRIPTEN__
#include <QWebSocket>
#else
#include <QTcpSocket>
#endif
*/

#include <QTcpSocket>

SocketManager::SocketManager(QObject *parent) :
    QObject(parent)
{
    qDebug() << "SocketManager created";

    /*
#ifdef __EMSCRIPTEN__
    socket = new QWebSocket("localhost:6667");
    connect(socket, SIGNAL(binaryMessageReceived(QByteArray)), this, SLOT(Read_Data(QByteArray)));
#else
    socket = new QTcpSocket();
    QObject::connect(socket, SIGNAL(hostFound()), this, SLOT(Success_HostLookup()));
    QObject::connect(socket, SIGNAL(readyRead()), this, SLOT(Read_Data()));
#endif
*/

    socket = new QTcpSocket();
    QObject::connect(socket, SIGNAL(hostFound()), this, SLOT(Success_HostLookup()));
    QObject::connect(socket, SIGNAL(readyRead()), this, SLOT(Read_Data()));

    connect(socket, SIGNAL(connected()), this, SLOT(Success_Connected()));
    connect(socket, SIGNAL(disconnected()), this, SLOT(Success_Disconnected()));
    connect(socket, SIGNAL(errorOccurred(QAbstractSocket::SocketError)), this, SLOT(Error_Occurred(QAbstractSocket::SocketError)));
    connect(socket, SIGNAL(bytesWritten(qint64)), this, SLOT(Bytes_Written(qint64)));
    connect(socket, SIGNAL(connected()), this, SLOT(Success_Connected()));


}


void SocketManager::ServerConnect()
{
    qDebug() << "ServerConnect started";
    QString l_host = QString();
    quint16 l_port;

#if defined(QT_DEBUG)
#if defined(Q_OS_ANDROID)
    l_host.append("10.0.2.2");
#elif defined(Q_OS_IOS)
    l_host.append("localhost");
#else
    l_host.append("localhost");
#endif
    l_port = 6667;
#else
    l_host.append("localhost");
    l_port = 6667;
#endif

    socket->connectToHost(l_host, l_port);
    if (socket->waitForConnected(1000))
        qDebug("SocketManager::ServerConnect(): Waiting for connected success");
    /*
#ifndef __EMSCRIPTEN__
    socket->connectToHost(l_host, l_port);
    if (socket->waitForConnected(1000))
        qDebug("SocketManager::ServerConnect(): Waiting for connected success");
#else
    l_host.append("ws://localhost");
    l_port = 6667;
    QUrl l_url(l_host);
    l_url.port(l_port);
    socket->open(l_url);
#endif
    */

    qDebug("SocketManager connected successfully");
}

/*
#ifndef __EMSCRIPTEN__
void SocketManager::Success_HostLookup()
{
    qDebug() << "Host lookup successful!!";
}
#endif
*/

void SocketManager::Success_HostLookup()
{
    qDebug() << "Host lookup successful!!";
}

QString SocketManager::nick() const
{
    return m_nick;
}

void SocketManager::setNick(const QString &newNick)
{
    m_nick = newNick;
}

void SocketManager::Success_Connected()
{
    qDebug() << "Successfully connected to server: " << socket->peerName() << socket->peerAddress().toString() << ":" << socket->peerPort();
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

//#ifndef __EMSCRIPTEN__
void SocketManager::Read_Data()
{
    // TODO
    // Run in separate thread for i/o

    QByteArray res = socket->readAll();
    qDebug() << "SocketManager::Read_Data() RES: " << res.toStdString();
    emit Display_Data(res);
    //m_result = res;

}
//#else
/*
void SocketManager::Read_Data(const QByteArray &data)
{
    qDebug() << "SocketManager::Read_Data start RES: " << data.toStdString();
}
*/
//#endif


void SocketManager::Write_Data(const QByteArray &data)
{
    //#ifdef __EMSCRIPTEN__
    //qint64 l_bytes_written = socket->sendBinaryMessage(data);
    //#else
    //QString l_final_data = QString(":%1@%2 %3 \r\n").arg(l_nickname).arg(m_socketManager->socket->localAddress().toString()).arg(l_msg);
    qint64 l_bytes_written = socket->write(data);
    //#endif
    if (l_bytes_written == -1) {
        qDebug() << "Error writing to socket";
    }

    qDebug() << l_bytes_written << " bytes written to socket";

}

void SocketManager::Debug_Send(QByteArray data)
{
    qDebug() << "Debug_Send " << data.toStdString();
}

void SocketManager::PRIVMSG(QByteArray data)
{
    qDebug() << "PRIVMSG start: " << data.toStdString();
    QString l_nick = m_nick.isEmpty() ? "*" : m_nick;
    QString l_msg = QString("PRIVMSG %1 %2").arg(l_nick).arg(data);
    qDebug() << "PRIVMSG cmd: " + l_msg;
}

///////////////
// GETTERS
///////////////

QByteArray SocketManager::Get_Message_Result() { return this->m_result; }
