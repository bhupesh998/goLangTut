package main

import (
"errors"
"fmt"
)

var ErrDivideByZero = errors.New("Division By Zero")

func Divide(a, b int)(int, error){
	if b==0 {
		return 0, ErrDivideByZero
	}
	return a/b, nil
}


func ErrorOne(){
	// a, b := 10, 0
	 a, b := 5, 6
	// a, b := 0, 0
	fmt.Println(ErrDivideByZero)

	result, err := Divide(a,b)
	if err != nil{
		fmt.Println(err)
		switch {
			case errors.Is(err, ErrDivideByZero): 
				fmt.Println("Divide By Zero Error")
			default : 
				fmt.Println("Unexpected Division Error %s\n", err)
		
		}
		return 
	}

	fmt.Printf("%d / %d = %d\n", a, b , result)
}