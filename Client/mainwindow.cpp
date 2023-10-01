#include "mainwindow.h"
#include "./ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    g_settings = new QSettings();


    m_tokenPkg = g_settings->value("tokenPkg").toString();
    qDebug() << "Tokenpkg == " << m_tokenPkg;

    //g_settings->setValue("tokenPkg", "HELLO TOKENPKG");

    m_socketManager = new SocketManager();
    m_socketManager->ServerConnect();

    ShowLoginPage();
    //this->layout()->addWidget(m_loginpage);
    //m_loginpage->show();
}

MainWindow::~MainWindow()
{
    delete ui;
}



void MainWindow::ShowRegisterPage()
{
    qDebug() << "MainWindow::ShowRegisterPage: ";
    if (!this->centralWidget()->layout()->isEmpty()) {
        qDebug() << "Layout not empty!: " << this->centralWidget()->layout()->activate();
        this->centralWidget()->layout()->itemAt(0)->widget()->deleteLater();
        //this->centralWidget()->layout()->removeWidget(this->centralWidget()->layout()->itemAt(0)->widget());
    }

    m_registerpage = new Registerpage(this, m_socketManager);
    connect(m_registerpage, SIGNAL(RegisterSuccess()), this, SLOT(ShowMainChatPage()));

    this->centralWidget()->layout()->addWidget(m_registerpage);
}

void MainWindow::ShowLoginPage()
{
    qDebug() << "MainWindow::ShowLoginPage";
    if (!this->centralWidget()->layout()->isEmpty()) {
        qDebug() << "Layout not empty!: " << this->centralWidget()->layout()->activate();
        this->centralWidget()->layout()->itemAt(0)->widget()->deleteLater();
        //this->centralWidget()->layout()->removeWidget(this->centralWidget()->layout()->itemAt(0)->widget());
    }

    m_loginpage = new LoginPage(this, m_socketManager);
    this->centralWidget()->layout()->addWidget(m_loginpage);
}

void MainWindow::ShowMainChatPage()
{
    qDebug() << "MainWindow::ShowMainChatPage";
    if (!this->centralWidget()->layout()->isEmpty()) {
        qDebug() << "Layout not empty!: " << this->centralWidget()->layout()->activate();
        this->centralWidget()->layout()->itemAt(0)->widget()->deleteLater();
        //this->centralWidget()->layout()->removeWidget(this->centralWidget()->layout()->itemAt(0)->widget());
    }

    m_mainChatWindow = new MainChatWindow(this);
    this->centralWidget()->layout()->addWidget(m_mainChatWindow);
}
