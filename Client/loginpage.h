#ifndef LOGINPAGE_H
#define LOGINPAGE_H

#include <QWidget>
#include "socketmanager.h"

namespace Ui {
class LoginPage;
}

class LoginPage : public QWidget
{
    Q_OBJECT

public:
    explicit LoginPage(QWidget *parent = nullptr);
    ~LoginPage();

    SocketManager *m_socketManager;

public slots:
    void Login();
    void Register();
    void ForgotPW();

private:
    Ui::LoginPage *ui;
};

#endif // LOGINPAGE_H
