package rpa

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-vgo/robotgo"
)

func Test_SearchImg(t *testing.T) {
	title := "pwsh"
	testImgPath := "./test/test.png"

	scaleFactors := []float64{
		1.0,
		1.25,
		1.5,
	}

	results := make([]SearchedData, len(scaleFactors))

	for i, scaleFactor := range scaleFactors {
		wX, wY, wW, _ := getBounds(title, scaleFactor)
		findImg := robotgo.CaptureScreen(wX+wW-150, wY+10, 100, 100)
		defer robotgo.FreeBitmap(findImg)
		robotgo.SaveBitmap(findImg, testImgPath)
		ch := SearchImg(title, testImgPath, ScaleFactor(scaleFactor))

		results[i] = <-ch

		if err := os.Remove(testImgPath); err != nil {
			t.Error(err)
		}
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
