package main

import (
	"fmt"
)


type Dog struct {
	Name string
	Height int
	Weight int
}

func (d *Dog) bark () bool {
	var condition bool = false;
	if d.Weight < 160 {
		condition = true;
		fmt.Println("Dog is barking")
		return condition
	}
	fmt.Println("Dog is not barking ")
	return condition
}


func CheckErr(err error){
	if err != nil {
		fmt.Println("error")
	}
}

func main() {

	dogObj := Dog {
		Name: "Shay",
		Height: 155,
		Weight: 140,
	}

	result := dogObj.bark()
	fmt.Println(result)
}

