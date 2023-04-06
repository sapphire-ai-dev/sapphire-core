package text

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestDirectoryDisplayText(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory

	d1Name, d21Name, d22Name := "d1", "d21", "d22"
	f11Name, f12Name, f221Name, f222Name := "f11", "f12", "f221", "f222"
	d1 := root.newDirectory(d1Name)
	d21 := d1.newDirectory(d21Name)
	d22 := d1.newDirectory(d22Name)
	f11 := d1.newFile(f11Name)
	f12 := d1.newFile(f12Name)
	f321 := d21.newFile(f221Name)
	f322 := d22.newFile(f222Name)

	txt := w.directoryText()
	assert.Contains(t, txt, d1Name)
	assert.Contains(t, txt, d21Name)
	assert.Contains(t, txt, d22Name)
	assert.Contains(t, txt, f11Name)
	assert.Contains(t, txt, f12Name)
	assert.Contains(t, txt, f221Name)
	assert.Contains(t, txt, f222Name)
	assert.Contains(t, txt, strconv.Itoa(d1.id()))
	assert.Contains(t, txt, strconv.Itoa(d21.id()))
	assert.Contains(t, txt, strconv.Itoa(d22.id()))
	assert.Contains(t, txt, strconv.Itoa(f11.id()))
	assert.Contains(t, txt, strconv.Itoa(f12.id()))
	assert.Contains(t, txt, strconv.Itoa(f321.id()))
	assert.Contains(t, txt, strconv.Itoa(f322.id()))
}

func TestFileDisplayText(t *testing.T) {
	Init()
	w := world.GetWorld().(*textWorld)
	root := w.rootDirectory
	f1Name, f2Name := "f1", "f2"
	f1 := root.newFile(f1Name)
	f2 := root.newFile(f2Name)
	assert.Equal(t, w.fileText(f1.id()), "\n")
	assert.Equal(t, w.fileText(f2.id()), "\n")

	line1 := "roses are red"
	line2 := "violets are blue"
	line3 := ""
	line4 := "line 3 is empty"
	lines := []string{line1, line2, line3, line4}
	for i, l := range lines {
		for j := range l {
			f1.lines[i].characters = append(f1.lines[i].characters, f1.lines[i].newCharacter(l[j:j+1]))
		}
		f1.appendLine(f1.newLine())
	}
	assert.Equal(t, w.fileText(f1.id()), strings.Join(lines, "\n")+"\n\n")
	assert.Equal(t, w.fileText(f2.id()), "\n")
}

func TestDisplayErrors(t *testing.T) {
	w := newTextWorld()
	assert.Empty(t, w.fileText(-1))

	display = nil
	w.displaySendFile(-1)
}
