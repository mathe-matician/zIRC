cmake_minimum_required(VERSION 3.16)

project(IRC_Client VERSION 0.1 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD_REQUIRED ON)

set(QT_QUICK_CONTROLS_CONF "qtquickcontrols2.conf")

find_package(Qt6 6.5 REQUIRED COMPONENTS Core Gui Quick NetworkAuth)

qt_standard_project_setup(REQUIRES 6.5)

qt_add_executable(appIRC_Client
    main.cpp
)
#
qt_add_qml_module(appIRC_Client
    URI IRC_Client
    VERSION 1.0
    QML_FILES
        Main.qml
    SOURCES
        backend.h backend.cpp
        socketmanager.h socketmanager.cpp
    RESOURCES
        "qtquickcontrols2.conf"
        QML_FILES register.qml
        QML_FILES chatview.qml
        QML_FILES MessageModel.qml
        QML_FILES popupError.qml
        RESOURCES resc.qrc
        SOURCES googlesso.h googlesso.cpp
        SOURCES authgen.h authgen.cpp

)

# Qt for iOS sets MACOSX_BUNDLE_GUI_IDENTIFIER automatically since Qt 6.1.
# If you are developing for iOS or macOS you should consider setting an
# explicit, fixed bundle identifier manually though.
set_target_properties(appIRC_Client PROPERTIES
#    MACOSX_BUNDLE_GUI_IDENTIFIER com.example.appIRC_Client
    MACOSX_BUNDLE_BUNDLE_VERSION ${PROJECT_VERSION}
    MACOSX_BUNDLE_SHORT_VERSION_STRING ${PROJECT_VERSION_MAJOR}.${PROJECT_VERSION_MINOR}
    MACOSX_BUNDLE TRUE
    WIN32_EXECUTABLE TRUE
)

target_link_libraries(appIRC_Client
    PRIVATE Qt6::Core
    PRIVATE Qt6::Gui
    PRIVATE Qt6::Quick
    PRIVATE Qt6::NetworkAuth
)

include(GNUInstallDirs)
install(TARGETS appIRC_Client
    BUNDLE DESTINATION .
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)
