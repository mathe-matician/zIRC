#include "mainchatwindow.h"
#include "ui_mainchatwindow.h"

#include <QTreeWidgetItem>
#include <QFileSystemModel>
#include <QJsonDocument>
#include <QJsonObject>
#include <QCborValue>
#include <QCborMap>
//#include <QValidator>

MainChatWindow::MainChatWindow(QWidget *parent, SocketManager *a_socketManager) :
    QWidget(parent),
    ui(new Ui::MainChatWindow)
{
    ui->setupUi(this);
    m_vLayout = ui->verticalLayout;
    m_channelSelectionWindow = new ChatChannelSelectionWindow();
    m_vLayout->insertWidget(0, m_channelSelectionWindow);
    ui->treeWidget->setFixedWidth(160);

    m_socketManager = a_socketManager;

    // user the event filter defined in this object
    this->installEventFilter(this);

    // Create new item (top level item)
    QTreeWidgetItem *topLevelItem = new QTreeWidgetItem(ui->treeWidget);
    // Add it on our tree as the top item.
    ui->treeWidget->addTopLevelItem(topLevelItem);
    // Set text for item
    topLevelItem->setText(0,"Channels");
    // Create new item and add as child item
    QTreeWidgetItem *item=new QTreeWidgetItem(topLevelItem);
    // Set text for item
    item->setText(0,"General");

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

bool MainChatWindow::m_loadState(SaveFormat saveFormat)
{
    QFile loadFile(saveFormat == Json ? "state.json" : "state.dat");

    if (!loadFile.open(QIODevice::ReadOnly)) {
        qWarning("Couldn't open save file.");
        return false;
    }

    QByteArray saveData = loadFile.readAll();

    QJsonDocument loadDoc(saveFormat == Json
                              ? QJsonDocument::fromJson(saveData)
                              : QJsonDocument(QCborValue::fromCbor(saveData).toMap().toJsonObject()));

    //read(loadDoc.object());

    QTextStream(stdout) << "Loaded save for " << loadDoc["player"]["name"].toString()
                        << " using " << (saveFormat != Json ? "CBOR" : "JSON") << "...\n";
    return true;
}

const bool MainChatWindow::m_saveState(SaveFormat saveFormat)
{
    QFile saveFile(saveFormat == Json ? "save.json" : "save.dat");

    if (!saveFile.open(QIODevice::WriteOnly)) {
        qWarning("Couldn't open save file.");
        return false;
    }

    //QJsonObject gameObject = toJson();
    //saveFile.write(saveFormat == Json ? QJsonDocument(gameObject).toJson()
                                      //: QCborValue::fromJsonValue(gameObject).toCbor());

    return true;
}

void MainChatWindow::MenuItemDoubleClicked(QTreeWidgetItem *a_item, int column)
{
    qDebug() << "MainChatWindow::MenuItemDoubleClicked start";
    if (a_item != nullptr) {
        QString l_selectedItem = a_item->text(column);
        qDebug() << "Text == " << l_selectedItem;

        if (l_selectedItem == "Channels") {
            qDebug() << "Channels Selected";

            if (!m_vLayout->layout()->isEmpty() && m_vLayout->layout()->itemAt(0)->widget()->objectName() != "ChatChannelSelectionWindow") {
                qDebug() << "Layout not empty!: ";
                if (m_vLayout->layout()->itemAt(0)->widget()->objectName() == "ChatInboxWindow") {
                    m_inboxWindow = nullptr;
                }

                if (m_vLayout->layout()->itemAt(0)->widget()->objectName() == "ChatView") {
                    qDebug() << "l_chatbox != nullptr";
                    m_chatBox = nullptr;
                    m_vLayout->layout()->itemAt(1)->widget()->deleteLater();
                    m_vLayout->layout()->removeWidget(m_vLayout->layout()->itemAt(1)->widget());
                }

                m_vLayout->layout()->itemAt(0)->widget()->deleteLater();
                m_vLayout->layout()->removeWidget(m_vLayout->layout()->itemAt(0)->widget());

                if (m_channelSelectionWindow == nullptr) {
                    m_channelSelectionWindow = new ChatChannelSelectionWindow();
                }

                m_vLayout->insertWidget(0, m_channelSelectionWindow);
                m_currentChan = "";
            }
        } else if (l_selectedItem == "Inbox") {
            qDebug() << "Inbox Selected";

            if (!m_vLayout->layout()->isEmpty() && m_vLayout->layout()->itemAt(0)->widget()->objectName() != "ChatInboxWindow") {
                qDebug() << "Layout not empty!: ";
                qDebug() << "ACtive objec tname = " << m_vLayout->layout()->itemAt(0)->widget()->objectName();
                if (m_vLayout->layout()->itemAt(0)->widget()->objectName() == "ChatChannelSelectionWindow") {
                    m_channelSelectionWindow = nullptr;
                }

                if (m_vLayout->layout()->itemAt(0)->widget()->objectName() == "ChatView") {
                    // if ChatView exists, ChatBox will always exist as well, so remove the ChatBox
                    qDebug() << "l_chatbox != nullptr";
                    m_chatBox = nullptr;
                    m_vLayout->layout()->itemAt(1)->widget()->deleteLater();
                    m_vLayout->layout()->removeWidget(m_vLayout->layout()->itemAt(1)->widget());
                }

                m_vLayout->layout()->itemAt(0)->widget()->deleteLater();
                m_vLayout->layout()->removeWidget(m_vLayout->layout()->itemAt(0)->widget());

                if (m_inboxWindow == nullptr) {
                    m_inboxWindow = new ChatInboxWindow();
                }

                m_vLayout->insertWidget(0, m_inboxWindow);
                m_currentChan = "";
            }
        } else {
            // else some sub menu was selected so open up a view for that
            qDebug() << "Other selected: " << l_selectedItem;

            if (l_selectedItem == m_currentChan) {
                qDebug() << "Channel already selected";
                return;
            }

            m_currentChan = l_selectedItem;

            QString l_activeWidget = m_vLayout->layout()->itemAt(0)->widget()->objectName();
            if (l_activeWidget == "ChatChannelSelectionWindow") {
                m_channelSelectionWindow = nullptr;
            } else if (l_activeWidget == "ChatInboxWindow") {
                m_inboxWindow = nullptr;
            }

            m_vLayout->layout()->itemAt(0)->widget()->deleteLater();

            // TODO
            // Insert chat view of proposed chat
            // Get cached chat from database
            // get chat values from server if no cache exists
            m_chatView = new ChatView();
            m_vLayout->insertWidget(0, m_chatView);

            // TODO
            // nullout ChatView in other views above

            if (m_vLayout->layout()->findChild<QWidget *>(QString("ChatBox")) == nullptr) {
                m_chatBox = new ChatBox();
                m_vLayout->insertWidget(1, m_chatBox);
            }
        }
    } else {
        qDebug() << "Nothing selected";
    }
}
