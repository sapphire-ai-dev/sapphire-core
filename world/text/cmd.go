package text

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
)

const (
	CmdTypeGetRootDirectoryId = iota
	CmdTypeCreateDirectory
	CmdTypeCreateFile
	CmdTypeAddCharacter
	CmdTypeMoveActor
)

func (w *textWorld) Cmd(args ...any) {
	cmdType := world.GetArg[int](0, true, 0, args)
	cmdOptions := map[int]func(...any){
		CmdTypeGetRootDirectoryId: w.cmdGetRootDirectoryId,
		CmdTypeCreateDirectory:    w.cmdCreateDirectory,
		CmdTypeCreateFile:         w.cmdCreateFile,
		CmdTypeAddCharacter:       w.cmdAddCharacter,
		CmdTypeMoveActor:          w.moveActor,
	}

	if option, seen := cmdOptions[cmdType]; seen {
		option(args[1:]...)
	} else {
		panic(world.ErrInvalidCmd)
	}
}

func (w *textWorld) cmdGetRootDirectoryId(args ...any) {
	world.SetOutArg[int](0, true, w.rootDirectory.i, args)
}

func (w *textWorld) cmdCreateDirectory(args ...any) {
	itemId := world.GetArg[int](0, true, 0, args)
	newDirectoryName := world.GetArg[string](1, true, "", args)
	d := w.locateDirectory(itemId)
	if d != nil {
		world.SetOutArg[int](2, false, d.newDirectory(newDirectoryName).i, args)
	}
}

func (w *textWorld) cmdCreateFile(args ...any) {
	itemId := world.GetArg[int](0, true, 0, args)
	fileName := world.GetArg[string](1, true, "", args)
	d := w.locateDirectory(itemId)
	if d != nil {
		world.SetOutArg[int](2, false, d.newFile(fileName).i, args)
	}
}

func (w *textWorld) cmdAddCharacter(args ...any) {
	itemId := world.GetArg[int](0, true, 0, args)
	lineNum := world.GetArg[int](1, true, 0, args)
	charNum := world.GetArg[int](2, true, 0, args)
	charId := world.GetArg[int](3, true, 0, args)
	f := w.locateFile(itemId)
	if f == nil || lineNum < 0 || lineNum >= len(f.lines) || charNum < 0 || charNum > len(f.lines[lineNum].characters) {
		return
	}

	c, charSeen := pressKeyCmds[charId]
	if !charSeen {
		return
	}

	currLine := f.lines[lineNum]
	left, right := currLine.characters[:charNum], currLine.characters[charNum:]
	currLine.characters = append(append(left, currLine.newCharacter(c)), right...)
}

func (w *textWorld) moveActor(args ...any) {
	actorId := world.GetArg[int](0, true, 0, args)
	itemId := world.GetArg[int](1, true, 0, args)
	w.actors[actorId].currItemId = itemId
}

func (w *textWorld) locateDirectory(itemId int) *directory {
	i, itemSeen := w.items[itemId]
	if !itemSeen {
		return nil
	}

	d, isDirectory := i.(*directory)
	if !isDirectory {
		return nil
	}

	return d
}

func (w *textWorld) locateFile(itemId int) *file {
	i, itemSeen := w.items[itemId]
	if !itemSeen {
		return nil
	}

	f, isFile := i.(*file)
	if !isFile {
		return nil
	}

	return f
}
