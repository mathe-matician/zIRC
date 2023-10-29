#ifndef SOCKETMANAGER_H
#define SOCKETMANAGER_H

#include <QObject>
#include <QQmlEngine>
#include <QTcpSocket>
#include <QtQml>

class SocketManager : public QObject
{   
    Q_OBJECT
    QML_ELEMENT
public:
    explicit SocketManager(QObject *parent = nullptr);

    QTcpSocket *socket;

    QString nick() const;
    void setNick(const QString &newNick);


public slots:
    Q_INVOKABLE void ServerConnect();
    void Success_Connected();
    void Success_Disconnected();
    void Error_Occurred(QAbstractSocket::SocketError socketError);
    void Bytes_Written(qint64 bytes);
    void Read_Data();
    void Debug_Send(QByteArray data);
    void PRIVMSG(QByteArray data);
    Q_INVOKABLE void Write_Data(const QByteArray &data);

    QByteArray Get_Message_Result();
    void Success_HostLookup();


signals:
    void Display_Data(QByteArray);

private:
    QByteArray m_result;
    QString m_nick;
};

#endif // SOCKETMANAGER_H
