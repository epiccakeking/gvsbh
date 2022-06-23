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

const PulsarSlideTicks = Tickrate / 2

var PulsarSprite = Resource("res/Pulsar.png")

type Pulsar struct {
	targetY, x, y float64
	health        int
	slideTicks    int
	returnTimer   int
	spin          float64
	shotCooldown  int
}

func NewPulsar(x, y float64) *Pulsar {
	return &Pulsar{
		x:           x,
		y:           y,
		targetY:     y,
		health:      10,
		returnTimer: Tickrate * 3,
	}
}

func (p *Pulsar) Position() (x, y float64) {
	return p.x, p.y
}

func (p *Pulsar) Size() float64 {
	return 5
}

func (p *Pulsar) Team() Team {
	return EnemyTeam
}

func (p *Pulsar) Draw(screen *ebiten.Image) {
	PulsarSprite.Draw(screen, p.x, p.y, p.spin)
}

func (p *Pulsar) Logic(g *Level) {
	p.spin += 2 * math.Pi / float64(Tickrate)
	if p.returnTimer > 0 {
		if p.slideTicks < PulsarSlideTicks {
			p.slideTicks++
		} else {
			p.returnTimer--
			if p.shotCooldown > 0 {
				p.shotCooldown--
			} else {
				g.Entities[&Bullet{x: p.x, y: p.y, team: p.Team(), orientation: p.spin, speed: .25}] = struct{}{}
				g.Entities[&Bullet{x: p.x, y: p.y, team: p.Team(), orientation: p.spin + math.Pi/2, speed: .25}] = struct{}{}
				g.Entities[&Bullet{x: p.x, y: p.y, team: p.Team(), orientation: p.spin + math.Pi, speed: .25}] = struct{}{}
				g.Entities[&Bullet{x: p.x, y: p.y, team: p.Team(), orientation: p.spin + math.Pi*3/2, speed: .25}] = struct{}{}
				p.shotCooldown = Tickrate / 5
			}
		}
	} else {
		p.slideTicks--
		if p.slideTicks <= 0 {
			g.RemoveEntity(p)
		}
	}
	p.y = p.targetY * float64(p.slideTicks) / float64(PulsarSlideTicks)
}

func (p *Pulsar) Hurt(g *Level, damage int) {
	p.health -= damage
	if p.health <= 0 {
		g.RemoveEntity(p)
		g.Score += PulsarScore
	}
}
