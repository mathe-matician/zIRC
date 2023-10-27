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

    Rectangle {
        id: image_placeholder
        color: "#95e295"
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.topMargin: 30
        width: 100
        height: 100
    }

    TextField {
        id: email_input
        anchors.top: image_placeholder.bottom
        anchors.horizontalCenter: image_placeholder.horizontalCenter
        anchors.topMargin: 30
        width: parent.width / 2
        placeholderText: "email"
        cursorVisible: true
        leftPadding: 4
        topPadding: 4
        echoMode: TextInput.Normal
        maximumLength: 256
    }

    TextField {
        id: password_input
        placeholderText: "password"
        anchors.top: email_input.bottom
        anchors.left: email_input.left
        anchors.topMargin: 10
        width: parent.width / 2
        cursorVisible: true
        leftPadding: 4
        topPadding: 2
        echoMode: TextInput.Password
        maximumLength: 256
    }

    TextField {
        id: password_confirm_input
        placeholderText: "confirm password"
        anchors.top: password_input.bottom
        anchors.left: password_input.left
        anchors.topMargin: 10
        width: parent.width / 2
        cursorVisible: true
        leftPadding: 4
        topPadding: 2
        echoMode: TextInput.Password
        maximumLength: 256

        Connections {
            target: password_confirm_input
            onEditingFinished: password_input.text !== password_confirm_input.text ? console.log("PASSWORDS DONT MATCH") : console.log("passwords match!")
        }
    }

    Button {
        id: register_btn
        anchors.top: password_confirm_input.bottom
        anchors.left: password_confirm_input.left
        anchors.topMargin: 10
        width: login_btn.width
        text: qsTr("Register")
        onClicked: {
            pageLoader.setSource("chatview.qml",
                                 {x: register_window.x, y: register_window.y})
        }
    }
}
