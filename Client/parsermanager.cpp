#include "parsermanager.h"
#include <QDebug>

ParserManager::ParserManager()
{

}

ParserManager::ParserManager(QByteArray data)
{
    qDebug() << "ParserManager::ParserManager start";
    qDebug() << "Data: " << data.toStdString();

    // Messages should be prefixed with :
    // e.g. a server response should look like:
    // "undefined undefined0 400 :Unknown error occurred"
}

//~ParserManager();
