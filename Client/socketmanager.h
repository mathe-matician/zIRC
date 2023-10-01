#ifndef SOCKETMANAGER_H
#define SOCKETMANAGER_H

/*
#ifdef __EMSCRIPTEN__
#include <QWebSocket>
#else
#include <QTcpSocket>
#endif
*/

#include <QTcpSocket>
#include <QObject>


class SocketManager : public QObject
{
    Q_OBJECT
public:
/*
#ifdef __EMSCRIPTEN__
    QWebSocket *socket;
#else
    QTcpSocket *socket;
#endif
*/
    QTcpSocket *socket;

    SocketManager();

public slots:
    void ServerConnect();
    void Success_Connected();
    void Success_Disconnected();
    void Error_Occurred(QAbstractSocket::SocketError socketError);
    void Bytes_Written(qint64 bytes);
    void Read_Data();
/*
#ifndef __EMSCRIPTEN__
    void Read_Data();
#else
    void Read_Data(const QByteArray &data);
#endif
*/
    void Write_Data(const QByteArray &data);

    QByteArray Get_Message_Result();
    void Success_HostLookup();

signals:
    void Display_Data(QByteArray);

private:
    QByteArray m_result;

};

#endif // SOCKETMANAGER_H
