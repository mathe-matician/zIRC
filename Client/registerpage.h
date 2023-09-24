#ifndef REGISTERPAGE_H
#define REGISTERPAGE_H

#include <QWidget>

namespace Ui {
class Registerpage;
}

class Registerpage : public QWidget
{
    Q_OBJECT

public:
    explicit Registerpage(QWidget *parent = nullptr);
    ~Registerpage();

private:
    Ui::Registerpage *ui;
};

#endif // REGISTERPAGE_H
