package main

import (
	"fmt"
	"log"

	doublearray "github.com/kampersanda/doublearray-go"
)

func main() {
	keys := []string{
		"Aru", "Bocci", "Kai", "Kako", "Nako", "Nakosuke", "Sotca",
	}
	values := []int{
		1224, 505, 505, 611, 504, 504, 731,
	}

	da, err := doublearray.Build(keys, values)
	if err != nil {
		log.Fatal(err)
	}

	value, found := da.Lookup("Bocci")
	fmt.Printf("Bocci -> %d (found = %t)\n", value, found)

	value, found = da.Lookup("Peko")
	fmt.Printf("Peko -> ? (found = %t)\n", found)

	decKeys, decValues := da.Enumerate()
	for i := 0; i < len(decKeys); i++ {
		fmt.Printf("%s -> %d\n", decKeys[i], decValues[i])
	}
}
