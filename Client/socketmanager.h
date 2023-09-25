#ifndef SOCKETMANAGER_H
#define SOCKETMANAGER_H

#ifdef __EMSCRIPTEN__
#include <QWebSocket>
#else
#include <QTcpSocket>
#endif

#include <QObject>


class SocketManager : public QObject
{
    Q_OBJECT
public:
#ifdef __EMSCRIPTEN__
    QWebSocket *socket;
#else
    QTcpSocket *socket;
#endif

    SocketManager();

public slots:
    void ServerConnect();
    void Success_Connected();
    void Success_Disconnected();
    void Error_Occurred(QAbstractSocket::SocketError socketError);
    void Bytes_Written(qint64 bytes);
    void Read_Data();
    void Write_Data(const QByteArray &data);

#ifdef __EMSCRIPTEN__

#else
    void Success_HostLookup();
#endif

private:

};

#endif // SOCKETMANAGER_H
