#ifndef AUTHGEN_H
#define AUTHGEN_H

#include <QOAuth2AuthorizationCodeFlow>
#include <QByteArray>

class AuthGen : public QOAuth2AuthorizationCodeFlow
{
public:
    explicit AuthGen(QObject *parent = nullptr);
    QByteArray genRandom(quint8 length);
};

#endif // AUTHGEN_H
