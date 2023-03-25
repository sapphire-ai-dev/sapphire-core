package text

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
)

const (
	itemTypeRoot      = "[itemType]"
	itemTypeDirectory = "[directory]"
	itemTypeFile      = "[file]"
	itemDirection     = "[itemDirection]"
	contentTypeRoot   = "[contentType]"
	contentTypeLine   = "[line]"
	lineDirection     = "[lineDirection]"
	charDirection     = "[lineDirection]"
)

type item interface {
	id() int
	parent() item
	name() string
	itemImg(itemDelta int) *world.Image
	dirImgs(actorPos int) []*world.Image
	fileImgs(cursorLine, cursorChar int) []*world.Image
}

type abstractItem struct {
	w    *textWorld
	self item
	i    int
	p    item
	n    string
}

func (a *abstractItem) id() int {
	return a.i
}

func (a *abstractItem) parent() item {
	return a.p
}

func (a *abstractItem) name() string {
	return a.n
}

func itemImg(id int, name string, itemType string, itemDelta int) *world.Image {
	itemDir := world.TernaryZro
	if itemDelta > 0 {
		itemDir = world.TernaryPos
	} else if itemDelta < 0 {
		itemDir = world.TernaryNeg
	}

	return &world.Image{
		Id:   id,
		Name: name,
		Permanent: []*world.Info{
			{
				Labels: []string{world.InfoLabelObservable, itemTypeRoot, itemType},
				Value:  nil,
			},
		},
		Transient: []*world.Info{
			{
				Labels: []string{world.InfoLabelObservable, itemDirection, itemDir},
				Value:  itemDelta,
			},
		},
	}
}

func (w *textWorld) newAbstractItem(self, parent item, name string, out **abstractItem) {
	*out = &abstractItem{
		w:    w,
		self: self,
		i:    world.NewUnitId(),
		p:    parent,
		n:    name,
	}

	w.items[(*out).i] = self
}

type directory struct {
	*abstractItem
	content []item
}

func (d *directory) itemImg(itemDelta int) *world.Image {
	return itemImg(d.i, d.n, itemTypeDirectory, itemDelta)
}

func (d *directory) dirImgs(actorPos int) []*world.Image {
	var result []*world.Image

	if d.p != nil {
		result = append(result, itemImg(d.p.id(), d.p.name(), itemTypeDirectory, actorPos))
	}

	i := len(result)
	for _, elem := range d.content {
		result = append(result, elem.itemImg(actorPos-i))
		i++
	}

	return result
}

func (d *directory) fileImgs(_, _ int) []*world.Image {
	return []*world.Image{}
}

func (d *directory) newDirectory(name string) *directory {
	result := &directory{
		content: []item{},
	}

	d.w.newAbstractItem(result, d, name, &result.abstractItem)
	d.content = append(d.content, result)
	return result
}

type file struct {
	*abstractItem
	lines []*line
}

func (f *file) itemImg(itemDelta int) *world.Image {
	return itemImg(f.i, f.n, itemTypeFile, itemDelta)
}

func (f *file) dirImgs(actorPos int) []*world.Image {
	var result []*world.Image

	if f.p != nil {
		result = append(result, itemImg(f.p.id(), f.p.name(), itemTypeDirectory, actorPos))
	}

	return result
}

func (f *file) fileImgs(cursorLine, cursorChar int) []*world.Image {
	var result []*world.Image

	for i, l := range f.lines {
		result = append(result, l.img(cursorLine-i, cursorChar)...)
	}

	return result
}

func (d *directory) newFile(name string) *file {
	result := &file{}
	result.lines = []*line{result.newLine()}

	d.w.newAbstractItem(result, d, name, &result.abstractItem)
	d.content = append(d.content, result)
	return result
}

type line struct {
	id         int
	parent     *file
	characters []*character
}

func (l *line) img(cursorLineDelta, cursorChar int) []*world.Image {
	var result []*world.Image
	cursorDir := world.TernaryZro
	if cursorLineDelta > 0 {
		cursorDir = world.TernaryPos
	} else if cursorLineDelta < 0 {
		cursorDir = world.TernaryNeg
	}

	result = append(result, &world.Image{
		Id:   l.id,
		Name: "",
		Permanent: []*world.Info{
			{
				Labels: []string{world.InfoLabelObservable, contentTypeRoot, contentTypeLine},
				Value:  nil,
			},
		},
		Transient: []*world.Info{
			{
				Labels: []string{world.InfoLabelObservable, lineDirection, cursorDir},
				Value:  cursorLineDelta,
			},
		},
	})

	for j, c := range l.characters {
		result = append(result, c.img(j-cursorChar))
	}

	return result
}

func (f *file) newLine() *line {
	result := &line{
		id:         world.NewUnitId(),
		parent:     f,
		characters: []*character{},
	}

	f.lines = append(f.lines, result)
	return result
}

type character struct {
	id     int
	parent *line
	shape  string
}

func (c *character) img(cursorCharDelta int) *world.Image {
	cursorDir := world.TernaryZro
	if cursorCharDelta > 0 {
		cursorDir = world.TernaryPos
	} else if cursorCharDelta < 0 {
		cursorDir = world.TernaryNeg
	}

	return &world.Image{
		Id:   c.id,
		Name: "",
		Permanent: []*world.Info{
			{
				Labels: []string{world.InfoLabelObservable, contentTypeRoot, c.shape},
				Value:  nil,
			},
		},
		Transient: []*world.Info{
			{
				Labels: []string{world.InfoLabelObservable, charDirection, cursorDir},
				Value:  cursorCharDelta,
			},
		},
	}
}

func (l *line) newCharacter(shape string) *character {
	result := &character{
		id:     world.NewUnitId(),
		parent: l,
		shape:  shape,
	}

	return result
}
