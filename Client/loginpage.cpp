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
    //m_socketManager->ServerConnect();

    m_loginBTN = ui->BTNlogin;
    m_registerBTN = ui->BTNregister;
    m_forgotpwBTN = ui->BTNforgotpw;

    //connect(m_loginBTN, SIGNAL(clicked()), m_socketManager, SLOT(ServerConnect()));
    connect(m_loginBTN, SIGNAL(clicked()), m_socketManager, SLOT(ServerConnect()));
    connect(m_registerBTN, SIGNAL(clicked()), this->parentWidget(), SLOT(ShowRegisterPage()));
    connect(m_forgotpwBTN, SIGNAL(clicked()), this, SLOT(ForgotPW()));

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
}

void LoginPage::Register()
{
    qDebug() << "Register";
}

void LoginPage::ForgotPW()
{
    qDebug() << "ForgotPW";
}

