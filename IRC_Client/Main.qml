//import QtQuick
//import QtQuick.Window

//Window {
//    width: 640
//    height: 480
//    visible: true
//    title: qsTr("Hello World")
//}
import QtQuick
import QtQuick.Layouts
import QtQuick.Controls
import QtQuick.Controls.Material
import IRC_Client


ApplicationWindow {
    id: login_window
    width: 360
    height: 520
    visible: true
    title: qsTr("Login")

    SocketManager {
        id: socketmanager
        Component.onCompleted: {
            socketmanager.ServerConnect();
        }
    }

//    BackEnd {
//        id: myBackend
//        onNumberEmitted: (num) => {
//                            console.log(num)
//                            randNum.text = num
//                         }
//        Component.onCompleted: {
//            myBackend.generateNumber(1,100);
//        }
//    }

    Loader {
        id: pageLoader
    }

//    Image {
//        id: login_logo
//        anchors.horizontalCenter: parent.horizontalCenter
//        source: ":/qt-logo.png"
//        width: 100
//        height: 100
//    }

    Rectangle {
        id: test
        color: "#95e295"
        anchors.horizontalCenter: parent.horizontalCenter
        width: 100
        height: 100
    }

    TextField {
        id: email_input
        anchors.top: test.bottom
        anchors.horizontalCenter: test.horizontalCenter
        placeholderText: "email"
        cursorVisible: true
        leftPadding: 4
        topPadding: 4
        echoMode: TextInput.Normal
    }

    TextField {
        id: password_input
        placeholderText: "password"
        anchors.top: email_input.bottom
        anchors.left: email_input.left
        cursorVisible: true
        leftPadding: 4
        topPadding: 2
        echoMode: TextInput.Password
    }

    RoundButton {
        id: forgot_password_btn
        text: qsTr("?")
        width: 30
        height: 30
        anchors.left: password_input.right
        anchors.verticalCenter: password_input.verticalCenter
        onClicked: {
            popup.openWithContent("Password Reset Sent", "info")
        }
    }

    Button {
        id: login_btn
        anchors.top: password_input.bottom
        anchors.left: password_input.left
        width: password_input.width
        text: qsTr("Login")
        onClicked: {
            console.log("Login btn clicked")
            var err = "";
            if (email_input.text === "") {
                console.log("Need email input!")
                err += "Need Email\n"
            }

            if (password_input.text === "") {
                console.log("Need password")
                err += "Need password"
            }

            if (err !== "") {
                popup.openWithContent(err, "error")
            } else {
                // TODO what IRC command is this? LOGIN or AUTHENTICATE? Probably AUTHENTICATE
                const l_msg = `LOGIN ${email_input.text} ${password_input.text}`
                const l_final_data = `:${nickname}@${localAddr} ${l_msg} \r\n`
//                    socketmanager.Write_Data(l_final_data.toString("utf8"))
            }
        }
    }

    Button {
        id: register_btn
        anchors.top: login_btn.bottom
        anchors.left: login_btn.left
        width: login_btn.width
        text: qsTr("Register")
        onClicked: {
            login_window.hide()
            pageLoader.setSource("register.qml",
                                 {x: login_window.x, y: login_window.y})
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

//        Label {
//            id: randNum
//            text: ""
//        }
}
