package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

type ReadState int

const (
	INACTIVE   ReadState = 0 // In this state, the only input that will advance the state is "do()".
	ACTIVE     ReadState = 1 // In this state, the program will actively be looking for operation commands.
	FIRST_ARG  ReadState = 2 // In this state, the program is still reading the first argument of the mult() command.
	SECOND_ARG ReadState = 3 // In this state, the program is still reading the second argument of the mult() command.
)

const bufferLen = 7

type StateInternal struct {
	state  ReadState
	buffer []rune
	count  int
	first  int
	second int
}

func ResetState(state *StateInternal) {
	state.state = ACTIVE
	state.buffer = make([]rune, bufferLen)
	state.count = 0
	state.first = 0
	state.second = 0
}

func ShiftStateBuffer(input rune, state *StateInternal) {
	for i := 0; i < bufferLen-1; i++ {
		state.buffer[i] = state.buffer[i+1]
	}

	state.buffer[bufferLen-1] = input
}

func HandleInactive(input rune, state *StateInternal) {
	ShiftStateBuffer(input, state)

	if input == ')' && string(state.buffer[bufferLen-4:]) == "do()" {
		state.state = ACTIVE
	}
}

func HandleActive(input rune, state *StateInternal) {
	ShiftStateBuffer(input, state)

	if input == '(' && string(state.buffer[bufferLen-4:]) == "mul(" {
		state.state = FIRST_ARG
	} else if input == ')' && string(state.buffer) == "don't()" {
		state.state = INACTIVE
	}
}

func HandleFirstArg(input rune, state *StateInternal) {
	if unicode.IsDigit(input) {
		/*
			When a number is intered, we either shift store it in the buffer,
			or if we've exceded the digit length we just reset the state.
		*/
		state.count++
		if state.count <= 3 {
			ShiftStateBuffer(input, state)
		} else {
			ResetState(state)
		}
	} else if input == ',' {
		/*
			A comma input means that we've reached the end of the digits for the first arg.
			Whether we advance to the next state depends on the length of the first arg.
		*/
		if state.count > 0 && state.count <= 3 {
			numPosition := bufferLen - state.count
			state.first, _ = strconv.Atoi(string(state.buffer[numPosition:]))

			state.state = SECOND_ARG
			state.count = 0
		} else {
			ResetState(state)
		}
	} else {
		// All other inputs indicate an invalid operation.
		ResetState(state)
	}
}

func HandleSecondArg(input rune, state *StateInternal) int {
	if unicode.IsDigit(input) {
		/*
			When a number is intered, we either shift store it in the buffer,
			or if we've exceded the digit length we just reset the state.
		*/
		state.count++
		if state.count <= 3 {
			ShiftStateBuffer(input, state)
		} else {
			ResetState(state)
		}
	} else if input == ')' {
		/*
			A close parenthesis input means that we've reached the end of the digits for the second arg.
			If the second arg is valid (has the correct number of digits), we can return the operation result.
		*/
		if state.count > 0 && state.count <= 3 {
			numPosition := bufferLen - state.count
			state.second, _ = strconv.Atoi(string(state.buffer[numPosition:]))

			mult := state.first * state.second
			ResetState(state)

			return mult
		} else {
			ResetState(state)
		}
	} else {
		// All other inputs indicate an invalid operation.
		ResetState(state)
	}

	return 0
}

func HandleComplete(input rune, state *StateInternal) int {
	mult := state.first * state.second
	ResetState(state)

	return mult
}

func SwitchState(input rune, state *StateInternal) int {
	switch state.state {
	case INACTIVE:
		HandleInactive(input, state)
	case ACTIVE:
		HandleActive(input, state)
	case FIRST_ARG:
		HandleFirstArg(input, state)
	case SECOND_ARG:
		return HandleSecondArg(input, state)
	}

	return 0
}

func main() {
	var state StateInternal
	ResetState(&state)

	file, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	total := 0
	for {
		input, _, err := reader.ReadRune()
		input = rune(input)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			total += SwitchState(input, &state)
		}
	}

	fmt.Printf("Operation total %d\n", total)
}
