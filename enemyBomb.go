/*
This file is part of gvsbh.

gvsbh is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

gvsbh is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with gvsbh. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var EnemyBombSprite = Resource("res/EnemyBomb.png")

type EnemyBomb struct {
	x, y   float64
	health int
}

func NewBomb(x, y float64) *EnemyBomb {
	return &EnemyBomb{
		x: x, y: y,
		health: 3,
	}
}

func (b *EnemyBomb) Position() (x, y float64) {
	return float64(b.x), float64(b.y)
}

func (b *EnemyBomb) Size() float64 {
	return 5
}

func (b *EnemyBomb) Team() Team {
	return EnemyTeam
}

func (b *EnemyBomb) Draw(screen *ebiten.Image) {
	EnemyBombSprite.Draw(screen, float64(b.x), float64(b.y), 0)
}

func (b *EnemyBomb) Logic(g *Level) {
	b.y += .5
	if b.y > screenHeight {
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
			g.Entities[&Bullet{x: b.x, y: b.y, team: NeitherTeam, orientation: float64(i) * (2 * math.Pi / 10), speed: .5}] = struct{}{}
		}
	}
}
