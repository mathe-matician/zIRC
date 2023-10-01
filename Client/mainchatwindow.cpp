#include "mainchatwindow.h"
#include "ui_mainchatwindow.h"

MainChatWindow::MainChatWindow(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::MainChatWindow)
{
    ui->setupUi(this);
}

MainChatWindow::~MainChatWindow()
{
    delete ui;
}
