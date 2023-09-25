#include "startpage.h"
#include "ui_startpage.h"
#include <QDebug>

StartPage::StartPage(QSettings *g_settings, QWidget *parent) :
    QWidget(parent),
    ui(new Ui::StartPage)
{
    ui->setupUi(this);

    m_tokenPkg = g_settings->value("tokenPkg").toString();
    qDebug() << "Tokenpkg == " << m_tokenPkg;

    //g_settings->setValue("tokenPkg", "HELLO TOKENPKG");

    m_socketManager = new SocketManager();

    m_loginpage = new LoginPage(this);
    m_registerpage = new Registerpage(this);

    m_registerpage->hide();
    m_loginpage->show();
}

StartPage::~StartPage()
{
    delete ui;
}

void StartPage::ShowRegisterPage()
{
    qDebug() << "StartPage::ShowRegisterPage";
    if (m_loginpage->isVisible()) {
        m_loginpage->hide();
    }
    m_registerpage->show();
}

void StartPage::ShowLoginPage()
{
    qDebug() << "StartPage::ShowLoginPage";
    if (m_registerpage->isVisible()) {
        m_registerpage->hide();
    }
    m_loginpage->show();
}
