package strategy

import (
	"fmt"
	"testing"
)

func Test_getTopKActionIDs(t *testing.T) {
	d := []int{5, 3, 8, 10, 9, 2, 3, 8, 7, 6}
	ids := getTopKActionIDs(d, 3)
	fmt.Println(ids)
	if len(ids) != 2 {
		t.Error("getTopKActionIDs incorrect")
	}
}
