#ifndef MAINCHATWINDOW_H
#define MAINCHATWINDOW_H

#include <QWidget>
#include <QTreeWidgetItem>
#include <QVBoxLayout>

#include "socketmanager.h"
#include "chatchannelselectionwindow.h"
#include "chatinboxwindow.h"

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

    QVBoxLayout *m_vLayout;
    ChatChannelSelectionWindow *m_channelSelectionWindow;
    ChatInboxWindow *m_inboxWindow = nullptr;

public slots:
    void MenuItemDoubleClicked(QTreeWidgetItem *a_item, int column);

private:
    Ui::MainChatWindow *ui;
    SocketManager *m_socketManager;
};

#endif // MAINCHATWINDOW_H
