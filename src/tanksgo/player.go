package main

import (
	"net"
)

type Player struct {
	Conn      net.Conn
	Name      string
	Health    int64
	Color     int
	Bombs     int
	Tank      Tank
}

func (p *Player) initPlayer(id int) {
	switch id {
	case 0:
		initX, initY := 1, 1
		p.Tank.Borders = Rectangle{[4]Point{
			{initX, initY},
			{initX + tankWidth - 1, initY},
			{initX + tankWidth - 1, initY + tankHeight - 1},
			{initX, initY + tankHeight - 1}}}
		p.Tank.Direction = RIGHT
	case 1:
		initX, initY := mapWidth-nameTableWidth-tankWidth, 1
		p.Tank.Borders = Rectangle{[4]Point{
			{initX, initY},
			{initX + tankWidth - 1, initY},
			{initX + tankWidth, initY + tankHeight - 1},
			{initX, initY + tankHeight - 1}}}
		p.Tank.Direction = DOWN
	case 2:
		initX, initY := mapWidth-nameTableWidth-tankWidth, mapHeight-tankHeight-1
		p.Tank.Borders = Rectangle{[4]Point{
			{initX, initY},
			{initX + tankWidth - 1, initY},
			{initX + tankWidth - 1, initY + tankHeight - 1},
			{initX, initY + tankHeight - 1}}}
		p.Tank.Direction = LEFT
	case 3:
		initX, initY := 1, mapHeight-tankHeight-1
		p.Tank.Borders = Rectangle{[4]Point{
			{initX, initY},
			{initX + tankWidth - 1, initY},
			{initX + tankWidth - 1, initY + tankHeight - 1},
			{initX, initY + tankHeight - 1}}}
		p.Tank.Direction = UP
	case 4:
		initX, initY := (mapWidth-nameTableWidth)/2, mapHeight/2
		p.Tank.Borders = Rectangle{[4]Point{
			{initX, initY},
			{initX + tankWidth - 1, initY},
			{initX + tankWidth - 1, initY + tankHeight - 1},
			{initX, initY + tankHeight - 1}}}
		p.Tank.Direction = DOWN
	}
	p.Bombs = maxBombs
	p.Color = RED + id
}

func (player *Player) readDirection(round *Round) {
	if initTelnet(player.Conn) != nil {
		return
	}

	for {
		if player.Health <= 0 {
			return
		}
		direction := make([]byte, 1)

		// Read all possible bytes and try to find a sequence of:
		// ESC [ cursor_key
		escpos := 0
		for {
			_, err := player.Conn.Read(direction)
			if err != nil {
				player.Health = 0
				return
			}

			// Check if telnet want to negotiate something
			if escpos == 0 && direction[0] == 255 {
				readTelnet(player.Conn)
			} else if escpos == 0 && direction[0] == 3 {
				// Ctrl+C
				player.Health = 0
				return
			} else if escpos == 0 && direction[0] == 32 {
				// Space
				if player.Bombs > 0 {
					player.Bombs--
				}
			} else if escpos == 0 && direction[0] == 27 {
				escpos = 1
			} else if escpos == 1 && direction[0] == 91 {
				escpos = 2
			} else if escpos == 2 {
				break
			}
		}
		switch direction[0] {
		case 68:
			// Left
			if player.Tank.Direction != LEFT {
				player.Tank.Direction = LEFT
			} else {
				player.Tank.moveTo(LEFT, player.willBeCrash(LEFT, round))
			}
		case 67:
			// Right
			if player.Tank.Direction != RIGHT {
				player.Tank.Direction = RIGHT
			} else {
				player.Tank.moveTo(RIGHT, player.willBeCrash(RIGHT, round))
			}
		case 65:
			// Up
			if player.Tank.Direction != UP {
				player.Tank.Direction = UP
			} else {
				player.Tank.moveTo(UP, player.willBeCrash(UP, round))
			}
		case 66:
			// Down
			if player.Tank.Direction != DOWN {
				player.Tank.Direction = DOWN
			} else {
				player.Tank.moveTo(DOWN, player.willBeCrash(DOWN, round))
			}
		}
	}
}

func (player *Player) willBeCrash(direction int, round *Round) bool {
	result := false
	tmpTank := player.Tank.nextTo(direction)
	// hit wall
	for _, point := range tmpTank.Borders.Points {
		if point.X < 1 || point.X >= mapWidth-nameTableWidth || point.Y < 1 || point.Y >= mapHeight-1 {
			player.Health -= crashDamage
			result = true
		}
	}
	// hit another tank
	for _, opponent := range round.Players {
		if player.Name == opponent.Name {
			continue
		}
		if tmpTank.Borders.intersects(&opponent.Tank.Borders) {
			player.Health -= crashDamage
			opponent.Health -= ramDamage
			result = true
		}
	}
	return result
}

func (player *Player) writeToThePlayer(message []byte, clean bool, go_home bool) {
	if clean {
		_, err := player.Conn.Write(clear)
		if err != nil {
			// Kick user if connection got lost
			player.Health = 0
			return
		}
	}

	if go_home {
		_, err := player.Conn.Write(home)
		if err != nil {
			player.Health = 0
			return
		}
	}

	_, err := player.Conn.Write(message)
	if err != nil {
		player.Health = 0
		return
	}
}
