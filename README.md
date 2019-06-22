# doublearray-go

This package implements double-array minimal-prefix trie.

## Usage

```go
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
```

will produce

```
Bocci -> 505 (found = true)
Peko -> ? (found = false)
Aru -> 1224
Bocci -> 505
Kai -> 505
Kako -> 611
Nako -> 504
Nakosuke -> 504
Sotca -> 731
```
