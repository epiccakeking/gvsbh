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

var BulletSprite = Resource("res/Bullet.png")

type Bullet struct {
	position    f32.Point
	team        Team
	orientation float32
	speed       float32
}

func (b *Bullet) Position() f32.Point {
	return b.position
}
func (b *Bullet) Size() float32 {
	return 1
}
func (b *Bullet) Team() Team {
	return b.team
}
func (b *Bullet) Draw(ops *op.Ops) {
	defer op.Affine(f32.Affine2D{}.Rotate(f32.Point{}, b.orientation).Offset(
		b.position.Sub(f32.Point{X: float32(BulletSprite.Size().X / 2), Y: float32(BulletSprite.Size().Y / 2)}),
	)).Push(ops).Pop()
	BulletSprite.Add(ops)
	paint.PaintOp{}.Add(ops)
}
func (b *Bullet) Logic(g *Level) {
	const damage = 1
	b.position.Y -= float32(math.Cos(float64(b.orientation))) * b.speed
	b.position.X += float32(math.Sin(float64(b.orientation))) * b.speed
	if OOB(b) {
		delete(g.Entities, b)
	}
	for e := range g.Entities {
		if e.Team() != b.team && collides(b, e) {
			if e, ok := e.(Damageable); ok {
				e.Hurt(g, damage)
				delete(g.Entities, b)
			}
		}
	}
}
