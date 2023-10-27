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

    MouseArea {
        id: base_mousearea
        anchors.fill: parent
        onClicked: forceActiveFocus()
    }

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
        placeholderText: email_input.text.length === 0 ? "<font color=\"red\">*</font>email" : "email"
        cursorVisible: true
        leftPadding: 4
        topPadding: 4
        echoMode: TextInput.Normal
        maximumLength: 256
    }

    TextField {
        id: password_input
        placeholderText: password_input.text.length === 0 ? "<font color=\"red\">*</font>password" : "password"
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
        placeholderText: password_confirm_input.text.length === 0 ? "<font color=\"red\">*</font>confirm password" : "confirm password"
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
            function onEditingFinished() {
                if (password_input.text !== password_confirm_input.text) {
                    //password_confirm_input.placeholderText = "<font color=\"red\">confirm password</font>"
                    console.log("PASSWORDS DO NOT MATCH AND CANT BE EMPTY")
                } else {
                    console.log(`passwords match: ${password_input.text} and ${password_confirm_input.text}`)
                    console.log(`passwords match: ${typeof(password_input.text)} and ${password_confirm_input.text.length}`)
                }
            }
        }
    }

    Button {
        id: register_btn
        anchors.top: password_confirm_input.bottom
        anchors.left: password_confirm_input.left
        anchors.topMargin: 10
        width: login_btn.width
        text: qsTr("Register")
        enabled: email_input.text.length !== 0 && password_input.text.length !== 0 && password_confirm_input.text.length !== 0 ? true : false
        onClicked: {
            if (password_input.text.length === 0
                || password_confirm_input.text.length === 0) {
                console.log("PASSWORDS CANT BE EMPTY")
                popup.openWithContent("Password fields are empty", "error")
                return
            }

            if (email_input.text.length === 0) {
                console.log("EMAIL CANT BE EMPTY")
                popup.openWithContent("Email field is empty", "error")
                return
            }

            if (password_input.text !== password_confirm_input.text) {
                console.log("PASSWORDS DO NOT MATCH")

                popup.openWithContent("Passwords do not match", "error")
            } else {
                console.log(`passwords match: ${password_input.text} and ${password_confirm_input.text}`)
                console.log(`passwords match: ${typeof(password_input.text)} and ${password_confirm_input.text.length}`)
                pageLoader.setSource("chatview.qml",
                                     {x: register_window.x, y: register_window.y})
            }
        }
    }

    Popup {
        property string popup_text
        property string border_color
        property string popup_title

        id: popup
        width: 200
        height: 200
        modal: true
        focus: true
        // centers
        anchors.centerIn: Overlay.overlay

        function openWithContent(text, type) {
            popup.popup_text = text
            if (type === "error") {
               popup.border_color = "red"
               popup.popup_title = "Error"
            } else {
                popup.border_color = "white"
                popup.popup_title = "Alert"
            }

            popup.open()
        }

        contentItem: Item{
            ColumnLayout {
                Text {
                    id: content_title
                    text: "<h1>" + popup.popup_title + "</h1>"
                }

                Text {
                    id: content_body
                    text: popup.popup_text
                }
            }
        }

        background: Rectangle {
            id: background
            color: "white"
            border.color: popup.border_color
            border.width: 3
            radius: 7
        }

        closePolicy: Popup.CloseOnPressOutside
    }
}
