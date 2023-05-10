package action

func (a Action) CanBeSupliedBy(ps2 PokerSet2) bool {
	return ps2.Covers(a.PokerSet2())
}

func ParseAction(ipt string) (bool, Action) {
	//先将ipt做成statMap , 然后看看能不能一手出完
	pkCnt := len(ipt)
	if ipt == "P" || ipt == "p" {
		return true, [2]uint64{}
	}
	if ok, ps2 := ParsePokerSet2(ipt); ok {
		actionList := ps2.possibleActionsWithoutContext()
		if len(actionList) == 1 {
			return true, actionList[0]
		} else {
			for _, a := range actionList {
				if a.PokerCount() == pkCnt {
					return true, a
				}
			}
		}
	}
	return false, Action{}

}

// seq: [2]uint64{base , len}
func getSubSeqsWithMinLen(seq [2]uint64, minLen uint64) [][2]uint64 {
	subSeqs := [][2]uint64{}
	base, totalLen := seq[0], seq[1]
	var start, len uint64
	for start = 0; start < totalLen-minLen+1; start++ {
		for len = minLen; len < totalLen-start+1; len++ {
			subSeqs = append(subSeqs, [2]uint64{base + start, len})
		}
	}
	return subSeqs
}

// seq: [2]uint64{base , len}
func getSubSeqsWithFixLen(seq [2]uint64, fixLen uint64) [][2]uint64 {
	subSeqs := [][2]uint64{}
	base, totalLen := seq[0], seq[1]
	var start uint64
	for start = 0; start < totalLen-fixLen+1; start++ {
		subSeqs = append(subSeqs, [2]uint64{base + start, fixLen})
	}
	return subSeqs
}

// 每个牌张只能取一次
func getAffsCouple(ps PokerSet2, length int) []PokerSet2 {
	if length == 0 {
		return []PokerSet2{}
	}
	var t PokerSet2 = 0b1111
	ps1 := ps
	validCards := []PokerSet2{}
	var psAll PokerSet2
	for power := 0; power < 15; power++ {
		cnt := int(t & ps1)
		if cnt >= 2 {
			var tmp PokerSet2 = 0b10 << (power * 4)
			validCards = append(validCards, tmp)
			psAll = psAll.CombineWith(tmp)
		}
		ps1 >>= 4
	}
	cCnt := len(validCards)

	if length > cCnt {
		return []PokerSet2{}
	} else if cCnt == length {
		return []PokerSet2{psAll}
	}
	//如果length大于cCnt的一半则取不要的部分
	reverse := false
	if length > cCnt/2 {
		reverse = true
		length = cCnt - length
	}
	var affs []PokerSet2
	if length == 1 {
		affs = append(affs, validCards...)
	} else {
		lenAffsMap := make([][]PokerSet2, length) //0 , 1 , 2 , ... length
		lenAffsMap[0] = append(lenAffsMap[0], 0)
		for _, power := range validCards {
			//take add lenaff wont bigger the length
			for l := length - 1; l >= 0; l-- {
				if l == length-1 {
					for _, item := range lenAffsMap[l] {
						affs = append(affs, item.CombineWith(power))
					}
				} else {
					for _, item := range lenAffsMap[l] {
						lenAffsMap[l+1] = append(lenAffsMap[l+1], item.CombineWith(power))
					}
				}
			}
		}
	}
	if reverse {
		res := []PokerSet2{}
		for _, item := range affs {
			res = append(res, psAll.Subtract(item))
		}
		return res
	} else {
		return affs
	}
}

// psMain为主牌的PokerSet，加入配牌不得与主牌相同的逻辑
func getAffsSingle(ps, psMain PokerSet2, length int) []PokerSet2 {
	if length == 0 {
		return []PokerSet2{}
	}
	var t PokerSet2 = 0b1111
	for power := 0; power < 15; power++ {
		cnt := int(t & psMain)
		if cnt > 0 {
			ps &= ^(0b1111 << (power << 2))
		}
		psMain >>= 4

	}
	pkCnt := ps.PokerCount()

	if length > pkCnt {
		return []PokerSet2{}
	} else if pkCnt == length {
		return []PokerSet2{ps}
	}
	//如果length大于pkCnt的一半则取不要的部分
	reverse := false
	if length > pkCnt/2 {
		reverse = true
		length = pkCnt - length
	}
	var affs []PokerSet2
	ps1 := ps
	if length == 1 {
		for power := 0; power < 15; power++ {
			if t&ps1 > 0 {
				affs = append(affs, 1<<(power<<2))
			}
			ps1 >>= 4
		}
	} else {
		lenAffsMap := make([][]PokerSet2, length) //0 , 1 , 2 , ... length
		lenAffsMap[0] = append(lenAffsMap[0], 0)
		for power := 0; power < 15; power++ {
			cnt := int(t & ps1)
			if cnt == 4 { //cant take all 4
				cnt = 3
			}
			ps1 >>= 4
			tmp := make([][]PokerSet2, length) //store result this power gens
			for take := 1; take <= cnt && take <= length; take++ {
				psTake := PokerSet2(take) << (power << 2)
				if take == length { //short cut
					affs = append(affs, psTake)
					continue
				}
				//take add lenaff wont bigger the length
				for l := 0; l+take <= length; l++ {
					if l+take == length {
						for _, item := range lenAffsMap[l] {
							affs = append(affs, item.CombineWith(psTake))
						}
					} else {
						for _, item := range lenAffsMap[l] {
							tmp[take+l] = append(tmp[take+l], item.CombineWith(psTake))
						}
					}
				}
			}
			//merge tmp to lenAffsMap
			for i := 1; i < length; i++ {
				lenAffsMap[i] = append(lenAffsMap[i], tmp[i]...)
			}
		}
	}
	if reverse {
		res := []PokerSet2{}
		for _, item := range affs {
			res = append(res, ps.Subtract(item))
		}
		return res
	} else {
		return affs
	}
}

func (ps PokerSet2) haveWZ() bool {
	return ps.Covers(getWZPS())
}

// 如果PS可以以以下几种牌形中的一种一手出完，则只会返回这一种牌形
// 单牌 ， 对子 ， 三不带，三带一，三带二，炸弹，单顺，双顺，飞机（飞机，三带均为在不拆炸弹的前提下）
func (ps PokerSet2) possibleActionsWithoutContext() []Action {
	pkCnt, sCnt, cCnt, tCnt, bCnt, sList, cList, tList, bList := ps.parse()
	pkCntu64 := uint64(pkCnt)

	if ps == getWZPS() {
		return []Action{getWZ()}
	}
	//解析一手牌
	if pkCnt == 1 {
		return []Action{NewActionWithPokerSet2(Single, 1, 0, sList[0][0], uint64(ps))}
	}
	if sCnt == cCnt {
		if pkCnt == 2 { //只有一对
			return []Action{NewActionWithPokerSet2(Couple, 1, 0, cList[0][0], uint64(ps))}
		} else if tCnt == 0 && len(cList) == 1 && cList[0][1] >= 3 { //只有双顺
			return []Action{NewActionWithPokerSet2(CL, cList[0][1], 0, cList[0][0], uint64(ps))}
		}
	}
	if cCnt == 0 && len(sList) == 1 && sList[0][1] >= 5 { //只有一个单顺
		return []Action{NewActionWithPokerSet2(SL, sList[0][1], 0, sList[0][0], uint64(ps))}
	}
	if sCnt == bCnt { //全炸
		res := []Action{}
		for _, b := range bList {
			res = append(res, NewActionWithoutAff(Bomb, 1, b))
		}
		return res
	}
	if bCnt == 0 && tCnt == 1 && pkCnt <= 4 { //三带
		return []Action{NewActionWithPokerSet2(Triple, tList[0][1], uint64(pkCnt-3), tList[0][0], uint64(ps))}
	}
	if bCnt == 0 && tCnt > 1 { //飞机
		if len(tList) == 1 {
			seqLen := tList[0][1]
			if pkCntu64 == seqLen*3 || pkCntu64 == seqLen*4 {
				return []Action{NewActionWithPokerSet2(TL, seqLen, pkCntu64-seqLen*3, tList[0][0], uint64(ps))}
			}
			//去掉一个作配牌
			if pkCntu64 == (seqLen-1)*4 {
				return []Action{NewActionWithPokerSet2(TL, (seqLen - 1), pkCntu64-(seqLen-1)*3, tList[0][0], uint64(ps))}
			}
		} else { //取最长的，其它的做配牌
			var maxTList [2]uint64 = tList[0]
			for _, item := range tList {
				if item[1] > maxTList[1] {
					maxTList = item
				}
			}
			seqLen := maxTList[1]
			if pkCntu64 == seqLen*3 || pkCntu64 == seqLen*4 {
				return []Action{NewActionWithPokerSet2(TL, seqLen, pkCntu64-seqLen*3, maxTList[0], uint64(ps))}
			}
		}
	}
	//如果不符合一手牌
	res := make([]Action, 0, 200)
	if ps.haveWZ() {
		res = append(res, getWZ())
	}
	//单牌和单顺
	for _, item := range sList {
		var i uint64
		for i = item[0]; i < item[0]+item[1]; i++ {
			res = append(res, NewActionWithoutAff(Single, 1, i))
		}
		if item[1] >= 5 {
			seqs := getSubSeqsWithMinLen(item, 5)
			for _, seq := range seqs {
				res = append(res, NewActionWithoutAff(SL, seq[1], seq[0]))
			}
		}
	}
	//对子和双顺
	for _, item := range cList {
		var i uint64
		for i = item[0]; i < item[0]+item[1]; i++ {
			res = append(res, NewActionWithoutAff(Couple, 1, i))
		}
		if item[1] >= 2 {
			seqs := getSubSeqsWithMinLen(item, 3)
			for _, seq := range seqs {
				res = append(res, NewActionWithoutAff(CL, seq[1], seq[0]))
			}
		}
	}

	var actionsNeedAff []Action

	//得到本脾可能三不带和飞机不带
	for _, item := range tList {
		var i uint64
		for i = item[0]; i < item[0]+item[1]; i++ {
			a := NewActionWithoutAff(Triple, 1, i)
			res = append(res, a)
			actionsNeedAff = append(actionsNeedAff, a)
		}
		if item[1] >= 2 {
			seqs := getSubSeqsWithMinLen(item, 2)
			for _, seq := range seqs {
				a := NewActionWithoutAff(TL, seq[1], seq[0])
				res = append(res, a)
				actionsNeedAff = append(actionsNeedAff, a)
			}
		}
	}
	//得到本脾可能炸弹（四不带）
	for _, item := range bList {
		a := NewActionWithoutAff(Bomb, 1, item)
		res = append(res, a)
		actionsNeedAff = append(actionsNeedAff, a)
	}

	for _, a := range actionsNeedAff {
		//var affLen int
		actionType := a.ActionType()
		//只能四带二
		if actionType == Bomb {
			affs := getAffsSingle(ps.Subtract(a.PokerSet2()), a.PokerSet2(), 2)
			//add affs
			for _, aff := range affs {
				res = append(res, a.SetAff(aff))
			}
			affs = getAffsCouple(ps.Subtract(a.PokerSet2()), 2)
			//add affs
			for _, aff := range affs {
				res = append(res, a.SetAff(aff))
			}
		} else if actionType == Triple {
			affs := getAffsSingle(ps.Subtract(a.PokerSet2()), a.PokerSet2(), 1)
			//add affs
			for _, aff := range affs {
				res = append(res, a.SetAff(aff))
			}
			affs = getAffsCouple(ps.Subtract(a.PokerSet2()), 1)
			//add affs
			for _, aff := range affs {
				res = append(res, a.SetAff(aff))
			}
		} else if actionType == TL {
			seqLen := a.SeqLen()
			affs := getAffsSingle(ps.Subtract(a.PokerSet2()), a.PokerSet2(), seqLen)
			//add affs
			for _, aff := range affs {
				res = append(res, a.SetAff(aff))
			}
			affs = getAffsCouple(ps.Subtract(a.PokerSet2()), seqLen)
			//add affs
			for _, aff := range affs {
				res = append(res, a.SetAff(aff))
			}
		}
	}
	return res
}

func (ps PokerSet2) possibleActionsWithContext(c Action) []Action {
	pkCnt, _, _, _, _, sList, cList, tList, bList := ps.parse()
	actionType := c.ActionType()
	if actionType == Pass {
		panic("actionType musnt be Pass")
	}
	res := []Action{}
	if ps.haveWZ() {
		res = append(res, getWZ())
	}
	if c.IsBomb() { //find bigger bombs
		for _, b := range bList {
			keyCard := c.KeyCard()
			if b > keyCard {
				res = append(res, NewActionWithoutAff(Bomb, 1, b))
			}
		}
		return res
	}
	//add bombs , for whatever other actiontype
	for _, b := range bList {
		res = append(res, NewActionWithoutAff(Bomb, 1, b))
	}
	//not enough
	if pkCnt < c.PokerCount() {
		return res
	}
	// process actiontype seperately
	keyCard := c.KeyCard()
	switch actionType {
	case Single:
		for _, item := range sList {
			if keyCard < item[1]+item[0]-1 {
				t := keyCard + 1
				if t < item[0] {
					t = item[0]
				}
				for t < item[1]+item[0] {
					res = append(res, NewActionWithoutAff(Single, 1, t))
					t++
				}
			}
		}
	case Couple:
		for _, item := range cList {
			if keyCard < item[1]+item[0]-1 {
				t := keyCard + 1
				if t < item[0] {
					t = item[0]
				}
				for t < item[1]+item[0] {
					res = append(res, NewActionWithoutAff(Couple, 1, t))
					t++
				}
			}
		}
	case Triple:
		affLen := c.AffLen()
		coupleAff := false
		if affLen == 2 {
			coupleAff = true
		}
		actionsNeedAff := []Action{}
		for _, item := range tList {
			if keyCard < item[1]+item[0]-1 {
				t := keyCard + 1
				if t < item[0] {
					t = item[0]
				}
				for t < item[1]+item[0] {
					if affLen == 0 {
						res = append(res, NewActionWithoutAff(Triple, 1, t))
					} else {
						actionsNeedAff = append(actionsNeedAff, NewActionWithoutAff(Triple, 1, t))
					}
					t++
				}
			}
		}
		if affLen != 0 {
			for _, a := range actionsNeedAff {
				var affs []PokerSet2
				if coupleAff {
					affs = getAffsCouple(ps.Subtract(a.PokerSet2()), 1)
				} else {
					affs = getAffsSingle(ps.Subtract(a.PokerSet2()), a.PokerSet2(), 1)
				}
				for _, aff := range affs {
					res = append(res, a.SetAff(aff))
				}
			}
		}
	case BW:
		affLen := c.AffLen()
		if affLen == 0 {
			panic("BW actionType , but afflen is 0")
		}
		coupleAff := false
		if affLen == 4 {
			coupleAff = true
		}
		actionsNeedAff := []Action{}
		for _, item := range bList {
			if keyCard < item {
				actionsNeedAff = append(actionsNeedAff, NewActionWithoutAff(Bomb, 1, item))
			}
		}
		for _, a := range actionsNeedAff {
			var affs []PokerSet2
			if coupleAff {
				affs = getAffsCouple(ps.Subtract(a.PokerSet2()), 2)
			} else {
				affs = getAffsSingle(ps.Subtract(a.PokerSet2()), a.PokerSet2(), 2)
			}
			for _, aff := range affs {
				res = append(res, a.SetAff(aff))
			}
		}
	case SL:
		seqLen := c.SeqLen()
		for _, item := range sList {
			if item[0] > keyCard && item[1] >= uint64(seqLen) {
				seqs := getSubSeqsWithFixLen(item, uint64(seqLen))
				for _, seq := range seqs {
					res = append(res, NewActionWithoutAff(SL, seq[1], seq[0]))
				}

			} else if item[0] <= keyCard && item[0]+item[1] >= keyCard+1+uint64(seqLen) {
				//trim item
				item[1] = item[1] - (keyCard + 1 - item[0])
				item[0] = keyCard + 1
				seqs := getSubSeqsWithFixLen(item, uint64(seqLen))
				for _, seq := range seqs {
					res = append(res, NewActionWithoutAff(SL, seq[1], seq[0]))
				}
			}
		}
	case CL:
		seqLen := c.SeqLen()
		for _, item := range cList {
			if item[0] > keyCard && item[1] >= uint64(seqLen) {
				seqs := getSubSeqsWithFixLen(item, uint64(seqLen))
				for _, seq := range seqs {
					res = append(res, NewActionWithoutAff(CL, seq[1], seq[0]))
				}

			} else if item[0] <= keyCard && item[0]+item[1] >= keyCard+1+uint64(seqLen) {
				//trim item
				item[1] = item[1] - (keyCard + 1 - item[0])
				item[0] = keyCard + 1
				seqs := getSubSeqsWithFixLen(item, uint64(seqLen))
				for _, seq := range seqs {
					res = append(res, NewActionWithoutAff(CL, seq[1], seq[0]))
				}
			}
		}
	case TL:
		seqLen := c.SeqLen()
		affLen := c.AffLen()
		coupleAff := false
		if affLen == seqLen*2 {
			coupleAff = true
		}
		actionsNeedAff := []Action{}
		for _, item := range tList {
			if item[0] > keyCard && item[1] >= uint64(seqLen) {
				seqs := getSubSeqsWithFixLen(item, uint64(seqLen))
				for _, seq := range seqs {
					actionsNeedAff = append(actionsNeedAff, NewActionWithoutAff(TL, seq[1], seq[0]))
				}

			} else if item[0] <= keyCard && item[0]+item[1] >= keyCard+1+uint64(seqLen) {
				//trim item
				item[1] = item[1] - (keyCard + 1 - item[0])
				item[0] = keyCard + 1
				seqs := getSubSeqsWithFixLen(item, uint64(seqLen))
				for _, seq := range seqs {
					actionsNeedAff = append(actionsNeedAff, NewActionWithoutAff(TL, seq[1], seq[0]))
				}
			}
		}
		for _, a := range actionsNeedAff {
			var affs []PokerSet2
			if coupleAff {
				affs = getAffsCouple(ps.Subtract(a.PokerSet2()), seqLen)
			} else {
				affs = getAffsSingle(ps.Subtract(a.PokerSet2()), a.PokerSet2(), seqLen)
			}
			for _, aff := range affs {
				res = append(res, a.SetAff(aff))
			}
		}
	}
	return res
}

func (ps PokerSet2) PossibleActionsWithContext(c Action) []Action {
	if c.ActionType() == Pass {
		return ps.possibleActionsWithoutContext()
	} else {
		actionList := ps.possibleActionsWithContext(c)
		actionList = append(actionList, NewActionWithoutAff(Pass, 0, 0))
		return actionList
	}
}
