#include "authgen.h"


AuthGen::AuthGen(QObject *parent)
    : QOAuth2AuthorizationCodeFlow(parent)
{}

QByteArray AuthGen::genRandom(quint8 length)
{
    // genered here because we aren't deriving a class for google sso
    // and we want to use this auth method which is protected otherwise and in accessible.
    return this->generateRandomString(length);
}
