package game

import (
	"fmt"
	"os"

	"github.com/conggova/doudizhu-robot/pkg/action"
)

func baseCentAndDizhu(cents [3]int) (int, int) {
	dizhu := 0
	baseCent := 0
	for i := 0; i < 3; i++ {
		if cents[i] > baseCent {
			baseCent = cents[i]
			dizhu = i
		}
	}
	return baseCent, dizhu
}

// 记录一次出牌
func (g *Game) putAction(playerId int, a action.Action) {
	if g.stageFlag != 4 {
		panic("must startPlay before u can putAction")
	}
	if a.IsBomb() {
		g.bombCnt += 1
	}
	_, dizhu := baseCentAndDizhu(g.cents)
	if a.ActionType() != action.Pass {
		if playerId == dizhu {
			//更新未出底牌,  假定玩家优先出底牌里面的牌
			g.dipaiRemain = g.dipaiRemain.SubtractBestEffort(a.PokerSet2())
			g.dizhuPlayCnt += 1
		} else {
			g.nongminPlayCnt += 1
		}
	}
}

// 计算每个玩家的收益
func (g *Game) calcuProfit() {
	baseCent, dizhu := baseCentAndDizhu(g.cents)
	beishu := baseCent * (1 + g.bombCnt)
	//春天与反春
	if (g.winner == dizhu && g.nongminPlayCnt == 0) || (g.winner != dizhu && g.dizhuPlayCnt == 1) {
		beishu += baseCent
	}
	if g.winner != dizhu {
		beishu = -beishu
	}
	g.profits[dizhu] += 2 * beishu * 5 // 5 per cent
	g.profits[(dizhu+2)%3] -= beishu * 5
	g.profits[(dizhu+1)%3] -= beishu * 5
}

func fwrite(str string, onlyfile bool) {
	if !onlyfile {
		fmt.Print(str)
	}
	fout, err := os.OpenFile("gameLog.txt", os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		err := fout.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	fout.WriteString(str)
}
