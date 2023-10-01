#include "mainchatwindow.h"
#include "ui_mainchatwindow.h"

#include <QTreeWidgetItem>
#include <QFileSystemModel>
//#include <QValidator>

MainChatWindow::MainChatWindow(QWidget *parent, SocketManager *a_socketManager) :
    QWidget(parent),
    ui(new Ui::MainChatWindow)
{
    ui->setupUi(this);
    m_vLayout = ui->verticalLayout;
    // insertWidget(index)
    m_channelSelectionWindow = new ChatChannelSelectionWindow();
    m_vLayout->insertWidget(0, m_channelSelectionWindow);
    ui->treeWidget->setFixedWidth(160);

    m_socketManager = a_socketManager;

    // user the event filter defined in this object
    this->installEventFilter(this);

    // chatbox
    ui->chatBox->setFixedHeight(40);

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

    QTreeWidgetItem *inbox = new QTreeWidgetItem(ui->treeWidget);
    // Add it on our tree as the top item.
    ui->treeWidget->addTopLevelItem(inbox);
    // Set text for item
    inbox->setText(0,"Inbox");


    //itemDoubleClicked(QTreeWidgetItem *item, int column)
    //connect(ui->treeWidget, SIGNAL(itemDoubleClicked(QTreeWidgetItem*, int)), this, SLOT(MenuItemDoubleClicked(QTreeWidgetItem*, int)));
    connect(ui->treeWidget, SIGNAL(itemClicked(QTreeWidgetItem*, int)), this, SLOT(MenuItemDoubleClicked(QTreeWidgetItem*, int)));
}

MainChatWindow::~MainChatWindow()
{
    delete ui;
}

bool MainChatWindow::eventFilter(QObject* obj, QEvent* event)
{
    //qDebug() << "eventFilter::event = " << event->type();

    if (this->focusWidget() && this->focusWidget()->objectName() == "chatBox") {
        qDebug() << "FOCUS == chatBox";
        if (event->type()==QEvent::KeyRelease || event->type()==QEvent::KeyPress) {
            QKeyEvent* key = static_cast<QKeyEvent*>(event);

            if ((key->modifiers() & Qt::ShiftModifier) && (key->key() == Qt::Key_Enter || key->key() == Qt::Key_Return)) {
                qDebug() << "Shift + Enter pressed - create newline";
            } else if ((key->key()==Qt::Key_Enter) || (key->key()==Qt::Key_Return)) {
                qDebug() << "KeyRelease was enter or return";
                event->ignore();
                if (!ui->chatBox->document()->isEmpty()) {
                    qDebug() << "chatbox NOT EMPTY";
                    m_socketManager->Debug_Send(ui->chatBox->document()->toPlainText().toUtf8());
                    ui->chatBox->clear();
                } else {
                    qDebug() << "chatbox EMPTY";
                }
            }
        }
    }

    /*
    if (event->type()==QEvent::KeyPress) {
        QKeyEvent* key = static_cast<QKeyEvent*>(event);
        if ( (key->key()==Qt::Key_Enter) || (key->key()==Qt::Key_Return) ) {
            //Enter or return was pressed
            qDebug() << "Enter or return was pressed";

            if (this->focusWidget() && this->focusWidget()->objectName() == "chatBox") {
                qDebug() << "chatbox focus AND enter pressed";

                if (!ui->chatBox->document()->isEmpty()) {
                    qDebug() << "chatbox empty second";
                }
            }

        } else {
            qDebug() << "Another key was pressed";
            return QObject::eventFilter(obj, event);
        }
        return true;
    } else {
        return QObject::eventFilter(obj, event);
    }
*/
    return false;
}

void MainChatWindow::MenuItemDoubleClicked(QTreeWidgetItem *a_item, int column)
{
    qDebug() << "MainChatWindow::MenuItemDoubleClicked start";
    if (a_item != nullptr) {
        QString l_selectedItem = a_item->text(column);
        qDebug() << "Text == " << l_selectedItem;

        if (l_selectedItem == "Channels") {
            qDebug() << "Channels Selected";

            if (!m_vLayout->layout()->isEmpty()) {
                // && m_vLayout->layout()->widget()->objectName() != ...
                // to avoid allocs
                qDebug() << "Layout not empty!: ";
                m_vLayout->layout()->itemAt(0)->widget()->deleteLater();
            }

            if (m_channelSelectionWindow == nullptr) {
                m_channelSelectionWindow = new ChatChannelSelectionWindow();
            }

            m_vLayout->insertWidget(0, m_channelSelectionWindow);

        } else if (l_selectedItem == "Inbox") {
            qDebug() << "Inbox Selected";

            if (!m_vLayout->layout()->isEmpty()) {
                qDebug() << "Layout not empty!: ";
                m_vLayout->layout()->itemAt(0)->widget()->deleteLater();
            }

            if (m_inboxWindow == nullptr) {
                m_inboxWindow = new ChatInboxWindow();
            }

            m_vLayout->insertWidget(0, m_inboxWindow);
        }
    } else {
        qDebug() << "Nothing selected";
    }
}
