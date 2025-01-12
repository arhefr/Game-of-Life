package life

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Отрисовывает игровое состояние мира
func (w *World) Print(background *ebiten.Image, backgroundColor color.RGBA) {
	R, G, B := backgroundColor.R, backgroundColor.G, backgroundColor.B

	for y, row := range w.Cells {
		for x, cell := range row {
			mutation := (x%2 == 0 && y%2 != 0) || (x%2 != 0 && y%2 == 0)
			if cell {
				if mutation {
					ebitenutil.DrawRect(background, float64(x), float64(y), 1, 1, color.RGBA{227, 233, 229, 0xff})
				} else {
					ebitenutil.DrawRect(background, float64(x), float64(y), 1, 1, color.RGBA{255, 255, 255, 0xff})
				}
			} else if mutation {
				ebitenutil.DrawRect(background, float64(x), float64(y), 1, 1, color.RGBA{R - 5, G - 5, B - 5, 0xff})
			}
		}
	}
}
