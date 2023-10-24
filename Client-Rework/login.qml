import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

ApplicationWindow {
    id: login_window
    width: 360
    height: 520
    visible: true
    title: qsTr("Login")

    Loader {
        id: pageLoader
    }

    ColumnLayout {
        anchors.fill: parent

        RowLayout {
            Label {
                id: email_label
                Layout.fillWidth: true
                text: qsTr("Email")
            }
            TextField {
                id: email_input
                Layout.fillWidth: true
                cursorVisible: true
                leftPadding: 4
                topPadding: 2
                echoMode: TextInput.Normal
            }
        }

        RowLayout {
            Label {
                id: password_label
                Layout.fillWidth: true
                text: qsTr("Password")
            }
            TextField {
                id: password_input
                Layout.fillWidth: true
                cursorVisible: true
                leftPadding: 4
                topPadding: 2
                echoMode: TextInput.Password
            }
        }

        Button {
            id: next_page_btn
            text: qsTr("Login")
            onClicked: console.log("Login btn clicked")
        }

        Button {
            id: register_btn
            text: qsTr("Register")
            onClicked: pageLoader.source = "register.qml"
        }

        Item {
            // spacer item
            Layout.fillWidth: true
            Layout.fillHeight: true
            Rectangle { anchors.fill: parent; color: "#ffaaaa" } // to visualize the spacer
        }
    }
}
