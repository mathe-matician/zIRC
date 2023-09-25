#include "registerpage.h"
#include "ui_registerpage.h"

Registerpage::Registerpage(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::Registerpage)
{
    ui->setupUi(this);

    connect(ui->BTN_back, SIGNAL(clicked()), this->parentWidget(), SLOT(ShowLoginPage()));
}

Registerpage::~Registerpage()
{
    delete ui;
}
