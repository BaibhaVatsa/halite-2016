package main

import (
	"hlt"
	// "fmt"
	// "math/rand"
)

func main() {
	conn, gameMap := hlt.NewConnection()
	conn.SendName("MyBot")
	for {
		var moves hlt.MoveSet
		gameMap = conn.GetFrame()
		for y := 0; y < gameMap.Height; y++ {
			for x := 0; x < gameMap.Width; x++ {
				loc := hlt.NewLocation(x, y)
				if gameMap.GetSite(loc, hlt.STILL).Owner == conn.PlayerTag {
					moves = append(moves, assignMove(loc, gameMap, conn.PlayerTag))
				}
			}
		}
		conn.SendFrame(moves)
	}
}

type Movement struct {
	site hlt.Site
	direction hlt.Direction
}

func getNeighbors(loc hlt.Location, gameMap hlt.GameMap) []Movement {
	var moves []Movement
	for _, dir := range hlt.Directions {
		moves = append(moves, Movement{gameMap.GetSite(loc, dir), dir})
	}
	return moves
}

func assignMove(loc hlt.Location, gameMap hlt.GameMap, selfId int) hlt.Move {
	site := gameMap.GetSite(loc, hlt.STILL)
	neighborStack := getNeighbors(loc, gameMap)
	minStrength := Movement{site, hlt.STILL}
	minStrength.site.Strength = 256

	for _, move := range neighborStack {
		if move.site.Owner != selfId && move.site.Strength < site.Strength {
			return hlt.Move{
				Location: loc,
				Direction: move.direction,
			}
		}
		
		if move.site.Owner == selfId && minStrength.site.Strength > move.site.Strength {
			minStrength = move
		}
	}
	
	if minStrength.site.Strength != 256 {
		return hlt.Move{
			Location:  loc,
			Direction: minStrength.direction,
		}
	} else {
		return hlt.Move{
			Location: loc,
			Direction: hlt.STILL,
		}
	}
	
}