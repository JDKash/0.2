package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	MaxNumber       = 10
	MaxOutputLength = 40
)

type OperationType int

const (
	Add OperationType = iota
	Subtract
	Multiply
	Divide
)

type Operation struct {
	LeftStr  string
	RightStr string
	Number   int
	OpType   OperationType
}

func parseInput(input string) (Operation, error) {
	trimmedInput := strings.TrimSpace(input)
	if trimmedInput == "" {
		panic("пустая строка ввода")
	}

	regex := regexp.MustCompile(`^"([^"]{1,10})"\s([+\-*/])\s(?:("([^"]{1,10})")|(\d+))$`)
	matches := regex.FindStringSubmatch(trimmedInput)
	if matches == nil {
		panic("неправильный формат ввода")
	}

	leftStr := matches[1]
	opType, err := determineOperationType(matches[2])
	if err != nil {
		panic(err.Error())
	}

	var rightStr string
	var number int
	if matches[3] != "" {
		rightStr = matches[4]
		if opType == Multiply || opType == Divide {
			panic("операция требует числового второго операнда")
		}
	} else {
		number, err = strconv.Atoi(matches[5])
		if err != nil || number < 1 || number > MaxNumber {
			panic("число вне допустимого диапазона")
		}
		if opType == Add || opType == Subtract {
			panic("операция требует строкового второго операнда")
		}
	}

	return Operation{
		LeftStr:  leftStr,
		RightStr: rightStr,
		Number:   number,
		OpType:   opType,
	}, nil
}

func determineOperationType(opSymbol string) (OperationType, error) {
	switch opSymbol {
	case "+":
		return Add, nil
	case "-":
		return Subtract, nil
	case "*":
		return Multiply, nil
	case "/":
		return Divide, nil
	default:
		return 0, fmt.Errorf("неподдерживаемая операция: %s", opSymbol)
	}
}

func calculate(operation Operation) (string, error) {
	switch operation.OpType {
	case Add:
		return handleStringOverflow(operation.LeftStr + operation.RightStr), nil
	case Subtract:
		return handleStringOverflow(strings.Replace(operation.LeftStr, operation.RightStr, "", -1)), nil
	case Multiply:
		return handleStringOverflow(strings.Repeat(operation.LeftStr, operation.Number)), nil
	case Divide:
		if operation.Number == 0 {
			panic("деление на ноль")
		}
		partLength := len(operation.LeftStr) / operation.Number
		if partLength == 0 {
			return "", nil
		}
		return handleStringOverflow(operation.LeftStr[:partLength]), nil
	default:
		panic("неизвестная операция")
	}
}

func handleStringOverflow(result string) string {
	if len(result) > MaxOutputLength {
		return result[:MaxOutputLength] + "..."
	}
	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("String Calculator")
	fmt.Println("---------------------")
	for {
		fmt.Print("-> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка при вводе:", err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}

		operation, err := parseInput(input)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			continue
		}

		result, err := calculate(operation)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			continue
		}

		fmt.Println("Результат:", result)
	}
}
