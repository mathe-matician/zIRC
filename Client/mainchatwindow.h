#ifndef MAINCHATWINDOW_H
#define MAINCHATWINDOW_H

#include <QWidget>

namespace Ui {
class MainChatWindow;
}

class MainChatWindow : public QWidget
{
    Q_OBJECT

public:
    explicit MainChatWindow(QWidget *parent = nullptr);
    ~MainChatWindow();

private:
    Ui::MainChatWindow *ui;
};

#endif // MAINCHATWINDOW_H
