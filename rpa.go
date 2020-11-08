package rpa

import (
	"fmt"
	"os"
	"time"

	"github.com/TakahiroSato/imgconv"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
)

// SearchImgOpts : SearchImg関数のオプション引数
type SearchImgOpts struct {
	Tolerance   float64
	IsSaveImg   bool
	ScaleFactor float64
	IsGrayScale bool
}

type option func(*SearchImgOpts)

// OptTolerance : 一致度の許容値（デフォルトは0.01 違いを1%まで許容する....多分)
func OptTolerance(v float64) option {
	return func(o *SearchImgOpts) {
		o.Tolerance = v
	}
}

// OptIsSaveImg : スクリーンショットを保存するかどうか（デフォルトはしないfalse）
func OptIsSaveImg(v bool) option {
	return func(o *SearchImgOpts) {
		o.IsSaveImg = v
	}
}

// OptScaleFactor : 画面の拡大率（デフォルトは1.5）
func OptScaleFactor(v float64) option {
	return func(o *SearchImgOpts) {
		o.ScaleFactor = v
	}
}

// OptIsGrayScale : グレースケールで検索するかどうか（デフォルトはしないfalse）
func OptIsGrayScale(v bool) option {
	return func(o *SearchImgOpts) {
		o.IsGrayScale = v
	}
}

// SearchImg : 指定タイトルのウィンドウから指定した画像(.png)を探しその情報を返却する
func SearchImg(title string, imgPath string, opts ...option) <-chan SearchedData {
	o := SearchImgOpts{
		Tolerance:   0.01,
		IsSaveImg:   false,
		ScaleFactor: 1.5,
		IsGrayScale: false,
	}

	for _, opt := range opts {
		opt(&o)
	}
	ch := make(chan SearchedData)

	go func() {
		wX, wY, wW, wH := getBounds(title, o.ScaleFactor)

		printPreTime("Finding " + imgPath + " ...")
		var x, y int
		for i := 0; i < 5; i++ {
			x, y = innerSearchImg(wX, wY, wW, wH, imgPath, o.Tolerance, o.IsGrayScale)
			if x > 0 || y > 0 {
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

// private functions

func innerSearchImg(wX, wY, wW, wH int, imgPath string, tolerance float64, isGrayScale bool) (int, int) {
	refRect := robotgo.CaptureScreen(wX, wY, wW, wH)
	defer robotgo.FreeBitmap(refRect)

	_genTmpFilePath := func() string {
		tmpFileName, _ := MakeRandomStr(10)
		tmpFilePath := fmt.Sprintf("./tmp/%s.bmp", tmpFileName)

		return tmpFilePath
	}

	if isGrayScale {
		refImgName, _ := MakeRandomStr(10)
		refImgPath := fmt.Sprintf("./tmp/%s.png", refImgName)
		robotgo.SaveBitmap(refRect, refImgPath)
		grayRefImgPath := _genTmpFilePath()
		imgconv.ToGrayScale(LoadImage(refImgPath)).SaveAsBmp(grayRefImgPath)
		grayFindImgPath := _genTmpFilePath()
		imgconv.ToGrayScale(LoadImage(imgPath)).SaveAsBmp(grayFindImgPath)
		bit := robotgo.OpenBitmap(grayFindImgPath, 2) // 2 = bitmap
		sbit := robotgo.OpenBitmap(grayRefImgPath, 2) // 2 = bitmap
		defer func() {
			robotgo.FreeBitmap(bit)
			robotgo.FreeBitmap(sbit)
			removeFilePaths := []string{refImgPath, grayRefImgPath, grayFindImgPath}
			for {
				if len(removeFilePaths) == 0 {
					break
				}
				var removedPaths []string
				for _, path := range removeFilePaths {
					err := os.Remove(path)
					if err == nil {
						removedPaths = append(removedPaths, path)
					}
				}
				for _, p := range removedPaths {
					i, _ := FindIndexFromStringSlice(removeFilePaths, p)
					removeFilePaths = DeleteElementFromStringSlice(removeFilePaths, i)
				}
			}
		}()
		return robotgo.FindBitmap(bit, sbit, tolerance)
	}

	return robotgo.FindPic(imgPath, refRect, tolerance)
}

func printPreTime(str string) {
	now := time.Now()
	fmt.Println(fmt.Sprintf("[%02d/%02d %02d:%02d:%02d] ", now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()) + str)
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
