#include "registerpage.h"
#include "ui_registerpage.h"

Registerpage::Registerpage(QWidget *parent, SocketManager *a_socketManager) :
    QWidget(parent),
    ui(new Ui::Registerpage)
{
    ui->setupUi(this);

    m_socketManager = a_socketManager;

    ui->INPUT_password->setEchoMode(QLineEdit::Password);

    connect(ui->BTN_back, SIGNAL(clicked()), this->parentWidget(), SLOT(ShowLoginPage()));
    connect(ui->BTN_register, SIGNAL(clicked()), this, SLOT(Register()));
}

Registerpage::~Registerpage()
{
    delete ui;
}

void Registerpage::Register()
{
    qDebug() << "Registerpage::Register Start";
    QString l_email = ui->INPUT_email->text();
    QString l_nickname = ui->INPUT_nick->text();
    QString l_password = ui->INPUT_password->text();

    QString l_msg = QString("REGISTER %1 %2 %3").arg(l_nickname).arg(l_email).arg(l_password);

    qDebug() << "REGISTER sending: " << l_msg;

    QString l_final_data = QString(":%1@%2 %3 \r\n").arg(l_nickname).arg(m_socketManager->socket->localAddress().toString()).arg(l_msg);
    qDebug() << "Before write";
    m_socketManager->Write_Data(l_final_data.toUtf8());
    qDebug() << "After write";
}
