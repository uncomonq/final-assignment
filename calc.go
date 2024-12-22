package main

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidExpression   = errors.New("invalid expression")
	ErrMismatchedParentheses = errors.New("mismatched parentheses")
	ErrDivisionByZero      = errors.New("division by zero")
)

func Calc(expression string) (float64, error) {
	if strings.TrimSpace(expression) == "" {
		return 0, ErrInvalidExpression
	}

	if !isValidExpression(expression) {
		return 0, ErrInvalidExpression
	}

	tokens := tokenize(expression)
	postfix, err := infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}
	return evaluatePostfix(postfix)
}

func isValidExpression(expr string) bool {
	for _, char := range expr {
		if !(char >= '0' && char <= '9') && !strings.ContainsRune("+-*/() ", char) {
			return false
		}
	}
	return true
}

func tokenize(expr string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, char := range expr {
		switch char {
		case ' ':
			continue
		case '+', '-', '*', '/', '(', ')':
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(char))
		default:
			currentToken.WriteRune(char)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	var oper []string

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			oper = append(oper, token)
		} else if token == ")" {
			for len(oper) > 0 && oper[len(oper)-1] != "(" {
				output = append(output, oper[len(oper)-1])
				oper = oper[:len(oper)-1]
			}
			if len(oper) == 0 {
				return nil, ErrMismatchedParentheses
			}
			oper = oper[:len(oper)-1]
		} else if isOperator(token) {
			for len(oper) > 0 && precedence(oper[len(oper)-1]) >= precedence(token) {
				output = append(output, oper[len(oper)-1])
				oper = oper[:len(oper)-1]
			}
			oper = append(oper, token)
		} else {
			return nil, ErrInvalidExpression
		}
	}

	for len(oper) > 0 {
		if oper[len(oper)-1] == "(" {
			return nil, ErrMismatchedParentheses
		}
		output = append(output, oper[len(oper)-1])
		oper = oper[:len(oper)-1]
	}

	return output, nil
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				stack = append(stack, a/b)
			default:
				return 0, ErrInvalidExpression
			}
		} else {
			return 0, ErrInvalidExpression
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}

func isNumber(token string) bool {
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return true
	}
	return false
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}
