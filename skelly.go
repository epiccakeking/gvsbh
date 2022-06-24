package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var SkellySprite = Resource("res/SkellyHead.png")

type Skelly struct {
	x, y   float64
	team   Team
	health int
}

// Skelly is made of multiple entities, so it returns all parts.
func NewSkelly(x, y float64) []Entity {
	head := &Skelly{x: x, y: y, team: EnemyTeam, health: 100}
	return []Entity{
		head,
		&SkellyArm{Parent: head, offsetX: -20, offsetY: 10, health: 100},
		&SkellyArm{Parent: head, offsetX: 20, offsetY: 10, health: 100},
	}
}

func (s *Skelly) Position() (x, y float64) {
	return s.x, s.y
}
func (s *Skelly) Size() float64 {
	return 10
}
func (s *Skelly) Team() Team {
	return s.team
}
func (s *Skelly) Draw(screen *ebiten.Image) {
	SkellySprite.Draw(screen, s.x, s.y, 0)
}
func (s *Skelly) Logic(g *Level) {
}

func (s *Skelly) Hurt(g *Level, damage int) {
	s.health -= damage
}

var SkellyArmRightSprite = Resource("res/SkellyArmRight.png")

type SkellyArm struct {
	Parent           *Skelly
	offsetX, offsetY float64
	orientation      float64
	health           int
	shotTimer        int
}

func (a *SkellyArm) Position() (x, y float64) {
	pX, pY := a.Parent.Position()
	return pX + a.offsetX + a.Size()*math.Sin(a.orientation), pY + a.offsetY - a.Size()*math.Cos(a.orientation)
}
func (a *SkellyArm) Size() float64 {
	return 10
}
func (a *SkellyArm) Team() Team {
	return a.Parent.Team()
}
func (a *SkellyArm) Draw(screen *ebiten.Image) {
	x, y := a.Position()
	SkellyArmRightSprite.Draw(screen, x, y, a.orientation)
}
func (a *SkellyArm) Logic(g *Level) {
	nearestDistanceSquared := math.Inf(1)
	// x and y are calculated manually because Position moves based on orientation
	x, y := a.Parent.Position()
	x += a.offsetX
	y += a.offsetY
	var nearest Entity
	for e := range g.Entities {
		if e.Team() != a.Team() {
			if _, ok := e.(Damageable); !ok {
				continue
			}
			dX, dY := e.Position()
			dX -= x
			dY -= y
			if dSquared := dX*dX + dY*dY; dSquared < nearestDistanceSquared {
				nearestDistanceSquared = dSquared
				nearest = e
			}
		}
	}
	if nearest == nil {
		return
	}
	dX, dY := nearest.Position()
	dX -= x
	dY -= y
	if dY == 0 {
		if dX < 0 {
			a.orientation = -math.Pi / 2
		} else {
			a.orientation = math.Pi / 2
		}
	} else {
		a.orientation = -math.Atan(dX / dY)
		if dY > 0 {
			a.orientation += math.Pi
		}
	}
	if a.shotTimer <= 0 {
		a.shotTimer = Tickrate
		g.AddEntity(&Bullet{x: x, y: y, team: a.Team(), speed: 1, orientation: a.orientation})
	} else {
		a.shotTimer--
	}
}

func (a *SkellyArm) Hurt(g *Level, damage int) {
	a.health -= damage
}
