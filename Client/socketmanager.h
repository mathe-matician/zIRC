#ifndef SOCKETMANAGER_H
#define SOCKETMANAGER_H

#include <QTcpSocket>
#include <QObject>


class SocketManager : public QTcpSocket
{
    Q_OBJECT
public:
    SocketManager();

public slots:
    void ServerConnect();
    void Success_HostLookup();
    void Success_Connected();
    void Success_Disconnected();
    void Error_Occurred(QAbstractSocket::SocketError socketError);
    void Bytes_Written(qint64 bytes);
    void Read_Data();

private:

};

#endif // SOCKETMANAGER_H
