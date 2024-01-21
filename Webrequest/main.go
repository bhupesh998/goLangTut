package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const url="https://jioface.in/api/module-s/ping"

func main(){
	fmt.Println("web request")

	resp , err := http.Get(url)
	if err!=nil {
		panic(err)
	}
	fmt.Printf("response is of type %T\n", resp)
	
	dataBytes , err :=ioutil.ReadAll(resp.Body)
	if err!=nil {
		panic(err)
	}

	fmt.Println("Response is ", string(dataBytes))
	defer resp.Body.Close()

}