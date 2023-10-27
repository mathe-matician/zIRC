#ifndef LOGINPAGE_H
#define LOGINPAGE_H

#include <QtQml>
#include "socketmanager.h"

namespace Ui {
class LoginPage;
}

class LoginPage : public QObject
{
    Q_OBJECT
    QML_ELEMENT
public:
    explicit LoginPage(QObject *parent = nullptr, SocketManager *a_socketManager = nullptr);
//    ~LoginPage();

    SocketManager *m_socketManager;

public slots:
    void Login();
    void Register();
    void ForgotPW();

private:
    Ui::LoginPage *ui;
};

#endif // LOGINPAGE_H
