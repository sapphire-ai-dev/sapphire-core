package text

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectoryConstructor(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	assert.Equal(t, w, root.w)
	assert.Equal(t, root, root.self)
	assert.Equal(t, root.n, root.name())
	assert.Equal(t, root.i, root.id())
	assert.Equal(t, root.p, root.parent())
	assert.Equal(t, root, w.items[root.id()])
	assert.Empty(t, root.name())
	assert.Nil(t, root.parent())
	assertEmptyNotNil(t, root.content)

	name := "abc"
	d := root.newDirectory(name)
	assert.Equal(t, d, d.self)
	assert.Equal(t, name, d.name())
	assert.Equal(t, root, d.parent())
	assert.NotEqual(t, root.id(), d.id())
	assertEmptyNotNil(t, d.content)
}

func TestFileConstructor(t *testing.T) {
	root := newTextWorld().rootDirectory
	name := "abc"
	f := root.newFile(name)
	assert.Equal(t, name, f.name())
	assert.Equal(t, root, f.parent())
	assert.NotEqual(t, root.id(), f.id())
	assert.Len(t, f.lines, 1)
}

func TestDirectoryFileImage(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	actorId, _ := w.NewActor()
	actor := w.actors[actorId]
	assert.Empty(t, root.fileImgs(actor.cursorLine, actor.cursorChar))
}

func TestDirectoryImageFormat(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	actorId, _ := w.NewActor()

	dName := "abc"
	d := root.newDirectory(dName)
	imgs := w.Look(actorId)
	assert.Len(t, imgs, 1)
	assert.Equal(t, imgs[0].Name, dName)
	assert.Equal(t, imgs[0].Id, d.id())
	assert.Contains(t, imgs[0].Permanent[0].Labels, world.InfoLabelObservable)
	assert.Contains(t, imgs[0].Permanent[0].Labels, itemTypeRoot)
	assert.Contains(t, imgs[0].Permanent[0].Labels, itemTypeDirectory)
	assert.Nil(t, imgs[0].Permanent[0].Value)
	assert.Contains(t, imgs[0].Transient[0].Labels, world.InfoLabelObservable)
	assert.Contains(t, imgs[0].Transient[0].Labels, itemDirection)
	assert.Contains(t, imgs[0].Transient[0].Labels, world.TernaryZro)
	assert.Zero(t, imgs[0].Transient[0].Value)
}

func TestFileImageFormat(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	actorId, _ := w.NewActor()

	fName := "abc"
	f := root.newFile(fName)
	imgs := w.Look(actorId)
	assert.Len(t, imgs, 1)
	assert.Equal(t, imgs[0].Name, fName)
	assert.Equal(t, imgs[0].Id, f.id())
	assert.Contains(t, imgs[0].Permanent[0].Labels, world.InfoLabelObservable)
	assert.Contains(t, imgs[0].Permanent[0].Labels, itemTypeRoot)
	assert.Contains(t, imgs[0].Permanent[0].Labels, itemTypeFile)
	assert.Nil(t, imgs[0].Permanent[0].Value)
	assert.Contains(t, imgs[0].Transient[0].Labels, world.InfoLabelObservable)
	assert.Contains(t, imgs[0].Transient[0].Labels, itemDirection)
	assert.Contains(t, imgs[0].Transient[0].Labels, world.TernaryZro)
	assert.Zero(t, imgs[0].Transient[0].Value)
}

func TestDirectoryImage(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	actorId, _ := w.NewActor()

	d1Name, d2Name, d31Name, d32Name := "d1", "d2", "d31", "d32"
	d1 := root.newDirectory(d1Name)
	d2 := d1.newDirectory(d2Name)
	d31 := d2.newDirectory(d31Name)
	d32 := d2.newDirectory(d32Name)

	w.actors[actorId].currItemId = d2.id()
	imgs := w.Look(actorId)
	imgMap := map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 3)
	assert.Equal(t, imgMap[d1.id()].Name, d1Name)
	assert.Contains(t, imgMap[d1.id()].Transient[0].Labels, world.TernaryZro)
	assert.Equal(t, imgMap[d1.id()].Transient[0].Value, 0)
	assert.Equal(t, imgMap[d31.id()].Name, d31Name)
	assert.Contains(t, imgMap[d31.id()].Transient[0].Labels, world.TernaryNeg)
	assert.Equal(t, imgMap[d31.id()].Transient[0].Value, -1)
	assert.Equal(t, imgMap[d32.id()].Name, d32Name)
	assert.Contains(t, imgMap[d32.id()].Transient[0].Labels, world.TernaryNeg)
	assert.Equal(t, imgMap[d32.id()].Transient[0].Value, -2)

	w.actors[actorId].cursorItem = 1
	imgs = w.Look(actorId)
	imgMap = map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 3)
	assert.Equal(t, imgMap[d1.id()].Name, d1Name)
	assert.Contains(t, imgMap[d1.id()].Transient[0].Labels, world.TernaryPos)
	assert.Equal(t, imgMap[d1.id()].Transient[0].Value, 1)
	assert.Equal(t, imgMap[d31.id()].Name, d31Name)
	assert.Contains(t, imgMap[d31.id()].Transient[0].Labels, world.TernaryZro)
	assert.Equal(t, imgMap[d31.id()].Transient[0].Value, 0)
	assert.Equal(t, imgMap[d32.id()].Name, d32Name)
	assert.Contains(t, imgMap[d32.id()].Transient[0].Labels, world.TernaryNeg)
	assert.Equal(t, imgMap[d32.id()].Transient[0].Value, -1)
}

func TestFileImage(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	actorId, _ := w.NewActor()

	dName, f1Name, f2Name := "d", "f1", "f2"
	d := root.newDirectory(dName)
	f1 := d.newFile(f1Name)
	f2 := d.newFile(f2Name)

	w.actors[actorId].currItemId = d.id()
	imgs := w.Look(actorId)
	imgMap := map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 3)
	assert.Equal(t, imgMap[root.id()].Name, root.name())
	assert.Contains(t, imgMap[root.id()].Transient[0].Labels, world.TernaryZro)
	assert.Equal(t, imgMap[root.id()].Transient[0].Value, 0)
	assert.Equal(t, imgMap[f1.id()].Name, f1Name)
	assert.Contains(t, imgMap[f1.id()].Transient[0].Labels, world.TernaryNeg)
	assert.Equal(t, imgMap[f1.id()].Transient[0].Value, -1)
	assert.Equal(t, imgMap[f2.id()].Name, f2Name)
	assert.Contains(t, imgMap[f2.id()].Transient[0].Labels, world.TernaryNeg)
	assert.Equal(t, imgMap[f2.id()].Transient[0].Value, -2)

	w.actors[actorId].currItemId = f1.id()
	imgs = w.Look(actorId)
	imgMap = map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 2)
	assert.Equal(t, imgMap[d.id()].Name, dName)
	assert.Contains(t, imgMap[d.id()].Transient[0].Labels, world.TernaryZro)
	assert.Equal(t, imgMap[d.id()].Transient[0].Value, 0)
}

func TestLineConstructor(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	f := root.newFile("")
	l := f.newLine()
	assert.Equal(t, f, l.parent)
	assert.Empty(t, l.characters)

	l2 := f.newLine()
	assert.NotEqual(t, l.id, l2.id)
}

func TestLineImageFormat(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	actorId, _ := w.NewActor()
	f := root.newFile("")
	l := f.lines[0]
	w.actors[actorId].currItemId = f.id()
	imgs := w.Look(actorId)
	imgMap := map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 2)
	assert.Equal(t, imgMap[l.id].Name, "")
	assert.Contains(t, imgMap[l.id].Permanent[0].Labels, world.InfoLabelObservable)
	assert.Contains(t, imgMap[l.id].Permanent[0].Labels, contentTypeRoot)
	assert.Contains(t, imgMap[l.id].Permanent[0].Labels, contentTypeLine)
	assert.Nil(t, imgMap[l.id].Permanent[0].Value)
	assert.Contains(t, imgMap[l.id].Transient[0].Labels, world.InfoLabelObservable)
	assert.Contains(t, imgMap[l.id].Transient[0].Labels, lineDirection)
	assert.Contains(t, imgMap[l.id].Transient[0].Labels, world.TernaryZro)
	assert.Zero(t, imgMap[l.id].Transient[0].Value)
}

func TestLineImage(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	actorId, _ := w.NewActor()
	f := root.newFile("")
	l1 := f.lines[0]
	l2 := f.newLine()
	l3 := f.newLine()
	w.actors[actorId].currItemId = f.id()
	imgs := w.Look(actorId)
	imgMap := map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 4)
	assert.Contains(t, imgMap[l1.id].Transient[0].Labels, world.TernaryZro)
	assert.Equal(t, imgMap[l1.id].Transient[0].Value, 0)
	assert.Len(t, imgMap, 4)
	assert.Contains(t, imgMap[l2.id].Transient[0].Labels, world.TernaryNeg)
	assert.Equal(t, imgMap[l2.id].Transient[0].Value, -1)
	assert.Len(t, imgMap, 4)
	assert.Contains(t, imgMap[l3.id].Transient[0].Labels, world.TernaryNeg)
	assert.Equal(t, imgMap[l3.id].Transient[0].Value, -2)

	w.actors[actorId].cursorLine = 1
	imgs = w.Look(actorId)
	imgMap = map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 4)
	assert.Contains(t, imgMap[l1.id].Transient[0].Labels, world.TernaryPos)
	assert.Equal(t, imgMap[l1.id].Transient[0].Value, 1)
	assert.Contains(t, imgMap[l2.id].Transient[0].Labels, world.TernaryZro)
	assert.Equal(t, imgMap[l2.id].Transient[0].Value, 0)
	assert.Contains(t, imgMap[l3.id].Transient[0].Labels, world.TernaryNeg)
	assert.Equal(t, imgMap[l3.id].Transient[0].Value, -1)
}

func TestCharacterConstructor(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	f := root.newFile("")
	l := f.newLine()
	shape := pressKeyCmds[PressKeyCmd0]
	c := l.newCharacter(shape)
	assert.Equal(t, l, c.parent)
	assert.Equal(t, shape, c.shape)

	c2 := l.newCharacter(shape)
	assert.NotEqual(t, c.id, c2.id)
}

func TestCharImageFormat(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	actorId, _ := w.NewActor()
	f := root.newFile("")
	l := f.lines[0]
	shape := pressKeyCmds[PressKeyCmd0]
	c := l.newCharacter(shape)
	l.characters = append(l.characters, c)

	w.actors[actorId].currItemId = f.id()
	imgs := w.Look(actorId)
	imgMap := map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 3)
	assert.Equal(t, imgMap[c.id].Name, "")
	assert.Contains(t, imgMap[c.id].Permanent[0].Labels, world.InfoLabelObservable)
	assert.Contains(t, imgMap[c.id].Permanent[0].Labels, contentTypeRoot)
	assert.Contains(t, imgMap[c.id].Permanent[0].Labels, shape)
	assert.Nil(t, imgMap[c.id].Permanent[0].Value)
	assert.Contains(t, imgMap[c.id].Transient[0].Labels, world.InfoLabelObservable)
	assert.Contains(t, imgMap[c.id].Transient[0].Labels, charDirection)
	assert.Contains(t, imgMap[c.id].Transient[0].Labels, world.TernaryZro)
	assert.Zero(t, imgMap[c.id].Transient[0].Value)
}

func TestCharImage(t *testing.T) {
	w := newTextWorld()
	root := w.rootDirectory
	actorId, _ := w.NewActor()
	f := root.newFile("")
	l := f.lines[0]
	shape1, shape2, shape3 := pressKeyCmds[PressKeyCmd1], pressKeyCmds[PressKeyCmd2], pressKeyCmds[PressKeyCmd3]
	c1 := l.newCharacter(shape1)
	c2 := l.newCharacter(shape2)
	c3 := l.newCharacter(shape3)
	l.characters = append(l.characters, c1, c2, c3)
	w.actors[actorId].currItemId = f.id()
	imgs := w.Look(actorId)
	imgMap := map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 5)
	assert.Contains(t, imgMap[c1.id].Transient[0].Labels, world.TernaryZro)
	assert.Equal(t, imgMap[c1.id].Transient[0].Value, 0)
	assert.Contains(t, imgMap[c2.id].Transient[0].Labels, world.TernaryPos)
	assert.Equal(t, imgMap[c2.id].Transient[0].Value, 1)
	assert.Contains(t, imgMap[c3.id].Transient[0].Labels, world.TernaryPos)
	assert.Equal(t, imgMap[c3.id].Transient[0].Value, 2)

	w.actors[actorId].cursorChar = 1
	imgs = w.Look(actorId)
	imgMap = map[int]*world.Image{}
	for _, img := range imgs {
		imgMap[img.Id] = img
	}

	assert.Len(t, imgMap, 5)
	assert.Contains(t, imgMap[c1.id].Transient[0].Labels, world.TernaryNeg)
	assert.Equal(t, imgMap[c1.id].Transient[0].Value, -1)
	assert.Contains(t, imgMap[c2.id].Transient[0].Labels, world.TernaryZro)
	assert.Equal(t, imgMap[c2.id].Transient[0].Value, 0)
	assert.Contains(t, imgMap[c3.id].Transient[0].Labels, world.TernaryPos)
	assert.Equal(t, imgMap[c3.id].Transient[0].Value, 1)
}
