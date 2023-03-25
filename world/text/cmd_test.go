package text

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextWorldCmd(t *testing.T) {
	Init()
	assert.PanicsWithError(t, world.ErrInvalidArgs.Error(), func() {
		world.Cmd()
	})
	assert.PanicsWithError(t, world.ErrInvalidCmd.Error(), func() {
		world.Cmd(-1)
	})
}

func TestCmdGetRootDirectoryId(t *testing.T) {
	w := newTextWorld()
	out := 0
	w.Cmd(CmdTypeGetRootDirectoryId, &out)
	assert.Equal(t, out, w.rootDirectory.i)
}

func TestCmdCreateDirectory(t *testing.T) {
	w := newTextWorld()
	assert.Empty(t, w.rootDirectory.content)
	dName := "dName"
	w.Cmd(CmdTypeCreateDirectory, w.rootDirectory.i, dName)
	assert.Equal(t, dName, w.rootDirectory.content[0].(*directory).name())

	out := 0
	w.Cmd(CmdTypeCreateDirectory, w.rootDirectory.i, dName, &out)
	assert.Equal(t, out, w.rootDirectory.content[1].(*directory).id())

	w.Cmd(CmdTypeCreateDirectory, w.rootDirectory.i-1, dName, &out)
}

func TestCmdCreateFile(t *testing.T) {
	w := newTextWorld()
	assert.Empty(t, w.rootDirectory.content)
	fName := "fName"
	w.Cmd(CmdTypeCreateFile, w.rootDirectory.i, fName)
	assert.Equal(t, fName, w.rootDirectory.content[0].(*file).name())

	out := 0
	w.Cmd(CmdTypeCreateFile, w.rootDirectory.i, fName, &out)
	assert.Equal(t, out, w.rootDirectory.content[1].(*file).id())

	w.Cmd(CmdTypeCreateFile, out, fName, &out)
}

func TestCmdTypeCharacter(t *testing.T) {
	w := newTextWorld()
	f := w.rootDirectory.newFile("")
	assert.Empty(t, f.lines[0].characters)

	w.Cmd(CmdTypeAddCharacter, f.id(), 0, 0, PressKeyCmd1)
	assert.Equal(t, f.lines[0].characters[0].shape, pressKeyCmds[PressKeyCmd1])

	w.Cmd(CmdTypeAddCharacter, w.rootDirectory.id()-1, 0, 0, PressKeyCmd1)
	w.Cmd(CmdTypeAddCharacter, w.rootDirectory.id(), 0, 0, PressKeyCmd1)
	w.Cmd(CmdTypeAddCharacter, f.id(), -1, 0, 0)
	w.Cmd(CmdTypeAddCharacter, f.id(), 0, -1, 0)
	w.Cmd(CmdTypeAddCharacter, f.id(), 0, 0, -1)
}
