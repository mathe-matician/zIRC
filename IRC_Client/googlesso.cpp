// GoogleSSO.cpp
#include "googlesso.h"
#include <QString>
#include <QDir>
#include <QUrl>
#include <QOAuthHttpServerReplyHandler>
#include <QDesktopServices>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonArray>

#include "authgen.h"

// Get these from https://console.developers.google.com/apis/credentials
#define CLIENT_ID "CLIENT-ID"
#define CLIENT_SECRET "CLIENT-SECRET"
#define AUTH_URI "https://accounts.google.com/o/oauth2/auth"
#define TOKEN_URI "https://oauth2.googleapis.com/token"
#define REDIRECT_URI "http://127.0.0.1:8080/"

GoogleSSO::GoogleSSO(QObject *parent) : QObject(parent) {
}

GoogleSSO::~GoogleSSO() {
    delete this->google;
}

// Invoked externally to initiate
void GoogleSSO::authenticate() {
    this->google = new QOAuth2AuthorizationCodeFlow(this);
    // https://developers.google.com/identity/openid-connect/openid-connect#sendauthrequest
    // use OIDC and OAuth2
    this->google->setScope("openid email");

    m_state = QByteArray(AuthGen().genRandom(quint8(255)));
    this->google->setState(m_state);

    connect(this->google, &QOAuth2AuthorizationCodeFlow::authorizeWithBrowser, [=](QUrl url) {
        QUrlQuery query(url);

        query.addQueryItem("prompt", "consent");      // Param required to get data everytime
        query.addQueryItem("access_type", "offline"); // Needed for Refresh Token (as AccessToken expires shortly)
        query.addQueryItem("login_hint", "zach@syllogi.io"); //
        url.setQuery(query);

        // TODO
        // check for valid browsers that Google auth can use IF we are using Google SSO.
        QDesktopServices::openUrl(url);
    });

    //const QUrl redirectUri(redirectUris[0].toString());
    //const auto port = static_cast<quint16>(redirectUri.port());

    this->google->setAuthorizationUrl(QUrl(AUTH_URI));
    this->google->setAccessTokenUrl(QUrl(TOKEN_URI));
    this->google->setClientIdentifier(CLIENT_ID);
    this->google->setClientIdentifierSharedKey(CLIENT_SECRET);

    this->google->setModifyParametersFunction([](QAbstractOAuth::Stage stage, QMultiMap<QString, QVariant> *parameters) -> void {
        // Percent-decode the "code" parameter so Google can match it
        if (stage == QAbstractOAuth::Stage::RequestingAccessToken) {
            QByteArray code = parameters->value("code").toByteArray();
            parameters->replace("code", QUrl::fromPercentEncoding(code));
        }
    });

    auto replyHandler = new QOAuthHttpServerReplyHandler(8080, this);
    this->google->setReplyHandler(replyHandler);

    connect(this->google, &QOAuth2AuthorizationCodeFlow::granted, [=](){
        qDebug() << __FUNCTION__ << __LINE__ << "Access Granted!";

        if (this->google->extraTokens().value("state").toString() == m_state) {
            // TODO
            // error handling if they do not match - something seriously went wrong!
            qDebug() << "STATES MATCH YO!";
        }


        qDebug() << "Extra token id_token: " << this->google->extraTokens().value("id_token");

        /*
            authuser
            hd
            id_token
            prompt
            scope
            state
        */
//        for (auto i : this->google->extraTokens().keys()) {
//            qDebug() << "Extra token key: " << i;
//        }
//        auto reply = this->google->get(QUrl("https://www.googleapis.com/plus/v1/people/me"));
//        connect(reply, &QNetworkReply::finished, [reply](){
//            qDebug() << "REQUEST FINISHED. Error? " << (reply->error() != QNetworkReply::NoError);
//            qDebug() << reply->readAll();
//        });
    });

    this->google->grant();
}
