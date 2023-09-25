#include "loginpage.h"
#include "ui_loginpage.h"

#include <QLineEdit>
#include <QDebug>

LoginPage::LoginPage(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::LoginPage)
{
    ui->setupUi(this);

    m_socketManager = new SocketManager();
    m_socketManager->ServerConnect();

    connect(ui->BTNlogin, SIGNAL(clicked()), this, SLOT(Login()));
    connect(ui->BTNregister, SIGNAL(clicked()), this->parentWidget(), SLOT(ShowRegisterPage()));
    connect(ui->BTNforgotpw, SIGNAL(clicked()), this, SLOT(ForgotPW()));

    ui->INPUTpassword->setEchoMode(QLineEdit::Password);

/*
// Only for EMSCRIPTEN webassembly stuff
#ifndef __EMSCRIPTEN__
    ui->LABEL_logo->setPixmap(QPixmap(":/images/qt_logo_green_128x128px.png"));
#else
    ui->LABEL_logo->setText("Using emscripten");
#endif
*/
}

LoginPage::~LoginPage()
{
    delete ui;
}

void LoginPage::Login()
{
    qDebug() << "Login";
    QString l_email = ui->INPUTemail->text();
    QString l_password = ui->INPUTpassword->text();
    qDebug() << "Email: " << l_email << " Password: " << l_password;

    if (l_email.isEmpty()) {
        qDebug() << "Email is empty";
    }

    if (l_password.isEmpty()) {
        qDebug() << "Password is empty";
    }
    /*
    if (client) {
        // console.log(`CLIENT UID: '${clientUID}', clientUIDExists: ${clientUIDExists}`);
        if (tokenExists) {
            client.write(`tokenPkg::${tokenpkg} :${nickname}@${host} ${line} \r\n`);
        } else if (clientUIDExists) {
            client.write(`clientUID::${clientUID} :${nickname}@${host} ${line} \r\n`)
        } else {
            client.write(`:${nickname}@${host} ${line} \r\n`)
        }
    }
*/
    // TODO what IRC command is this? LOGIN or AUTHENTICATE? Probably AUTHENTICATE
    QString l_msg = QString("LOGIN %1 %2").arg(l_email).arg(l_password);

    // TODO check if tokenPkg exists in state
    // TODO get nickname from local state OR have the user pass it in with email
    QString l_final_data = QString(":%1@%2 %3 \r\n").arg("nickname").arg(m_socketManager->socket->localAddress().toString()).arg(l_msg);

    m_socketManager->Write_Data(l_final_data.toUtf8());
}

void LoginPage::Register()
{
    qDebug() << "Register";
}

void LoginPage::ForgotPW()
{
    qDebug() << "ForgotPW";
}

