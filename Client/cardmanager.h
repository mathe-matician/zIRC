#ifndef CARDMANAGER_H
#define CARDMANAGER_H

#include <QObject>

#include "socketmanager.h"

class CardManager
{
    Q_OBJECT
public:
    CardManager(SocketManager *a_socketManager);
};

#endif // CARDMANAGER_H
