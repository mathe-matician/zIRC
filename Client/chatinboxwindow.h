#ifndef CHATINBOXWINDOW_H
#define CHATINBOXWINDOW_H

#include <QWidget>

namespace Ui {
class ChatInboxWindow;
}

class ChatInboxWindow : public QWidget
{
    Q_OBJECT

public:
    explicit ChatInboxWindow(QWidget *parent = nullptr);
    ~ChatInboxWindow();

private:
    Ui::ChatInboxWindow *ui;
};

#endif // CHATINBOXWINDOW_H
