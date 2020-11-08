package rpa

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/TakahiroSato/imgconv"
	"github.com/go-vgo/robotgo"
)

func getWindowTitles() ([]string, error) {
	titles := []string{}
	testWindowTitles, err := os.Open("./test/test_window_titles.txt")
	defer testWindowTitles.Close()
	if err != nil {
		return titles, err
	}

	scanner := bufio.NewScanner(testWindowTitles)
	for scanner.Scan() {
		titles = append(titles, scanner.Text())
	}
	return titles, nil
}

func Test_SearchImg(t *testing.T) {
	testImgPath := "./tmp/test.png"

	scaleFactors := []float64{
		1.0,
		1.25,
		1.5,
	}

	testFunc := func(title string) {
		results := make([]SearchedData, len(scaleFactors))

		for i, scaleFactor := range scaleFactors {
			func() {
				wX, wY, wW, wH := getBounds(title, scaleFactor)
				screenshot := robotgo.CaptureScreen(wX, wY, wW, wH)
				defer robotgo.FreeBitmap(screenshot)
				robotgo.SaveBitmap(screenshot, fmt.Sprintf("./screenshot/test_%d.png", i))
				findImg := robotgo.CaptureScreen(wX+wW-110, wY+wH-110, 100, 100)
				defer robotgo.FreeBitmap(findImg)
				robotgo.SaveBitmap(findImg, testImgPath)
				ch := SearchImg(title, testImgPath, OptScaleFactor(scaleFactor))

				results[i] = <-ch

				if err := os.Remove(testImgPath); err != nil {
					t.Error(err)
				}
			}()
		}

		testResult := false
		for i, r := range results {
			if r.Ok {
				fmt.Printf("OK %0.2f [x=%d, y=%d]\n", scaleFactors[i], r.X, r.Y)
				testResult = true
			} else {
				fmt.Printf("NG %0.2f\n", scaleFactors[i])
			}
		}

		if !testResult {
			t.Error("NG")
		}
	}

	titles, err := getWindowTitles()
	if err != nil {
		t.Error(err)
	}
	for _, title := range titles {
		testFunc(title)
	}
}

func Test_DragAndDrop(t *testing.T) {
	SearchedData{
		Ok: true,
		X:  500,
		Y:  500,
	}.DragAndDrop(500, -50, -100, 100)
}

func Test_ToGrayScale(t *testing.T) {
	imgconv.ToGrayScale(LoadImage("./screenshot/test_2.png")).SaveAsBmp("./test/gray.bmp")
}

func Test_innerSearchImg_GrayScale(t *testing.T) {
	titles, err := getWindowTitles()
	if err != nil {
		t.Error(err)
	}
	scaleFactor := 1.5
	for _, title := range titles {
		func() {
			wX, wY, wW, wH := getBounds(title, scaleFactor)
			findImg := robotgo.CaptureScreen(wX+500, wY+500, 250, 250)
			defer robotgo.FreeBitmap(findImg)
			findImgPath := "./tmp/test.png"
			robotgo.SaveBitmap(findImg, findImgPath)
			defer os.Remove(findImgPath)
			x, y := innerSearchImg(wX, wY, wW, wH, findImgPath, 0.01, true)
			x = multiply(x, 1/scaleFactor)
			y = multiply(y, 1/scaleFactor)
			_wX := multiply(wX, 1/scaleFactor)
			_wY := multiply(wY, 1/scaleFactor)
			robotgo.MoveMouse(_wX+x, _wY+y)
		}()
	}
}
