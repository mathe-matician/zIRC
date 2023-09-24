#ifndef VERIFYEMAILPAGE_H
#define VERIFYEMAILPAGE_H

#include <QWidget>

namespace Ui {
class VerifyEmailPage;
}

class VerifyEmailPage : public QWidget
{
    Q_OBJECT

public:
    explicit VerifyEmailPage(QWidget *parent = nullptr);
    ~VerifyEmailPage();

private:
    Ui::VerifyEmailPage *ui;
};

#endif // VERIFYEMAILPAGE_H
