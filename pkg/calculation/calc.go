package calc

import (
	"fmt"
	"regexp"
	"strconv"
)

func tokenize(input string) []string {
	re := regexp.MustCompile(`\d+|\+|\-|\*|\/|\(|\)`)
	tokens := re.FindAllString(input, -1)
	return tokens
}

func findClosingBracket(tokens []string, start int) int {
	count := 1
	for i := start + 1; i < len(tokens); i++ {
		if tokens[i] == "(" {
			count++
		} else if tokens[i] == ")" {
			count--
			if count == 0 {
				return i
			}
		}
	}
	return -1
}

func evaluateExpression(tokens []string) ([]string, error) {
	processedTokens := []string{}
	i := 0
	for i < len(tokens) {
		if tokens[i] == "(" {
			end := findClosingBracket(tokens, i)
			if end == -1 {
				return []string{}, fmt.Errorf("no closing bracket found")
			}

			innerExpression := tokens[i+1 : end]
			evaluated, err := evaluateExpression(innerExpression)
			if err != nil {
				return []string{}, fmt.Errorf("error: %v", err)
			}
			innerResult, err := evaluateHighPriority(evaluated)
			if err != nil {
				return []string{}, fmt.Errorf("error: %v", err)
			}
			result, err := evaluateLowPriority(innerResult)
			if err != nil {
				return []string{}, fmt.Errorf("error: %v", err)
			}

			processedTokens = append(processedTokens, fmt.Sprintf("%f", result))
			i = end + 1
		} else {
			processedTokens = append(processedTokens, tokens[i])
			i++
		}
	}
	return processedTokens, nil
}

func evaluateHighPriority(tokens []string) ([]string, error) {
	i := 0
	processedTokens := []string{}

	for i < len(tokens) {
		if tokens[i] == "*" || tokens[i] == "/" {
			if len(processedTokens) == 0 {
				return []string{}, fmt.Errorf("error: missing left operand for operator %s", tokens[i])
			}
			if i+1 >= len(tokens) {
				return []string{}, fmt.Errorf("error: missing right operand for operator %s", tokens[i])
			}

			left, err1 := strconv.ParseFloat(processedTokens[len(processedTokens)-1], 64)
			right, err2 := strconv.ParseFloat(tokens[i+1], 64)

			if err1 != nil || err2 != nil {
				return []string{}, fmt.Errorf("error while parsing number: %v %v", err1, err2)
			}

			var result float64
			if tokens[i] == "*" {
				result = left * right
			} else {
				result = left / right
			}

			processedTokens[len(processedTokens)-1] = fmt.Sprintf("%f", result)
			i += 2
		} else {
			processedTokens = append(processedTokens, tokens[i])
			i++
		}
	}

	return processedTokens, nil
}

func evaluateLowPriority(tokens []string) (float64, error) {
	if len(tokens) == 0 {
		return -1, fmt.Errorf("error: no tokens to evaluate")
	}

	result, err := strconv.ParseFloat(tokens[0], 64)
	if err != nil {
		return -1, fmt.Errorf("error while parsing number: %v", err)
	}

	i := 1
	for i < len(tokens) {
		operator := tokens[i]

		if i+1 >= len(tokens) {
			return -1, fmt.Errorf("error: missing right operand for operator %s", operator)
		}

		right, err := strconv.ParseFloat(tokens[i+1], 64)
		if err != nil {
			return -1, fmt.Errorf("error while parsing number: %v", err)
		}

		if operator == "+" {
			result += right
		} else if operator == "-" {
			result -= right
		}
		i += 2
	}

	return result, nil
}

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	tokens, err := evaluateExpression(tokens)
	if err != nil {
		return 0, fmt.Errorf("error: %v", err)
	}
	tokens, err = evaluateHighPriority(tokens)
	if err != nil {
		return 0, fmt.Errorf("error: %v", err)
	}
	result, err := evaluateLowPriority(tokens)
	if err != nil {
		return 0, fmt.Errorf("error: %v", err)
	}

	return result, nil
}

func main() {
	input := "78*(2485+48*(45+7))"
	fmt.Println(Calc(input))
}
