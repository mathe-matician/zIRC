#include "chatinboxwindow.h"
#include "ui_chatinboxwindow.h"

ChatInboxWindow::ChatInboxWindow(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::ChatInboxWindow)
{
    ui->setupUi(this);
}

ChatInboxWindow::~ChatInboxWindow()
{
    delete ui;
}
