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

var EnemyBombSprite = Resource("res/EnemyBomb.png")

type EnemyBomb struct {
	position f32.Point
	health   int
}

func NewBomb(position f32.Point) *EnemyBomb {
	return &EnemyBomb{
		position: position,
		health:   3,
	}
}

func (b *EnemyBomb) Position() f32.Point {
	return b.position
}

func (b *EnemyBomb) Size() float32 {
	return 5
}

func (b *EnemyBomb) Team() Team {
	return EnemyTeam
}

func (b *EnemyBomb) Draw(ops *op.Ops) {
	defer op.Affine(f32.Affine2D{}.Offset(
		b.position.Sub(f32.Point{X: float32(EnemyBombSprite.Size().X / 2), Y: float32(EnemyBombSprite.Size().Y / 2)}),
	)).Push(ops).Pop()
	EnemyBombSprite.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func (b *EnemyBomb) Logic(g *Level) {
	b.position.Y += .5
	if b.position.Y > float32(screenHeight) {
		delete(g.Entities, b)
	}
	for e := range g.Entities {
		if e.Team() != b.Team() && collides(b, e) {
			if e, ok := e.(Damageable); ok {
				e.Hurt(g, 10)
				delete(g.Entities, b)
			}
		}
	}
}

func (b *EnemyBomb) Hurt(g *Level, damage int) {
	b.health -= damage
	if b.health <= 0 {
		delete(g.Entities, b)
		g.Score += BombScore
		// Spawn bullets as fragments
		for i := 0; i < 20; i++ {
			g.Entities[&Bullet{position: b.position, team: NeitherTeam, orientation: float32(i) * (2 * math.Pi / 10), speed: .5}] = struct{}{}
		}
	}
}
