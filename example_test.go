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
		1, 2, 3, 4, 5, 6, 7,
	}

	da, _ = doublearray.Build(keys, values)

	value, found := da.Lookup("Bocci")
	fmt.Println(value, found)

	value, found = da.Lookup("Peko")
	fmt.Println(value, found)

	// Output:
	// 2 true
	// 0 false
}

func Example_prefixLookup() {
	keys, values := da.PrefixLookup("Nakosuke")
	fmt.Println(keys)
	fmt.Println(values)

	// Output:
	// [Nako Nakosuke]
	// [5 6]
}

func Example_predictiveLookup() {
	keys, values := da.PredictiveLookup("Ka")
	fmt.Println(keys)
	fmt.Println(values)

	// Output:
	// [Kai Kako]
	// [3 4]
}

func Example_enumerate() {
	keys, values := da.PredictiveLookup("")
	fmt.Println(keys)
	fmt.Println(values)

	// Output:
	// [Aru Bocci Kai Kako Nako Nakosuke Sotca]
	// [1 2 3 4 5 6 7]
}
