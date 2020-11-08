package rpa

import (
	"crypto/rand"
	"image"
	"image/png"
	"os"
)

// MakeRandomStr : ランダム文字列生成
// 参考 : https://qiita.com/RyotaNakaya/items/7d269525a288c4b3ecda
func MakeRandomStr(digit uint32) (string, error) {
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

// DeleteElementFromStringSlice : スライスから指定インデックスの要素を削除する
// 参考 : https://gawawa124.hatenablog.com/entry/2015/04/09/190000
func DeleteElementFromStringSlice(s []string, i int) []string {
	s = append(s[:i], s[i+1:]...)
	//新しいスライスを用意することがポイント
	n := make([]string, len(s))
	copy(n, s)
	return n
}

// FindIndexFromStringSlice : ストリングスライスから検索ワードで最初に見つかった場所のインデックスを返す
func FindIndexFromStringSlice(s []string, searchTerm string) (int, bool) {
	for i, str := range s {
		if str == searchTerm {
			return i, true
		}
	}
	return -1, false
}

// LoadImage : 画像ファイル読み込み(pngのみ)
func LoadImage(path string) image.Image {
	imgFile, _ := os.Open(path)
	defer imgFile.Close()
	img, _ := png.Decode(imgFile)

	return img
}
