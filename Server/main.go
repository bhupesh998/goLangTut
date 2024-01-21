package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main(){

	fmt.Println("Server in goLang")
	//performGet()
	//performPost()
	performFormRequest()
}


func performGet(){
	const myUrl="http://jioface.in/api/module-s/ping"

	resp, err := http.Get(myUrl)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("Status Code :", resp.StatusCode)
	fmt.Println("Status :", resp.Status)
	fmt.Println("Content Length :", resp.ContentLength)

	 content , _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("Response body is", string(content))

	var respBuilder strings.Builder
	byteCount , _ := respBuilder.Write(content)

	fmt.Println("Byte Count" , byteCount)
	fmt.Println("Response Through builder is " , respBuilder.String())
	
}


func performPost(){
	const postUrl="https://dummyjson.com/posts/add"

	requestBody := strings.NewReader(`
	{
		"title": "I am in love with someone.",
		"userId": "7",
	}
	`)

	resp , err := http.Post(postUrl , "application/json" , requestBody)
	  if err != nil {
		panic(err)
	  }

	  defer resp.Body.Close()	  

	  
	fmt.Println("Status Code :", resp.StatusCode)
	fmt.Println("Status :", resp.Status)
	fmt.Println("Content Length :", resp.ContentLength)
	
	content , _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response body is", string(content))
}

func performFormRequest(){
	const postUrl="https://dummyjson.com/posts/add"

	data := url.Values{}
	data.Add("firstname", "bhupesh")
	data.Add("lastname", "choudhary")
	data.Add("userId", "5")

	resp , err := http.PostForm(postUrl , data)
	  if err != nil {
		panic(err)
	  }

	  defer resp.Body.Close()	  

	  
	fmt.Println("Status Code :", resp.StatusCode)
	fmt.Println("Status :", resp.Status)
	fmt.Println("Content Length :", resp.ContentLength)
	
	content , _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response body is", string(content))

}
