/*
This file is part of gvsbh.

gvsbh is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

gvsbh is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with gvsbh. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

const Speed = .5     // Speed is the speed of the ship (affected by Tickrate)
const Tickrate = 240 // Number of game ticks per second

// Dimensions of the screen in sprite pixels
const (
	screenHeight = 200
	screenWidth  = 100
)

func main() {
	game := Level{
		CustomLogic: NewLevel1Logic(),
		Entities: map[Entity]struct{}{
			NewPlayer(f32.Point{X: screenWidth / 2, Y: screenHeight - 10}): {},
		},
		entityLock: new(sync.Mutex),
	}
	go func() {
		th := material.NewTheme(gofont.Collection())
		ops := new(op.Ops)
		for e := range app.NewWindow(app.Size(500, 1000)).Events() {
			switch e := e.(type) {
			case system.DestroyEvent:
				return
			case system.FrameEvent:
				gtx := layout.NewContext(ops, e)
				key.InputOp{
					Tag:  gameTag,
					Keys: key.NameLeftArrow + "|" + key.NameRightArrow + "|" + key.NameUpArrow + "|" + key.NameDownArrow + "|" + key.NameSpace + "|P",
				}.Add(ops)
				for _, e := range e.Queue.Events(gameTag) {
					// Spaghetti switches
					switch e := e.(type) {
					case key.Event:
						switch e.Name {
						case key.NameLeftArrow:
							switch e.State {
							case key.Press:
								game.MovementVector.X = -Speed
							case key.Release:
								if game.MovementVector.X == -Speed {
									game.MovementVector.X = 0
								}
							}
						case key.NameRightArrow:
							switch e.State {
							case key.Press:
								game.MovementVector.X = Speed
							case key.Release:
								if game.MovementVector.X == Speed {
									game.MovementVector.X = 0
								}
							}
						case key.NameUpArrow:
							switch e.State {
							case key.Press:
								game.MovementVector.Y = -Speed
							case key.Release:
								if game.MovementVector.Y == -Speed {
									game.MovementVector.Y = 0
								}

							}
						case key.NameDownArrow:
							switch e.State {
							case key.Press:
								game.MovementVector.Y = Speed
							case key.Release:
								if game.MovementVector.Y == Speed {
									game.MovementVector.Y = 0
								}
							}
						case key.NameSpace:
							switch e.State {
							case key.Press:
								game.Shooting = true
							case key.Release:
								game.Shooting = false
							}
						case "P":
							if e.State == key.Press {
								game.Paused = !game.Paused
							}
						}
					}
				}
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(material.H6(th, fmt.Sprintf("Score: %d", game.Score)).Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						defer clip.Rect{Max: gtx.Constraints.Max}.Push(ops).Pop()
						paint.ColorOp{Color: color.NRGBA{A: 0xff}}.Add(ops)
						paint.PaintOp{}.Add(ops)
						return layout.Center.Layout(gtx, game.Layout)
					}),
				)
				if !game.Paused {
					op.InvalidateOp{}.Add(ops)
				}
				e.Frame(gtx.Ops)
			}
		}
	}()
	go func() {
		for range time.Tick(time.Second / Tickrate) {
			game.Logic()
		}
	}()
	app.Main()
}

var gameTag = new(struct{})

type Level struct {
	GameTime       float64
	CustomLogic    func(g *Level)
	CurrentTimer   int
	Entities       map[Entity]struct{}
	entityLock     *sync.Mutex
	MovementVector f32.Point
	Score          int64
	Shooting       bool
	Paused         bool
}

func (g *Level) Layout(gtx layout.Context) layout.Dimensions {
	scale := float32(gtx.Constraints.Max.X) / float32(screenWidth)
	if s := float32(gtx.Constraints.Max.Y) / float32(screenHeight); s < scale {
		scale = s
	}
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	defer op.Affine(f32.Affine2D{}.Scale(f32.Point{X: 0, Y: 0}, f32.Point{X: scale, Y: scale})).Push(gtx.Ops).Pop()
	defer clip.Rect{Max: image.Point{X: screenWidth, Y: screenHeight}}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color.NRGBA{0xff, 0xff, 0xff, 0xff}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	for e := range g.Entities {
		e.Draw(gtx.Ops)
	}
	return layout.Dimensions{Size: image.Point{X: int(float32(screenWidth) * scale), Y: int(float32(screenHeight) * scale)}}
}

func (g *Level) Logic() {
	if g.Paused {
		return
	}
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	g.GameTime += 1
	g.CustomLogic(g)
	for e := range g.Entities {
		e.Logic(g)
	}
}

type SpawnTimer struct {
	Time float64
	Entity
}

type Entity interface {
	Position() f32.Point
	Size() float32 // Radius of the entity
	Draw(ops *op.Ops)
	Team() Team
	Logic(*Level) // Perform any game logic the entity has
}

func collides(a, b Entity) bool {
	distance := a.Position().Sub(b.Position())
	requiredDistance := a.Size() + b.Size()
	return distance.X*distance.X+distance.Y*distance.Y < requiredDistance*requiredDistance
}

// Check if the entity is out of bounds
func OOB(e Entity) bool {
	return e.Position().X < 0 || e.Position().X > screenWidth || e.Position().Y < 0 || e.Position().Y > screenHeight
}

type Team uint8

const (
	PlayerTeam Team = iota
	EnemyTeam
	NeitherTeam // For hurting both
)

type Damageable interface {
	Hurt(*Level, int)
}
