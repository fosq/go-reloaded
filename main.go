package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Main function
func main() {
	arguments := os.Args[1:]
	var oldString string
	var newString []string

	// Check if arguments more than 0
	if len(arguments) == 0 {
		bytes, _ := ioutil.ReadAll(os.Stdin)
		fmt.Println(bytes)
	} else if len(arguments) > 0 {
		// Read first argument and open it's file
		file, err := os.Open(arguments[0])

		// Goes through all the functions to correct the text, and write to file
		if err == nil {
			fileBytes, _ := file.Stat()
			arr := make([]byte, fileBytes.Size())
			file.Read(arr)
			oldString = string(arr)
			newString = stringManipulation(oldString)
			newString = fixPunctuation(newString)
			justString := strings.Join(newString, " ")
			justString = strings.TrimRight(justString, " ")
			file.Close()
			output, err2 := os.Create(arguments[1])
			if err2 != nil {
				fmt.Println(err2.Error())
				os.Exit(1)
			}
			defer output.Close()
			_, err3 := output.Write([]byte(justString))

			if err3 != nil {
				fmt.Println(err3.Error())
				os.Exit(1)
			}

		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

// Function for converting a to an, (hex), (bin), (up), (low), (cap) and the latter 3 number variants
func stringManipulation(argument string) []string {
	var strArray []string
	strArray = strings.Fields(argument)
	// Switch cases
	for i := range strArray {
		switch strArray[i] {
		case "(cap)":
			strArray[i-1] = strings.Title(strArray[i-1])
			remove(strArray, i)
		case "(up)":
			strArray[i-1] = strings.ToUpper(strArray[i-1])
			remove(strArray, i)
		case "(low)":
			strArray[i-1] = strings.ToLower(strArray[i-1])
			remove(strArray, i)
		case "(cap,":
			// Remove ")" from 6)
			nextWord := strings.TrimSuffix(strArray[i+1], ")")
			// Convert string to int
			nextWordCount, _ := strconv.Atoi(nextWord)
			// Uppercases the number of previous words from current index
			for j := 0; j <= nextWordCount; j++ {
				strArray[i-j] = strings.Title(strArray[i-j])
			}
			// Remove (cap, 6) from array of word, and if it is, then corrects its
			remove(strArray, i)
			remove(strArray, i)

		case "(low,":
			// Remove ")" from 6)
			nextWord := strings.TrimSuffix(strArray[i+1], ")")
			// Convert string to int
			nextWordCount, _ := strconv.Atoi(nextWord)
			// Lowercases the number of previous words from current index
			for j := 0; j <= nextWordCount; j++ {
				strArray[i-j] = strings.ToLower(strArray[i-j])
			}
			remove(strArray, i)
			remove(strArray, i)

		// Replace word with decimanl version
		case "(hex)":
			num, err3 := strconv.ParseInt(strArray[i-1], 16, 64)

			if err3 != nil {
				fmt.Println(err3.Error())
				os.Exit(1)
			}
			strArray[i-1] = fmt.Sprint(num)
			remove(strArray, i)

		// Replace word with decimal version
		case "(bin)":
			num, err3 := strconv.ParseInt(strArray[i-1], 2, 64)

			if err3 != nil {
				fmt.Println(err3.Error())
				os.Exit(1)
			}
			strArray[i-1] = fmt.Sprint(num)
			remove(strArray, i)

		// Change a to an when next word starts with a vowel
		case "a":
			nextWord := strArray[i+1]
			nextLetter := nextWord[:1]
			switch nextLetter {
			case "a", "e", "i", "o", "u", "A", "E", "I", "O", "U":
				strArray[i] = "an"
			}
		}
	}
	return strArray
}

// When finding lonely punctuations, or punctuations at the start of a word, corrects them to the appropriate position
func fixPunctuation(strArray []string) []string {
	var firstFound bool
	for i := range strArray {
		// Check if punctuation in word, and check if NOT any letters or numbers in word
		if strings.ContainsAny(strArray[i], ".,!?:;") && !strings.ContainsAny(strArray[i], "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789") {

			// If lonely punctuations are found, put them at the end of previous word
			strArray[i-1] = strArray[i-1] + strArray[i]
			remove(strArray, i)

			// Should've used switch cases, oopsie
			// Check if comma at start of word, and if it is, then corrects it
		} else if strings.HasPrefix(strArray[i], ",") {
			strArray[i] = strArray[i][1:]
			strArray[i-1] = strArray[i-1] + ","

			// Check if dot at start of word, and if it is, then corrects it
		} else if strings.HasPrefix(strArray[i], ".") {
			strArray[i] = strArray[i][1:]
			strArray[i-1] = strArray[i-1] + "."

			// Check if exclamation mark at start of word, and if it is, then corrects it
		} else if strings.HasPrefix(strArray[i], "!") {
			strArray[i] = strArray[i][1:]
			strArray[i-1] = strArray[i-1] + "!"

			// Check if question mark at start of word, and if it is, then corrects it
		} else if strings.HasPrefix(strArray[i], "?") {
			strArray[i] = strArray[i][1:]
			strArray[i-1] = strArray[i-1] + "?"

			// Check if colon at start of word, and if it is, then corrects it
		} else if strings.HasPrefix(strArray[i], ":") {
			strArray[i] = strArray[i][1:]
			strArray[i-1] = strArray[i-1] + ":"

			// Check if semicolon at start of word, and if it is, then corrects it
		} else if strings.HasPrefix(strArray[i], ";") {
			strArray[i] = strArray[i][1:]
			strArray[i-1] = strArray[i-1] + ";"
		}

		// Fix `` punctuations
		if strArray[i] == "'" {
			// Check for first ` and put it at the start of next word
			if firstFound == false {
				strArray[i+1] = strArray[i] + strArray[i+1]
				remove(strArray, i)
				firstFound = true

				// Check for second ` and put it at the end of previous word
			} else if firstFound == true {
				strArray[i-1] = strArray[i-1] + strArray[i]
				remove(strArray, i)
				firstFound = false
			}
		}
	}
	// If all checks done, return string array
	return strArray
}

// Removes arguments from array by index
func remove(array []string, i int) []string {
	copy(array[i:], array[i+1:])
	array[len(array)-1] = ""
	array = array[:len(array)-1]

	return array
}
