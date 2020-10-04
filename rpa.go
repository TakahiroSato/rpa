package rpa

import (
	"fmt"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
)

var scaleFactor float64

func init() {

	// 高DPIディスプレイ(4K)とかで、見づらいからWin10の設定で拡大率を100%から変更してると、robotgoでとってくる諸々の座標がずれちゃう

	// GetScaleFactorForMonitorとGetScaleFactorForDevice試したけど何故かずれる（150がほしいのに140が来る）

	// hwnd := robotgo.FindWindow("dollfro")
	// hMonitor := win.MonitorFromWindow(hwnd, 2)
	//libShcore := windows.NewLazySystemDLL("Shcore.dll")
	// GetScaleFactorForMonitor := libShcore.NewProc("GetScaleFactorForMonitor")
	// var s uint
	// syscall.Syscall(GetScaleFactorForMonitor.Addr(), 2, uintptr(hMonitor), uintptr(unsafe.Pointer(&s)), 0)
	// scaleFactor = float64(s) / 100.0

	// getScaleFactorForDevice := libShcore.NewProc("GetScaleFactorForDevice")
	// s, _, _ := syscall.Syscall(getScaleFactorForDevice.Addr(), 1, 2, 0, 0)

	// 自分が150%に設定してるので、1.5固定で一旦設定
	scaleFactor = 1.5
}

// SearchImgAndClickOpts : SearchImgAndClick関数のオプション引数
type SearchImgAndClickOpts struct {
	Tolerance float64
	IsSaveImg bool
}

type option func(*SearchImgAndClickOpts)

// Tolerance : 一致度の許容値（デフォルトは0.01 違いを1%まで許容する....多分)
func Tolerance(v float64) option {
	return func(o *SearchImgAndClickOpts) {
		o.Tolerance = v
	}
}

// IsSaveImg : スクリーンショットを保存するかどうか（デフォルトはしないfalse）
func IsSaveImg(v bool) option {
	return func(o *SearchImgAndClickOpts) {
		o.IsSaveImg = v
	}
}

// SearchImgAndClick : 指定タイトルのウィンドウから指定した画像(.png)を探しその場所をクリックする
func SearchImgAndClick(title string, imgPath string, opts ...option) <-chan bool {
	o := SearchImgAndClickOpts{
		Tolerance: 0.01,
		IsSaveImg: false,
	}

	for _, opt := range opts {
		opt(&o)
	}
	ch := make(chan bool)

	go func() {
		hwnd := robotgo.FindWindow(title)
		if hwnd == 0 {
			fmt.Fprintf(os.Stderr, "Could not find window")
			os.Exit(1)
		}

		var pid uint32
		win.GetWindowThreadProcessId(hwnd, &pid)

		wX, wY, wW, wH := robotgo.GetBounds(int32(pid))

		multiply := func(n int, mag float64) int {
			return int(float64(n) * mag)
		}
		wX = multiply(wX, scaleFactor)
		wY = multiply(wY, scaleFactor)
		wW = multiply(wW, scaleFactor)
		wH = multiply(wH, scaleFactor)

		printPreTime := func(str string) {
			now := time.Now()
			fmt.Println(fmt.Sprintf("[%02d/%02d %02d:%02d:%02d] ", now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()) + str)
		}

		var x, y int
		// 1分くらい探して見つからなかったら無いと判断して終わる。
		for i := 0; i < 6; i++ {
			printPreTime("Searching " + imgPath + " ...")
			refRect := robotgo.CaptureScreen(wX, wY, wW, wH)
			x, y = robotgo.FindPic(imgPath, refRect, o.Tolerance)
			if x <= 0 || y <= 0 {
				time.Sleep(time.Second * 10)
			} else {
				break
			}
		}

		if x > 0 && y > 0 {
			if o.IsSaveImg {
				go saveImg(title, wX, wY, wW, wH)
			}

			x = multiply(x, 1/scaleFactor)
			y = multiply(y, 1/scaleFactor)
			_wX := multiply(wX, 1/scaleFactor)
			_wY := multiply(wY, 1/scaleFactor)

			robotgo.MoveMouse(_wX+x, _wY+y)
			robotgo.MouseClick("left", true)
			printPreTime(imgPath + " is found!")
			ch <- true
		} else {
			printPreTime(imgPath + " is not found...")
			ch <- false
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
