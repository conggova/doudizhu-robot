package strategy

import (
	"fmt"
	"testing"

	"github.com/conggova/doudizhu-robot/pkg/action"
)

func Test_RandDistCards(t *testing.T) {
	pss := RandDistCards()
	a := pss[0].CombineWith(pss[1]).CombineWith(pss[2]).CombineWith(pss[3])
	fmt.Println(pss)
	fmt.Println(a.PokerSet2())
	fmt.Println(action.PokerSet2(0x114444444444444))
	if a.PokerSet2() != 0x114444444444444 {
		t.Error("RandDistCards incorrect")
	}
	if pss[0].PokerCount() != 17 || pss[1].PokerCount() != 17 || pss[2].PokerCount() != 17 || pss[3].PokerCount() != 3 {
		t.Error("RandDistCards incorrect2")
	}
}

func Test_RandDistCards2(t *testing.T) {
	p1, p2 := RandDistCards2(0x102200000, 1, 4)
	if p1.CombineWith(p2) != 0x102200000 || p1.PokerCount() != 1 || p2.PokerCount() != 4 {
		t.Error("RandDistCards2 incorrect")
	}
}
