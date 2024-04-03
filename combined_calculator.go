package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Определение констант
const (
	MaxStringLength = 10
	MaxNumber       = 10
)

// OperationType Определение типа операции
type OperationType int

const (
	Add      OperationType = iota // Сложение
	Subtract                      // Вычитание
	Multiply                      // Умножение
	Divide                        // Деление
)

// Operation Определение структуры операции
type Operation struct {
	LeftStr  string        // Левый операнд (строка)
	RightStr string        // Правый операнд (строка, может быть пустым)
	Number   int           // Операнд (число, может быть 0)
	OpType   OperationType // Тип операции
}

// Функция для валидации и обработки строки, заключенной в кавычки
func validateAndProcessString(str string) (string, error) {
	// Проверка, что строка заключена в двойные кавычки
	if !strings.HasPrefix(str, "\"") || !strings.HasSuffix(str, "\"") {
		return "", errors.New("строка должна быть заключена в двойные кавычки")
	}

	// Удаление кавычек
	processedStr := str[1 : len(str)-1]

	// Проверка длины строки
	if len(processedStr) > MaxStringLength {
		return "", errors.New("длина строки превышает максимально допустимую")
	}

	return processedStr, nil
}

// Функция для определения типа операции
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
		return 0, errors.New("неподдерживаемая операция")
	}
}

// Функция для обработки второго операнда в зависимости от типа операции
func processSecondOperand(operand string, opType OperationType) (string, int, error) {
	switch opType {
	case Add, Subtract:
		// Для сложения и вычитания второй операнд должен быть строкой
		processedStr, err := validateAndProcessString(operand)
		if err != nil {
			return "", 0, err
		}
		return processedStr, 0, nil
	case Multiply, Divide:
		// Для умножения и деления второй операнд должен быть числом
		number, err := strconv.Atoi(operand)
		if err != nil {
			return "", 0, errors.New("не удалось преобразовать второй операнд в число")
		}
		// Проверка диапазона числа
		if number < 1 || number > MaxNumber {
			return "", 0, errors.New("число вне допустимого диапазона")
		}
		return "", number, nil
	default:
		// Для полноты
		return "", 0, errors.New("неизвестный тип операции")
	}
}

/* Парсим ввод
func parseInput(input string) (Operation, error) {
	trimmedInput := strings.TrimSpace(input)
	if trimmedInput == "" {
		return Operation{}, errors.New("пустая строка ввода")
	}

	parts := strings.Fields(trimmedInput)
	if len(parts) != 3 {
		return Operation{}, errors.New("ввод должен содержать ровно три элемента")
	}

	leftStr, err := validateAndProcessString(parts[0])
	if err != nil {
		return Operation{}, err
	}

	opType, err := determineOperationType(parts[1])
	if err != nil {
		return Operation{}, err
	}

	rightStr, number, err := processSecondOperand(parts[2], opType)
	if err != nil {
		return Operation{}, err
	}

	return Operation{
		LeftStr:  leftStr,
		RightStr: rightStr,
		Number:   number,
		OpType:   opType,
	}, nil
}
*/

// с учетом пробелов
func parseInput(input string) (Operation, error) {
	trimmedInput := strings.TrimSpace(input)
	if trimmedInput == "" {
		return Operation{}, errors.New("пустая строка ввода")
	}

	regex := regexp.MustCompile(`^"([^"]+)"\s([+\-*/])\s(.+)$`)
	matches := regex.FindStringSubmatch(trimmedInput)

	if matches == nil || len(matches) != 4 {
		return Operation{}, errors.New("ввод должен содержать ровно три элемента")
	}

	leftStr, err := validateAndProcessString("\"" + matches[1] + "\"")
	if err != nil {
		return Operation{}, err
	}

	opType, err := determineOperationType(matches[2])
	if err != nil {
		return Operation{}, err
	}

	rightStr, number, err := processSecondOperand(matches[3], opType)
	if err != nil {
		return Operation{}, err
	}

	return Operation{
		LeftStr:  leftStr,
		RightStr: rightStr,
		Number:   number,
		OpType:   opType,
	}, nil
}

// Функция для выполнения вычислений
func calculate(operation Operation) (string, error) {
	switch operation.OpType {
	case Add:
		// Сложение строк
		return handleStringOverflow(operation.LeftStr + operation.RightStr), nil
	case Subtract:
		// Вычитание строки из строки
		return handleStringOverflow(strings.Replace(operation.LeftStr, operation.RightStr, "", -1)), nil
	case Multiply:
		// Умножение строки на число
		return handleStringOverflow(strings.Repeat(operation.LeftStr, operation.Number)), nil
	case Divide:
		// Деление строки на число
		if operation.Number == 0 {
			return "", errors.New("деление на ноль")
		}
		partLength := len(operation.LeftStr) / operation.Number
		if partLength == 0 {
			return "", nil
		}
		return handleStringOverflow(operation.LeftStr[:partLength]), nil
	default:
		return "", errors.New("неизвестная операция")
	}
}

// Функция для обработки переполнения строки
func handleStringOverflow(result string) string {
	const maxOutputLength = 40
	if len(result) > maxOutputLength {
		return result[:maxOutputLength] + "..."
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
