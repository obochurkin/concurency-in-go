package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var subArraysCount int = 4

func sortRoutine(sli []int, ch chan []int) {
	t := time.Now().UnixMicro()
	fmt.Println("Slice to sort: ", sli, t)
	sort.Ints(sli)
	ch <- sli
}

/*
user input serializer
*/
func serializeInput(input *string) []int {
  var sli []int
  inputStrings := strings.Split(*input, ",")
  for _, v := range inputStrings {
    number, err := strconv.Atoi(strings.TrimSpace(v))

		if err != nil {
			fmt.Println("Got an error", err)
      panic(err)
		}
    sli = append(sli, int(number))
  }

  return sli
}

func main() {
	ch := make(chan []int, 4) //initialize channel that can buffer 5 int arrays
	var sorted []int
	scanner := bufio.NewScanner(os.Stdin)

  fmt.Println("Pls provide comma separated string of integers. (ex. 5, 2, 4, 1, 4)")

  scanner.Scan()
	input := scanner.Text()
  intArray := serializeInput(&input)
	chunkSize := len(intArray)/subArraysCount

	for i := 0; i < subArraysCount; i++ {
		var chunk []int
		if (i+1) == subArraysCount {
			chunk = intArray[i*chunkSize:]
		} else {
			chunk = intArray[i*chunkSize : (i+1)*chunkSize]
		}
		
		go sortRoutine(chunk, ch)
		temp := <-ch
		sorted = append(sorted, temp...)
	}

	sort.Ints(sorted)
	t := time.Now().UnixMicro()
	fmt.Printf("sorted array: %d %d\n", sorted, t)
}
