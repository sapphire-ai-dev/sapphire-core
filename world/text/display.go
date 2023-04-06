package text

import (
	"encoding/json"
	"strconv"
)

const indentTab = "  "

func (w *textWorld) directoryText() string {
	return w.directoryTextHelper("", w.rootDirectory)
}

func (w *textWorld) directoryTextHelper(indent string, i item) string {
	d, isD := i.(*directory)
	iType := "[D]"
	if !isD {
		iType = "[F]"
	}

	result := indent + iType + " " + i.name() + " " + strconv.Itoa(i.id()) + "\n"
	if isD {
		for _, c := range d.content {
			result += w.directoryTextHelper(indent+indentTab, c)
		}
	}

	return result
}

func (w *textWorld) fileText(id int) string {
	f := w.locateFile(id)
	if f == nil {
		return ""
	}

	result := ""
	for _, l := range f.lines {
		for _, c := range l.characters {
			result += c.shape
		}
		result += "\n"
	}

	return result
}

func (w *textWorld) displaySendDirectory() {
	if display == nil {
		return
	}

	data := &directoryDisplayData{Directory: w.directoryText()}
	bytes, err := json.Marshal(data)
	if err == nil {
		display.Send(bytes)
	}
}

func (w *textWorld) displaySendFile(id int) {
	if display == nil {
		return
	}

	data := &fileDisplayData{File: w.fileText(id)}
	bytes, err := json.Marshal(data)
	if err == nil {
		display.Send(bytes)
	}
}

type directoryDisplayData struct {
	Directory string `json:"directory"`
}

type fileDisplayData struct {
	File string `json:"file"`
}
