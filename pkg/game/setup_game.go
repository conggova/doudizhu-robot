package game

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/conggova/doudizhu-robot/pkg/player"
	"github.com/conggova/doudizhu-robot/pkg/strategy"
)

// 人机对战
func ManAndRobot() {
	var ipt string
	gameOpenCard := false
	fmt.Println("您希望游戏是明牌吗？（大家都能看到彼此的牌）")
	fmt.Print("Y for yes , others for no :")
	ipt = "N"
	fmt.Scanln(&ipt)
	if ipt == "Y" {
		gameOpenCard = true
	}

	othersAreParterners := false
	othersCom := false
	fmt.Println("您希望另外两个玩家是同伙吗？（利益共同体）")
	fmt.Print("Y for yes , others for no :")
	ipt = "N"
	fmt.Scanln(&ipt)
	if ipt == "Y" {
		othersAreParterners = true
		if !gameOpenCard {
			fmt.Println("您希望他们能够私下通信吗?（能相互看牌）")
			fmt.Print("Y for yes , others for no :")
			ipt = "N"
			fmt.Scanln(&ipt)
			if ipt == "Y" {
				othersCom = true
			}
		}
	}

	gameTotal := 0
	fmt.Print("您想玩几局？")
	fmt.Scanf("%d\n", &gameTotal)
	headPlayer := rand.Intn(3)
	manTotalProfit := 0
	for i := 0; i < gameTotal; i++ {
		manProfit := manAndRobotPlayOneGame(gameOpenCard, othersAreParterners, othersCom, headPlayer)
		fmt.Println("这一局您的收益是 ", manProfit)
		manTotalProfit += manProfit
		fmt.Println("目前您的总收益是 ", manTotalProfit)
		fmt.Println("")
		headPlayer = (headPlayer + 1) % 3
	}
	fmt.Println("你一共玩了 ", gameTotal, " 局牌， 总收益是 ", manTotalProfit, " 。 ")
	fmt.Println("")
	fmt.Println("")

	fmt.Println("您要再玩几局吗？")
	fmt.Print("Y for yes , others for no :")
	ipt = "N"
	fmt.Scanln(&ipt)
	if ipt == "Y" {
		ManAndRobot()
	}
}

func manAndRobotPlayOneGame(gameOpenCard bool, robotsArePartners bool, robotsCom bool, headPlayer int) int {
	var player1 player.Player = player.NewManualPlayer(0)
	var player2 player.Player = player.NewStrategyPlayer(1,
		strategy.TwoStageSimulateStrategy{
			Simulator:               &Game{},
			DistCardTimesWhenRecall: 100,
			GamesPerDistWhenRecall:  1,
			StrategyWhenRecall:      strategy.RandomStrategy{},
			RecallNum:               20,
			DistCardTimesWhenRank:   10,
			GamesPerDistWhenRank:    2,
			StrategyWhenRank:        strategy.SimulateStrategy{DistCardTimes: 5, SimulateTimesPerDist: 1, Simulator: &Game{}, StrategyWhenSimulate: strategy.RandomStrategy{}},
		},
		strategy.NoCheat, strategy.ShareInterest, nil)

	var player3 player.Player = player.NewStrategyPlayer(2,
		strategy.UCB1SimulateStrategy{
			Simulator:            &Game{},
			StrategyWhenSimulate: strategy.RandomStrategy{},
			AvgTryCnt:            50,
		},
		strategy.NoCheat, strategy.ShareInterest, nil)

	if robotsArePartners {
		player2.SetCheatFlag(strategy.WithNexter)
		player3.SetCheatFlag(strategy.WithPreer)
		if robotsCom {
			player2.SetCheatMethod(strategy.CommInSecret)
			player3.SetCheatMethod(strategy.CommInSecret)
		}
		player2.SetPartner(player3)
		player3.SetPartner(player2)
	}

	game := &Game{players: [3]player.Player{player1, player2, player3}, openCard: gameOpenCard, keepLog: true}
	game.RunGameWithRandomBeginning(headPlayer)
	manProfit := game.GetProfits()[0]
	return manProfit
}

// 机器人大战
func RobotFight() {
	headPlayer := rand.Intn(3)
	totalProfits := [3]int{}
	for i := 1; i < 10000; i++ {
		fmt.Println("第 ", i, " 局")
		profits := robotFightPlayOneGame(false, headPlayer)
		for i := 0; i < 3; i++ {
			totalProfits[i] += profits[i]
		}
		fmt.Print("\n")
		fmt.Print("总收益：\n")
		fmt.Print("	刘备 ")
		fmt.Print(strconv.Itoa(totalProfits[0]))
		fmt.Print("\n")
		fmt.Print("	关羽 ")
		fmt.Print(strconv.Itoa(totalProfits[1]))
		fmt.Print("\n")
		fmt.Print("	张飞 ")
		fmt.Print(strconv.Itoa(totalProfits[2]))
		fmt.Print("\n")
		fmt.Print("\n")
		fmt.Print("回车继续... ")
		var ipt string
		fmt.Scanln(&ipt)
		headPlayer = (headPlayer + 1) % 3
	}
}

func robotFightPlayOneGame(opencard bool, headPlayer int) [3]int {
	var player1 player.Player = player.NewStrategyPlayer(0,
		strategy.TwoStageSimulateStrategy{
			Simulator:               &Game{},
			DistCardTimesWhenRecall: 100,
			GamesPerDistWhenRecall:  1,
			StrategyWhenRecall:      strategy.RandomStrategy{},
			RecallNum:               20,
			DistCardTimesWhenRank:   10,
			GamesPerDistWhenRank:    2,
			StrategyWhenRank:        strategy.SimulateStrategy{DistCardTimes: 5, SimulateTimesPerDist: 1, Simulator: &Game{}, StrategyWhenSimulate: strategy.RandomStrategy{}},
		},
		strategy.NoCheat, strategy.ShareInterest, nil)
	var player2 player.Player = player.NewStrategyPlayer(1,
		strategy.UCB1SimulateStrategy{
			Simulator:            &Game{},
			StrategyWhenSimulate: strategy.RandomStrategy{},
			AvgTryCnt:            50,
		},
		strategy.NoCheat, strategy.ShareInterest, nil)
	var player3 player.Player = player.NewStrategyPlayer(2,
		strategy.UCB1SimulateStrategy{
			Simulator:            &Game{},
			StrategyWhenSimulate: strategy.RandomStrategy{},
			AvgTryCnt:            50,
		},
		strategy.NoCheat, strategy.ShareInterest, nil)
	game := &Game{players: [3]player.Player{player1, player2, player3}, openCard: opencard, keepLog: true}
	game.RunGameWithRandomBeginning(headPlayer)
	return game.GetProfits()
}
