package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main(){
	fmt.Println("Working with files in GO")

	content := "I am being written to a file"

	file , err := os.Create("./file.txt")
	if err !=nil {
		panic(err)
	}

	length , err := io.WriteString(file , content)
	if err != nil{
		panic(err)
	}

	fmt.Println("Length is %d", length)

	defer file.Close()
	readFile("./file.txt")
}

func readFile(fileName string){

	byteData , err := ioutil.ReadFile(fileName)
	if err != nil{
		panic(err)
	}

	fmt.Println("Byte data " , byteData)
	fmt.Println("Byte data " , string(byteData))
}