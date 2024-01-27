package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version ="1.0.0"

type (
	Logger interface{
		Fatal(string , ...interface{})
		Error(string , ...interface{})
		Warn(string , ...interface{})
		Info(string , ...interface{})
		Debug(string , ...interface{})
		Trace(string , ...interface{})
	}

	Driver struct{
		mutex sync.Mutex
		mutexes map[string]*sync.Mutex
		dir string
		log Logger
	}
)

type Options struct{
	Logger
}

func New(dir string, option *Options)(*Driver, error){

	dir = filepath.Clean(dir)

	opts := Options{}
	if option != nil {
		opts = *option
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir : dir,
		mutexes: make(map[string]*sync.Mutex),
		log : opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (Database Already Exist)\n ", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Createing the Database at '%s'...\n", dir)

	return &driver, os.MkdirAll(dir, 0755)

}



func (d *Driver) Write(collection , resource string, v interface{}) error {

	if collection == ""{
		return fmt.Errorf("Missing Collection Details")
	}

	if resource == ""{
		return fmt.Errorf("Missing Resource -Unable to save record")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	finalPath := filepath.Join(dir, resource+".json")
	tempPath := finalPath+".tmp"

	if err:= os.MkdirAll(dir , 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b=append(b, byte('\n'))

	if err := ioutil.WriteFile(tempPath, b , 0644); err != nil{
		return err
	}

	return os.Rename(tempPath, finalPath)
}

func (d *Driver)  Read(collection , resource string, v interface{}) error{
	if collection == ""{
		return fmt.Errorf("Missing Collection Details")
	}

	if resource == ""{
		return fmt.Errorf("Missing Resource -Unable to save record")
	}

	record := filepath.Join(d.dir, collection, resource)
	if _, err := stat(record); err != nil {
		return err
	}

	b , err := ioutil.ReadFile(record + ".json")
	if  err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}


func (d *Driver)  ReadAll(collection string) ([]string , error){
	if collection == ""{
		return nil , fmt.Errorf("Missing Collection Details")
	}

	

	dir := filepath.Join(d.dir, collection)
	if _, err := stat(dir); err != nil {
		return nil , err
	}

	files , err := ioutil.ReadDir(dir)
	if  err != nil {
		return nil, err
	}

	var records []string

	for _, file := range files{
	b , err :=ioutil.ReadFile(filepath.Join(dir, file.Name()) )
	if  err != nil {
		return nil, err
	}

	records = append(records, string(b))
	}

	return records, nil
}


func (d *Driver) Delete(collection , resource string) error{
	if collection == ""{
		return   fmt.Errorf("Missing Collection Details")
	}

	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)
	switch fi , err := stat(dir); {
	case fi == nil , err != nil :
		return fmt.Errorf("Unable to find the file or directory %v\n", path)

	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	case fi.Mode().IsRegular(): 
		return os.RemoveAll(dir + ".json")

	
	}

	return nil

}

func (d *Driver)  getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()

m , ok := 	d.mutexes[collection]

	if !ok{
		m= &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string)(fi os.FileInfo, err error){
	if fi, err = os.Stat(path); os.IsNotExist(err){
		fi, err = os.Stat(path + ".json")
	}
	return 
}



type User struct{
	Name string
	Age json.Number
	Contact string
	Company string
	Address Address
}

type Address struct{
	City string
	State string
	Country string
	Pincode json.Number
}

func main()  {
	fmt.Println("GO OWN DATABASE ")
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	employees := []User{
		{"john", "23" , "34567890", "Abc ltd", Address{"ghansoli", "Maharashtra", "India", "471010"} },
		{"Mru", "23" , "34567990", "RbK llp", Address{"Satara", "Maharashtra", "India", "481010"} },
		{"Bhu", "25" , "34568090", "Mai ltd", Address{"Indore", "Madhya Pradesh", "India", "450010"} },
	}

	for _, val := range employees{
		db.Write("users", val.Name, User{
			Name : val.Name,
	Age: val.Age ,
	Contact : val.Contact,
	Company : val.Company,
	Address : val.Address,
		})


	}

	records , err := db.ReadAll("users")

	if err != nil {
		fmt.Println("Error", err)
	}
	
	fmt.Println(records)

	allUsers := []User{}

	for _, val := range records{
		employeeFound := User{}
		if err := json.Unmarshal([]byte(val), &employeeFound); err != nil{
			fmt.Println("Error", err)
		}
		allUsers = append(allUsers, employeeFound)
	}

	fmt.Println(allUsers)

	if err := db.Delete("users", "john"); err != nil {
		fmt.Println("Error", err)
	}

	// if err:= db.Delete("users", ""); err != nil {
	// 	fmt.Println("Error", err)
	// }

}