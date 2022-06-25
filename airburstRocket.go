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

var AirburstRocketSprite = Resource("res/AirburstRocket.png")

type AirburstRocket struct {
	x, y        float64
	team        Team
	orientation float64
	speed       float64
}

func (r *AirburstRocket) Position() (x, y float64) {
	return r.x, r.y
}
func (r *AirburstRocket) Size() float64 {
	return 20 // Fake oversized radius for airbursting
}
func (r *AirburstRocket) Team() Team {
	return r.team
}
func (r *AirburstRocket) Draw(s *ebiten.Image) {
	AirburstRocketSprite.Draw(s, r.x, r.y, r.orientation)
}
func (r *AirburstRocket) Logic(g *Level) {
	r.y -= math.Cos(r.orientation) * r.speed
	r.x += math.Sin(r.orientation) * r.speed
	if OOB(r) {
		delete(g.Entities, r)
	}
	for e := range g.Entities {
		if e.Team() != r.team && collides(r, e) {
			// Airburst if near a damagable entity
			if _, ok := e.(Damageable); ok {
				for i := -4; i < 5; i++ {
					g.Entities[&Bullet{
						x: r.x, y: r.y,
						team:        r.team,
						orientation: r.orientation + math.Pi/20*float64(i),
						speed:       r.speed,
					}] = struct{}{}
				}
				delete(g.Entities, r)
			}
		}
	}
}
