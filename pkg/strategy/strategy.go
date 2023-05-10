package strategy

import (
	"github.com/conggova/doudizhu-robot/pkg/action"
)

type TCheatFlag int
type TCheatMethod int

const (
	CommInSecret  TCheatMethod = iota //私下通信，互知底牌
	ShareInterest                     //利益共通
)

const (
	NoCheat TCheatFlag = iota
	WithPreer
	WithNexter
)

type PlayContext struct {
	DizhuFlag            int //0 self 1 nexter 2 preer
	PreerAction          action.Action
	NexterAction         action.Action
	PreerPkCnt           int
	NexterPkCnt          int
	DizhuDipaiRemain     action.PokerSet2
	DizhuPlayCnt         int
	NongminPlayCnt       int
	GameOpenCard         bool             //叫分时可用
	OthersRemainPokerSet action.PokerSet2 //必有
	PreerRemainPokerSet  action.PokerSet2 //GameOpenCard为true才有有效值 ， 叫分时可用
	NexterRemainPokerSet action.PokerSet2 //GameOpenCard为true才有有效值 ， 叫分时可用
}

type CallContext struct {
	StageFlag            int              //0 首叫 ， 1 2 。。。
	PreerCent            int              //上家叫分 如有
	NexterCent           int              //下家叫分 如有
	GameOpenCard         bool             //叫分时可用
	PreerRemainPokerSet  action.PokerSet2 //GameOpenCard为true才有有效值 ， 叫分时可用
	NexterRemainPokerSet action.PokerSet2 //GameOpenCard为true才有有效值 ， 叫分时可用
}

type StrategyPlayContext struct {
	RemainPokerSet action.PokerSet2 //
	CheatFlag      TCheatFlag
	CheatMethod    TCheatMethod
	OpenCard4Me    bool //为True时PreerRemainPokerSet和NexterRemainPokerSet必有有效值
	PlayContext
}

type StrategyCallContext struct {
	RemainPokerSet  action.PokerSet2 //
	CheatFlag       TCheatFlag
	CheatMethod     TCheatMethod
	PartnerPokerSet action.PokerSet2 //通信做弊时有可值
	CallContext
}

type GameBuildInfo struct {
	Strategys   [3]Strategy
	OpenCard    bool
	Partners    [2]int
	CheatMethod TCheatMethod
}

type GameRestoreInfo struct {
	StageFlag        int              //牌局到了哪个阶段 0 未开始  1 首叫完成 2 第二家叫分完成  3 第三家叫分完成  4正在打牌
	Cents            [3]int           //各家叫的分, 隐含地主信息
	LastAction       action.Action    //最后出的一手牌是什么
	LLastAction      action.Action    //上上手牌
	OnTurn           int              //该谁了
	DizhuDipaiRemain action.PokerSet2 //地主还没出的底牌
	DizhuPlayCnt     int
	NongminPlayCnt   int
	BombCnt          int
	RemainPokerSets  [3]action.PokerSet2 //各家余牌
}

// 为不让Strategy不直接依赖Game（会出现循环依赖），让Game实现Simulator接口
type Simulator interface {
	BuildGame(GameBuildInfo) Simulator //通过定义在接口里的构造函数，其实实现了循环依赖
	RestoreGame(GameRestoreInfo)
	RunRestoredGame()
	GetProfits() [3]int
}

type Strategy interface {
	MakeDecisionWhenCall(c StrategyCallContext) int
	MakeDecisionWhenPlay(c StrategyPlayContext) action.Action
}
