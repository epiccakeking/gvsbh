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

const PulsarSlideTicks = Tickrate / 2

var PulsarSprite = Resource("res/Pulsar.png")

type Pulsar struct {
	targetPosition f32.Point
	health         int
	slideTicks     int
	returnTimer    int
	spin           float32
	shotCooldown   int
}

func NewPulsar(targetPosition f32.Point) *Pulsar {
	return &Pulsar{
		targetPosition: targetPosition,
		health:         10,
		returnTimer:    Tickrate * 3,
	}
}

func (p *Pulsar) Position() f32.Point {
	return f32.Point{X: p.targetPosition.X, Y: p.targetPosition.Y * float32(p.slideTicks) / float32(PulsarSlideTicks)}
}

func (p *Pulsar) Size() float32 {
	return 5
}

func (p *Pulsar) Team() Team {
	return EnemyTeam
}

func (p *Pulsar) Draw(ops *op.Ops) {
	defer op.Affine(f32.Affine2D{}.Rotate(f32.Point{X: float32(PulsarSprite.Size().X) / 2, Y: float32(PulsarSprite.Size().Y) / 2}, p.spin).Offset(
		p.Position().Sub(f32.Point{X: float32(PulsarSprite.Size().X) / 2, Y: float32(PulsarSprite.Size().Y) / 2}),
	)).Push(ops).Pop()
	PulsarSprite.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func (p *Pulsar) Logic(g *Level) {
	p.spin += float32(2 * math.Pi / float64(Tickrate))
	if p.returnTimer > 0 {
		if p.slideTicks < PulsarSlideTicks {
			p.slideTicks++
		} else {
			p.returnTimer--
			if p.shotCooldown > 0 {
				p.shotCooldown--
			} else {
				g.Entities[&Bullet{position: p.targetPosition, team: p.Team(), orientation: p.spin, speed: .25}] = struct{}{}
				g.Entities[&Bullet{position: p.targetPosition, team: p.Team(), orientation: p.spin + math.Pi/2, speed: .25}] = struct{}{}
				g.Entities[&Bullet{position: p.targetPosition, team: p.Team(), orientation: p.spin + math.Pi, speed: .25}] = struct{}{}
				g.Entities[&Bullet{position: p.targetPosition, team: p.Team(), orientation: p.spin + math.Pi*3/2, speed: .25}] = struct{}{}
				p.shotCooldown = Tickrate / 5
			}
		}
	} else {
		p.slideTicks--
		if p.slideTicks <= 0 {
			delete(g.Entities, p)
		}
	}
}

func (p *Pulsar) Hurt(g *Level, damage int) {
	p.health -= damage
	if p.health <= 0 {
		delete(g.Entities, p)
		g.Score += PulsarScore
	}
}
