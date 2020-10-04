package rpa

import (
	"fmt"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
)

// SearchedData : 見つけた場所の情報
type SearchedData struct {
	Ok bool
	X  int
	Y  int
}

// Move : 見つめた場所にマウスカーソルを移動させる
func (d SearchedData) Move() {
	if d.Ok {
		robotgo.MoveMouse(d.X, d.Y)
	}
}

// Click : 見つけた場所をクリックする
func (d SearchedData) Click() {
	if d.Ok {
		robotgo.MoveMouse(d.X, d.Y)
		robotgo.MouseClick("left", true)
	}
}

// SearchImgOpts : SearchImg関数のオプション引数
type SearchImgOpts struct {
	Tolerance   float64
	IsSaveImg   bool
	ScaleFactor float64
}

type option func(*SearchImgOpts)

// Tolerance : 一致度の許容値（デフォルトは0.01 違いを1%まで許容する....多分)
func Tolerance(v float64) option {
	return func(o *SearchImgOpts) {
		o.Tolerance = v
	}
}

// IsSaveImg : スクリーンショットを保存するかどうか（デフォルトはしないfalse）
func IsSaveImg(v bool) option {
	return func(o *SearchImgOpts) {
		o.IsSaveImg = v
	}
}

// ScaleFactor : 画面の拡大率（デフォルトは1.5）
func ScaleFactor(v float64) option {
	return func(o *SearchImgOpts) {
		o.ScaleFactor = v
	}
}

// SearchImg : 指定タイトルのウィンドウから指定した画像(.png)を探しその情報を返却する
func SearchImg(title string, imgPath string, opts ...option) <-chan SearchedData {
	o := SearchImgOpts{
		Tolerance:   0.01,
		IsSaveImg:   false,
		ScaleFactor: 1.5,
	}

	for _, opt := range opts {
		opt(&o)
	}
	ch := make(chan SearchedData)

	go func() {
		wX, wY, wW, wH := getBounds(title, o.ScaleFactor)

		printPreTime := func(str string) {
			now := time.Now()
			fmt.Println(fmt.Sprintf("[%02d/%02d %02d:%02d:%02d] ", now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()) + str)
		}

		var x, y int
		// 5回探して見つからなかったら無いことにする
		printPreTime("Finding " + imgPath + " ...")
		for i := 0; i < 5; i++ {
			result := func() bool {
				refRect := robotgo.CaptureScreen(wX, wY, wW, wH)
				defer robotgo.FreeBitmap(refRect)
				x, y = robotgo.FindPic(imgPath, refRect, o.Tolerance)
				if x <= 0 || y <= 0 {
					time.Sleep(time.Second * 1)
					return false
				}
				return true
			}()
			if result {
				break
			}
		}

		if x > 0 && y > 0 {
			printPreTime(imgPath + " is found!")
			if o.IsSaveImg {
				go saveImg(title, wX, wY, wW, wH)
			}

			x = multiply(x, 1/o.ScaleFactor)
			y = multiply(y, 1/o.ScaleFactor)
			_wX := multiply(wX, 1/o.ScaleFactor)
			_wY := multiply(wY, 1/o.ScaleFactor)

			ch <- SearchedData{
				Ok: true,
				X:  x + _wX,
				Y:  y + _wY,
			}
		} else {
			printPreTime(imgPath + " is not found...")
			ch <- SearchedData{
				Ok: false,
				X:  -1,
				Y:  -1,
			}
		}

		close(ch)
	}()

	return ch
}

// TODO: 保存先を指定出来るようにする
func saveImg(title string, wX, wY, wW, wH int) {
	bitmap := robotgo.CaptureScreen(wX, wY, wW, wH)
	// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
	defer robotgo.FreeBitmap(bitmap)
	now := time.Now()
	fileName := fmt.Sprintf("%d_%02d_%02d_%02d_%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	fileName = title + "_" + fileName
	robotgo.SaveBitmap(bitmap, "./screenshot/"+fileName+".png")
}

func multiply(n int, mag float64) int {
	return int(float64(n) * mag)
}

func getBounds(title string, scaleFactor float64) (int, int, int, int) {
	hwnd := robotgo.FindWindow(title)
	if hwnd == 0 {
		fmt.Fprintf(os.Stderr, "Could not find window")
		os.Exit(1)
	}

	var pid uint32
	win.GetWindowThreadProcessId(hwnd, &pid)

	wX, wY, wW, wH := robotgo.GetBounds(int32(pid))
	wX = multiply(wX, scaleFactor)
	wY = multiply(wY, scaleFactor)
	wW = multiply(wW, scaleFactor)
	wH = multiply(wH, scaleFactor)
	return wX, wY, wW, wH
}
