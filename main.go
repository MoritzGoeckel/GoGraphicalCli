package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	//"io"
	//"math/rand"
	"os"
	"strings"
	"time"
)

func clear() {
	fmt.Println("\033[2J")
}

var end = false
var pressedKeys = ""

// https://godoc.org/github.com/pborman/ansi
const (
	Green = "\033[32m"
	Red   = "\033[31m"
	Black = "\033[30m"
	Blue  = "\033[34m"
	White = "\033[37m"

	Block = '\u2588'
)

type Pixel struct {
	color string
	char  rune
}

func main() {
	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		panic("stdin/stdout should be terminal")
	}

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	clear()

	width, height, _ := terminal.GetSize(0)
	go runDrawingLoop(width, height)
	go runInputLoop()
	go runUpdateLoop()

	for end == false {
		time.Sleep(100 * time.Millisecond)
	}

	clear()
}

func runDrawingLoop(width int, height int) {
	for {
		frame := draw(width, height)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				pixel := frame[y][x]
				//if pixel == nil {
				//	fmt.Print(string(Black) + string(Block) + "\033[0m")
				//} else {
				fmt.Print(string(pixel.color) + string(pixel.char) + "\033[0m")
				//}
			}
			fmt.Print("\r\n")
		}
		time.Sleep(120 * time.Millisecond)
	}
}

func runUpdateLoop() {
	for {
		doUpdate(100, pressedKeys)
		pressedKeys = ""
		time.Sleep(120 * time.Millisecond)
	}
}

func runInputLoop() {
	reader := bufio.NewReader(os.Stdin)
	char := '0'

	for {
		char, _, _ = reader.ReadRune()
		//Save in keys table
		pressedKeys += string(char)
		if char == 'X' {
			quit()
		}
	}
}

func quit() {
	end = true
}

var px = 5
var py = 29

func doUpdate(eTime int, pressedKeys string) {
	if py == 28 {
		py = 29
	}

	if strings.ContainsRune(pressedKeys, 'd') {
		px++
	}

	if strings.ContainsRune(pressedKeys, 'a') {
		px--
	}

	if strings.ContainsRune(pressedKeys, ' ') {
		py--
	}
}

func draw(width int, height int) [][]Pixel {
	frame := make([][]Pixel, height)
	for y := 0; y < height; y++ {
		frame[y] = make([]Pixel, width)
		for x := 0; x < width; x++ {
			if y == 30 {
				frame[y][x] = Pixel{Green, Block}
			} else {
				frame[y][x] = Pixel{Black, Block}
			}
		}
	}

	frame[py][px] = Pixel{Green, Block}

	return frame
}
