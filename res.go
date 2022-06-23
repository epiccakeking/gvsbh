/*
This file is part of gvsbh.

gvsbh is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

gvsbh is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with gvsbh. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed res
var res embed.FS

type Sprite struct{ *ebiten.Image }

func Resource(path string) Sprite {
	f, err := res.Open(path)
	if err != nil {
		panic(err)
	}
	i, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return Sprite{ebiten.NewImageFromImage(i)}
}

// Draw
func (s Sprite) Draw(screen *ebiten.Image, x, y float64, orientation float64) {
	op := &ebiten.DrawImageOptions{}
	sizeX, sizeY := s.Size()
	op.GeoM.Translate(float64(-sizeX/2), float64(-sizeY/2))
	op.GeoM.Rotate(orientation)
	op.GeoM.Translate(x, y)
	op.GeoM.Scale(scale, scale)
	screen.DrawImage(s.Image, op)
}
