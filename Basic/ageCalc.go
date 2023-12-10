package main

import (
	"fmt"
	"bufio"
	"os"
	// "strconv"
)

func main1()  {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	fmt.Printf("you typed : %q", input)
}