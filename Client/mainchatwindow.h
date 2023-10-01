#ifndef MAINCHATWINDOW_H
#define MAINCHATWINDOW_H

#include <QWidget>

#include "socketmanager.h"

namespace Ui {
class MainChatWindow;
}

class MainChatWindow : public QWidget
{
    Q_OBJECT

public:
    explicit MainChatWindow(QWidget *parent = nullptr, SocketManager *a_socketManager = nullptr);
    ~MainChatWindow();
    bool eventFilter(QObject *object, QEvent *event);

private:
    Ui::MainChatWindow *ui;
    SocketManager *m_socketManager;
};

#endif // MAINCHATWINDOW_H
