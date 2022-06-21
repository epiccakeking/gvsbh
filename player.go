/*
This file is part of gvsbh.

gvsbh is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

gvsbh is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with gvsbh. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"
)

var PlayerSprite = Resource("res/Player.png")

type Player struct {
	position     f32.Point
	health       int
	shotCooldown int
}

func NewPlayer(position f32.Point) *Player {
	return &Player{
		position: position,
		health:   100,
	}
}

func (p *Player) Position() f32.Point {
	return p.position
}

func (p *Player) Size() float32 {
	return 6
}

func (p *Player) Team() Team {
	return PlayerTeam
}

func (p *Player) Draw(ops *op.Ops) {
	defer op.Affine(f32.Affine2D{}.Offset(
		p.position.Sub(f32.Point{X: float32(PlayerSprite.Size().X / 2), Y: float32(PlayerSprite.Size().Y / 2)}),
	)).Push(ops).Pop()
	PlayerSprite.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func (p *Player) Logic(g *Level) {
	p.position.X += g.MovementVector.X
	if p.position.X < 0 {
		p.position.X = 0
	} else if p.position.X > float32(screenWidth) {
		p.position.X = float32(screenWidth)
	}
	if p.position.Y < 0 {
		p.position.Y = 0
	} else if p.position.Y > float32(screenHeight) {
		p.position.Y = float32(screenHeight)
	}

	p.position.Y += g.MovementVector.Y
	if p.shotCooldown > 0 {
		p.shotCooldown--
	}
	if g.Shooting && p.shotCooldown == 0 {
		g.Entities[&AirburstRocket{
			position: f32.Point{X: p.position.X, Y: p.position.Y - 5},
			team:     p.Team(),
			speed:    1,
		}] = struct{}{}
		p.shotCooldown = 50
	}
}

func (p *Player) Hurt(g *Level, damage int) {
	p.health -= damage
	if p.health <= 0 {
		delete(g.Entities, p)
	}
}
