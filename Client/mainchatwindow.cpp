#include "mainchatwindow.h"
#include "ui_mainchatwindow.h"
#include <QTreeWidgetItem>
#include <QFileSystemModel>

MainChatWindow::MainChatWindow(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::MainChatWindow)
{
    ui->setupUi(this);
    ui->treeWidget->setFixedWidth(160);

    // Create new item (top level item)
    QTreeWidgetItem *topLevelItem = new QTreeWidgetItem(ui->treeWidget);
    // Add it on our tree as the top item.
    ui->treeWidget->addTopLevelItem(topLevelItem);
    // Set text for item
    topLevelItem->setText(0,"Channels");
    // Create new item and add as child item
    QTreeWidgetItem *item=new QTreeWidgetItem(topLevelItem);
    // Set text for item
    item->setText(0,"#General");

    QTreeWidgetItem *dms = new QTreeWidgetItem(ui->treeWidget);
    // Add it on our tree as the top item.
    ui->treeWidget->addTopLevelItem(dms);
    // Set text for item
    topLevelItem->setText(1,"dms");
}

MainChatWindow::~MainChatWindow()
{
    delete ui;
}
