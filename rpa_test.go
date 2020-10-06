package rpa

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-vgo/robotgo"
)

func Test_SearchImg(t *testing.T) {
	title := "rpa_test.go - rpa - Visual Studio Code"
	testImgPath := "./test/test.png"

	scaleFactors := []float64{
		1.0,
		1.25,
		1.5,
	}

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

func Test_DragAndDrop(t *testing.T) {
	SearchedData{
		Ok: true,
		X:  500,
		Y:  500,
	}.DragAndDrop(500, -50, -100, 100)
}

func Test_ToGrayScale(t *testing.T) {
	ToGrayScale("./screenshot/test_1.png")
}
