package main

import (
	// "fmt"
	// "log"
	// "time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	"github.com/pa001024/elang"
	"github.com/pa001024/elang/zhconv"
)

const (
	Mode1 byte = iota
	Mode2
)

type ElangGuiData struct {
	Src  string
	Dst  string
	Mode byte
}
type ElangGui struct {
	mw    *walk.MainWindow
	mSrc  *walk.TextEdit
	mTar  *walk.TextEdit
	mMode *walk.GroupBox

	Data *ElangGuiData
}

func NewElangGui() (this *ElangGui, err error) {
	this = &ElangGui{Data: &ElangGuiData{"我学会了新的姿势", "厘侥氏阻仟议徊米", Mode1}}
	var db *walk.DataBinder
	update := func() {
		src := this.mSrc.Text()
		switch this.Data.Mode {
		case Mode1:
			src = zhconv.EncodeString(src)
			src = elang.EncodeString(src)
			src = zhconv.EncodeString(src)
		case Mode2:
			src = zhconv.DecodeString(src)
			src = elang.DecodeString(src)
			src = zhconv.EncodeString(src)
		}
		this.mTar.SetText(src)
	}
	wmodel := MainWindow{
		AssignTo: &this.mw,
		Title:    "鹅语",
		MinSize:  Size{512, 450},
		Layout:   VBox{},
		DataBinder: DataBinder{
			AssignTo:   &db,
			DataSource: this.Data,
		},
		Children: []Widget{
			TextEdit{
				Text:     Bind("Src"),
				AssignTo: &this.mSrc,
			},
			RadioButtonGroupBox{
				Title:      "模式",
				AssignTo:   &this.mMode,
				Layout:     HBox{},
				DataMember: "Mode",
				Buttons: []RadioButton{
					{Text: "模式一", Value: Mode1, OnClicked: update},
					{Text: "模式二", Value: Mode2, OnClicked: update},
				},
			},
			TextEdit{
				Text:     Bind("Dst"),
				AssignTo: &this.mTar,
			},
		},
	}
	if err = wmodel.Create(); err != nil {
		walk.MsgBox(nil, "错误", "创建窗口失败", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	if icon, err2 := walk.NewIconFromResource("ICON_MAIN"); err2 == nil {
		this.mw.SetIcon(icon)
	}
	this.mSrc.TextChanged().Attach(update)
	db.SetAutoSubmit(true)
	return
}

func main() {
	if w, err := NewElangGui(); err == nil {
		w.mw.Run()
	}
}
