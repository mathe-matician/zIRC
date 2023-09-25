#include "mainwindow.h"

#include <QApplication>

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    MainWindow w;
    w.show();

    QCoreApplication::setOrganizationName("Syllogi");
    QCoreApplication::setOrganizationDomain("syllogi.io");
    QCoreApplication::setApplicationName("IRC");

    return a.exec();
}
