#ifndef CHATCHANNELSELECTIONWINDOW_H
#define CHATCHANNELSELECTIONWINDOW_H

#include <QWidget>

namespace Ui {
class ChatChannelSelectionWindow;
}

class ChatChannelSelectionWindow : public QWidget
{
    Q_OBJECT

public:
    explicit ChatChannelSelectionWindow(QWidget *parent = nullptr);
    ~ChatChannelSelectionWindow();

private:
    Ui::ChatChannelSelectionWindow *ui;
};

#endif // CHATCHANNELSELECTIONWINDOW_H
