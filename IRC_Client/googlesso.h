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

signals:
    void gotToken(const QString& token);

private:
    QOAuth2AuthorizationCodeFlow * google;
};

#endif // GOOGLESSO_H
