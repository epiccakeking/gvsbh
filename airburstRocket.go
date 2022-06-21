/*
This file is part of gvsbh.

gvsbh is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

gvsbh is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with gvsbh. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"
)

var AirburstRocketSprite = Resource("res/AirburstRocket.png")

type AirburstRocket struct {
	position    f32.Point
	team        Team
	orientation float32
	speed       float32
}

func (b *AirburstRocket) Position() f32.Point {
	return b.position
}
func (b *AirburstRocket) Size() float32 {
	return 20 // Fake oversized radius for airbursting
}
func (b *AirburstRocket) Team() Team {
	return b.team
}
func (b *AirburstRocket) Draw(ops *op.Ops) {
	defer op.Affine(f32.Affine2D{}.Rotate(f32.Point{}, b.orientation).Offset(
		b.position.Sub(f32.Point{X: float32(AirburstRocketSprite.Size().X / 2), Y: float32(AirburstRocketSprite.Size().Y / 2)}),
	)).Push(ops).Pop()
	AirburstRocketSprite.Add(ops)
	paint.PaintOp{}.Add(ops)
}
func (b *AirburstRocket) Logic(g *Level) {
	b.position.Y -= float32(math.Cos(float64(b.orientation))) * b.speed
	b.position.X += float32(math.Sin(float64(b.orientation))) * b.speed
	if OOB(b) {
		delete(g.Entities, b)
	}
	for e := range g.Entities {
		if e.Team() != b.team && collides(b, e) {
			// Airburst if near a damagable entity
			if _, ok := e.(Damageable); ok {
				for i := -4; i < 5; i++ {
					g.Entities[&Bullet{
						position:    b.position,
						team:        b.team,
						orientation: b.orientation + float32(math.Pi/20)*float32(i),
						speed:       b.speed,
					}] = struct{}{}
				}
				delete(g.Entities, b)
			}
		}
	}
}
