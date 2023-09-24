#ifndef SOCKETMANAGER_H
#define SOCKETMANAGER_H

#include <QTcpSocket>


class SocketManager : public QTcpSocket
{
public:
    SocketManager();

public slots:
    void ServerConnect();

private:

};

#endif // SOCKETMANAGER_H
