package strategy

// UCB1 introduce https://blog.csdn.net/conggova/article/details/79498761

import (
	"math"

	"github.com/conggova/doudizhu-robot/pkg/action"
)

// UCB1算法的变形
type UCB1SimulateStrategy struct {
	Simulator            Simulator
	StrategyWhenSimulate Strategy
	AvgTryCnt            int //
}

func (s UCB1SimulateStrategy) MakeDecisionWhenCall(c StrategyCallContext) int {
	return simulateCall(c, s.Simulator, SimulateStrategy{
		DistCardTimes:        5,
		SimulateTimesPerDist: 1,
		Simulator:            s.Simulator,
		StrategyWhenSimulate: RandomStrategy{},
	}, 10, 1, false)
}

func (s UCB1SimulateStrategy) MakeDecisionWhenPlay(c StrategyPlayContext) (actionTaken action.Action) {
	contextAction := c.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = c.NexterAction
	}
	possibleActionList := c.RemainPokerSet.PossibleActionsWithContext(contextAction)
	if len(possibleActionList) == 1 {
		return possibleActionList[0]
	}
	return ucb1SimulateChooseAction(possibleActionList, c, s.Simulator, s.StrategyWhenSimulate, s.AvgTryCnt)
}

func ucb1SimulateChooseAction(possibleActionList []action.Action, c StrategyPlayContext,
	game Simulator, strategyWhenSimulate Strategy, avgTryCnt int) action.Action {

	partners := [2]int{}
	if c.CheatFlag == WithNexter {
		partners = [2]int{0, 1}
	} else if c.CheatFlag == WithPreer {
		partners = [2]int{0, 2}
	}
	simulator := game.BuildGame(GameBuildInfo{
		Strategys:   [3]Strategy{strategyWhenSimulate, strategyWhenSimulate, strategyWhenSimulate},
		OpenCard:    c.GameOpenCard,
		Partners:    partners,
		CheatMethod: c.CheatMethod,
	})

	actionCnt := len(possibleActionList)
	totalProfits := make([]int, actionCnt)
	tryCnt := make([]int, actionCnt)
	totalTryCnt := 0
	var totalProfitABSSum float64

	tryAction := func(id int, redistCard bool) {
		if !c.OpenCard4Me && redistCard {
			c.PreerRemainPokerSet, c.NexterRemainPokerSet = RandDistCards2(c.OthersRemainPokerSet, c.PreerPkCnt, c.NexterPkCnt)
		}
		tryCnt[id] += 1
		totalTryCnt += 1
		a := possibleActionList[id]
		dizhuPlayCnt := c.DizhuPlayCnt
		nongminPlayCnt := c.NongminPlayCnt
		dizhuDipaiRemain := c.DizhuDipaiRemain
		if a.ActionType() != action.Pass {
			if c.DizhuFlag == 0 {
				dizhuDipaiRemain = dizhuDipaiRemain.SubtractBestEffort(a.PokerSet2())
				dizhuPlayCnt += 1
			} else {
				nongminPlayCnt += 1
			}
		}
		bombCnt := 0
		if a.IsBomb() {
			bombCnt += 1
		}
		cents := [3]int{}      //根据地主是谁伪造叫分
		cents[c.DizhuFlag] = 3 //分数对于收益对比来说没有意义
		//simulator中当前player为0号
		simulator.RestoreGame(GameRestoreInfo{
			StageFlag:        4,
			Cents:            cents,
			LastAction:       a,
			LLastAction:      c.PreerAction,
			OnTurn:           1,
			DizhuDipaiRemain: dizhuDipaiRemain,
			DizhuPlayCnt:     dizhuPlayCnt,
			NongminPlayCnt:   nongminPlayCnt,
			BombCnt:          bombCnt,
			RemainPokerSets: [3]action.PokerSet2{
				c.RemainPokerSet.Subtract(a.PokerSet2()),
				c.NexterRemainPokerSet,
				c.PreerRemainPokerSet,
			},
		})
		simulator.RunRestoredGame()
		profits := simulator.GetProfits()
		profit := profits[0]
		if c.CheatFlag == WithNexter { //with nexter
			profit += profits[1]
		} else if c.CheatFlag == WithPreer { //with preer
			profit += profits[2]
		}
		totalProfits[id] += profit
		totalProfitABSSum += math.Abs(float64(profit))
	}
	//init try every action once
	if !c.OpenCard4Me { //如果看不到另外两家的牌
		c.PreerRemainPokerSet, c.NexterRemainPokerSet = RandDistCards2(c.OthersRemainPokerSet, c.PreerPkCnt, c.NexterPkCnt)
	}
	for i := 0; i < actionCnt; i++ {
		tryAction(i, false)
	}
	getNextTry := func() int {
		//找出上界最大的选择
		absMean := totalProfitABSSum / float64(totalTryCnt)
		// 上界为 mean +  absMean * (2ln(n)/ni)**0.5
		ln_n2 := math.Log(float64(totalTryCnt)) * 2
		selected := -1
		maxUpBound := 0.0
		for i := 0; i < actionCnt; i++ {
			upBound := float64(totalProfits[i])/float64(tryCnt[i]) + absMean*math.Pow(ln_n2/float64(tryCnt[i]), 0.5)
			if selected == -1 {
				selected = i
				maxUpBound = upBound
			} else {
				if upBound > maxUpBound {
					maxUpBound = upBound
					selected = i
				}
			}
		}
		return selected
	}
	for i := 0; i < actionCnt*avgTryCnt; i++ {
		tryAction(getNextTry(), true)
	}
	//返回均值最大的
	maxId := -1
	maxAvg := 0
	for i := 0; i < actionCnt; i++ {
		avg := totalProfits[i] / tryCnt[i]
		if maxId == -1 {
			maxId = i
			maxAvg = avg
		} else {
			if avg > maxAvg {
				maxAvg = avg
				maxId = i
			}
		}
	}
	return possibleActionList[maxId]
}
