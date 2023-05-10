package game

import (
	"fmt"

	"github.com/conggova/doudizhu-robot/pkg/action"
	"github.com/conggova/doudizhu-robot/pkg/player"
	"github.com/conggova/doudizhu-robot/pkg/strategy"
)

// Game要满足Simulator接口
type Game struct {
	//基础信息
	players  [3]player.Player //玩家
	openCard bool             //这局游戏是否明牌
	//牌局相关信息（为了打牌需要）
	stageFlag   int              //牌局到了哪个阶段 0 未开始  1 首叫完成 2 第二家叫分完成  3 第三家叫分完成  4正在打牌
	cents       [3]int           //各家叫分(含地主与底分信息)
	lastAction  action.Action    //最后一手牌
	llastAction action.Action    //倒数第二手牌
	onTurn      int              //该谁行动
	dipaiRemain action.PokerSet2 //未出现的底牌
	//统计信息（计账需要）
	winner         int    //赢家
	bombCnt        int    //炸弹数
	dizhuPlayCnt   int    //地主出牌手数
	nongminPlayCnt int    //农民出牌手数
	profits        [3]int //记录收益
	end            bool   //游戏已经结束
	//附加信息（日志需要）
	headPlayer int  //头家
	keepLog    bool //如果KeepLog会写日志到文件
}

// 重置牌局有关信息
func (g *Game) resetGame() {
	g.stageFlag = 0
	g.cents = [3]int{}
	g.lastAction = [2]uint64{}
	g.llastAction = [2]uint64{}
	g.onTurn = 0
	g.dipaiRemain = 0
	g.winner = 0
	g.bombCnt = 0
	g.dizhuPlayCnt = 0
	g.nongminPlayCnt = 0
	g.end = false
	g.profits = [3]int{}
}

// 返回指定玩家的收益
func (g Game) GetProfits() [3]int {
	if !g.end {
		panic("game is still ongoing , u cant get profits")
	}
	return g.profits
}

func (g *Game) RunGameWithRandomBeginning(headPlayer int) {
	//随机发牌
	pokerSets := strategy.RandDistCards()
	pokerSet2s := [4]action.PokerSet2{}
	for i := 0; i < 4; i++ {
		pokerSet2s[i] = pokerSets[i].PokerSet2()
	}
	g.dipaiRemain = pokerSet2s[3]
	//g.dipai = pokerSet2s[3]
	for i := 0; i < 3; i++ {
		g.players[i].SetRemainPokerSet(pokerSet2s[i])
	}
	g.headPlayer = headPlayer
	g.onTurn = headPlayer
	g.runGame()
}

func (g *Game) runGame() {
	nameMap := [3]string{"刘备", "关羽", "张飞"}
	centMap := [4]string{"不叫", "叫一分", "叫二分", "叫三分"}
	stageMap := [3]string{"首家叫分", "二家叫分", "尾家叫分"}
	if g.keepLog {
		fwrite("\n\n----------------牌局开始---------------\n", false)
		if g.openCard {
			fwrite("这一局牌明牌进行\n", false)
			for i := 0; i < 3; i++ {
				fwrite(fmt.Sprintf("%s 手牌为 %s\n", nameMap[i], g.players[i].GetRemainPokerSet().String()), false)
			}
			fwrite(fmt.Sprintf("三张底牌为 %s\n", g.dipaiRemain.String()), false)
		} else {
			for i := 0; i < 3; i++ {
				fwrite(fmt.Sprintf("%s 手牌为 %s\n", nameMap[i], g.players[i].GetRemainPokerSet().String()), true)
			}
			fwrite(fmt.Sprintf("三张底牌为 %s\n", g.dipaiRemain.String()), true)
		}
		fwrite("开始叫分\n", false)
	}

	if g.stageFlag < 3 {
		t := g.onTurn
		for i := g.stageFlag; i < 3; i++ {
			g.cents[t] = g.players[t].Call(strategy.CallContext{
				StageFlag:            g.stageFlag,
				PreerCent:            g.cents[(t+2)%3],
				NexterCent:           g.cents[(t+1)%3],
				GameOpenCard:         g.openCard,
				PreerRemainPokerSet:  g.players[(t+2)%3].GetRemainPokerSet(),
				NexterRemainPokerSet: g.players[(t+1)%3].GetRemainPokerSet(),
			})
			if g.keepLog {
				fwrite(fmt.Sprintf("%s \n%s %s \n\n", stageMap[i], nameMap[t], centMap[g.cents[t]]), false)
			}
			g.stageFlag += 1
			t = (t + 1) % 3
		}
	}
	if g.stageFlag == 3 {
		baseCent, dizhu := baseCentAndDizhu(g.cents)
		if baseCent == 0 {
			g.end = true
			if g.keepLog {
				fwrite("无人叫分，此局以流局结束\n----------------牌局结束---------------\n\n", false)
			}
			return
		}
		if g.dipaiRemain.PokerCount() != 3 {
			panic("dipaiRemain is not three")
		}
		g.players[dizhu].SetRemainPokerSet(g.players[dizhu].GetRemainPokerSet().CombineWith(g.dipaiRemain))
		g.onTurn = dizhu
		if g.keepLog {
			fwrite(fmt.Sprintf("叫分结束，地主为 %s , 这一局的底分为 %d \n\n开始打牌\n", nameMap[dizhu], baseCent), false)
		}
		g.stageFlag = 4
	}
	g.enterPlayLoop()
}

func (g *Game) enterPlayLoop() {
	if g.stageFlag != 4 {
		panic("not start play yet")
	}
	_, dizhu := baseCentAndDizhu(g.cents)
	nameMap := [3]string{"刘备", "关羽", "张飞"}
	nameMap[dizhu] += "*"
	for !(g.players[0].GetRemainPokerSet() == 0 || g.players[1].GetRemainPokerSet() == 0 || g.players[2].GetRemainPokerSet() == 0) {
		var a action.Action
		currentPlayer := g.onTurn
		//如果是明牌，会给Player另外两家的牌形信息
		dizhuFlag := 0
		if dizhu == (currentPlayer+1)%3 {
			dizhuFlag = 1
		} else if dizhu == (currentPlayer+2)%3 {
			dizhuFlag = 2
		}

		if g.openCard {
			preerPS := g.players[(currentPlayer+2)%3].GetRemainPokerSet()
			nexterPS := g.players[(currentPlayer+1)%3].GetRemainPokerSet()
			a = g.players[currentPlayer].Play(strategy.PlayContext{
				DizhuFlag:            dizhuFlag,
				PreerAction:          g.lastAction,
				NexterAction:         g.llastAction,
				PreerPkCnt:           preerPS.PokerCount(),
				NexterPkCnt:          nexterPS.PokerCount(),
				DizhuDipaiRemain:     g.dipaiRemain,
				GameOpenCard:         g.openCard,
				OthersRemainPokerSet: preerPS.CombineWith(nexterPS),
				PreerRemainPokerSet:  preerPS,
				NexterRemainPokerSet: nexterPS,
			})

		} else {
			preerPS := g.players[(currentPlayer+2)%3].GetRemainPokerSet()
			nexterPS := g.players[(currentPlayer+1)%3].GetRemainPokerSet()
			a = g.players[currentPlayer].Play(strategy.PlayContext{
				DizhuFlag:            dizhuFlag,
				PreerAction:          g.lastAction,
				NexterAction:         g.llastAction,
				PreerPkCnt:           preerPS.PokerCount(),
				NexterPkCnt:          nexterPS.PokerCount(),
				DizhuDipaiRemain:     g.dipaiRemain,
				GameOpenCard:         g.openCard,
				OthersRemainPokerSet: preerPS.CombineWith(nexterPS),
				PreerRemainPokerSet:  0,
				NexterRemainPokerSet: 0,
			})
		}
		g.putAction(currentPlayer, a)
		if g.keepLog {
			fwrite(fmt.Sprintf("%s 出 %s    ", nameMap[currentPlayer], a.String()), false)
			if g.openCard {
				fwrite(fmt.Sprintf("他手里还剩 %s     ", g.players[currentPlayer].GetRemainPokerSet()), false)
			} else {
				fwrite(fmt.Sprintf("他手里还剩 %d 张     ", g.players[currentPlayer].GetRemainPokerSet().PokerCount()), false)
				fwrite(fmt.Sprintf(", 他手里还剩 %s     ", g.players[currentPlayer].GetRemainPokerSet()), true)
			}
			if g.onTurn == dizhu {
				fwrite(fmt.Sprintf("地主还没出的底牌为 %s\n\n", g.dipaiRemain.String()), false)
			} else {
				fwrite("\n\n", false)
			}
		}
		////转到下家
		g.llastAction = g.lastAction
		g.lastAction = a
		g.onTurn = (currentPlayer + 1) % 3
	}
	g.winner = -1
	for i := 0; i < 3; i++ {
		if g.players[i].GetRemainPokerSet() == 0 {
			g.winner = i
		}
	}

	if g.winner == -1 {
		panic("end game with no player empty his cards")
	}
	g.end = true
	g.calcuProfit()
	if g.keepLog {
		fwrite("此局结束 ， 各玩家晾牌\n", false)
		for i := 0; i < 3; i++ {
			fwrite(fmt.Sprintf("%s 余牌为 %s\n", nameMap[i], g.players[i].GetRemainPokerSet().String()), false)
		}
		if g.winner == dizhu {
			fwrite("此局地主取胜\n", false)
		} else {
			fwrite("此局农民取胜\n", false)
		}
		fwrite(fmt.Sprintf("各方收益：\n	刘备 %d\n	关羽 %d\n	张飞 %d\n", g.profits[0], g.profits[1], g.profits[2]), false)
		fwrite("----------------牌局结束---------------\n", false)
	}
}
