#ifndef CHATBOX_H
#define CHATBOX_H

#include <QWidget>

#include "socketmanager.h"

namespace Ui {
class ChatBox;
}

class ChatBox : public QWidget
{
    Q_OBJECT

public:
    explicit ChatBox(QWidget *parent = nullptr, SocketManager *a_socketManager = nullptr);
    ~ChatBox();

    bool eventFilter(QObject *object, QEvent *event);

    SocketManager *m_socketManager;

private:
    Ui::ChatBox *ui;
};

#endif // CHATBOX_H
