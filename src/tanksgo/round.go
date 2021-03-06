package main

import (
	"sync"
	"time"
	"fmt"
)

type Round struct {
	Players         []*Player
	Bombs           []*Bomb
	FrameBuffer     Symbols
	State			int
	countPlayers 	int
	sync.Mutex
}

func (round *Round) generateMap() {
	for row := 0; row < mapHeight; row++ {
		for column := 0; column < mapWidth; column++ {
			var char []byte
			if (row == 0 || row == mapHeight-1) && column < mapWidth-2 {
				char = []byte("─")
			} else if column == 0 || column == mapWidth-3 || column == mapWidth-nameTableWidth {
				char = []byte("│")
			} else if column == mapWidth-2 {
				char = []byte("\r")
			} else if column == mapWidth-1 {
				char = []byte("\n")
			} else {
				char = []byte(" ")
			}
			round.FrameBuffer[row*mapWidth+column] = Symbol{0, char}
		}
	}
}

func (round *Round) writeToAllPlayers(message []byte, clean bool, go_home bool) {
	for i := range round.Players {
		go round.Players[i].writeToThePlayer(message, clean, go_home)
	}
}

func (round *Round) gameLogic() {
	for i := range round.Players {
		go round.Players[i].readDirection(round)
	}
	go round.moveBombs()
}

func (round *Round) moveBombs() {
	for {
		for _, b := range round.Bombs {
			if b.Active {
				b.move()
				if !b.Active {
					continue
				}
				for _, p := range round.Players {
					if p.Tank.Borders.pointInside(b.Point) {
						if p.Health > 0 {
							p.Health -= bombDamage
						}
						b.Active = false
					}
				}
			}
		}
		time.Sleep(1.0 / framesPerSecond * 1000 * time.Millisecond)
	}
}

func (round *Round) applyNames(lineBetweenPlayersInBar int) {
	for line, player := range round.Players {
		for i, char := range []byte(player.Name) {
			round.FrameBuffer[(line*lineBetweenPlayersInBar+1)*mapWidth+(mapWidth-nameTableWidth+1)+i] = Symbol{player.Color, []byte{char}}
		}
		round.FrameBuffer[(line*lineBetweenPlayersInBar+1)*mapWidth+(mapWidth-nameTableWidth+1)+len(player.Name)] = Symbol{RESET, []byte{':'}}
	}
}

func (round *Round) applyUserData(activeFrameBuffer []Symbol, lineBetweenPlayersInBar int) {
	for num, player := range round.Players {
		if player.Health <= 0 {
			player.Color = BOLD
		}
		// Apply health
		health := []byte(fmt.Sprintf("Health: %3d", player.Health))
		for i, char := range health {
			// +1 because health is next line after the name
			activeFrameBuffer[((num*lineBetweenPlayersInBar+1)+1)*mapWidth+(mapWidth-3)-len(health)+i] = Symbol{player.Color, []byte{char}}
		}

		// Apply the amount of bombs to the bar
		bombs := []byte(fmt.Sprintf("Bombs: %4d", round.Players[num].Bombs))
		for i, char := range bombs {
			// +2 because "bombs" is next line after the name
			activeFrameBuffer[((num*lineBetweenPlayersInBar+2)+1)*mapWidth+(mapWidth-3)-len(bombs)+i] = Symbol{player.Color, []byte{char}}
		}
	}
}

func (round *Round) applyTanks(activeMap []Symbol) {
	defer func() {
		recover()
	}()

	for _, player := range round.Players {
		charPosX, charPosY := 0, 0
		for i := 0; i < len(tanks[player.Tank.Direction]); i++ {
			var chars []byte
			if tanks[player.Tank.Direction][i] == byte('\n') {
				charPosY++
				charPosX = 0
				continue
			} else if tanks[player.Tank.Direction][i] == 226 {
				/*
				 This means extended ASCII is used. After 226 2 bytes must follow
				*/
				chars = []byte{tanks[player.Tank.Direction][i], tanks[player.Tank.Direction][i+1], tanks[player.Tank.Direction][i+2]}
				i += 2
			} else if tanks[player.Tank.Direction][i] == 194 {
				/*
				 This means extended ASCII is used. After 194 1 bytes must follow
				*/
				chars = []byte{tanks[player.Tank.Direction][i], tanks[player.Tank.Direction][i+1]}
				i++
			} else {
				chars = []byte{tanks[player.Tank.Direction][i]}
			}
			activeMap[(player.Tank.Borders.Points[LEFTUP].Y+charPosY)*mapWidth+player.Tank.Borders.Points[LEFTUP].X+charPosX] = Symbol{player.Color, chars}
			charPosX++
		}
	}
}

func (round *Round) applyBombs(activeFrameBuffer []Symbol) {
	if round.State != STARTED {
		return
	}
	for _, b := range round.Bombs {
		if b.Active {
			activeFrameBuffer[b.Point.Y*mapWidth+b.Point.X] = Symbol{BOLD, []byte(bomb)}
		}
	}
}

func (round *Round) applyGetReady(activeFrameBuffer []Symbol, getReadyCounter *int) {
	if round.State == INIT {
		getReady := "GET READY!"
		if *getReadyCounter == 0 {
			fmt.Println("Round has started!")
			round.State = STARTED
		} else if *getReadyCounter <= framesPerSecond*1 {
			getReady += " 1"
		} else if *getReadyCounter <= framesPerSecond*2 {
			getReady += " 2"
		} else if *getReadyCounter <= framesPerSecond*3 {
			getReady += " 3"
		}

		for i, char := range []byte(getReady) {
			activeFrameBuffer[mapWidth*(mapHeight/2-2)+mapWidth/2-len(getReady)/2+i] = Symbol{GREEN, []byte{char}}
		}
		*getReadyCounter--
	}
}

func (round *Round) checkGameOver(activeFrameBuffer Symbols) bool {
	deadPlayers := 0
	winnersName := ""

	for _, p := range round.Players {
		if p.Health <= 0 {
			deadPlayers++
		} else {
			winnersName = p.Name
		}
	}
	if round.countPlayers-deadPlayers <= 1 {
		fmt.Println("Round has finished!")
		winnerStr := "THE WINNER IS " + winnersName + "!!!"
		for i, char := range []byte(winnerStr) {
			activeFrameBuffer[mapWidth*(mapHeight/2-2)+mapWidth/2-len(winnerStr)/2+i] = Symbol{GREEN, []byte{char}}
		}
		round.writeToAllPlayers(activeFrameBuffer.symbolsToByte(), false, true)
		time.Sleep(5 * time.Second)
		round.State = FINISHED
		return true
	} else {
		return false
	}
}

func (round *Round) start() {
	for len(round.Players) != round.countPlayers {
		time.Sleep(1*time.Second)
	}

	lineBetweenPlayersInBar := mapHeight / len(round.Players)
	getReadyCounter := getReadyPause

	round.gameLogic()
	round.generateMap()
	round.applyNames(lineBetweenPlayersInBar)

	for {
		activeFrameBuffer := make(Symbols, len(round.FrameBuffer))
		copy(activeFrameBuffer, round.FrameBuffer)

		round.applyGetReady(activeFrameBuffer, &getReadyCounter)
		round.applyUserData(activeFrameBuffer, lineBetweenPlayersInBar)
		round.applyTanks(activeFrameBuffer)
		round.applyBombs(activeFrameBuffer)

		round.writeToAllPlayers(activeFrameBuffer.symbolsToByte(), false, true)

		if round.checkGameOver(activeFrameBuffer) {
			break
		}

		time.Sleep(1.0/framesPerSecond * 1000 * time.Millisecond)
	}
}
