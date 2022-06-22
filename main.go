/*
This file is part of gvsbh.

gvsbh is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

gvsbh is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with gvsbh. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"log"
	"sync"

	"gioui.org/io/key"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const Speed = .5     // Speed is the speed of the ship (affected by Tickrate)
const Tickrate = 240 // Number of game ticks per second

// Dimensions of the screen in sprite pixels
const (
	screenHeight = 200
	screenWidth  = 100
)

var scale float64 = 1
var touches, touchQueue []ebiten.TouchID

// Keybindings
const (
	LeftKey  = key.NameLeftArrow
	RightKey = key.NameRightArrow
	UpKey    = key.NameUpArrow
	DownKey  = key.NameDownArrow
	ShootKey = key.NameSpace
	PauseKey = "P"
)

func main() {
	game := &Level{
		CustomLogic: NewLevel1Logic(),
		Entities: map[Entity]struct{}{
			NewPlayer(screenWidth/2, screenHeight-10): {},
		},
		entityLock: new(sync.Mutex),
	}
	ebiten.SetWindowTitle("gvsbh")
	ebiten.SetMaxTPS(Tickrate)
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Level struct {
	GameTime             float64
	CustomLogic          func(g *Level)
	CurrentTimer         int
	Entities             map[Entity]struct{}
	entityLock           *sync.Mutex
	MovementX, MovementY float64
	Score                int64
	Shooting             bool
	Paused               bool
	// Touch related information
	UseTouch       bool
	TouchX, TouchY float64
}

func (g *Level) Draw(screen *ebiten.Image) {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	for e := range g.Entities {
		e.Draw(screen)
	}
}
func (g *Level) Layout(outsideWidth, outsideHeight int) (int, int) {
	// SCaling is done manually so it will be smooth
	newScale := float64(outsideWidth) / screenWidth
	if s := float64(outsideHeight) / screenHeight; s < newScale {
		scale = s
	} else {
		scale = newScale
	}
	return int(screenWidth * scale), int(screenHeight * scale)
}
func (g *Level) Update() (err error) {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.Paused = !g.Paused
	}
	if g.Paused {
		return
	}
	g.MovementX = 0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.MovementX -= Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.MovementX += Speed
	}
	g.MovementY = 0
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.MovementY -= Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.MovementY += Speed
	}
	g.Shooting = ebiten.IsKeyPressed(ebiten.KeySpace)
	touches = ebiten.AppendTouchIDs(touches[:0])
	touchQueue = inpututil.AppendJustPressedTouchIDs(touchQueue)
	if len(touches) > 0 {
		if len(touches) > 1 {
			g.Shooting = true
		}
		g.UseTouch = true
		var x, y int
		for len(touchQueue) > 0 {
			x, y = ebiten.TouchPosition(touchQueue[0])
			if x == 0 && y == 0 {
				touchQueue = touchQueue[1:]
			} else {
				break
			}
		}
		g.TouchX = float64(x) / scale
		g.TouchY = float64(y) / scale
	} else {
		g.UseTouch = false
	}
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	g.CustomLogic(g)
	for e := range g.Entities {
		e.Logic(g)
	}
	return
}

type SpawnTimer struct {
	Time float64
	Entity
}

type Entity interface {
	Position() (x float64, y float64)
	Size() float64 // Radius of the entity
	Draw(screen *ebiten.Image)
	Team() Team
	Logic(*Level) // Perform any game logic the entity has
}

func collides(a, b Entity) bool {
	x, y := a.Position()
	x2, y2 := b.Position()
	x -= x2
	y -= y2
	d := a.Size() + b.Size()
	return x*x+y*y < d*d
}

// Check if the entity is out of bounds
func OOB(e Entity) bool {
	x, y := e.Position()
	return x < 0 || x > screenWidth || y < 0 || y > screenHeight
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
