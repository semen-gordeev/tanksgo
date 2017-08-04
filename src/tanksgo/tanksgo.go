package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	_ "net/http/pprof"
	"os"
	"strings"
)

const framesPerSecond = 8
const getReadyPause = 400
const maxNameLength = 25
const maxParallelRounds = 100
const maxPlayersPerRound = 5
const minPlayersPerRound = 1
const maxRoundWaitingTimeSec = 5
const maxRoundRunningTimeSec = 600
const maxSpeed = 5

const bonusPoint = 5
const lowFactor = 50
const highFactor = 5

const mapWidth = 179
const mapHeight = 38
const nameTableWidth = 30
const tankWidth = 3
const tankHeight = 3
const maxBombs = 100

const colorPrefix = "\x1b["
const colorPostfix = "m"
const bonus = "\xE2\x99\xA5"
const bomb = "\xE2\x9C\xB3"

var (
	tanks = [4][]byte{}
	// http://www.isthe.com/chongo/tech/comp/ansi_escapes.html
	home  = []byte{27, 91, 72}
	clear = []byte{27, 91, 50, 74}
	// [20A[90D
	middle = []byte{27, 91, 50, 48, 65, 27, 91, 57, 48, 68}
	enter = []byte{27, 91, 49, 66}
)

// States of the round
const (
	COMPILING = 1 + iota
	WAITING
	STARTING
	RUNNING
	FINISHED
)

// Directions
const (
	LEFT = 0 + iota
	RIGHT
	UP
	DOWN
)

// Car Borders
const (
	LEFTUP = 0 + iota
	RIGHTUP
	RIGHTDOWN
	LEFTDOWN
)

// Damage points
const (
	DAMAGE_BACK  = 2
	DAMAGE_FRONT = 4
	DAMAGE_SIDE  = 6
)

// Colors
const (
	RESET = 0
	BOLD  = 1
	RED   = 31
	GREEN = 32
	//YELLOW  = 33
	//BLUE    = 34
	//MAGENTA = 35
)

type Point struct {
	X, Y int
}

type Symbol struct {
	Color int
	Char  []byte
}

type Symbols []Symbol

func getModel(fileName string) ([]byte, error) {
	fileStat, err := os.Stat(fileName)
	if err != nil {
		fmt.Printf("File with model %s does not exist: %v\n", fileName, err)
		return []byte{}, err
	}

	model := make([]byte, fileStat.Size())
	f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("Error while opening %s: %v\n", fileName, err)
		os.Exit(1)
	}
	defer f.Close()

	f.Read(model)
	return model, nil
}

func getPlayerData(conn net.Conn, splash []byte) (Player, error) {
	_, err := conn.Write(clear)
	if err != nil {
		return Player{}, errors.New("Communication error")
	}
	_, err = conn.Write(home)
	if err != nil {
		return Player{}, errors.New("Communication error")
	}
	_, err = conn.Write(splash)
	if err != nil {
		return Player{}, errors.New("Communication error")
	}
	_, err = conn.Write(middle)
	if err != nil {
		return Player{}, errors.New("Communication error")
	}

	io := bufio.NewReader(conn)

	line, err := io.ReadString('\n')
	if err != nil {
		return Player{}, errors.New("Communication error")
	}

	name := strings.Replace(strings.Replace(line, "\n", "", -1), "\r", "", -1)
	if name == "" {
		return Player{}, errors.New("Empty name")
	}
	if len(name) > maxNameLength {
		return Player{}, errors.New("Too long name")
	}

	return Player{Conn: conn, Name: name, Health: 100, Tank: Tank{}}, nil
}

func (symbols Symbols) symbolsToByte() []byte {
	var returnSlice []byte
	for _, symbol := range symbols {
		// Should be something like \x1b[31m^\x1b[0m for symbols with colors or ^ without
		if symbol.Color != RESET {
			returnSlice = append(returnSlice, []byte(colorPrefix+fmt.Sprintf("%d", symbol.Color)+colorPostfix)...)
		}
		returnSlice = append(returnSlice, symbol.Char...)
		if symbol.Color != RESET {
			returnSlice = append(returnSlice, []byte(colorPrefix+fmt.Sprintf("%d", RESET)+colorPostfix)...)
		}
	}
	return returnSlice
}

func prepare(conn net.Conn, splash []byte, round *Round) {
	p, err := getPlayerData(conn, splash)
	if err != nil {
		conn.Close()
		return
	}
	p.initPlayer(len(round.Players))
	round.Players = append(round.Players, p)
	p.writeToThePlayer([]byte("Waiting for other players to join...\n"), true, true)
	round.writeToAllPlayers([]byte(fmt.Sprintf("%s is added\n", p.Name)), false, false)
}

func main() {
	var modelsPath string
	var port, countUsers, users int

	flag.IntVar(&port, "p", 8000, "Port to listen")
	flag.IntVar(&countUsers, "c", 3, "Count of users")
	flag.StringVar(&modelsPath, "m", "./models", "Models of tank location")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		os.Exit(2)
	}
	defer l.Close()

	// Read sketches
	tanks[LEFT], _ = getModel(modelsPath + "/" + "tankLeft.txt")
	tanks[RIGHT], _ = getModel(modelsPath + "/" + "tankRight.txt")
	tanks[UP], _ = getModel(modelsPath + "/" + "tankUp.txt")
	tanks[DOWN], _ = getModel(modelsPath + "/" + "tankDown.txt")
	splash, _ := getModel(modelsPath + "/" + "splash.txt")

	var round = Round{}

	for users < countUsers {
		conn, err := l.Accept()
		users++
		fmt.Println("In total", users, "users connected")
		if err != nil {
			fmt.Println("Failed to accept request", err)
		}

		go prepare(conn, splash, &round)
	}

	round.start()
}
