package main


import "fmt"

type employee struct {
	firstname string
	lastname  string
	age       int
}

func main() {
	const Maxrities = 5
	var employeeid int = 1001
	var emp employee
	emp.firstname = "John"
	emp.lastname = "Doe"
	emp.age = 30
	fmt.Println("Employee ID:", employeeid)
	fmt.Println("Employee Name:", emp.firstname, emp.lastname)
	fmt.Println("Employee Age:", emp.age)

}

