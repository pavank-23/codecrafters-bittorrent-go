package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

func decodeString(bencodedString string, index int) (interface{}, error, int) {
	var firstColonIndex int
	for i := index; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[index : firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err, -1
	}

	return bencodedString[firstColonIndex + 1 : firstColonIndex + 1 + length], nil, firstColonIndex + length
}

func decodeInteger(bencodedString string, index int) (interface{}, error, int) {
	var endOfInteger int
	for i := index; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			endOfInteger = i
			break
		}
	}
	integer, err := strconv.Atoi(bencodedString[index + 1 : endOfInteger])
	if err != nil {
		return "", err, -1
	}
	return integer, nil, endOfInteger
}

func decodeList(bencodedString string, index int) (interface{}, error, int) {
	var decodedList []interface{}
	var i int = index + 1
	for i < len(bencodedString) && bencodedString[i] != 'e' {
		if bencodedString[i] == 'i' {
			result, err, index := decodeInteger(bencodedString, i)
			i = index
			if err != nil {
				return "", err, -1
			}
			decodedList = append(decodedList, result)
		} else if unicode.IsDigit(rune(bencodedString[i])) {
			result, err, index := decodeString(bencodedString, i)
			i = index
			if err != nil {
				return "", err, -1
			}
			decodedList = append(decodedList, result)
		} else if bencodedString[i] == 'l' {
			list, err, index := decodeList(bencodedString, i)
			if err != nil {
				return "", err, -1
			}
			decodedList = append(decodedList, list)
			i = index
		}
		i += 1
	}
	if len(decodedList) == 0 {
		return []interface{}{}, nil, i
	}
	return decodedList, nil, i
}

func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		result, err, _ := decodeString(bencodedString, 0)
		if err != nil {
			return "", err
		}
		return result, err
	} else if bencodedString[0] == 'i' {
		result, err, _ := decodeInteger(bencodedString, 0)
		if err != nil {
			return "", err
		}
		return result, err
	} else if bencodedString[0] == 'l' {
		result, err, _ := decodeList(bencodedString, 0)
		if err != nil {
			return "", err
		}
		return result, err
	} else {
		return "", fmt.Errorf("only strings, integers and lists are supported at the moment")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}