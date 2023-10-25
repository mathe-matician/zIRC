import QtQuick.Window 2.2
import QtQuick
import QtQuick.Layouts
import QtQuick.Controls
import QtQuick.Controls.Material

ApplicationWindow {
    id: window
    width: 360
    height: 520
    visible: true
    title: qsTr("zIRC")

    //! [orientation]
    readonly property bool inPortrait: window.width < window.height
    //! [orientation]

    ToolBar {
        id: overlayHeader

        z: 1
        width: parent.width
        parent: Overlay.overlay

        RowLayout {
            anchors.fill: parent

            ToolButton {
                id: left_menu
                text: qsTr("⋮")
                anchors.left: parent
                onClicked: drawer.opened ? drawer.close() : drawer.open()
            }
            Label {
                id: label
                text: "#General"
                anchors.centerIn: parent
                elide: Label.ElideRight
                horizontalAlignment: Qt.AlignHCenter
                verticalAlignment: Qt.AlignVCenter
                Layout.fillWidth: true
            }
            ToolButton {
                id: right_menu
                text: qsTr("⋮")
                anchors.right: parent
                onClicked: drawer_right.opened ? drawer_right.close() : drawer_right.open()
            }
        }
    }

    Drawer {
        id: drawer
        edge: Qt.LeftEdge

        y: overlayHeader.height
        width: window.width / 2
        height: window.height - overlayHeader.height

        modal: inPortrait
        interactive: inPortrait
        position: inPortrait ? 0 : 1
        visible: !inPortrait

        ListView {
            id: listView
            anchors.fill: parent

            headerPositioning: ListView.OverlayHeader
            header: Pane {
                id: header
                z: 2
                width: parent.width

                contentHeight: logo.height

                Image {
                    id: logo
                    width: parent.width
                    source: "images/qt-logo.png"
                    fillMode: implicitWidth > width ? Image.PreserveAspectFit : Image.Pad
                }

                MenuSeparator {
                    parent: header
                    width: parent.width
                    anchors.verticalCenter: parent.bottom
                    visible: !listView.atYBeginning
                }
            }

            footer: ItemDelegate {
                id: footer
                text: qsTr("Footer")
                width: parent.width

                MenuSeparator {
                    parent: footer
                    width: parent.width
                    anchors.verticalCenter: parent.top
                }
            }

            model: 10

            delegate: ItemDelegate {
                text: qsTr("Title %1").arg(index + 1)
                width: listView.width
            }

            ScrollIndicator.vertical: ScrollIndicator { }
        }
    }

    Drawer {
        id: drawer_right

        edge: Qt.RightEdge

        y: overlayHeader.height
        width: window.width / 2
        height: window.height - overlayHeader.height

        modal: inPortrait
        interactive: inPortrait
        position: inPortrait ? 0 : 1
        visible: !inPortrait

        ListView {
            id: listView_right
            anchors.fill: parent

            headerPositioning: ListView.OverlayHeader
            header: Pane {
                id: header_right
                z: 2
                width: parent.width

                contentHeight: logo.height

                Image {
                    id: logo_right
                    width: parent.width
                    source: "images/qt-logo.png"
                    fillMode: implicitWidth > width ? Image.PreserveAspectFit : Image.Pad
                }

                MenuSeparator {
                    parent: header_right
                    width: parent.width
                    anchors.verticalCenter: parent.bottom
                    visible: !listView_right.atYBeginning
                }
            }

            footer: ItemDelegate {
                id: footer_right
                text: qsTr("Footer")
                width: parent.width

                MenuSeparator {
                    parent: footer_right
                    width: parent.width
                    anchors.verticalCenter: parent.top
                }
            }

            model: 10

            delegate: ItemDelegate {
                text: qsTr("Title %1").arg(index + 1)
                width: listView_right.width
            }

            ScrollIndicator.vertical: ScrollIndicator { }
        }
    }


    Flickable {
        id: flickable

        anchors.fill: parent
        anchors.topMargin: overlayHeader.height
        anchors.leftMargin: !inPortrait ? drawer.width : undefined

        topMargin: 20
        bottomMargin: 20
        contentHeight: column.height

        Layout.fillWidth: true
        Layout.fillHeight: true

        ColumnLayout {
            id: column
            spacing: 20
            anchors.margins: 20
            anchors.left: parent.left
            anchors.right: parent.right

            Layout.fillWidth: true
            Layout.fillHeight: true


            ScrollView {
                width: parent.width
                height: 300
                Layout.fillWidth: true
                Layout.fillHeight: true

                Component {
                    id: messageDelegate
                    ItemDelegate {
                        width: parent.width
                        height: 40

                        Column {
                            Row {
                                Text { text: '<b>' + name + '</b>'}
                            }
                            Text {
                                text: msg
                            }
                        }
                        onClicked: console.log("clicked:", name)
                    }
                }

                ListView {
                    anchors.fill: parent
                    model: MessageModel {}
                    delegate: messageDelegate
                    //highlight: Rectangle { color: "lightsteelblue"; radius: 5 }
                    focus: true
                }
            }

            TextField {
                id: chat_input
                Layout.fillWidth: true
                Layout.fillHeight: true
                cursorVisible: true
                leftPadding: 4
                topPadding: 2
            }

        }

        ScrollIndicator.vertical: ScrollIndicator { }
    }
}
