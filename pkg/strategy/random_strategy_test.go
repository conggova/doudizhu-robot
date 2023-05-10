package strategy

import (
	"testing"

	"github.com/conggova/doudizhu-robot/pkg/action"
)

func TestRondomStrategy(t *testing.T) {
	s := RandomStrategy{}
	c := StrategyCallContext{
		RemainPokerSet:  0,
		CheatFlag:       0,
		CheatMethod:     0,
		PartnerPokerSet: 0,
		CallContext: CallContext{
			StageFlag:            2,
			PreerCent:            3,
			NexterCent:           0,
			GameOpenCard:         false,
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}
	cent := s.MakeDecisionWhenCall(c)
	if cent != 0 {
		t.Error("incorrect")
	}
	pc := StrategyPlayContext{
		RemainPokerSet: 0x323,
		CheatFlag:      0,
		CheatMethod:    0,
		OpenCard4Me:    false,
		PlayContext: PlayContext{
			DizhuFlag:   0,
			PreerAction: action.NewActionWithoutAff(action.Single, 1, 0),
			NexterAction: [2]uint64{
				0,
				0,
			},
			PreerPkCnt:           0,
			NexterPkCnt:          0,
			DizhuDipaiRemain:     0,
			DizhuPlayCnt:         0,
			NongminPlayCnt:       0,
			GameOpenCard:         false,
			OthersRemainPokerSet: 0,
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}
	a := s.MakeDecisionWhenPlay(pc)
	t.Errorf(a.String())
	//fmt.Println(a)
}

func TestRondomStrategy2(t *testing.T) {
	s := RandomStrategy2{3}
	c := StrategyCallContext{
		RemainPokerSet:  0,
		CheatFlag:       0,
		CheatMethod:     0,
		PartnerPokerSet: 0,
		CallContext: CallContext{
			StageFlag:            2,
			PreerCent:            3,
			NexterCent:           0,
			GameOpenCard:         false,
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}
	cent := s.MakeDecisionWhenCall(c)
	if cent != 0 {
		t.Error("incorrect")
	}
	pc := StrategyPlayContext{
		RemainPokerSet: 0x323,
		CheatFlag:      0,
		CheatMethod:    0,
		OpenCard4Me:    false,
		PlayContext: PlayContext{
			DizhuFlag:   0,
			PreerAction: action.NewActionWithoutAff(action.Pass, 0, 0),
			NexterAction: [2]uint64{
				0,
				0,
			},
			PreerPkCnt:           0,
			NexterPkCnt:          0,
			DizhuDipaiRemain:     0,
			DizhuPlayCnt:         0,
			NongminPlayCnt:       0,
			GameOpenCard:         false,
			OthersRemainPokerSet: 0,
			PreerRemainPokerSet:  0,
			NexterRemainPokerSet: 0,
		},
	}
	a := s.MakeDecisionWhenPlay(pc)
	t.Errorf(a.String())
	//fmt.Println(a)
}
