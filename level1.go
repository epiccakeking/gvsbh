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
	"math/rand"
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
				g.AddEntity(NewBomb(float64(rng.Int63())/math.MaxInt64*screenWidth, 0))
			}
			if tick == Tickrate*10 {
				phase++
				tick = 0
			}
		case 2:
			if tick == Tickrate {
				g.AddEntity(NewPulsar(25, 50))
			} else if tick >= Tickrate*2 {
				g.AddEntity(NewPulsar(screenWidth-25, 50))
				phase++
				tick = 0
			}
		case 3:
			numEnemies := 0
			for e := range g.Entities {
				if e.Team() == EnemyTeam {
					numEnemies++
				}
			}
			// Prevent advancing stage if there are still enemies
			if numEnemies != 0 {
				tick = 0
			}
			if tick > Tickrate {
				phase++
				tick = 0
			}

		case 4:
			phase++
			tick = 0
			for _, e := range NewSkelly(screenWidth/2, -30) {
				g.AddEntity(e)
			}
		}
		tick++
	}
}
