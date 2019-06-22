# doublearray-go [![GoDoc](https://godoc.org/github.com/kampersanda/doublearray-go?status.svg)](https://godoc.org/github.com/kampersanda/doublearray-go) [![Go Report Card](https://goreportcard.com/badge/github.com/kampersanda/doublearray-go)](https://goreportcard.com/report/github.com/kampersanda/doublearray-go) [![Build Status](https://travis-ci.org/kampersanda/doublearray-go.svg?branch=master)](https://travis-ci.org/kampersanda/doublearray-go)

Package `doublearray-go` implements double-array minimal-prefix trie.

A double array is a fast and compact data structure for representing a trie, which can efficiently implement a dictionary with string keys.
The main feature of `doublearray-go` is to apply, instead of a (plain) trie, a **minimal-prefix trie** which replaces non-branching node-to-leaf paths in a trie into strings.
The minimal-prefix trie can reduce many nodes in a trie and can implement space- and cache-efficient trie-based dictionaries.


## Install

```
$ go get github.com/kampersanda/doublearray-go
```

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

	fmt.Println("Statistics:")
	{
		fmt.Printf("- NumKeys: %d\n", da.NumKeys())
		fmt.Printf("- NumNodes: %d\n", da.NumNodes())
		fmt.Printf("- ArrayLen: %d\n", da.ArrayLen())
		fmt.Printf("- TailLen: %d\n", da.TailLen())
		fmt.Printf("- AllocBytes: %d\n", da.AllocBytes())
	}
}
```

The output will be

```
Exact Lookup for 'Bocci' and 'Peko':
- Bocci -> 2 (found = true)
- Peko -> ? (found = false)
Common Prefix Lookup for 'Nakosuke':
- Nako -> 5
- Nakosuke -> 6
Common Prefix Lookup for 'Ka':
- Kai -> 3
- Kako -> 4
Enumerate all keys:
- Aru -> 1
- Bocci -> 2
- Kai -> 3
- Kako -> 4
- Nako -> 5
- Nakosuke -> 6
- Sotca -> 7
Statistics:
- NumKeys: 7
- NumNodes: 14
- ArrayLen: 256
- TailLen: 45
- AllocBytes: 2093
```