#ifndef GOOGLESSO_H
#define GOOGLESSO_H

#include <QObject>
#include <QtQml>

#include <QNetworkReply>
#include <QOAuth2AuthorizationCodeFlow>

class GoogleSSO : public QObject {
    Q_OBJECT
    QML_ELEMENT
public:
    GoogleSSO(QObject *parent=nullptr);
    virtual ~GoogleSSO();

public slots:
    Q_INVOKABLE void authenticate();
    Q_INVOKABLE void firebase_basic_auth(QString a_mail, QString a_password);

signals:
    void gotToken(const QString& token);

private:
    QOAuth2AuthorizationCodeFlow * google;
    QByteArray m_state;

    const char* m_email;
    const char* m_password;
};

#endif // GOOGLESSO_H
