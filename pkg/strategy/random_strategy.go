package strategy

import (
	"math/rand"

	"github.com/conggova/doudizhu-robot/pkg/action"
)

// 完全随机出牌，不考虑明牌和作弊信息
type RandomStrategy struct {
}

// 随机叫分
func (s RandomStrategy) MakeDecisionWhenCall(c StrategyCallContext) int {
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
	t := rand.Intn(4)
	if t > curMax {
		return t
	} else {
		return 0
	}
}

// =================================================
// 随机策略出牌 ， 在不放走的情况下
// -------------------------------------------------
// 如果下家是一张，要考虑不放走，既保证在随机牌局中不存在放走包赔
func (s RandomStrategy) MakeDecisionWhenPlay(c StrategyPlayContext) (actionTaken action.Action) {
	contextAction := c.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = c.NexterAction
	}
	//defer func() { actionTaken.SetContext(contextAction) }()
	possibleActionList := c.RemainPokerSet.PossibleActionsWithContext(contextAction)
	if len(possibleActionList) == 1 {
		return possibleActionList[0]
	}
	return possibleActionList[rand.Intn(len(possibleActionList))]
}

// 完全随机出牌，不考虑明牌和作弊信息
// 随机选Times次，最终选张数最多的
type RandomStrategy2 struct {
	Times int
}

// 随机叫分
func (s RandomStrategy2) MakeDecisionWhenCall(c StrategyCallContext) int {
	return RandomStrategy{}.MakeDecisionWhenCall(c)
}

// =================================================
// 随机策略出牌 ，
func (s RandomStrategy2) MakeDecisionWhenPlay(c StrategyPlayContext) (actionTaken action.Action) {
	contextAction := c.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = c.NexterAction
	}
	//defer func() { actionTaken = actionTaken.SetContext(contextAction) }()
	possibleActionList := c.RemainPokerSet.PossibleActionsWithContext(contextAction)
	if len(possibleActionList) == 1 {
		return possibleActionList[0]
	}
	if contextAction.ActionType() == action.Pass { //可以出任意牌
		res := possibleActionList[rand.Intn(len(possibleActionList))]
		for i := 1; i < s.Times; i++ {
			t := possibleActionList[rand.Intn(len(possibleActionList))]
			if t.PokerCount() > res.PokerCount() {
				res = t
			}
		}
		return res
	} else {
		return possibleActionList[rand.Intn(len(possibleActionList))]
	}
}
