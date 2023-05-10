package game

import (
	"fmt"
	"testing"

	"github.com/conggova/doudizhu-robot/pkg/action"
	"github.com/conggova/doudizhu-robot/pkg/strategy"
)

func TestSimulateStrategyCallWithRandom(t *testing.T) {
	s := strategy.SimulateStrategy{
		DistCardTimes:        5,
		SimulateTimesPerDist: 1,
		Simulator:            &Game{},
		StrategyWhenSimulate: strategy.RandomStrategy{},
		Verbose:              true,
	}

	_, ps1 := action.ParsePokerSet2("56889TJQQKKAAA22X")
	//_, ps2 := action.ParsePokerSet2("3334456799TJJQKKA")
	//_, ps3 := action.ParsePokerSet2("34455667889TTJQ2D")
	//_, ps4 := action.ParsePokerSet2("772")

	c := strategy.StrategyCallContext{
		RemainPokerSet:  ps1,
		CheatFlag:       0,
		CheatMethod:     0,
		PartnerPokerSet: 0,
		CallContext: strategy.CallContext{
			StageFlag:            1,
			PreerCent:            1,
			NexterCent:           2,
			GameOpenCard:         false,
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}

	cent := s.MakeDecisionWhenCall(c)
	fmt.Println("cent : ", cent)
	if cent != -1 {
		t.Error("incorrect")
	}

}

func TestSimulateStrategyPlayWithRandom(t *testing.T) {
	s := strategy.SimulateStrategy{
		DistCardTimes:        100,
		SimulateTimesPerDist: 1,
		Simulator:            &Game{},
		StrategyWhenSimulate: strategy.RandomStrategy2{Times: 3},
		Verbose:              true,
	}
	// 500*5结果同样不稳定 ， 与500 * 1 相差不多 ， 前20都还行，可以用作召回

	_, ps1 := action.ParsePokerSet2("56889TJQQKKAAA22X")
	_, ps2 := action.ParsePokerSet2("3334456799TJJQKKA")
	_, ps3 := action.ParsePokerSet2("34455667889TTJQ2D")
	_, ps4 := action.ParsePokerSet2("772")

	pc := strategy.StrategyPlayContext{
		RemainPokerSet: ps1.CombineWith(ps4),
		CheatFlag:      0,
		CheatMethod:    0,
		OpenCard4Me:    false,
		PlayContext: strategy.PlayContext{
			DizhuFlag:            0,
			PreerAction:          action.NewActionWithoutAff(action.Pass, 0, 0),
			NexterAction:         action.NewActionWithoutAff(action.Pass, 0, 0),
			PreerPkCnt:           17,
			NexterPkCnt:          17,
			DizhuDipaiRemain:     ps4,
			DizhuPlayCnt:         0,
			NongminPlayCnt:       0,
			GameOpenCard:         false,
			OthersRemainPokerSet: ps2.CombineWith(ps3),
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}
	a := s.MakeDecisionWhenPlay(pc)
	fmt.Println("res : ", a)
	if a != action.NewActionWithoutAff(action.Pass, 0, 0) {
		t.Error("incorrect")
	}
}

func TestSimulateStrategyCallRecur(t *testing.T) {
	s := strategy.SimulateStrategy{
		DistCardTimes:        10,
		SimulateTimesPerDist: 1,
		Simulator:            &Game{},
		StrategyWhenSimulate: strategy.SimulateStrategy{
			DistCardTimes:        5,
			SimulateTimesPerDist: 1,
			Simulator:            &Game{},
			StrategyWhenSimulate: strategy.RandomStrategy2{Times: 3},
		},
		Verbose: true,
	}

	_, ps1 := action.ParsePokerSet2("56889TJQQKKAAA22X")
	//_, ps2 := action.ParsePokerSet2("3334456799TJJQKKA")

	//_, ps3 := action.ParsePokerSet2("34455667889TTJQ2D")
	/*		_, ps4 := action.ParsePokerSet2("772")
	 */
	c := strategy.StrategyCallContext{
		RemainPokerSet:  ps1,
		CheatFlag:       0,
		CheatMethod:     0,
		PartnerPokerSet: 0,
		CallContext: strategy.CallContext{
			StageFlag:            1,
			PreerCent:            1,
			NexterCent:           2,
			GameOpenCard:         false,
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}

	cent := s.MakeDecisionWhenCall(c)
	fmt.Println("cent : ", cent)
	if cent != -1 {
		t.Error("incorrect")
	}
}

func TestSimulateStrategyPlayRecur(t *testing.T) {
	s := strategy.SimulateStrategy{
		DistCardTimes:        10,
		SimulateTimesPerDist: 1,
		Simulator:            &Game{},
		StrategyWhenSimulate: strategy.SimulateStrategy{
			DistCardTimes:        5,
			SimulateTimesPerDist: 1,
			Simulator:            &Game{},
			StrategyWhenSimulate: strategy.RandomStrategy{},
		},
		Verbose: true,
	}

	_, ps1 := action.ParsePokerSet2("56889TJQQKKAAA22X")
	_, ps2 := action.ParsePokerSet2("3334456799TJJQKKA")
	_, ps3 := action.ParsePokerSet2("34455667889TTJQ2D")
	_, ps4 := action.ParsePokerSet2("772")

	pc := strategy.StrategyPlayContext{
		RemainPokerSet: ps1.CombineWith(ps4),
		CheatFlag:      0,
		CheatMethod:    0,
		OpenCard4Me:    false,
		PlayContext: strategy.PlayContext{
			DizhuFlag:            0,
			PreerAction:          action.NewActionWithoutAff(action.Pass, 0, 0),
			NexterAction:         action.NewActionWithoutAff(action.Pass, 0, 0),
			PreerPkCnt:           17,
			NexterPkCnt:          17,
			DizhuDipaiRemain:     ps4,
			DizhuPlayCnt:         0,
			NongminPlayCnt:       0,
			GameOpenCard:         false,
			OthersRemainPokerSet: ps2.CombineWith(ps3),
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}
	a := s.MakeDecisionWhenPlay(pc)
	fmt.Println("res : ", a)
	if a != action.NewActionWithoutAff(action.Pass, 0, 0) {
		t.Error("incorrect")
	}
}

func TestTwoStageSimulateStrategyPlay(t *testing.T) {
	s := strategy.TwoStageSimulateStrategy{
		Simulator:               &Game{},
		DistCardTimesWhenRecall: 100,
		GamesPerDistWhenRecall:  1,
		StrategyWhenRecall:      strategy.RandomStrategy{},
		RecallNum:               20,
		DistCardTimesWhenRank:   10,
		GamesPerDistWhenRank:    2,
		StrategyWhenRank: strategy.SimulateStrategy{
			DistCardTimes:        5,
			SimulateTimesPerDist: 1,
			Simulator:            &Game{},
			StrategyWhenSimulate: strategy.RandomStrategy{},
		},
	}
	_, ps1 := action.ParsePokerSet2("56889TJQQKKAAA22X")
	_, ps2 := action.ParsePokerSet2("3334456799TJJQKKA")
	_, ps3 := action.ParsePokerSet2("34455667889TTJQ2D")
	_, ps4 := action.ParsePokerSet2("772")

	pc := strategy.StrategyPlayContext{
		RemainPokerSet: ps1.CombineWith(ps4),
		CheatFlag:      0,
		CheatMethod:    0,
		OpenCard4Me:    false,
		PlayContext: strategy.PlayContext{
			DizhuFlag:            0,
			PreerAction:          action.NewActionWithoutAff(action.Pass, 0, 0),
			NexterAction:         action.NewActionWithoutAff(action.Pass, 0, 0),
			PreerPkCnt:           17,
			NexterPkCnt:          17,
			DizhuDipaiRemain:     ps4,
			DizhuPlayCnt:         0,
			NongminPlayCnt:       0,
			GameOpenCard:         false,
			OthersRemainPokerSet: ps2.CombineWith(ps3),
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}
	a := s.MakeDecisionWhenPlay(pc)
	fmt.Println("res : ", a)
	if a != action.NewActionWithoutAff(action.Pass, 0, 0) {
		t.Error("incorrect")
	}
}

func TestUCB1StrategyPlay(t *testing.T) {
	s := strategy.UCB1SimulateStrategy{
		Simulator: &Game{},
		StrategyWhenSimulate: strategy.SimulateStrategy{
			DistCardTimes:        1,
			SimulateTimesPerDist: 1,
			Simulator:            &Game{},
			StrategyWhenSimulate: strategy.RandomStrategy{},
		},
		AvgTryCnt: 25,
	}
	_, ps1 := action.ParsePokerSet2("56889TJQQKKAAA22X")
	_, ps2 := action.ParsePokerSet2("3334456799TJJQKKA")
	_, ps3 := action.ParsePokerSet2("34455667889TTJQ2D")
	_, ps4 := action.ParsePokerSet2("772")

	pc := strategy.StrategyPlayContext{
		RemainPokerSet: ps1.CombineWith(ps4),
		CheatFlag:      0,
		CheatMethod:    0,
		OpenCard4Me:    false,
		PlayContext: strategy.PlayContext{
			DizhuFlag:            0,
			PreerAction:          action.NewActionWithoutAff(action.Pass, 0, 0),
			NexterAction:         action.NewActionWithoutAff(action.Pass, 0, 0),
			PreerPkCnt:           17,
			NexterPkCnt:          17,
			DizhuDipaiRemain:     ps4,
			DizhuPlayCnt:         0,
			NongminPlayCnt:       0,
			GameOpenCard:         false,
			OthersRemainPokerSet: ps2.CombineWith(ps3),
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}
	a := s.MakeDecisionWhenPlay(pc)
	fmt.Println("res : ", a)
	if a != action.NewActionWithoutAff(action.Pass, 0, 0) {
		t.Error("incorrect")
	}
}
