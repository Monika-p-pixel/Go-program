package main
import (

"fmt"
"math")
func main() {
	var a,b int = 10, 3
	var result int

	result = a + b
	fmt.Println("Addition:", result)

	result = a - b
	fmt.Println("Subtraction:", result)

	result = a * b
	fmt.Println("Multiplication:", result)

	result = a / b
	fmt.Println("Division:", result)

	result = a % b
	fmt.Println("Modulus:", result)

	const p float64 = 22/7	
	fmt.Println(p)

	var maxInt int64 = 9223372036854775807
	fmt.Println(maxInt)

	maxInt = maxInt + 1
	fmt.Println(maxInt)

	var umaxint uint64 = 18446744073709551615
	fmt.Println(umaxint)

	umaxint = umaxint + 1
	fmt.Println(umaxint)

	var smallint float64 = 	1.0e-323
	fmt.Println(smallint)
	smallint = smallint/math.MaxFloat64
	fmt.Println(smallint)


}	