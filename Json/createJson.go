package main

import (
	"encoding/json"
	"fmt"
)

type course struct{
	Name  string 	`json:"coursename"`
	Price int
	Platform  string
	Password  string `json:"-"`
	Tags  []string	`json:"tags,omitempty"`


}

func main() {
	fmt.Println("Created Json function ")
	//EncodeJson()
	DecodeJson()
}


func EncodeJson(){
	myCourses :=[]course{
		{ "React BootCamp" , 1000 , "Udemy", "React",[]string{ "Frontend Development"} },
		{ "Angular BootCamp" , 1000 , "Linkedin", "Angul", nil },
		{ "Nodejs BootCamp" , 1000 , "Coursera", "Node",[]string{ "Backend Development"} },
	}

	finalJson , err := json.MarshalIndent(myCourses, "", " ")

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", finalJson)
	
}

func DecodeJson(){
	jsonWeb := []byte(`
	{
  "coursename": "React BootCamp",
  "Price": 1000,
  "Platform": "Udemy",
  "tags": ["Frontend Development"]
 }
`)

var lcoCourse  course

checkValid := json.Valid(jsonWeb)

if(checkValid){
	 json.Unmarshal(jsonWeb, &lcoCourse)
	 fmt.Printf("%#v\n", lcoCourse)
}else{
	fmt.Println("JSON IS INVALID")
}

var myData map[string]interface{}
json.Unmarshal(jsonWeb, &myData)

fmt.Printf("%#v\n", myData)

for key , value := range myData{
	fmt.Printf("key is %v and value is %v\n", key , value)
}

}