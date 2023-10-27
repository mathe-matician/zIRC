import QtQuick
import QtQuick.Layouts
import QtQuick.Controls
import QtQuick.Controls.Material

ApplicationWindow {
    id: register_window
    width: 360
    height: 520
    visible: true
    title: qsTr("Register")

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
                maximumLength: 256
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
                maximumLength: 256
            }
        }

        RowLayout {
            Label {
                id: password_confirm_label
                Layout.fillWidth: true
                text: qsTr("Confirm Password")
            }
            TextField {
                id: password_confirm_input
                Layout.fillWidth: true
                leftPadding: 4
                cursorVisible: true
                topPadding: 2
                echoMode: TextInput.Password
                maximumLength: 256
            }
        }

        Button {
            id: register_btn
            text: qsTr("Register")
            onClicked: {
                pageLoader.setSource("chatview.qml",
                                     {x: register_window.x, y: register_window.y})
            }
        }

        Item {
            // spacer item
            Layout.fillWidth: true
            Layout.fillHeight: true
            Rectangle { anchors.fill: parent; color: "#ffaaaa" } // to visualize the spacer
        }
    }
}
