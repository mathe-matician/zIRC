#include "registerpage.h"
#include "ui_registerpage.h"

Registerpage::Registerpage(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::Registerpage)
{
    ui->setupUi(this);
}

Registerpage::~Registerpage()
{
    delete ui;
}
