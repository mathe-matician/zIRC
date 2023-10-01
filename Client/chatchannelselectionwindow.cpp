#include "chatchannelselectionwindow.h"
#include "ui_chatchannelselectionwindow.h"

ChatChannelSelectionWindow::ChatChannelSelectionWindow(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::ChatChannelSelectionWindow)
{
    ui->setupUi(this);
}

ChatChannelSelectionWindow::~ChatChannelSelectionWindow()
{
    delete ui;
}
