#include "verifyemailpage.h"
#include "ui_verifyemailpage.h"

VerifyEmailPage::VerifyEmailPage(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::VerifyEmailPage)
{
    ui->setupUi(this);
}

VerifyEmailPage::~VerifyEmailPage()
{
    delete ui;
}
