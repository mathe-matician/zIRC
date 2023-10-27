#include "backend.h"
#include <QDebug>

BackEnd::BackEnd(QObject *parent) :
    QObject(parent)
{
}

void BackEnd::generateNumber(int min, int max)
{
    const int randNum = QRandomGenerator::global()->bounded(min,max);
    emit numberEmitted(randNum);
}
