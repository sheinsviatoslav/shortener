package hash

import "fmt"

func ExampleGenerator() {
	out1 := Generator(8)
	fmt.Println(out1)

	out2 := Generator(20)
	fmt.Println(out2)

	out3 := Generator(100)
	fmt.Println(out3)
}
