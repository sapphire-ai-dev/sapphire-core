package world

type langReceiveBuffer struct {
	id       *int
	messages map[int][]*LangMessage // clock -> messages

	// records where an actor left off last listen
	lastListenTime        int
	lastListenCount       int
	lastListenCommonCount int
}

func (b *langReceiveBuffer) addMessage(time int, speakerId, listenerId *int, body string) {
	if _, seen := b.messages[time]; !seen {
		b.messages[time] = []*LangMessage{}
	}

	b.messages[time] = append(b.messages[time], &LangMessage{
		Src:  speakerId,
		Dst:  listenerId,
		Body: body,
	})
}

func (w *AbstractWorld) newLangReceiveBuffer(id *int) {
	result := &langReceiveBuffer{id: id, messages: map[int][]*LangMessage{}}
	if id == nil {
		w.langReceiveBuffers[-1] = result
	} else {
		w.langReceiveBuffers[*id] = result
	}
}

type LangMessage struct {
	Src  *int
	Dst  *int
	Body string
}

func (w *AbstractWorld) Listen(id int) []*LangMessage {
	var result []*LangMessage
	buffer, seen := w.langReceiveBuffers[id]
	if !seen {
		return result
	}

	commonBuffer := w.langReceiveBuffers[-1]
	result = append(result, w.ListenHelper(buffer, commonBuffer, w.clock-1)...)
	result = append(result, w.ListenHelper(buffer, commonBuffer, w.clock)...)
	return result
}

func (w *AbstractWorld) ListenHelper(buffer, commonBuffer *langReceiveBuffer, time int) []*LangMessage {
	var result []*LangMessage
	if buffer.lastListenTime < time {
		buffer.lastListenTime = time
		buffer.lastListenCount = 0
		buffer.lastListenCommonCount = 0
	}

	result = append(result, buffer.messages[buffer.lastListenTime][buffer.lastListenCount:]...)
	result = append(result, commonBuffer.messages[buffer.lastListenTime][buffer.lastListenCommonCount:]...)
	buffer.lastListenCount = len(buffer.messages[buffer.lastListenTime])
	buffer.lastListenCommonCount = len(commonBuffer.messages[buffer.lastListenTime])
	return result
}

func (w *AbstractWorld) Speak(speakerId, listenerId *int, content string) {
	if speakerId == nil {
		panic(ErrInvalidArgs)
	}

	buffer := w.langReceiveBuffers[-1]
	if listenerId != nil {
		buffer = w.langReceiveBuffers[*listenerId]
	}

	buffer.addMessage(w.clock, speakerId, listenerId, content)
}
