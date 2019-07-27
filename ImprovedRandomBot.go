package main

import (
	"hlt"
	"math/rand"
)

func main() {
	conn, gameMap := hlt.NewConnection()
	conn.SendName("ImprovedRandomBot")
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

func assignMove(loc hlt.Location, gameMap hlt.GameMap, selfId int) hlt.Move {
	site := gameMap.GetSite(loc, hlt.STILL)

	for i := 0; i < 5; i += 1 {
		if gameMap.GetSite(loc, hlt.Direction(i)).Owner != selfId && gameMap.GetSite(loc, hlt.Direction(i)).Strength < site.Strength {
			return hlt.Move{
				Location: loc,
				Direction: hlt.Direction(i),
			}
		}
	}
	
	if site.Strength < 5 * site.Production {
		return hlt.Move{
			Location: loc,
			Direction: hlt.STILL,
		}
	} else {
		return hlt.Move{
			Location:  loc,
			Direction: hlt.Direction(func() int{ if rand.Intn(5) > 2 { return 1 } else { return 4 } }()),
		}
	}
}