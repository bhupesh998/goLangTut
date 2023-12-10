package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

// in struct the fields are in caps so that they can be public i.e they can be viewed outside of the modules
type book struct{
	ID string 		`json: "id"`  // when serizing the fields with json convert them to lowercase
	Title string    `json: "title"`
	Author string 	`json: "author"`
	Price int	`json: "price"`
	Quantity int `json: "quantity"`
}

var books = []book{
	{ID: "100", Title: "Batman" , Author: "Batman", Price: 200 , Quantity: 5},
	{ID: "200", Title: "He man" , Author: "He man", Price: 300, Quantity: 7},
	{ID: "300", Title: "Superman" , Author: "Superman", Price: 400, Quantity: 9},
}

func getBooks (c *gin.Context){
c.IndentedJSON(http.StatusOK, books)
}

func bookById (c *gin.Context){
	id := c.Param("id") 
	book , err := getBookById(id)

	if(err !=nil){
		c.IndentedJSON(http.StatusNotFound , gin.H{"message": "Book Not Found " + id  } )
		return
	}
	c.IndentedJSON(http.StatusOK, book)

	}

func getBookById(Id string) (*book , error){
	for i, b := range books{
		if b.ID == Id {
			return &books[i] , nil
		}
	}
	return nil, errors.New("Book not found for %s" + Id)
}

func checkOutBook(c *gin.Context){
	id , ok := c.GetQuery("id")

	if ok == false  {
		c.IndentedJSON(http.StatusBadRequest , gin.H{"message": "Missing Arguments"  } )
		return 
	}
	book, err := getBookById(id)

	if(err !=nil){
		c.IndentedJSON(http.StatusNotFound , gin.H{"message": "Book Not Found " + id  } )
		return
	}
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusNotFound , gin.H{"message": "Book Not Available " + book.Title } )
		return 
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK , book)
}

func returnBook(c *gin.Context){
	id , ok := c.GetQuery("id")

	if ok == false  {
		c.IndentedJSON(http.StatusBadRequest , gin.H{"message": "Missing Arguments"  } )
		return 
	}
	book, err := getBookById(id)

	if(err !=nil){
		c.IndentedJSON(http.StatusNotFound , gin.H{"message": "Book Not Found " + id  } )
		return
	}
	
	book.Quantity += 1
	c.IndentedJSON(http.StatusOK , book)

}

func createBook(c *gin.Context){
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return // returning will not return a response , its the bind json method that will return the response
	}

	books= append(books, newBook)
	c.IndentedJSON(http.StatusOK , newBook)
}

func main(){
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById )
	router.POST("/books",createBook)
	router.PATCH("/checkOut", checkOutBook)
	router.PATCH("/returnBook", returnBook)
	router.Run("localhost:8080")
}





