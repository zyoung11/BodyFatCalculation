package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type tallEntry struct {
	widget.Entry
}

func newTallEntry() *tallEntry {
	e := &tallEntry{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *tallEntry) MinSize() fyne.Size {
	s := e.Entry.MinSize()
	return fyne.NewSize(s.Width, 48)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("体脂率计算器")
	myWindow.Resize(fyne.NewSize(400, 680))

	tabs := container.NewAppTabs(
		container.NewTabItem("  男士  ", buildGenderTab(myWindow, true)),
		container.NewTabItem("  女士  ", buildGenderTab(myWindow, false)),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

func buildGenderTab(win fyne.Window, isMale bool) fyne.CanvasObject {
	waistPH := "例: 85"
	neckPH := "例: 38"
	heightPH := "例: 175"
	if !isMale {
		waistPH = "例: 70"
		neckPH = "例: 33"
		heightPH = "例: 162"
	}

	waistEntry := newTallEntry()
	waistEntry.SetPlaceHolder(waistPH)
	neckEntry := newTallEntry()
	neckEntry.SetPlaceHolder(neckPH)
	heightEntry := newTallEntry()
	heightEntry.SetPlaceHolder(heightPH)

	bodyFatText := canvas.NewText("", color.White)
	bodyFatText.TextSize = 64
	bodyFatText.TextStyle = fyne.TextStyle{Bold: true}
	bodyFatText.Alignment = fyne.TextAlignCenter

	levelLabel := widget.NewLabel("")
	levelLabel.Alignment = fyne.TextAlignCenter
	levelLabel.TextStyle = fyne.TextStyle{Bold: true}

	outer := container.NewVBox()

	formView := container.NewVBox(
		container.NewPadded(widget.NewLabelWithStyle("腰围 (cm)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		container.NewPadded(waistEntry),
		container.NewPadded(widget.NewLabelWithStyle("颈围 (cm)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		container.NewPadded(neckEntry),
		container.NewPadded(widget.NewLabelWithStyle("身高 (cm)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		container.NewPadded(heightEntry),
		layout.NewSpacer(),
	)

	var formViewWithBtn *fyne.Container

	clearBtn := widget.NewButton("清除", func() {
		waistEntry.SetText("")
		neckEntry.SetText("")
		heightEntry.SetText("")
		bodyFatText.Text = ""
		bodyFatText.Refresh()
		levelLabel.SetText("")
		outer.Objects = []fyne.CanvasObject{formViewWithBtn}
		outer.Refresh()
	})

	resultView := container.NewVBox(
		layout.NewSpacer(),
		container.NewPadded(container.NewHBox(layout.NewSpacer(), bodyFatText, layout.NewSpacer())),
		container.NewPadded(levelLabel),
		layout.NewSpacer(),
		container.NewPadded(clearBtn),
	)

	calcBtn := widget.NewButton("计算体脂率", func() {
		waist, err1 := strconv.ParseFloat(waistEntry.Text, 64)
		neck, err2 := strconv.ParseFloat(neckEntry.Text, 64)
		height, err3 := strconv.ParseFloat(heightEntry.Text, 64)

		if err1 != nil || err2 != nil || err3 != nil {
			dialog.ShowError(fmt.Errorf("请输入有效的数字"), win)
			return
		}
		if waist <= neck {
			dialog.ShowError(fmt.Errorf("腰围必须大于颈围"), win)
			return
		}
		if height <= 0 {
			dialog.ShowError(fmt.Errorf("身高必须为正数"), win)
			return
		}

		var bodyFat float64
		if isMale {
			bodyFat = 86.010*math.Log10(waist-neck) - 70.041*math.Log10(height) + 36.76
		} else {
			hipEst := waist * 1.15
			bodyFat = 163.205*math.Log10(waist+hipEst-neck) - 97.684*math.Log10(height) - 78.387
		}

		lv, lc := getBodyFatLevel(bodyFat, isMale)
		bodyFatText.Text = fmt.Sprintf("%.1f%%", bodyFat)
		bodyFatText.Color = lc
		bodyFatText.Refresh()
		levelLabel.SetText(lv)

		outer.Objects = []fyne.CanvasObject{resultView}
		outer.Refresh()
	})
	calcBtn.Importance = widget.HighImportance

	formViewWithBtn = container.NewVBox(
		formView,
		container.NewPadded(calcBtn),
	)

	outer.Objects = []fyne.CanvasObject{formViewWithBtn}

	return container.NewVScroll(outer)
}

func getBodyFatLevel(bodyFat float64, isMale bool) (string, color.Color) {
	if isMale {
		switch {
		case bodyFat < 6:
			return "体脂过低", color.NRGBA{R: 220, G: 60, B: 60, A: 255}
		case bodyFat < 14:
			return "运动员水平", color.NRGBA{R: 50, G: 190, B: 50, A: 255}
		case bodyFat < 18:
			return "标准体型", color.NRGBA{R: 50, G: 170, B: 50, A: 255}
		case bodyFat < 25:
			return "偏胖", color.NRGBA{R: 230, G: 160, B: 30, A: 255}
		default:
			return "肥胖", color.NRGBA{R: 220, G: 60, B: 60, A: 255}
		}
	}
	switch {
	case bodyFat < 14:
		return "体脂过低", color.NRGBA{R: 220, G: 60, B: 60, A: 255}
	case bodyFat < 21:
		return "运动员水平", color.NRGBA{R: 50, G: 190, B: 50, A: 255}
	case bodyFat < 25:
		return "标准体型", color.NRGBA{R: 50, G: 170, B: 50, A: 255}
	case bodyFat < 32:
		return "偏胖", color.NRGBA{R: 230, G: 160, B: 30, A: 255}
	default:
		return "肥胖", color.NRGBA{R: 220, G: 60, B: 60, A: 255}
	}
}
