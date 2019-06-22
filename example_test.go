package doublearray_test

import (
	"fmt"

	doublearray "github.com/kampersanda/doublearray-go"
)

var da *doublearray.DoubleArray

func Example() {
	keys := []string{
		"Aru", "Bocci", "Kai", "Kako", "Nako", "Nakosuke", "Sotca",
	}
	values := []int{
		1224, 505, 505, 611, 504, 504, 731,
	}

	da, _ = doublearray.Build(keys, values)

	value, found := da.Lookup("Bocci")
	fmt.Println(value, found)

	value, found = da.Lookup("Peko")
	fmt.Println(value, found)

	// Output:
	// 505 true
	// 0 false
}

func Example_enumerate() {
	keys, values := da.Enumerate()
	fmt.Println(keys)
	fmt.Println(values)

	// Output:
	// [Aru Bocci Kai Kako Nako Nakosuke Sotca]
	// [1224 505 505 611 504 504 731]
}
