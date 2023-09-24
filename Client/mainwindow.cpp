#include "mainwindow.h"
#include "./ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    ui->setupUi(this);

    m_startpage = new StartPage(this);
    m_startpage->show();
}

MainWindow::~MainWindow()
{
    delete ui;
}

