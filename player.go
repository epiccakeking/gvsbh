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

var PlayerSprite = Resource("res/Player.png")

type Player struct {
	x, y         float64
	health       int
	shotCooldown int
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		x: x, y: y,
		health: 100,
	}
}

func (p *Player) Position() (x, y float64) {
	return p.x, p.y
}

func (p *Player) Size() float64 {
	return 6
}

func (p *Player) Team() Team {
	return PlayerTeam
}

func (p *Player) Draw(screen *ebiten.Image) {
	PlayerSprite.Draw(screen, float64(p.x), float64(p.y), 0)
}

func (p *Player) Logic(g *Level) {
	if g.UseTouch {
		if g.TouchX < p.x {
			g.MovementX = -Speed
		} else if g.TouchX > p.x {
			g.MovementX = Speed
		}
		if g.TouchY < p.y {
			g.MovementY = -Speed
		} else if g.TouchY > p.y {
			g.MovementY = Speed
		}
	}
	p.x += g.MovementX
	if p.x < 0 {
		p.x = 0
	} else if p.x > screenWidth {
		p.x = screenWidth
	}

	p.y += g.MovementY
	if p.y < 0 {
		p.y = 0
	} else if p.y > screenHeight {
		p.y = screenHeight
	}

	if p.shotCooldown > 0 {
		p.shotCooldown--
	}
	if g.Shooting && p.shotCooldown == 0 {
		g.Entities[&Bullet{
			x: p.x + 5, y: p.y - 5,
			team:        p.Team(),
			orientation: math.Pi * .1,
			speed:       1,
		}] = struct{}{}
		g.Entities[&Bullet{
			x: p.x - 5, y: p.y - 5,
			team:        p.Team(),
			orientation: math.Pi * -.1,
			speed:       1,
		}] = struct{}{}
		p.shotCooldown = 20
	}
}

func (p *Player) Hurt(g *Level, damage int) {
	p.health -= damage
	if p.health <= 0 {
		delete(g.Entities, p)
	}
}
