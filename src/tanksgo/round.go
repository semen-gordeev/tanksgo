package main

import (
	"sync"
	"time"
)

type Round struct {
	Players         []Player
	//Bombs           []Bomb
	FrameBuffer     Symbols
	sync.Mutex
}

func (round *Round) generateMap() {
	// http://www.theasciicode.com.ar
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
	//for i := range round.Players {
	//	go round.Players[i].readDirection(round)
	//	go round.Players[i].checkPosition(round)
	//	go round.Players[i].checkSpeed(round)
	//	go round.Players[i].checkHealth(round)
	//	go round.Players[i].checkBomb(round)
	//}
}

func (round *Round) applyNames(lineBetweenPlayersInBar int) {
	for line, player := range round.Players {
		for i, char := range []byte(player.Name) {
			round.FrameBuffer[(line*lineBetweenPlayersInBar+1)*mapWidth+(mapWidth-nameTableWidth+1)+i] = Symbol{player.Color, []byte{char}}
		}
		round.FrameBuffer[(line*lineBetweenPlayersInBar+1)*mapWidth+(mapWidth-nameTableWidth+1)+len(player.Name)] = Symbol{RESET, []byte{':'}}
	}
}

// We start round only if more than 0 player is presented
func (round *Round) start() {
	lineBetweenPlayersInBar := mapHeight / len(round.Players)
	//getReadyCounter := getReadyPause / framesPerSecond

	//round.gameLogic()
	round.generateMap()
	round.applyNames(lineBetweenPlayersInBar)

	for {
		activeFrameBuffer := make(Symbols, len(round.FrameBuffer))
		copy(activeFrameBuffer, round.FrameBuffer)

		//round.applyUserData(activeFrameBuffer, lineBetweenPlayersInBar)
		//round.applyBonus(activeFrameBuffer)
		//round.applyBombs(activeFrameBuffer, lineBetweenPlayersInBar)
		//round.applyCars(activeFrameBuffer)
		//round.applyGetReady(activeFrameBuffer, &getReadyCounter)
		//
		//round.writeToAllPlayers(activeFrameBuffer.symbolsToByte(), false)
		//
		//round.checkGameOver(activeFrameBuffer)

		time.Sleep(1 % framesPerSecond * 100 * time.Millisecond)
	}
}
