#ifndef STARTPAGE_H
#define STARTPAGE_H

#include <QWidget>
#include <QSettings>
#include "loginpage.h"
#include "registerpage.h"
#include "socketmanager.h"

namespace Ui {
class StartPage;
}

class StartPage : public QWidget
{
    Q_OBJECT

public:
    explicit StartPage(QSettings *g_settings, QWidget *parent = nullptr);
    ~StartPage();

    LoginPage *m_loginpage;
    Registerpage *m_registerpage;
    SocketManager *m_socketManager;

public slots:
    void ShowLoginPage();
    void ShowRegisterPage();

private:
    Ui::StartPage *ui;
};

#endif // STARTPAGE_H
