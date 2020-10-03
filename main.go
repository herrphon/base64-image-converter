package main

import (
	"log"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)


func main() {
	mw := new(MyMainWindow)
	var openAction *walk.Action

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Image Base64 Converter",
		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						AssignTo:    &openAction,
						Text:        "Exit",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&Help",

				Items: []MenuItem{
					Action{
						Text:        "About",
						OnTriggered: mw.aboutAction_Triggered,
					},
				},
			},
		},
		Size:     Size{500, 800},
		Layout:   VBox{},
		OnDropFiles: func(files []string) {
			mw.textArea.SetText(strings.Join(files, "\r\n"))
		},
		Children: []Widget{
			Composite{
				AssignTo: &mw.hideableComposite,
				Layout:   HBox{SpacingZero: true, MarginsZero: true},
				Children: []Widget{
					Label{
						Text: "Enter image location below",
					},
				},
			},

			VSplitter{
				Children: []Widget{
					Composite{
						Layout: HBox{SpacingZero: true, MarginsZero: true},
						Children: []Widget{
							PushButton{
								Text: "...",
								OnClicked: func() {
									dlgFile := new(walk.FileDialog)

									if ok, err := dlgFile.ShowOpen(mw); err != nil {
										println(dlgFile.FilePath, err)
									} else if !ok {
										println(dlgFile.FilePath, err)
									}

									mw.inputFileLocationInput.SetText(dlgFile.FilePath)
								},
							},
							LineEdit{
								AssignTo:    &mw.inputFileLocationInput,
								ToolTipText: "e.g. img/exit.png",
								OnKeyDown: func(key walk.Key) {
									if key == walk.KeyReturn {
										mw.convertAction()
									}
								},
							},
							HSpacer{Size: 3},
							PushButton{
								Text:      "Convert",
								OnClicked: mw.convertAction,
							},
						},
					},
				},
			},

			VSpacer{Size: 3},

			TextEdit{
				AssignTo: &mw.textArea,
				OnMouseDown: func(x, y int, button walk.MouseButton) {
					// fmt.Printf("position: %d - %d \n", x, y)
				},
				ToolTipText: "Click to copy to clipboard.",
			},

			PushButton{
				Text: "Base64 to Image",
				OnClicked: func() {
					drawImage(mw.textArea, mw.outputImageView)
				},
			},

			VSpacer{Size: 3},

			ImageView{
				AssignTo:   &mw.outputImageView,
				Background: SolidColorBrush{Color: walk.RGB(222, 222, 222)},
				Margin:     40,
				Mode:       ImageViewModeIdeal,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

type MyMainWindow struct {
	*walk.MainWindow
	hideableComposite      *walk.Composite
	inputFileLocationInput *walk.LineEdit
	textArea               *walk.TextEdit
	outputImageView        *walk.ImageView
	visible                bool
}

func (mw *MyMainWindow) convertAction() {
	mw.toggleVisibility(mw.hideableComposite)

	fileLocation := mw.inputFileLocationInput.Text()
	base64, err := getBase64FromImage(fileLocation)
	if err != nil {
		mw.textArea.SetText(err.Error())
	} else {
		mw.textArea.SetText(base64)
		drawImage(mw.textArea, mw.outputImageView)
	}
}

func (mw *MyMainWindow) toggleVisibility(hideableComposite *walk.Composite) {
	children := hideableComposite.Children()
	if mw.visible {
		hideableComposite.SetBounds(walk.Rectangle{1, 1, 2, 2})
		// children.Clear()
		mw.visible = false
	} else {
		// Add something and remove it right away to make the hideableComposite reappear.
		label, _ := walk.NewLabel(hideableComposite)
		label.SetText("...")
		children.Add(label)
		children.RemoveAt(1)
		mw.visible = true
	}
}

func drawImage(textArea *walk.TextEdit, outputImageView *walk.ImageView) {
	outputText := textArea.Text()

	if img, err := getImageFromBase64(outputText); err != nil {
		outputImageView.SetImage(nil)
		println(err)
	} else {
		if ic, err := walk.NewIconFromImageForDPI(img, 96); err != nil {
			outputImageView.SetImage(nil)
			println(err)
		} else {
			outputImageView.SetImage(ic)
		}
	}
}

func (mw *MyMainWindow) aboutAction_Triggered() {
	walk.MsgBox(mw, "About", "Version 1.0.0", walk.MsgBoxIconInformation)
}
