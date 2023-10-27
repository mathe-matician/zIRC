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
                topPadding: 4
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

                if (err !== "")
                    popup.openWithContent(err, "error")
//                else
//                    // TODO what IRC command is this? LOGIN or AUTHENTICATE? Probably AUTHENTICATE
//                    const l_msg = `LOGIN ${email_input.text} ${password_input.text}`
//                    const l_final_data = `:${nickname}@${localAddr} ${l_msg} \r\n`
//                    socketmanager.Write_Data(l_final_data.toString("utf8"))
            }
        }

        Button {
            id: register_btn
            text: qsTr("Register")
            onClicked: {
                login_window.hide()
                pageLoader.setSource("register.qml",
                                     {x: login_window.x, y: login_window.y})
            }
        }

        Button {
            id: callfunc
            text: qsTr("Forgot Password")
//            onClicked: {
//                myBackend.generateNumber(1,100)
//            }
            onClicked: {
                popup.openWithContent("TUFF LUK", "error")
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
                        horizontalAlignment: Text.AlignHCenter
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
//            contentItem: popup.popup_content

            closePolicy: Popup.CloseOnPressOutside
        }

//        Label {
//            id: randNum
//            text: ""
//        }

        Item {
            // spacer item
            Layout.fillWidth: true
            Layout.fillHeight: true
            Rectangle { anchors.fill: parent; color: "#ffaaaa" } // to visualize the spacer
        }
    }
}
