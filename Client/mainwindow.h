#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QSettings>

#include "loginpage.h"
#include "registerpage.h"
#include "socketmanager.h"
#include "mainchatwindow.h"

QT_BEGIN_NAMESPACE
namespace Ui { class MainWindow; }
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();
    QSettings *g_settings;

    LoginPage *m_loginpage;
    Registerpage *m_registerpage;
    SocketManager *m_socketManager;
    MainChatWindow *m_mainChatWindow;

public slots:
    void ShowLoginPage();
    void ShowRegisterPage();
    void ShowMainChatPage();

private:
    void IsPageActive();

    QString m_tokenPkg;
    Ui::MainWindow *ui;

};
#endif // MAINWINDOW_H
