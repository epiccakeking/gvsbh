/*
This file is part of gvsbh.

gvsbh is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

gvsbh is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with gvsbh. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"math"
	"math/rand"

	"gioui.org/f32"
)

func NewLevel1Logic() func(g *Level) {
	rng := rand.NewSource(0)
	phase := 0
	var tick uint64 // Keep track of logic ticks
	return func(g *Level) {
		switch phase {
		case 0: // Pre-game waiting
			if tick > Tickrate*2 {
				phase++
				tick = 0
			}
		case 1:
			if tick%(Tickrate/3) == 0 {
				g.Entities[NewBomb(f32.Point{X: float32(rng.Int63()) / math.MaxInt64 * screenWidth})] = struct{}{}
			}
			if tick == Tickrate*10 {
				phase++
				tick = 0
			}
		case 2:
			if tick == Tickrate {
				g.Entities[NewPulsar(f32.Point{X: 25, Y: 50})] = struct{}{}
			} else if tick == Tickrate*2 {
				g.Entities[NewPulsar(f32.Point{X: screenWidth - 25, Y: 50})] = struct{}{}
			}
			if tick == Tickrate*10 {
				phase++
				tick = 0
			}
		}
		tick++
	}
}
