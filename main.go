package main

import (
	"log"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	mw                     *walk.MainWindow
	hideableComposite      *walk.Composite
	visible                = true
	inputFileLocationInput *walk.LineEdit
	textArea               *walk.TextEdit
	outputImageView        *walk.ImageView
)

func main() {
	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "Image Base64 Converter",
		Size:     Size{500, 800},
		Layout:   VBox{},
		OnDropFiles: func(files []string) {
			textArea.SetText(strings.Join(files, "\r\n"))
		},
		Children: []Widget{
			Composite{
				AssignTo: &hideableComposite,
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

									inputFileLocationInput.SetText(dlgFile.FilePath)
								},
							},
							LineEdit{
								AssignTo:    &inputFileLocationInput,
								ToolTipText: "e.g. img/exit.png",
								OnKeyDown: func(key walk.Key) {
									if key == walk.KeyReturn {
										convertAction()
									}
								},
							},
							HSpacer{Size: 3},
							PushButton{
								Text:      "Convert",
								OnClicked: convertAction,
							},
						},
					},
				},
			},

			VSpacer{Size: 3},

			TextEdit{AssignTo: &textArea},

			PushButton{
				Text: "Base64 to Image",
				OnClicked: func() {
					drawImage(textArea, outputImageView)
				},
			},
			VSpacer{Size: 3},

			ImageView{
				AssignTo:   &outputImageView,
				Background: SolidColorBrush{Color: walk.RGB(222, 222, 222)},
				Margin:     40,
				Mode:       ImageViewModeIdeal,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

func convertAction() {
	toggleVisibility(hideableComposite)

	fileLocation := inputFileLocationInput.Text()
	base64, err := getBase64FromImage(fileLocation)
	if err != nil {
		textArea.SetText(err.Error())
	} else {
		textArea.SetText(base64)
		drawImage(textArea, outputImageView)
	}
}

func toggleVisibility(hideableComposite *walk.Composite) {
	children := hideableComposite.Children()
	if visible {
		hideableComposite.SetBounds(walk.Rectangle{1, 1, 2, 2})
		// children.Clear()
		visible = false
	} else {
		// Add something and remove it right away to make the hideableComposite reappear.
		label, _ := walk.NewLabel(hideableComposite)
		label.SetText("...")
		children.Add(label)
		children.RemoveAt(1)
		visible = true
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
