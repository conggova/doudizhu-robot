package strategy

import (
	"fmt"
	"math/rand"

	"github.com/conggova/doudizhu-robot/pkg/action"
)

func RandDistCards() [4]action.PokerSet {
	var res [4]action.PokerSet
	//1代表已被选过了
	var choosedMark uint64
	//生成前两手
	pkCnt := 0
	for {
		randomIdx := uint64(1 << rand.Intn(54))
		if choosedMark&randomIdx == 0 { //not choosen yet
			choosedMark |= randomIdx
			pkCnt += 1
			if pkCnt == 17 {
				break
			}
		}
	}
	res[0] = action.PokerSet(choosedMark)

	pkCnt = 0
	for {
		randomIdx := uint64(1 << rand.Intn(54))
		if choosedMark&randomIdx == 0 { //not choosen yet
			choosedMark |= randomIdx
			pkCnt += 1
			if pkCnt == 17 {
				break
			}
		}
	}
	res[1] = action.PokerSet(choosedMark).Subtract(res[0])

	pkCnt = 0
	for {
		randomIdx := uint64(1 << rand.Intn(54))
		if choosedMark&randomIdx == 0 { //not choosen yet
			choosedMark |= randomIdx
			pkCnt += 1
			if pkCnt == 3 {
				break
			}
		}
	}
	res[3] = action.PokerSet(choosedMark).Subtract(res[0]).Subtract(res[1])
	res[2] = action.PokerSet(0b111111111111111111111111111111111111111111111111111111).Subtract(action.PokerSet(choosedMark))
	return res
}

// 切牌函数
func RandDistCards2(ps action.PokerSet2, preerNum int,
	nexterNum int) (action.PokerSet2, action.PokerSet2) {
	prePS, nextPS := cutCardsRandom(ps, preerNum, nexterNum)
	return prePS, nextPS

}

func cutCardsRandom(ps action.PokerSet2, preerNum int, nexterNum int) (action.PokerSet2, action.PokerSet2) {
	cards := []int8{} //每个元素一张牌
	for i := 0; i < 15; i++ {
		cnt := ps >> (i << 2) & 0b1111
		for j := 0; j < int(cnt); j++ {
			cards = append(cards, int8(i))
		}
	}
	if preerNum+nexterNum != len(cards) {
		panic(fmt.Sprintf("preerNum %d + nexterNum %d != len(cards) %d , %s ", preerNum, nexterNum, len(cards), ps))
	}
	//true代表已被选过了
	choosed := make([]bool, len(cards))
	var prePS action.PokerSet2
	pkCnt := 0
	for {
		randomIdx := rand.Intn(int(preerNum + nexterNum))
		if !choosed[randomIdx] {
			choosed[randomIdx] = true
			prePS += 1 << (cards[randomIdx] << 2)
			pkCnt += 1
			if pkCnt == preerNum {
				break
			}
		}
	}
	return prePS, ps.Subtract(prePS)
}
