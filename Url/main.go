package main

import (
	"fmt"
	"net/url"
)

const myUrl1 string = "https://jioface.in/api/migrate-s/ping"
const myUrl2 string = "https://lco.dev:3000/learn?coursename=reactjs&paymentid=65327843"

func main(){
	fmt.Println("Dealing with URL's")
	result , err := url.Parse(myUrl2)

	if err != nil {
		panic(err)
	}

	fmt.Println("result is ", result)
	fmt.Println(result.Scheme)
	fmt.Println(result.Host)
	fmt.Println(result.Path)
	fmt.Println(result.RawQuery)
	fmt.Println(result.Port())

	qParmas := result.Query()
	fmt.Printf("query Params Type is %T", qParmas)
	fmt.Println("query Params  is ", qParmas)
	fmt.Println("query Params  is ", qParmas["coursename"])


	// to construct a url
	partsOfUrl := &url.URL{
		Scheme: "https",
		Host: "lco.dev",
		Path: "/learn",
		RawQuery: "user=hitesh",
	}

	generatedUrl := partsOfUrl.String()
	fmt.Println("Generated Url is", generatedUrl)


}