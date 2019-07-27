package main

import (
	"hlt"
	"sort"
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
	ratio int
}

func getNeighbors(loc hlt.Location, gameMap hlt.GameMap) []Movement {
	var moves []Movement
	for _, dir := range hlt.Directions {
		site := gameMap.GetSite(loc, dir)
		moves = append(moves, Movement{site, dir, func() int {
			if site.Strength != 0 {
				return site.Production / site.Strength
			}
			return site.Production
		}()})
	}
	return moves
}

func nearestBorderDirection(loc hlt.Location, selfId int, gameMap hlt.GameMap) hlt.Direction {
	direction := hlt.Directions[0]
	maxDistance := gameMap.Width/2
	for _, d := range hlt.Directions {
		distance := 0
		current := loc
		site := gameMap.GetSite(current, d)
		for site.Owner == selfId && distance < maxDistance {
			distance += 1
			current = gameMap.GetLocation(current, d)
			site = gameMap.GetSite(current, d)
		}

		if distance < maxDistance {
			direction = d
			maxDistance = distance
		}

	}

	return direction

}

func assignMove(loc hlt.Location, gameMap hlt.GameMap, selfId int) hlt.Move {
	site := gameMap.GetSite(loc, hlt.STILL)
	neighborStack := getNeighbors(loc, gameMap)
	minStrength := Movement{site, hlt.STILL, 0}
	minStrength.site.Strength = 256

	if site.Strength <= 5*site.Production {
		return hlt.Move{
			Location: loc,
			Direction: hlt.STILL,
		}
	}

	var notNiceNeighbors []Movement
	notNiceNeighbors = func() []Movement {
		for _, move := range neighborStack {
			if move.site.Owner != selfId {
				notNiceNeighbors = append(notNiceNeighbors, move)
			}
		}
		return notNiceNeighbors
	}()

	sort.Slice(notNiceNeighbors, func(i int, j int) bool {
		return notNiceNeighbors[i].ratio > notNiceNeighbors[i].ratio
	})
	

	if len(notNiceNeighbors) != 0 && notNiceNeighbors[0].site.Strength <= site.Strength {
		return hlt.Move {
			Location: loc,
			Direction: notNiceNeighbors[0].direction,
		}
	}

	return hlt.Move {
		Location: loc,
		Direction: nearestBorderDirection(loc, selfId, gameMap),
	}
	
}