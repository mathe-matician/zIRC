#ifndef CHATVIEW_H
#define CHATVIEW_H

#include <QWidget>

namespace Ui {
class ChatView;
}

class ChatView : public QWidget
{
    Q_OBJECT

public:
    explicit ChatView(QWidget *parent = nullptr);
    ~ChatView();

private:
    Ui::ChatView *ui;
};

#endif // CHATVIEW_H
