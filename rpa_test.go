package rpa

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-vgo/robotgo"
)

func Test_SearchImgAndClick(t *testing.T) {
	title := "pwsh"
	testImgPath := "./test/test.png"

	scaleFactors := [7]float64{
		1.0,
		1.25,
		1.5,
		1.75,
		2.0,
		2.25,
		2.5,
	}

	results := make([]bool, 7)

	for i, scaleFactor := range scaleFactors {
		wX, wY, wW, _ := getBounds(title, scaleFactor)
		findImg := robotgo.CaptureScreen(wX+wW-150, wY+10, 100, 100)
		defer robotgo.FreeBitmap(findImg)
		robotgo.SaveBitmap(findImg, testImgPath)
		ch := SearchImgAndClick(title, testImgPath, ScaleFactor(scaleFactor))

		results[i] = <-ch

		if err := os.Remove(testImgPath); err != nil {
			t.Error(err)
		}
	}

	testResult := false
	for i, r := range results {
		if r {
			fmt.Printf("OK %0.2f\n", scaleFactors[i])
			testResult = true
		} else {
			fmt.Printf("NG %0.2f\n", scaleFactors[i])
		}
	}

	if !testResult {
		t.Error("NG")
	}
}
