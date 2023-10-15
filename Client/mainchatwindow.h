#ifndef MAINCHATWINDOW_H
#define MAINCHATWINDOW_H

#include <QWidget>
#include <QTreeWidgetItem>
#include <QVBoxLayout>

#include "socketmanager.h"
#include "chatchannelselectionwindow.h"
#include "chatinboxwindow.h"
#include "chatbox.h"
#include "chatview.h"

namespace Ui {
class MainChatWindow;
}

class MainChatWindow : public QWidget
{
    Q_OBJECT

public:
    explicit MainChatWindow(QWidget *parent = nullptr, SocketManager *a_socketManager = nullptr);
    ~MainChatWindow();

    enum SaveFormat { Json, Binary };

    bool m_loadState(SaveFormat saveFormat);
    const bool m_saveState(SaveFormat saveFormat);

    QVBoxLayout *m_vLayout;
    ChatChannelSelectionWindow *m_channelSelectionWindow = nullptr;
    ChatInboxWindow *m_inboxWindow = nullptr;
    ChatBox *m_chatBox = nullptr;
    ChatView *m_chatView = nullptr;

public slots:
    void MenuItemDoubleClicked(QTreeWidgetItem *a_item, int column);

private:
    Ui::MainChatWindow *ui;
    SocketManager *m_socketManager;
    QString m_currentChan;
};

#endif // MAINCHATWINDOW_H
