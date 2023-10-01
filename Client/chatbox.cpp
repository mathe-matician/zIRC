#include "chatbox.h"
#include "ui_chatbox.h"

ChatBox::ChatBox(QWidget *parent, SocketManager *a_socketManager) :
    QWidget(parent),
    ui(new Ui::ChatBox)
{
    ui->setupUi(this);

    qDebug() << "ChatBox created";
    this->setFixedHeight(70);

    m_socketManager = a_socketManager;

    this->installEventFilter(this);
}

ChatBox::~ChatBox()
{
    delete ui;
}

bool ChatBox::eventFilter(QObject* obj, QEvent* event)
{
    //qDebug() << "eventFilter::event = " << event->type();
    //qDebug() << "Focus = " << this->focusWidget();

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

    return false;
}
