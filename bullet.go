/*
Copyright 2022 epiccakeking

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

var BulletSprite = Resource("res/Bullet.png")

type Bullet struct {
	x, y        float64
	team        Team
	orientation float64
	speed       float64
}

func (b *Bullet) Position() (x, y float64) {
	return b.x, b.y
}
func (b *Bullet) Size() float64 {
	return 1
}
func (b *Bullet) Team() Team {
	return b.team
}
func (b *Bullet) Draw(screen *ebiten.Image) {
	BulletSprite.Draw(screen, b.x, b.y, b.orientation)
}
func (b *Bullet) Logic(g *Level) {
	const damage = 1
	b.y -= math.Cos(b.orientation) * b.speed
	b.x += math.Sin(b.orientation) * b.speed
	if OOB(b) {
		g.RemoveEntity(b)
	}
	for e := range g.Entities {
		if e.Team() != b.team && collides(b, e) {
			if e, ok := e.(Damageable); ok {
				e.Hurt(g, damage)
				g.RemoveEntity(b)
			}
		}
	}
}
