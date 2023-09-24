#include "startpage.h"
#include "ui_startpage.h"
#include <QDebug>

StartPage::StartPage(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::StartPage)
{
    ui->setupUi(this);

    m_loginpage = new LoginPage(this);
    m_registerpage = new Registerpage(this);
    //m_socketManager = new SocketManager();
    m_registerpage->hide();
    m_loginpage->show();
}

StartPage::~StartPage()
{
    delete ui;
}

void StartPage::ShowLoginPage()
{

}

void StartPage::ShowRegisterPage()
{
    qDebug() << "StartPage::ShowRegisterPage";
    if (m_loginpage->isEnabled()) {
        m_loginpage->hide();
    }
    m_registerpage->show();
}
