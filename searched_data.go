package rpa

import "github.com/go-vgo/robotgo"

// SearchedData : 見つけた場所の情報
type SearchedData struct {
	Ok bool
	X  int
	Y  int
}

// Move : 見つけた場所にマウスカーソルを移動させる
func (d SearchedData) Move(offsetX int, offsetY int) SearchedData {
	if d.Ok {
		robotgo.MoveMouse(d.X+offsetX, d.Y+offsetY)
	}
	return d
}

func (d SearchedData) click(offsetX int, offsetY int, double bool) SearchedData {
	if d.Ok {
		robotgo.MoveMouse(d.X+offsetX, d.Y+offsetY)
		robotgo.MouseClick("left", double)
	}
	return d
}

// Click : 見つけた場所をクリックする
func (d SearchedData) Click(offsetX int, offsetY int) SearchedData {
	return d.click(offsetX, offsetY, false)
}

// DoubleClick : 見つけた場所をダブルクリックする
func (d SearchedData) DoubleClick(offsetX int, offsetY int) SearchedData {
	return d.click(offsetX, offsetY, true)
}

// DragAndDrop : 見つけた場所から、指定量ドラッグアンドドロップする
func (d SearchedData) DragAndDrop(offsetX int, offsetY int, moveX int, moveY int) SearchedData {
	if d.Ok {
		srcX := d.X + offsetX
		srcY := d.Y + offsetY
		d.Move(offsetX, offsetY)
		robotgo.DragSmooth(srcX+moveX, srcY+moveY)
	}
	return d
}
