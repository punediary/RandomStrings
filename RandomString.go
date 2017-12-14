package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// Fill the genrated Random numbers in randomNumbers
// The generated Random numbers should be in the range of (0, oLen)
// The argument oLen here is the total number of lines in 'words_alpha'
func getRandomNumbers(nRandom, oLen int) []int {

	randomNumbers := make([]int, 0)

	for i := 0; i < nRandom; i++ {
		r := random(0, oLen)
		randomNumbers = append(randomNumbers, r)
	}

	return randomNumbers
}

func prepareFileoffsets() ([]int, error) {

	offset := 0
	offsetSlice := make([]int, 0)
	str := make([]byte, 0)

	file, err := os.Open("words_alpha")
	defer file.Close()

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	// Read Entire file line by line and Store the offsets of each line.
	// The Size of file is Big and we are not using all file contents everytime.
	//
	for err == nil {
		str, err = reader.ReadBytes('\n')
		offsetSlice = append(offsetSlice, offset)
		offset += len(str)
	}

	// In case not able to read entire File then through Error.
	//
	if err != io.EOF {
		return []int{}, errors.New("Could Not read Entire File")
	}

	return offsetSlice, nil
}

func printRandomStrings(randomNumbers, offsetSlice []int) {

	file, err := os.Open("words_alpha")

	defer file.Close()

	if err != nil {
		panic(err)
	}

	for _, v := range randomNumbers {

		file.Seek(int64(offsetSlice[v]), 0)

		reader := bufio.NewReader(file)

		str, err := reader.ReadBytes('\n')
		if err == nil {
			fmt.Print(string(str))
		}
	}
}

func main() {

	// Store the file offset of each line in 'offsetSlice'
	offsetSlice, err := prepareFileoffsets()

	if err != nil {
		panic(err)
	}

	nRandom := 1
	if len(os.Args) == 2 {
		n, err := strconv.Atoi(os.Args[1])
		if err == nil {
			nRandom = n
		}
	}

	// Using default random number generator in Golang "math/rand"
	// There are much better Psudo Random number generation algorithums like 'Mersenne Twister' (PRNG) algorithums.
	// The value of generated Random number has to be between 0 and (total number of lines in 'words_alpha')
	randomNumbers := getRandomNumbers(nRandom, len(offsetSlice))

	// We pick up a random number and seek to perticular offset in file.
	// Its Good if the random numbers are arranged in assending order so that Seek on file is Linear (start to end )and not Random.
	if len(randomNumbers) > 1 {
		// Golang internally usess quick sort, I think its good enough for present requirement.
		sort.Ints(randomNumbers)
	}

	// The generated Random numbers are in the range of 0, total lines of file
	// This API will get each random number goes to that perticular line and finds its offset using 'offsetSlice' and print that string
	printRandomStrings(randomNumbers, offsetSlice)
}
