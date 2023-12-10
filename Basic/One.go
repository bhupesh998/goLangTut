package main

import "fmt"

func main(){
	var name string = "hello World"
	var num uint16= 291 // explicit declaration
	num=num+3
	fmt.Println(name, num)

	var num2 = 334
	fmt.Println(num2)
	fmt.Printf("%T\n",num2)

	flo:= 5566.00
	fmt.Println(flo)
	fmt.Printf("%T\n",flo)

	// flo = "hello"- cannot use "hello" (untyped string constant) as float64 value in assignment

	var def uint64
	var bl bool
	var str string
	fmt.Println(def, bl, str)

	fmt.Printf("Hello %T %v\n", 10 , 10)
	var x string = fmt.Sprintf("Hello %T %v", 10 , 10)
	fmt.Println(x)


	fmt.Printf("Number in HexaDecimal system   -  %x\n",  9999 )
	fmt.Printf("Number in HexaDecimal system   -  %X\n",  9999 ) // to print letter in CAPS 

	fmt.Printf("Number :  %-13q is cool\n", "11")
	fmt.Printf("Number :  %13q is cool", "11")
}