package __random

import (
	"fmt"
	"sort"
	"testing"
)

func TestSearch(t *testing.T) {
	s := []string{"a", "b", "c", "d", "e"}

	id := sort.SearchStrings(s, "b")
	fmt.Println("b at", id)

	if id != 1 {
		t.Fatal("id was not 1")
	}

	id = sort.SearchStrings(s, "z")
	fmt.Println("z at", id)

	if id != len(s) {
		t.Fatal("id was not len(s)")
	}
}
