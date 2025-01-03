package main

import (
	"errors"
	"fmt"
)

type DivisonError struct {
	IntA int
	IntB int
	Msg  string
}

func (e *DivisonError) Error() string {
	return e.Msg
}

func DivideCustom(a, b int) (int, error) {
	if b == 0 {
		return 0, &DivisonError{
			IntA: a,
			IntB: b,
			Msg:  fmt.Sprintf("Cannot Divide by '%d' by Zero", a),
		}
	}
	return a/b, nil
}

func CustomDivideCall(){
	a, b := 10, 0
	// a, b := 5, 6
	// a, b := 0, 0

	result, err := DivideCustom(a,b)
	if err != nil{
		fmt.Println(err)
		var divErr *DivisonError
		switch {
			case errors.As(err, &divErr): 
				fmt.Printf("%d /%d is Not Mathematically Valid: %s\n", divErr.IntA, divErr.IntB, divErr.Error())
			default : 
				fmt.Println("Unexpected Division Error %s\n", err)
		
		}
		return 
	}

	fmt.Printf("%d / %d = %d\n", a, b , result)
}
