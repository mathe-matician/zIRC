#include "messagecard.h"
#include "ui_messagecard.h"

MessageCard::MessageCard(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::MessageCard)
{
    ui->setupUi(this);

    connect(ui->BTN, SIGNAL(clicked()), this, SLOT(deleteLater()));
}

MessageCard::~MessageCard()
{
    delete ui;
}
