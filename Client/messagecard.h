#ifndef MESSAGECARD_H
#define MESSAGECARD_H

#include <QWidget>

namespace Ui {
class MessageCard;
}

class MessageCard : public QWidget
{
    Q_OBJECT

public:
    explicit MessageCard(QWidget *parent = nullptr);
    ~MessageCard();

//public slots:
    //void

private:
    Ui::MessageCard *ui;
};

#endif // MESSAGECARD_H
