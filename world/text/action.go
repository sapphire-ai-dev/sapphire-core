package text

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
)

func (w *textWorld) validCursorItem(currDir *directory, pos *actorPos, cmd int) bool {
	dirSize := len(currDir.content)
	if currDir.parent() != nil {
		dirSize++
	}

	if pos.cursorItem < 0 || pos.cursorItem >= dirSize {
		return false
	}

	if pos.cursorItem == 0 && cmd == ChangeItemCmdUp {
		return false
	}

	if pos.cursorItem == len(currDir.content)-1 && cmd == ChangeItemCmdDown {
		return false
	}

	return true
}

func (w *textWorld) locateItem(actorId int, cmd int) (*actorPos, item) {
	pos, posSeen := w.actors[actorId]
	if !posSeen {
		return nil, nil
	}

	currItem, currItemSeen := w.items[pos.currItemId]
	if !currItemSeen {
		return nil, nil
	}

	if _, ok := currItem.(*file); ok {
		if pos.cursorItem != 0 {
			return nil, nil
		}

		if cmd != ChangeItemCmdEnter {
			return nil, nil
		}
	}

	return pos, currItem
}

func (w *textWorld) changeItemReady(actorId, cmd int) bool {
	pos, currItem := w.locateItem(actorId, cmd)
	if pos == nil {
		return false
	}

	if _, ok := currItem.(*file); ok {
		return true
	}

	currDir := currItem.(*directory)
	return w.validCursorItem(currDir, pos, cmd)
}

func (w *textWorld) changeItemStep(actorId, cmd int) {
	pos, currItem := w.locateItem(actorId, cmd)
	if pos == nil {
		return
	}

	if _, ok := currItem.(*file); ok {
		w.actors[actorId].cursorItem = currItem.parent().id()
		return
	}

	currDir := currItem.(*directory)
	if !w.validCursorItem(currDir, pos, cmd) {
		return
	}

	if cmd == ChangeItemCmdUp {
		w.actors[actorId].cursorItem--
		return
	}

	if cmd == ChangeItemCmdDown {
		w.actors[actorId].cursorItem++
		return
	}

	cursorItem := pos.cursorItem
	if currDir.parent() != nil {
		cursorItem--
	}

	if cursorItem == -1 {
		w.actors[actorId].currItemId = currDir.parent().id()
	} else {
		w.actors[actorId].currItemId = currDir.content[cursorItem].id()
	}

	w.actors[actorId].cursorItem = 0
}

func (w *textWorld) changeItemWrap(actorId, cmd int) *world.ActionInterface {
	if cmd < 0 || cmd >= ChangeItemCmdEnd {
		return nil
	}

	return &world.ActionInterface{
		Name: pressKeyCmds[cmd],
		Ready: func() bool {
			return w.changeItemReady(actorId, cmd)
		},
		Step: func() {
			w.changeItemStep(actorId, cmd)
		},
	}
}

func (w *textWorld) identifyFile(actorId int) (*file, *actorPos) {
	pos, posSeen := w.actors[actorId]
	if !posSeen {
		return nil, nil
	}

	currItem, currItemSeen := w.items[pos.currItemId]
	if !currItemSeen {
		return nil, nil
	}

	currFile, isFile := currItem.(*file)
	if !isFile {
		return nil, nil
	}

	if pos.cursorLine < 0 || pos.cursorLine >= len(currFile.lines) {
		return nil, nil
	}

	currLine := currFile.lines[pos.cursorLine]
	if pos.cursorChar < 0 || pos.cursorChar > len(currLine.characters) {
		return nil, nil
	}

	return currFile, pos
}

func (w *textWorld) pressKeyReady(actorId int) bool {
	currFile, _ := w.identifyFile(actorId)
	return currFile != nil
}

func (w *textWorld) pressKeyStep(actorId, cmd int) {
	currFile, pos := w.identifyFile(actorId)
	if currFile == nil {
		return
	}

	if val, seen := pressKeyCmds[cmd]; seen {
		currLine := currFile.lines[pos.cursorLine]
		left, right := currLine.characters[:pos.cursorChar], currLine.characters[pos.cursorChar:]
		currLine.characters = append(append(left, currLine.newCharacter(val)), right...)
		w.actors[actorId].cursorChar++
	}
}

func (w *textWorld) pressKeyWrap(actorId, cmd int) *world.ActionInterface {
	if cmd < 0 || cmd >= PressKeyCmdEnd {
		return nil
	}

	return &world.ActionInterface{
		Name: "key" + pressKeyCmds[cmd],
		Ready: func() bool {
			return w.pressKeyReady(actorId)
		},
		Step: func() {
			w.pressKeyStep(actorId, cmd)
		},
	}
}

func (w *textWorld) newActionInterfaces(actorId int) []*world.ActionInterface {
	var result []*world.ActionInterface
	for cmd := range changeItemCmds {
		result = append(result, w.changeItemWrap(actorId, cmd))
	}

	for cmd := range pressKeyCmds {
		result = append(result, w.pressKeyWrap(actorId, cmd))
	}

	return result
}

const (
	ChangeItemCmdUp = iota
	ChangeItemCmdDown
	ChangeItemCmdEnter
	ChangeItemCmdExec
	ChangeItemCmdEnd
)

var changeItemCmds = map[int]string{
	ChangeItemCmdUp:    "itemUp",
	ChangeItemCmdDown:  "itemDown",
	ChangeItemCmdEnter: "itemEnter",
	ChangeItemCmdExec:  "itemExec",
}

const (
	PressKeyCmd0 = iota
	PressKeyCmd1
	PressKeyCmd2
	PressKeyCmd3
	PressKeyCmd4
	PressKeyCmd5
	PressKeyCmd6
	PressKeyCmd7
	PressKeyCmd8
	PressKeyCmd9
	PressKeyCmdA
	PressKeyCmdB
	PressKeyCmdC
	PressKeyCmdD
	PressKeyCmdE
	PressKeyCmdF
	PressKeyCmdG
	PressKeyCmdH
	PressKeyCmdI
	PressKeyCmdJ
	PressKeyCmdK
	PressKeyCmdL
	PressKeyCmdM
	PressKeyCmdN
	PressKeyCmdO
	PressKeyCmdP
	PressKeyCmdQ
	PressKeyCmdR
	PressKeyCmdS
	PressKeyCmdT
	PressKeyCmdU
	PressKeyCmdV
	PressKeyCmdW
	PressKeyCmdX
	PressKeyCmdY
	PressKeyCmdZ
	PressKeyCmdShift0
	PressKeyCmdShift1
	PressKeyCmdShift2
	PressKeyCmdShift3
	PressKeyCmdShift4
	PressKeyCmdShift5
	PressKeyCmdShift6
	PressKeyCmdShift7
	PressKeyCmdShift8
	PressKeyCmdShift9
	PressKeyCmdMinus
	PressKeyCmdPlus
	PressKeyCmdUnderscore
	PressKeyCmdEqual
	PressKeyCmdLeftSquareBracket
	PressKeyCmdLeftCurlyBracket
	PressKeyCmdRightSquareBracket
	PressKeyCmdRightCurlyBracket
	PressKeyCmdSpace
	PressKeyCmdComma
	PressKeyCmdPeriod
	PressKeyCmdSlash
	PressKeyCmdShiftComma
	PressKeyCmdShiftPeriod
	PressKeyCmdShiftSlash
	PressKeyCmdBackSlash
	PressKeyCmdVertical
	PressKeyCmdBackspace
	PressKeyCmdEnter
	PressKeyCmdUp
	PressKeyCmdDown
	PressKeyCmdLeft
	PressKeyCmdRight
	PressKeyCmdEnd
)

var pressKeyCmds = map[int]string{
	PressKeyCmd1:                  "1",
	PressKeyCmd2:                  "2",
	PressKeyCmd3:                  "3",
	PressKeyCmd4:                  "4",
	PressKeyCmd5:                  "5",
	PressKeyCmd6:                  "6",
	PressKeyCmd7:                  "7",
	PressKeyCmd8:                  "8",
	PressKeyCmd9:                  "9",
	PressKeyCmdA:                  "a",
	PressKeyCmdB:                  "b",
	PressKeyCmdC:                  "c",
	PressKeyCmdD:                  "d",
	PressKeyCmd0:                  "0",
	PressKeyCmdE:                  "e",
	PressKeyCmdF:                  "f",
	PressKeyCmdG:                  "g",
	PressKeyCmdH:                  "h",
	PressKeyCmdI:                  "i",
	PressKeyCmdJ:                  "j",
	PressKeyCmdK:                  "k",
	PressKeyCmdL:                  "l",
	PressKeyCmdM:                  "m",
	PressKeyCmdN:                  "n",
	PressKeyCmdO:                  "o",
	PressKeyCmdP:                  "p",
	PressKeyCmdQ:                  "q",
	PressKeyCmdR:                  "r",
	PressKeyCmdS:                  "s",
	PressKeyCmdT:                  "t",
	PressKeyCmdU:                  "u",
	PressKeyCmdV:                  "v",
	PressKeyCmdW:                  "w",
	PressKeyCmdX:                  "x",
	PressKeyCmdY:                  "y",
	PressKeyCmdZ:                  "z",
	PressKeyCmdShift0:             "!",
	PressKeyCmdShift1:             "@",
	PressKeyCmdShift2:             "#",
	PressKeyCmdShift3:             "$",
	PressKeyCmdShift4:             "%",
	PressKeyCmdShift5:             "^",
	PressKeyCmdShift6:             "&",
	PressKeyCmdShift7:             "*",
	PressKeyCmdShift8:             "(",
	PressKeyCmdShift9:             ")",
	PressKeyCmdMinus:              "-",
	PressKeyCmdPlus:               "+",
	PressKeyCmdUnderscore:         "_",
	PressKeyCmdEqual:              "=",
	PressKeyCmdLeftSquareBracket:  "[",
	PressKeyCmdLeftCurlyBracket:   "{",
	PressKeyCmdRightSquareBracket: "]",
	PressKeyCmdRightCurlyBracket:  "}",
	PressKeyCmdSpace:              " ",
	PressKeyCmdComma:              ",",
	PressKeyCmdPeriod:             ".",
	PressKeyCmdSlash:              "/",
	PressKeyCmdShiftComma:         "<",
	PressKeyCmdShiftPeriod:        ">",
	PressKeyCmdShiftSlash:         "?",
	PressKeyCmdBackSlash:          "\\",
	PressKeyCmdVertical:           "|",
}

var specialKeyCmds = map[int]bool{
	PressKeyCmdBackspace: true,
	PressKeyCmdEnter:     true,
	PressKeyCmdUp:        true,
	PressKeyCmdDown:      true,
	PressKeyCmdLeft:      true,
	PressKeyCmdRight:     true,
}
