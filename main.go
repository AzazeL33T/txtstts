package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

var writer io.Writer = os.Stdout

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

var debugMode bool
var colorMode bool

func debug(format string, args ...interface{}) {
	if debugMode {
		fmt.Fprintf(os.Stderr, "DEBUG: "+format, args...)
		fmt.Fprintf(os.Stderr, "\n")
	}
}

func openFile(filename string) (*os.File, error) {
	debug("Opening file : %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		debug("File \"%s\" hasn`t open: %v", filename, err)
		return nil, err
	}
	debug("File \"%s\" has been opened", filename)
	return file, nil
}

func closeFile(file *os.File) error {
	debug("Closing file : %s", file.Name())
	err := file.Close()
	if err != nil {
		debug("File \"%s\" hasn`t close: %v", file.Name(), err)
		return err
	}
	debug("File \"%s\" been closed", file.Name())
	return nil
}

func printFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintln(writer, scanner.Text())
	}
	return scanner.Err()
}

func countCharsInFile(file *os.File) error {
	counter := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		counter += utf8.RuneCountInString(scanner.Text())
	}
	if colorMode {
		fmt.Fprintf(writer, Red+"txtstts: "+Reset+"Total characters in file: "+Green+"%d"+Reset+"\n", counter)
	} else {
		fmt.Fprintf(writer, "txtstts: Total characters in file: %d \n", counter)
	}
	return scanner.Err()
}

func countWordsInFile(file *os.File) error {
	counter := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		counter += len(strings.Fields(scanner.Text()))
	}
	if colorMode {
		fmt.Fprintf(writer, Red+"txtstts: "+Reset+"Total words in file: "+Green+"%d"+Reset+"\n", counter)
	} else {
		fmt.Fprintf(writer, "txtstts: Total words in file: %d \n", counter)
	}
	return scanner.Err()
}

func countUniqueWordsInFile(file *os.File) error {
	wordsMap := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, text := range strings.Fields(scanner.Text()) {
			text = strings.Trim(text, ".,!?;:()\"'-")
			wordsMap[strings.ToLower(text)] = struct{}{}
		}
	}
	if colorMode {
		fmt.Fprintf(writer, Red+"txtstts: "+Reset+"Total unique words in file: "+Green+"%d"+Reset+"\n", len(wordsMap))
	} else {
		fmt.Fprintf(writer, "txtstts: Total unique words in file: %d \n", len(wordsMap))
	}
	return scanner.Err()
}

func countLineInFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	counter := 0
	for scanner.Scan() {
		counter++
	}
	if colorMode {
		fmt.Fprintf(writer, Red+"txtstts: "+Reset+"Total lines in file: "+Green+"%d"+Reset+"\n", counter)
	} else {
		fmt.Fprintf(writer, "txtstts: Total lines in file: %d \n", counter)
	}
	return scanner.Err()
}

func countCommonWordsInFile(file *os.File, N int) error {
	type WordFreq struct {
		Word      string
		Frequency int
	}

	wordsMap := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, text := range strings.Fields(scanner.Text()) {
			wordsMap[strings.Trim(text, ".,!?;:()\"'-")]++
		}
	}

	wordFreqSlice := make([]WordFreq, 0, len(wordsMap))

	for word, freq := range wordsMap {
		wordFreqSlice = append(wordFreqSlice, WordFreq{word, freq})
	}

	sort.Slice(wordFreqSlice, func(i, j int) bool {
		return wordFreqSlice[i].Frequency < wordFreqSlice[j].Frequency
	})

	if colorMode {
		fmt.Fprintf(writer, Red+"txtstts: "+Reset+"The most common words in file: \n")
	} else {
		fmt.Fprintf(writer, "txtstts: The most common words in file \n")
	}

	for i := len(wordFreqSlice) - 1; i >= len(wordFreqSlice)-N && i >= 0; i-- {
		if colorMode {
			fmt.Fprintf(writer, "-%s: "+Green+"%d\n"+Reset, wordFreqSlice[i].Word, wordFreqSlice[i].Frequency)
		} else {
			fmt.Fprintf(writer, "-%s: %d\n", wordFreqSlice[i].Word, wordFreqSlice[i].Frequency)
		}
	}
	return scanner.Err()
}

func isPalindrome(word []rune) bool {
	if len(word) < 3 {
		return false
	}
	for i := 0; i < len(word); i++ {
		if word[i] != word[len(word)-i-1] {
			return false
		}
	}
	return true
}

func findPalindromes(file *os.File) error {
	palindromes := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, word := range strings.Fields(scanner.Text()) {
			if isPalindrome([]rune(strings.Trim(strings.ToLower(word), ".,!?;:()\"'-"))) == true {
				palindromes[word] = struct{}{}
			}
		}
	}

	if len(palindromes) <= 0 {
		if colorMode {
			fmt.Fprintf(writer, Red+"txtstts: "+Reset+"There is no palindromes")
		} else {
			fmt.Fprintf(writer, "txtstts: There is not palindromes")
		}
	}

	if len(palindromes) > 0 {
		if colorMode {
			fmt.Fprintf(writer, Red+"txtstts: "+Reset+"Found %d palindromes:\n", len(palindromes))
			for word := range palindromes {
				fmt.Fprintf(writer, Green+"- %s\n"+Reset, word)
			}
		} else {
			fmt.Fprintf(writer, "txtstts: Found %d palindromes:\n", len(palindromes))
			for word := range palindromes {
				fmt.Fprintf(writer, "- %s\n", word)
			}
		}
	}

	return scanner.Err()
}

func withFile(file *os.File, fn func(file *os.File) error) error {
	err := fn(file)
	if err != nil {
		return err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	return nil
}

func calculateAvgWordLen(file *os.File) error {
	lettersCounter, wordsCounter := 0, 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordsSlice := strings.Fields(scanner.Text())
		wordsCounter += len(wordsSlice)
		for _, text := range wordsSlice {
			lettersCounter += utf8.RuneCountInString(strings.Trim(text, ".,!?;:()\"'-"))
		}
	}

	if wordsCounter == 0 {
		if colorMode {
			fmt.Fprint(writer, Red+"txtstts: "+Reset+"There is no words in file. Average word lenght is "+Green+"0\n"+Reset)
		} else {
			fmt.Fprint(writer, "txtstts: There is no words in file. Average word length is 0\n")
		}
		return scanner.Err()
	}

	if colorMode {
		fmt.Fprintf(writer, Red+"txtstts: "+Reset+"Average word length: "+Green+"%.1f \n"+Reset, float32(lettersCounter)/float32(wordsCounter))
	} else {
		fmt.Fprintf(writer, "txtstts: Average word length: %.1f", float32(lettersCounter)/float32(wordsCounter))
	}
	return scanner.Err()
}

func collectAllData(file *os.File) error {
	type WordFreq struct {
		Word      string
		Frequency int
	}
	palindromes := make(map[string]struct{})
	wordsMap := make(map[string]int)
	charactersCounter, wordsCounter, linesCounter := 0, 0, 0
	lettersCounter := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		charactersCounter += utf8.RuneCountInString(scanner.Text())
		wordsCounter += len(strings.Fields(scanner.Text()))
		linesCounter++
		for _, text := range strings.Fields(scanner.Text()) {
			if isPalindrome([]rune(strings.Trim(strings.ToLower(text), ".,!?;:()\"'-"))) == true {
				palindromes[text] = struct{}{}
			}
			wordsMap[strings.Trim(text, ".,!?;:()\"'-)")]++
			lettersCounter += utf8.RuneCountInString(strings.Trim(text, ".,!?;:()\"'-)"))
		}
	}

	wordFreqSlice := make([]WordFreq, 0, len(wordsMap))

	for word, freq := range wordsMap {
		wordFreqSlice = append(wordFreqSlice, WordFreq{word, freq})
	}

	sort.Slice(wordFreqSlice, func(i, j int) bool {
		return wordFreqSlice[i].Frequency < wordFreqSlice[j].Frequency
	})

	wordsCounterFixed := func(wordsCounter int) float32 {
		if wordsCounter > 0 {
			return float32(wordsCounter)
		}
		return 1
	}

	avgWordLen := float32(lettersCounter) / wordsCounterFixed(wordsCounter)

	if colorMode {
		fmt.Fprintf(writer, Red+"txtstts: "+Reset+"Total characters in file: "+Green+"%d"+Reset+"\n"+
			Red+"txtstts: "+Reset+"Total words in file: "+Green+"%d"+Reset+"\n"+
			Red+"txtstts: "+Reset+"Total lines in file: "+Green+"%d"+Reset+"\n"+
			Red+"txtstts: "+Reset+"Average word length: "+Green+"%.1f"+Reset+"\n"+
			Red+"txtstts: "+Reset+"Total unique words in file: "+Green+"%d"+Reset+"\n", charactersCounter, wordsCounter, linesCounter, avgWordLen, len(wordsMap))
		fmt.Fprintf(writer, Red+"txtstts: "+Reset+"Top 5 most common words in file: \n")
		for i := len(wordFreqSlice) - 1; i >= len(wordFreqSlice)-5 && i >= 0; i-- {
			if colorMode {
				fmt.Fprintf(writer, "-%s: "+Green+"%d\n"+Reset, wordFreqSlice[i].Word, wordFreqSlice[i].Frequency)
			}
		}
		if len(palindromes) > 0 {
			fmt.Fprintf(writer, Red+"txtstts: "+Reset+"Found %d palindromes:\n", len(palindromes))
			for word := range palindromes {
				fmt.Fprintf(writer, Green+"- %s\n"+Reset, word)
			}
		} else {
			fmt.Fprintf(writer, Red+"txtstts: "+Reset+"There is no palindromes")
		}

	} else {
		fmt.Fprintf(writer, "txtstts: Total characters in file: %d \n"+
			"txtstts: Total words in file: %d \n"+
			"txtstts: Total lines in file: %d \n"+
			"txtstts: Total unique words in file: %d \n", charactersCounter, wordsCounter, linesCounter, len(wordsMap))
		for i := len(wordFreqSlice) - 1; i >= len(wordFreqSlice)-5 && i >= 0; i-- {
			fmt.Fprintf(writer, "-%s: %d\n", wordFreqSlice[i].Word, wordFreqSlice[i].Frequency)
		}
		if len(palindromes) > 0 {
			fmt.Fprintf(writer, "txtstts: Found %d palindromes:\n", len(palindromes))
			for word := range palindromes {
				fmt.Fprintf(writer, "- %s\n", word)
			}
		} else {
			fmt.Fprintf(writer, "txtstts: There is not palindromes")
		}
	}
	return scanner.Err()
}

func main() {
	printMode := flag.Bool("print", false, "Print file")
	countCharactersMode := flag.Bool("chars", false, "Count characters in file")
	countWordsMode := flag.Bool("words", false, "Count words in file")
	countLinesMode := flag.Bool("lines", false, "Count lines in file")
	countUniqueWords := flag.Bool("unique", false, "Count unique words in file")
	commonWordsMode := flag.Int("common", 0, "Count most common words in file (Top N)")
	averageWordLenght := flag.Bool("avg-len", false, "Calculates the average word lenght in file")
	displayPalindromes := flag.Bool("palindrome", false, "Shows palindromes is file")
	outputPath := flag.String("o", "", "Output file")
	displayAll := flag.Bool("all", false, "Uses all counter features (lines, characters, words)")
	debugFlag := flag.Bool("debug", false, "Turns on debug mode")
	colorFlag := flag.Bool("color", false, "Turns on the color mode (can work incorrectly on some terminals)")
	flag.Parse()
	debugMode = *debugFlag
	colorMode = *colorFlag

	if *outputPath != "" {
		file, err := os.Create(*outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output: %v \n", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = file
	}

	if flag.NArg() < 1 {
		debug("No filename provided")
		fmt.Fprintf(os.Stderr, "Error: file name required \n")
		os.Exit(1)
	}

	files := flag.Args()

	if *displayAll && (*countWordsMode || *countCharactersMode || *countLinesMode || *countUniqueWords || *commonWordsMode > 0 || *averageWordLenght || *displayPalindromes) {
		fmt.Fprint(os.Stderr, Red+"Warning: "+Reset+"Individual stats flags (-chars, -words, -lines, -unique) are ignored when -all is used!\n")
	}

	if *commonWordsMode < 0 {
		fmt.Fprintf(os.Stderr, "Non-negative and non-zero interger required to show most common words\n")
	}

	for _, filename := range files {
		file, err := openFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v \n", err)
			os.Exit(1)
		}

		if *colorFlag {
			fmt.Printf("\n"+Green+"========== %s ==========\n\n"+Reset, filename)
		} else {
			fmt.Printf("\n========== %s ==========\n\n", filename)
		}

		if *printMode {
			if err := withFile(file, printFile); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v \n", err)
			}
		}

		if !*displayAll && (*countWordsMode || *countCharactersMode || *countLinesMode || *countUniqueWords || *commonWordsMode > 0 || *averageWordLenght || *displayPalindromes) {
			if *countWordsMode {
				if err := withFile(file, countWordsInFile); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v \n", err)
				}
			}
			if *countLinesMode {
				if err := withFile(file, countLineInFile); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v \n", err)
				}
			}
			if *countCharactersMode {
				if err := withFile(file, countCharsInFile); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v \n", err)
				}
			}
			if *countUniqueWords {
				if err := withFile(file, countUniqueWordsInFile); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v \n", err)
				}
			}

			if *commonWordsMode > 0 {
				n := *commonWordsMode
				if err := withFile(file, func(file *os.File) error {
					return countCommonWordsInFile(file, n)
				}); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v \n", err)
				}
			}

			if *averageWordLenght {
				if err := withFile(file, calculateAvgWordLen); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v \n", err)
				}
			}

			if *displayPalindromes {
				if err := withFile(file, findPalindromes); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v \n", err)
				}
			}
		}
		if *displayAll {
			if err := withFile(file, collectAllData); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v \n", err)
			}
		}
		if err := closeFile(file); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v \n", err)
		}
	}
}
