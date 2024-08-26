package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/fatih/color"
)

func main() {
	length := flag.Int("len", 10, "The number of words")
	flag.Parse()

	if *length <= 0 {
		fmt.Println("The length must be greater than 0")
		os.Exit(1)
	}

	oldTermState := setTermRawMode()
	defer restoreTermMode(oldTermState)

	var currentInput string
	var startTime time.Time

	targetText := generateText(length)

	fmt.Println("Type the following text (Ctrl + C to exit) :")

	for len(currentInput) != len(targetText) {
		printText(targetText, currentInput)

		b := readByte()

		if b == '\x03' {
			fmt.Println("\nExiting...")
			os.Exit(0)
		} else if b == '\b' && len(currentInput) > 0 {
			currentInput = currentInput[:len(currentInput)-1]
		} else if b == '\x17' {
			for len(currentInput) > 0 && currentInput[len(currentInput)-1] != ' ' {
				currentInput = currentInput[:len(currentInput)-1]
			}
		} else {
			currentInput += string(b)

			if len(currentInput) == 1 {
				startTime = time.Now()
			}
		}

	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	wpm := calculateWpm(targetText, currentInput, duration)
	acc := calculateAccuracy(targetText, currentInput)

	fmt.Printf("\n\nWPM: %.2f\nAccuracy: %.2f%%\n", wpm, acc)
}

func generateText(length *int) string {
	dictFile, err := os.ReadFile("data/dict.json")

	if err != nil {
		fmt.Printf("An error happened while reading the dictionary file: %v", err)
		os.Exit(1)
	}

	var words []string
	if err := json.Unmarshal(dictFile, &words); err != nil {
		fmt.Printf("An error happened while unmarshalling the dictionary file: %v", err)
		os.Exit(1)
	}


	var text string

	for i := 0; i < *length; i++ {
		text += pickRandom(words[:])

		if i != *length - 1 {
			text += " "
		}
	}

	return text
}

func pickRandom(arr []string) string {
	return arr[rand.Intn(len(arr))]
}

func setTermRawMode() *term.State {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		fmt.Printf("An error happened while setting the terminal to raw mode: %v", err)
		os.Exit(1)
	}

	return oldState
}

func restoreTermMode(oldState *term.State) {
	term.Restore(int(os.Stdin.Fd()), oldState)
}

func readByte() byte {
	b := make([]byte, 1)
	_, err := os.Stdin.Read(b)

	if err != nil {
		fmt.Printf("An error happened while reading the user input: %v", err)
		os.Exit(1)
	}

	return b[0]
}

func printText(text, current string) {
	cursorColor := color.New(color.FgWhite).Add(color.Underline)
	defaultColor := color.New(color.FgWhite)
	correctColor := color.New(color.FgBlue)
	wrongColor := color.New(color.FgRed)

	for i, r := range text {
		char := string(r)
		if len(current) <= i {
			if len(current) == i {
				cursorColor.Print(char)
			} else {
				defaultColor.Print(char)
			}
		} else {
			if char == string(current[i]) {
				correctColor.Print(char)
			} else {
				if char == " " {
					wrongColor.Print("_")
				} else {
					wrongColor.Print(char)
				}
			}
		}
	}

	fmt.Print("\r")
}

func calculateWpm(target, typed string, d time.Duration) float64 {
	gross := float64(len(typed) / 5) / d.Minutes()
	mistakes := countMistakes(target, typed)
	net := gross - float64(mistakes) / d.Minutes()

	return math.Abs(net)
}

func calculateAccuracy(target, typed string) float64 {
	mistakes := countMistakes(target, typed)

	return 100 - (float64(mistakes) / float64(len(target))) * 100
}

func countMistakes(target, typed string) int {
	mistakes := 0

	for i, r := range target {
		if i >= len(typed) {
			break
		}

		if r != rune(typed[i]) {
			mistakes++
		}
	}

	return mistakes
}