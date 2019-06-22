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
		1, 2, 3, 4, 5, 6, 7,
	}

	da, err := doublearray.Build(keys, values)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exact Lookup for 'Bocci' and 'Peko':")
	{
		value, found := da.Lookup("Bocci")
		fmt.Printf("- Bocci -> %d (found = %t)\n", value, found)

		value, found = da.Lookup("Peko")
		fmt.Printf("- Peko -> ? (found = %t)\n", found)
	}

	fmt.Println("Common Prefix Lookup for 'Nakosuke':")
	{
		targetKeys, targetValues := da.PrefixLookup("Nakosuke")
		for i := 0; i < len(targetKeys); i++ {
			fmt.Printf("- %s -> %d\n", targetKeys[i], targetValues[i])
		}
	}

	fmt.Println("Common Prefix Lookup for 'Ka':")
	{
		targetKeys, targetValues := da.PredictiveLookup("Ka")
		for i := 0; i < len(targetKeys); i++ {
			fmt.Printf("- %s -> %d\n", targetKeys[i], targetValues[i])
		}
	}

	fmt.Println("Enumerate all keys:")
	{
		targetKeys, targetValues := da.PredictiveLookup("")
		for i := 0; i < len(targetKeys); i++ {
			fmt.Printf("- %s -> %d\n", targetKeys[i], targetValues[i])
		}
	}
}
