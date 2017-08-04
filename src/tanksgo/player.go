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
