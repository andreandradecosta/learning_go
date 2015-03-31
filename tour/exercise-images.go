package main

import (
    "golang.org/x/tour/pic"
    "image/color"
    "image"
)

type Image struct{
    H, W int
}

func (i *Image) ColorModel() color.Model {
    return color.RGBAModel
}

func (i *Image) Bounds() image.Rectangle {
    return image.Rect(0, 0, i.H, i.W)
}

func (i *Image) At(x, y int) color.Color {
    c := uint8(x ^ y)
    return color.RGBA{c, c, 255, 255}
}

func main() {
    m := &Image{256, 256}
    pic.ShowImage(m)
}
