package rpa

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"golang.org/x/image/bmp"
)

// ToGrayScale : グレースケールに変換
// 参考 : https://qiita.com/tenntenn/items/0471e5f494df82c3e825
func ToGrayScale(src string) string {
	srcImg, _ := os.Open(src)
	defer srcImg.Close()
	img, _ := png.Decode(srcImg)
	bounds := img.Bounds()
	dest := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.Gray16Model.Convert(img.At(x, y))
			gray, _ := c.(color.Gray16)
			dest.Set(x, y, gray)
		}
	}

	tmpFileName, _ := MakeRandomStr(10)
	tmpFilePath := fmt.Sprintf("./tmp/%s.bmp", tmpFileName)
	destImg, err := os.Create(tmpFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// ビット深度16のpngがrobotgoで読み込めないのでbmpでエンコード
	bmp.Encode(destImg, dest)
	return tmpFilePath
}
