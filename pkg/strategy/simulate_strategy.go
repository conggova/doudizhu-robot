package strategy

import (
	"sort"

	"github.com/conggova/doudizhu-robot/pkg/action"
)

// 仿真策略
type SimulateStrategy struct {
	DistCardTimes        int //进行发牌次数
	SimulateTimesPerDist int //每次发牌进行的simulate次数
	Simulator            Simulator
	StrategyWhenSimulate Strategy //进行仿真时使用的策略出牌策略
	Verbose              bool
}

func (s SimulateStrategy) MakeDecisionWhenCall(c StrategyCallContext) int {
	return simulateCall(c, s.Simulator, s.StrategyWhenSimulate, s.DistCardTimes, s.SimulateTimesPerDist, s.Verbose)
}

func (s SimulateStrategy) MakeDecisionWhenPlay(c StrategyPlayContext) (actionTaken action.Action) {
	contextAction := c.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = c.NexterAction
	}
	possibleActionList := c.RemainPokerSet.PossibleActionsWithContext(contextAction)
	if len(possibleActionList) == 1 {
		return possibleActionList[0]
	}
	return simulateReduceAction(possibleActionList, c, s.Simulator, s.StrategyWhenSimulate, s.DistCardTimes, s.SimulateTimesPerDist, 1, s.Verbose)[0]
}

func simulateReduceAction(possibleActionList []action.Action, c StrategyPlayContext,
	game Simulator, strategyWhenSimulate Strategy, distCardTimes, simulateTimesPerDist, keep int, verbose bool) []action.Action {
	actionProfitDict := make([]int, len(possibleActionList))
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
	for distIdx := 0; distIdx < distCardTimes; distIdx++ {
		if !c.OpenCard4Me { //如果看不到另外两家的牌
			c.PreerRemainPokerSet, c.NexterRemainPokerSet = RandDistCards2(c.OthersRemainPokerSet, c.PreerPkCnt, c.NexterPkCnt)
		}
		for smIdx := 0; smIdx < simulateTimesPerDist; smIdx++ {
			for actionId, a := range possibleActionList {
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
				actionProfitDict[actionId] += profits[0]
				if c.CheatFlag == WithNexter { //with nexter
					actionProfitDict[actionId] += profits[1]
				} else if c.CheatFlag == WithPreer { //with preer
					actionProfitDict[actionId] += profits[2]
				}
			}
		}
	}
	if keep == 1 {
		return []action.Action{possibleActionList[getMaxProfitIdx(actionProfitDict)]}
	}
	//找出总收益最大的那个action
	ids := getTopKActionIDs(actionProfitDict, keep)
	res := make([]action.Action, len(ids))
	for idx, id := range ids {
		res[idx] = possibleActionList[id]
	}
	return res
}

func simulateCall(c StrategyCallContext, game Simulator, strategyWhenSimulate Strategy, distCardTimes, simulateTimesPerDist int, verbose bool) int {
	//stageFlag must in 0 1 2
	if c.StageFlag > 2 {
		panic("stageFlag > 2 ")
	}
	curMax := 0
	if c.StageFlag == 1 {
		curMax = c.PreerCent
	} else if c.StageFlag == 2 {
		curMax = c.PreerCent
		if curMax < c.NexterCent {
			curMax = c.NexterCent
		}
	}
	possibleCalls := []int{0}
	for i := curMax + 1; i <= 3; i++ {
		possibleCalls = append(possibleCalls, i)
	}
	if len(possibleCalls) == 1 {
		return possibleCalls[0]
	}
	profitDict := make([]int, len(possibleCalls))
	//通过随机仿真进行选择
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
	for distIdx := 0; distIdx < distCardTimes; distIdx++ {
		if !c.GameOpenCard {
			if c.CheatFlag == WithNexter && c.CheatMethod == CommInSecret {
				c.NexterRemainPokerSet = c.PartnerPokerSet
				//切上家牌与底牌
				_, c.PreerRemainPokerSet = RandDistCards2(action.PokerSet2(0x114444444444444).Subtract(c.RemainPokerSet).Subtract(c.PartnerPokerSet), 3, 17)
			} else if c.CheatFlag == WithPreer && c.CheatMethod == CommInSecret {
				c.PreerRemainPokerSet = c.PartnerPokerSet
				//切下家牌与底牌
				_, c.NexterRemainPokerSet = RandDistCards2(action.PokerSet2(0x114444444444444).Subtract(c.RemainPokerSet).Subtract(c.PartnerPokerSet), 3, 17)
			} else {
				//切上家下家与底牌
				c.PreerRemainPokerSet, _ = RandDistCards2(action.PokerSet2(0x114444444444444).Subtract(c.RemainPokerSet), 17, 20)
				_, c.NexterRemainPokerSet = RandDistCards2(action.PokerSet2(0x114444444444444).Subtract(c.RemainPokerSet).Subtract(c.PreerRemainPokerSet), 3, 17)
			}
		}
		for smIdx := 0; smIdx < simulateTimesPerDist; smIdx++ {
			for actionId, a := range possibleCalls {
				//simulator中当前player为0号
				simulator.RestoreGame(GameRestoreInfo{
					StageFlag: c.StageFlag + 1,
					Cents: [3]int{
						a,
						c.NexterCent,
						c.PreerCent,
					},
					OnTurn:           1,
					DizhuDipaiRemain: action.PokerSet2(0x114444444444444).Subtract(c.RemainPokerSet).Subtract(c.PreerRemainPokerSet).Subtract(c.NexterRemainPokerSet),
					DizhuPlayCnt:     0,
					NongminPlayCnt:   0,
					BombCnt:          0,
					RemainPokerSets: [3]action.PokerSet2{
						c.RemainPokerSet,
						c.NexterRemainPokerSet,
						c.PreerRemainPokerSet,
					},
				})
				simulator.RunRestoredGame()
				profits := simulator.GetProfits()
				profitDict[actionId] += profits[0]
				if c.CheatFlag == WithNexter { //with nexter
					profitDict[actionId] += profits[1]
				} else if c.CheatFlag == WithPreer { //with preer
					profitDict[actionId] += profits[2]
				}
			}
		}
	}
	//找出总收益最大的那个action
	res := possibleCalls[getMaxProfitIdx(profitDict)]
	return res
}

func getMaxProfitIdx(profitDict []int) int {
	maxProfitActionId := 0
	maxProfit := profitDict[0]
	for actionId, profit := range profitDict {
		if profit > maxProfit {
			maxProfit = profit
			maxProfitActionId = actionId
		}
	}
	return maxProfitActionId
}

func getTopKActionIDs(profitDict []int, k int) []int {
	a := make([][2]int, len(profitDict))
	for i := 0; i < len(profitDict); i++ {
		a[i] = [2]int{profitDict[i], i}
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i][0] > a[j][0]
	})
	res := []int{}
	for i := 0; i < k && i < len(a); i++ {
		res = append(res, a[i][1])
	}
	return res
}
