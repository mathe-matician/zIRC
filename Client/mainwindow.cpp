#include "mainwindow.h"
#include "./ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    g_settings = new QSettings();

    m_startpage = new StartPage(g_settings, this);
    m_startpage->show();
}

MainWindow::~MainWindow()
{
    delete ui;
}

