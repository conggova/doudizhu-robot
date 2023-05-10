package player

import (
	"fmt"

	"github.com/conggova/doudizhu-robot/pkg/action"
	"github.com/conggova/doudizhu-robot/pkg/strategy"
)

// 人类玩家,不参与作弊
type ManualPlayer struct {
	BasePlayer
}

func NewManualPlayer(playerId int) *ManualPlayer {
	return &ManualPlayer{BasePlayer{Id: playerId}}
}

func (p *ManualPlayer) Call(pc strategy.CallContext) int {
	fmt.Print("您的角色是刘备 ，现在该您叫分。")
	fmt.Println("您的手牌是 ", p.remainPokerSet)
	maxCent := pc.PreerCent
	if pc.NexterCent > maxCent {
		maxCent = pc.NexterCent
	}
	if maxCent == 3 {
		fmt.Print("当前已叫满三分，您只能不叫。回车继续......")
		var appLine string
		fmt.Scanln(&appLine)
		return 0
	}
	centMap := map[string]int{"0": 0, "1": 1, "2": 2, "3": 3}

	for {
		var ipt string
		fmt.Print("请输入（0不叫，1叫一分，2叫二分，3叫三分）:")
		fmt.Scanln(&ipt)
		if cent, ok := centMap[ipt]; ok {

			if cent != 0 && cent <= maxCent {
				fmt.Print("叫分不符合规则. 请重试!")
			} else {
				fmt.Print("您的输入是:", cent, " , R可重新输入 , 其它则确认 :")
				var ynInput string
				fmt.Scanln(&ynInput)
				confirm := true
				if ynInput == "R" {
					confirm = false
				}

				//只有输入正确 ， 并且确认的情况下
				if confirm {
					return cent
				} else {
					fmt.Print("已重置。")
				}
			}
		} else {
			fmt.Print("输入不符合规则，请重新输入。 ")
		}
	}
}

func (p *ManualPlayer) Play(pc strategy.PlayContext) (actionTaken action.Action) {
	var contextAction = pc.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = pc.NexterAction
	}
	defer func() {
		p.remainPokerSet = p.remainPokerSet.Subtract(actionTaken.PokerSet2())
	}()
	//有走必走
	possibleActionList := p.remainPokerSet.PossibleActionsWithContext(contextAction)
	//如果只有一种选择
	fmt.Print("您的角色是刘备。现在该您出牌。")
	if len(possibleActionList) == 1 {
		fmt.Print("您的手牌是 ", p.remainPokerSet, " , 唯一的选择是 : ", possibleActionList[0], " , 回车继续......")
		var appLine string
		fmt.Scanln(&appLine)
		return possibleActionList[0]
	}

	for {
		fmt.Print("您的手牌是 ", p.remainPokerSet, "  ,（输入H打印帮助信息）:")
		var ipt string
		fmt.Scanln(&ipt)

		//打印帮助
		if ipt == "H" {
			printHelpInfo()
			continue
		}

		//检查输入
		ok, checkResult := checkInput(ipt, p.remainPokerSet, contextAction)
		//如果是过 或者输入有效
		if ok {
			fmt.Println("您的输入是 :", checkResult, "  R可重新输入 , 其它则确认 :")
			var ynInput string
			fmt.Scanln(&ynInput)
			confirm := true
			if ynInput == "R" {
				confirm = false
			}

			//只有输入正确 ， 并且确认的情况下
			if confirm {
				return checkResult
			} else {
				fmt.Print("已重置。")
			}

		} else { //输入不符合规则
			fmt.Print("输入不符合规则，请重新输入。 ")
		}
	}
}

func checkInput(ipt string, ps action.PokerSet2, contextAction action.Action) (bool, action.Action) {
	//先检查输入是否符合规则
	ok, a := action.ParseAction(ipt)
	if !ok {
		return false, action.Action{}
	}
	//再检查输入是否符合背景
	if !a.PlayableInContext(contextAction) {
		return false, action.Action{}
	}
	//再检查是否能提供
	if !a.CanBeSupliedBy(ps) {
		return false, action.Action{}
	}
	return true, a
}

func printHelpInfo() {
	fmt.Println("#--------------------------------------#")
	fmt.Println("帮助信息 ：")
	fmt.Println("每张牌的表示：3456789TJQKA2XD，特别注意T代表10,X代表小王，D代表大王")
	fmt.Println("以上字母不分大小写")
	fmt.Println("各种牌形的介绍如下：")
	fmt.Println("过牌(不要): P")
	fmt.Println("单牌: J(一个J)")
	fmt.Println("对子: TT(对10)")
	fmt.Println("三条: 999(三个9)")
	fmt.Println("三带一: 999K(三个9带K)")
	fmt.Println("三带对: QQQ33(三个Q带对3)")
	fmt.Println("炸弹: KKKK")
	fmt.Println("四带二: 777789(四个7带8和9)")
	fmt.Println("四带两对: 77778899(四个7带对8和对9)")
	fmt.Println("顺子: 3456789TJQKA(顺子3到A)")
	fmt.Println("双顺: 667788(双顺6到8)")
	fmt.Println("飞机不带: 333444(333444)")
	fmt.Println("飞机带单: 33344456(333444带56)")
	fmt.Println("飞机带对: 3334445566(333444带5566)")
	fmt.Println("#--------------------------------------#")
}
