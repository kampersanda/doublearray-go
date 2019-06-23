package doublearray

import (
	"fmt"
)

// DoubleArray implements double-array minimal-prefix trie.
type DoubleArray struct {
	array    []node
	tail     []byte
	numKeys  int
	numNodes int
}

// Build returns a DoubleArray object built from sorted key strings and associated values.
// Key duplication and empty key are not allowed.
// NULL character byte(0) must not be included since it is used for the terminator.
func Build(keys []string, values []int) (*DoubleArray, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("keys must not be empty")
	}
	if len(keys) != len(values) {
		return nil, fmt.Errorf("The size of keys must be equal to that of values")
	}

	b := builder{keys: keys, values: values}
	b.init()
	err := b.arrange(0, len(keys), 0, 0)
	if err != nil {
		return nil, err
	}
	b.finish()

	numNodes := 1 // 1 is for the root
	for i := 1; i < len(b.array); i++ {
		if b.array[i].check >= 0 {
			numNodes++
		}
	}

	return &DoubleArray{array: b.array, tail: b.tail, numKeys: len(keys), numNodes: numNodes}, nil
}

// NumKeys returns the number of keys stored.
func (da *DoubleArray) NumKeys() int {
	return da.numKeys
}

// NumNodes returns the number of nodes.
func (da *DoubleArray) NumNodes() int {
	return da.numNodes
}

// ArrayLen returns the length of BASE/CHECK array
func (da *DoubleArray) ArrayLen() int {
	return len(da.array)
}

// TailLen returns the length of TAIL array
func (da *DoubleArray) TailLen() int {
	return len(da.tail)
}

// AllocBytes returns the allocated size in bytes.
func (da *DoubleArray) AllocBytes() int {
	return da.ArrayLen()*8 + da.TailLen()
}

// Lookup returns the associated value with the given key if found.
func (da *DoubleArray) Lookup(key string) (int, bool) {
	npos := 0
	depth := 0

	for ; depth < len(key); depth++ {
		if da.array[npos].base < 0 {
			break
		}
		cpos := da.array[npos].base ^ int(key[depth])
		if da.array[cpos].check != npos {
			return 0, false
		}
		npos = cpos
	}

	if da.array[npos].base >= 0 {
		cpos := da.array[npos].base // ^ int(terminator)
		if da.array[cpos].check != npos {
			return 0, false
		}
		return da.array[cpos].base, true
	}

	tpos := -da.array[npos].base
	for ; depth < len(key); depth++ {
		if da.tail[tpos] != key[depth] {
			return 0, false
		}
		tpos++
	}

	if da.tail[tpos] != terminator {
		return 0, false
	}
	return da.getValue(tpos + 1), true
}

// PrefixLookup returns the keys and associated values included as prefixes of the given key.
func (da *DoubleArray) PrefixLookup(key string) ([]string, []int) {
	keys := make([]string, 0)
	values := make([]int, 0)

	npos := 0
	depth := 0

	for ; depth < len(key); depth++ {
		if da.array[npos].base < 0 {
			break
		}

		base := da.array[npos].base

		if da.array[base].check == npos {
			keys = append(keys, key[:depth])
			values = append(values, da.array[base].base)
		}

		cpos := base ^ int(key[depth])
		if da.array[cpos].check != npos {
			return keys, values
		}
		npos = cpos
	}

	base := da.array[npos].base

	if base >= 0 {
		if da.array[base].check == npos {
			keys = append(keys, key[:depth])
			values = append(values, da.array[base].base)
		}
		return keys, values
	}

	tpos := -base
	for ; depth < len(key); depth++ {
		if da.tail[tpos] != key[depth] {
			return keys, values
		}
		tpos++
	}

	if da.tail[tpos] == terminator {
		keys = append(keys, key[:depth])
		values = append(values, da.getValue(tpos+1))
	}

	return keys, values
}

// PredictiveLookup returns the keys and associated values starting with prefixes of the given key.
func (da *DoubleArray) PredictiveLookup(key string) ([]string, []int) {
	keys := make([]string, 0, da.numKeys)
	values := make([]int, 0, da.numKeys)

	npos := 0
	depth := 0

	for ; depth < len(key); depth++ {
		if da.array[npos].base < 0 {
			return keys, values
		}

		cpos := da.array[npos].base ^ int(key[depth])
		if da.array[cpos].check != npos {
			return keys, values
		}
		npos = cpos
	}

	keys, values = da.enumerate(npos, depth, []byte(key), keys, values)

	return keys, values
}

func (da *DoubleArray) getValue(tpos int) int {
	return int(da.tail[tpos]) | int(da.tail[tpos+1])<<8 | int(da.tail[tpos+2])<<16 | int(da.tail[tpos+3])<<24
}

func (da *DoubleArray) enumerate(npos int, depth int, decoded []byte, keys []string, values []int) ([]string, []int) {
	if da.array[npos].base < 0 {
		tpos := -da.array[npos].base
		for da.tail[tpos] != byte(0) {
			decoded = append(decoded, da.tail[tpos])
			tpos++
		}
		keys = append(keys, string(decoded))
		values = append(values, da.getValue(tpos+1))
		return keys, values
	}

	base := da.array[npos].base
	cpos := base // ^ int(terminator)

	if da.array[cpos].check == npos {
		keys = append(keys, string(decoded))
		values = append(values, da.array[cpos].base)
	}

	for c := 1; c < 256; c++ {
		decoded = decoded[:depth]
		cpos = da.array[npos].base ^ c
		if da.array[cpos].check == npos {
			decoded = append(decoded, byte(c))
			keys, values = da.enumerate(cpos, depth+1, decoded, keys, values)
		}
	}

	return keys, values
}

const (
	terminator = byte(0)
)

type node struct {
	base, check int
}

type builder struct {
	array  []node
	tail   []byte
	keys   []string
	values []int
}

func (b *builder) init() {
	capa := 256
	for capa < len(b.keys) {
		capa <<= 1
	}

	array := make([]node, 256, capa)
	tail := make([]byte, 1)

	for i := 1; i < 256; i++ {
		array[i].base = -(i + 1)
		array[i].check = -(i - 1)
	}
	array[255].base = -1
	array[1].check = -255
	array[0].check = 1 // head empty

	b.array = array
	b.tail = tail
}

func (b *builder) finish() {
	b.array[0].check = -1
}

func (b *builder) enlarge() {
	oldLen := len(b.array)
	newLen := oldLen + 256

	for i := oldLen; i < newLen; i++ {
		b.array = append(b.array, node{base: -(i + 1), check: -(i - 1)})
	}

	if b.array[0].check == 0 {
		b.array[oldLen].check = -(newLen - 1) // prev
		b.array[newLen-1].base = -oldLen      // next
		b.array[0].check = oldLen
	} else {
		empHead := b.array[0].check
		empTail := -b.array[empHead].check
		b.array[oldLen].check = -empTail
		b.array[empTail].base = -oldLen
		b.array[empHead].check = -(newLen - 1)
		b.array[newLen-1].base = -empHead
	}
}

func (b *builder) fix(npos int) {
	next := -b.array[npos].base
	prev := -b.array[npos].check
	b.array[next].check = -prev
	b.array[prev].base = -next

	if npos == b.array[0].check {
		if next == npos {
			b.array[0].check = 0
		} else {
			b.array[0].check = next
		}
	}
}

func (b *builder) arrange(bpos, epos, depth, npos int) error {
	if bpos+1 == epos {
		b.array[npos].base = -len(b.tail)
		for ; depth < len(b.keys[bpos]); depth++ {
			if b.keys[bpos][depth] == terminator {
				return fmt.Errorf("keys must not include NULL terminator byte(0)")
			}
			b.tail = append(b.tail, b.keys[bpos][depth])
		}
		b.tail = append(b.tail, terminator)

		val := b.values[bpos]
		for i := 0; i < 4; i++ {
			b.tail = append(b.tail, byte(val%256))
			val >>= 8
		}
		return nil
	}

	edges := make([]byte, 0)
	isTerminate := len(b.keys[bpos]) == depth

	if isTerminate {
		bpos++
		if len(b.keys[bpos]) == depth {
			return fmt.Errorf("Key duplication is not allowed")
		}
		edges = append(edges, terminator)
	}

	c := b.keys[bpos][depth]
	for i := bpos + 1; i < epos; i++ {
		c2 := b.keys[i][depth]
		if c != c2 {
			if c2 < c {
				return fmt.Errorf("keys must be sorted in lex order")
			}
			if c == terminator {
				return fmt.Errorf("keys must not include NULL terminator byte(0)")
			}
			edges = append(edges, c)
			c = c2
		}
	}

	if c == terminator {
		return fmt.Errorf("keys must not include NULL terminator byte(0)")
	}
	edges = append(edges, c)

	base := b.xcheck(edges)
	if len(b.array) <= base {
		b.enlarge()
	}

	b.array[npos].base = base
	for _, c := range edges {
		cpos := base ^ int(c)
		b.fix(cpos)
		b.array[cpos].check = npos
	}

	if isTerminate {
		cpos := base // ^ int(terminator)
		b.array[cpos].base = b.values[bpos-1]
	}

	i := bpos
	c = b.keys[bpos][depth]
	for j := bpos + 1; j < epos; j++ {
		c2 := b.keys[j][depth]
		if c != c2 {
			err := b.arrange(i, j, depth+1, base^int(c))
			if err != nil {
				return err
			}
			i = j
			c = c2
		}
	}
	return b.arrange(i, epos, depth+1, base^int(c))
}

func (b *builder) xcheck(edges []byte) int {
	empHead := b.array[0].check

	if empHead == 0 {
		return len(b.array) ^ int(edges[0])
	}

	i := empHead
	for {
		base := i ^ int(edges[0])
		if b.isTarget(base, edges) {
			return base
		}
		i = -b.array[i].base
		if i == empHead {
			break
		}
	}

	return len(b.array) ^ int(edges[0])
}

func (b *builder) isTarget(base int, edges []byte) bool {
	for _, c := range edges {
		i := base ^ int(c)
		if b.array[i].check >= 0 {
			return false
		}
	}
	return true
}
