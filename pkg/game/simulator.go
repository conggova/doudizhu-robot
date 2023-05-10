package game

import (
	"github.com/conggova/doudizhu-robot/pkg/player"
	"github.com/conggova/doudizhu-robot/pkg/strategy"
)

// 载入Player信息
func (*Game) BuildGame(gbi strategy.GameBuildInfo) strategy.Simulator {
	g := &Game{players: [3]player.Player{
		player.NewStrategyPlayer(0, gbi.Strategys[0], strategy.NoCheat, strategy.ShareInterest, nil),
		player.NewStrategyPlayer(1, gbi.Strategys[1], strategy.NoCheat, strategy.ShareInterest, nil),
		player.NewStrategyPlayer(2, gbi.Strategys[2], strategy.NoCheat, strategy.ShareInterest, nil)},
		openCard: gbi.OpenCard, keepLog: false}
	//设置Player与打牌无关的数据
	for i := 0; i <= 2; i++ {
		//有两方是一伙
		if gbi.Partners != [2]int{} && (i == gbi.Partners[0] || i == gbi.Partners[1]) {
			//calcu cheatFlag
			partnerId := gbi.Partners[0]
			if i == gbi.Partners[0] {
				partnerId = gbi.Partners[1]
			}
			partnerFlag := strategy.WithPreer
			if (i+1)%3 == partnerId {
				partnerFlag = strategy.WithNexter
			}
			g.players[i].SetCheatFlag(partnerFlag)
			g.players[i].SetCheatMethod(gbi.CheatMethod)
			g.players[i].SetPartner(g.players[partnerId])
		} else {
			g.players[i].SetCheatFlag(strategy.NoCheat)
		}
	}
	return g
}

// 只载入牌局
func (g *Game) RestoreGame(gri strategy.GameRestoreInfo) {
	g.resetGame()
	//g.keepLog = true
	g.stageFlag = gri.StageFlag
	g.cents = gri.Cents
	g.lastAction = gri.LastAction
	g.llastAction = gri.LLastAction
	g.onTurn = gri.OnTurn
	g.dipaiRemain = gri.DizhuDipaiRemain
	g.dizhuPlayCnt = gri.DizhuPlayCnt
	g.nongminPlayCnt = gri.NongminPlayCnt
	g.bombCnt = gri.BombCnt
	for i := 0; i < 3; i++ {
		g.players[i].SetRemainPokerSet(gri.RemainPokerSets[i])
	}
}

// 指定上下文开始游戏
func (g *Game) RunRestoredGame() {
	g.runGame()
}
