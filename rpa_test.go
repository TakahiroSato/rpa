package rpa

import (
	"os"
	"testing"

	"github.com/go-vgo/robotgo"
)

func Test_SearchImgAndClick(t *testing.T) {
	title := "pwsh"
	testImgPath := "./test/test.png"
	wX, wY, wW, _ := getBounds(title)
	findImg := robotgo.CaptureScreen(wX+wW-150, wY+10, 100, 100)
	defer robotgo.FreeBitmap(findImg)
	robotgo.SaveBitmap(findImg, testImgPath)
	ch := SearchImgAndClick(title, testImgPath)

	result := <-ch
	if !result {
		t.Error("NG")
	}

	if err := os.Remove(testImgPath); err != nil {
		t.Error(err)
	}
}
