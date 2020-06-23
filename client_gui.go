package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"go_nat_git/Service"
	"go_nat_git/client"
)

func main() {
	a := app.New()
	var connectData = make([]string, 4)
	themes := theme.LightTheme()
	w := a.NewWindow("NAT GO")
	a.Settings().SetTheme(themes)
	sizeS := fyne.NewSize(800, 600)
	w.Resize(sizeS)
	resource, _ := fyne.LoadResourceFromPath("./nat.png")
	w.SetIcon(resource)
	var vBox = widget.NewVBox()
	addContent(vBox)
	var localServer *widget.Entry = widget.NewEntry()
	var localPort *widget.Entry = widget.NewEntry()
	var remoteServer *widget.Entry = widget.NewEntry()
	var remotePort *widget.Entry = widget.NewEntry()
	var clientConfig = Service.GetClientConfig()
	if len(clientConfig.LocalAddress) > 0 {
		localServer.SetText(clientConfig.LocalAddress)
	}
	if len(clientConfig.LocalPort) > 0 {
		localPort.SetText(clientConfig.LocalPort)
	}
	if len(clientConfig.RemoteAddress) > 0 {
		remoteServer.SetText(clientConfig.RemoteAddress)
	}
	if len(clientConfig.RemotePort) > 0 {
		remotePort.SetText(clientConfig.RemotePort)
	}
	addLocalServer(vBox, &connectData, localServer, localPort)
	addRemoteServer(vBox, &connectData, remoteServer, remotePort)
	addStartButton(vBox, a, &connectData, localServer, localPort, remoteServer, remotePort, w)
	w.SetContent(vBox)
	w.ShowAndRun()
}

func addContent(box *widget.Box) {
	var alignment = fyne.TextAlign(1)
	var style = fyne.TextStyle{}
	style.Bold = true
	style.Italic = true
	style.Monospace = true
	var label = widget.NewLabelWithStyle("Nat Go", alignment, style)
	box.Append(label)
}

func addLocalServer(box *widget.Box, connectData *[]string, localServer *widget.Entry, localPort *widget.Entry) {
	localServer.SetPlaceHolder("Default 127.0.0.1")
	formItem := widget.NewFormItem("LocalAddr:", localServer)
	localPort.SetPlaceHolder("Please set local port")
	formItem2 := widget.NewFormItem("LocalPort:", localPort)
	form := widget.NewForm(formItem, formItem2)
	var group = widget.NewGroup("Local Server", form)

	box.Append(group)
}

func addRemoteServer(box *widget.Box, connectData *[]string, remoteServer *widget.Entry, remotePort *widget.Entry) {
	remoteServer.OnChanged = func(strings string) {
	}
	remoteServer.SetPlaceHolder("Default 127.0.0.1")
	formItem := widget.NewFormItem("RemoteAddr:", remoteServer)
	remotePort.SetPlaceHolder("Please set remote port")
	formItem2 := widget.NewFormItem("RemotePort:", remotePort)
	form := widget.NewForm(formItem, formItem2)
	var group = widget.NewGroup("Remote Server", form)
	box.Append(group)
}

func addStartButton(box *widget.Box, a fyne.App, connectData *[]string, localServer *widget.Entry, localPort *widget.Entry, remoteServer *widget.Entry, remotePort *widget.Entry, w fyne.Window) *widget.Button {
	var button = widget.NewButton("Start", func() {
		var localAddr = localServer.Text
		var localPort = localPort.Text
		var remoteAddr = remoteServer.Text
		var remotePort = remotePort.Text
		var notificationStr = localAddr + ":" + localPort + "\n" + remoteAddr + ":" + remotePort
		notification := fyne.NewNotification("Start Client", notificationStr)
		Service.SetClientConfig(localAddr, localPort, remoteAddr, remotePort)
		a.SendNotification(notification)
		go client.Run()
		w.Hide()
	})
	box.Append(button)
	return button
}
