package rpa

import (
	"crypto/rand"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// ToGrayScale : グレースケールに変換
// 参考 : https://qiita.com/tenntenn/items/0471e5f494df82c3e825
func ToGrayScale(src string) {
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

	tmpFileName, _ := makeRandomStr(10)
	destImg, err := os.Create(fmt.Sprintf("./tmp/%s.png", tmpFileName))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	png.Encode(destImg, dest)
}

// 参考 : https://qiita.com/RyotaNakaya/items/7d269525a288c4b3ecda
func makeRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}
