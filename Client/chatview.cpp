#include "chatview.h"
#include "ui_chatview.h"

ChatView::ChatView(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::ChatView)
{
    ui->setupUi(this);
    ui->plainTextEdit->setReadOnly(true);
}

ChatView::~ChatView()
{
    delete ui;
}
