#ifndef REGISTERPAGE_H
#define REGISTERPAGE_H

#include <QWidget>
#include "socketmanager.h"

namespace Ui {
class Registerpage;
}

class Registerpage : public QWidget
{
    Q_OBJECT

public:
    explicit Registerpage(QWidget *parent = nullptr, SocketManager *a_socketManager = nullptr);
    ~Registerpage();

    SocketManager *m_socketManager;

public slots:
    void Register();

private:
    Ui::Registerpage *ui;
};

#endif // REGISTERPAGE_H
