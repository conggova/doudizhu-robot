package strategy

import (
	"github.com/conggova/doudizhu-robot/pkg/action"
)

// 先召回，再排序
type TwoStageSimulateStrategy struct {
	Simulator               Simulator
	DistCardTimesWhenRecall int
	GamesPerDistWhenRecall  int
	StrategyWhenRecall      Strategy
	RecallNum               int
	DistCardTimesWhenRank   int
	GamesPerDistWhenRank    int
	StrategyWhenRank        Strategy
}

func (s TwoStageSimulateStrategy) MakeDecisionWhenCall(c StrategyCallContext) int {
	return simulateCall(c, s.Simulator, SimulateStrategy{
		DistCardTimes:        5,
		SimulateTimesPerDist: 1,
		Simulator:            s.Simulator,
		StrategyWhenSimulate: RandomStrategy{},
	}, 10, 1, false)
}

func (s TwoStageSimulateStrategy) MakeDecisionWhenPlay(c StrategyPlayContext) (actionTaken action.Action) {
	contextAction := c.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = c.NexterAction
	}
	possibleActionList := c.RemainPokerSet.PossibleActionsWithContext(contextAction)
	if len(possibleActionList) == 1 {
		return possibleActionList[0]
	}
	if len(possibleActionList) > s.RecallNum {
		possibleActionList = simulateReduceAction(possibleActionList, c, s.Simulator, s.StrategyWhenRecall,
			s.DistCardTimesWhenRecall, s.GamesPerDistWhenRecall, s.RecallNum, false)
	}
	return simulateReduceAction(possibleActionList, c, s.Simulator, s.StrategyWhenRank, s.DistCardTimesWhenRank, s.GamesPerDistWhenRank, 1, false)[0]
}
