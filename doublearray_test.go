package doublearray

import (
	"bufio"
	"os"
	"testing"
)

func TestDoubleArray(t *testing.T) {
	fp, err := os.Open("words.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()

	keys := make([]string, 0)
	values := make([]int, 0)

	scanner := bufio.NewScanner(fp)
	for v := 0; scanner.Scan(); v++ {
		keys = append(keys, scanner.Text())
		values = append(values, v)
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	da, err := Build(keys, values)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(keys); i++ {
		v, found := da.Lookup(keys[i])
		if !found {
			t.Fatalf("keys[i] = %s is not found", keys[i])
		}
		if v != values[i] {
			t.Fatalf("v = %d, but values[i] = %d", v, values[i])
		}
	}

	decKeys, decValues := da.Enumerate()
	if len(decKeys) != len(keys) {
		t.Fatalf("len(decKeys) = %d, but len(keys) = %d", len(decKeys), len(keys))
	}

	for i := 0; i < len(keys); i++ {
		if decKeys[i] != keys[i] {
			t.Fatalf("decKeys[i] = %s, but keys[i] = %s", decKeys[i], keys[i])
		}
		if decValues[i] != values[i] {
			t.Fatalf("decValues[i] = %d, but values[i] = %d", decValues[i], values[i])
		}
	}
}
