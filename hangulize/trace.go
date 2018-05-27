package hangulize

import "fmt"

// Trace is emitted when a replacement occurs.  It is used for tracing of
// Hangulize pipeline internal.
type Trace struct {
	step string
	why  string
	word string
}

func (t *Trace) String() string {
	return fmt.Sprintf("[%s] %#v %s", t.step, t.word, t.why)
}

type Tracer struct {
	traces   []Trace
	lastWord string
}

func (tr *Tracer) Traces() []Trace {
	return tr.traces
}

func (tr *Tracer) trace(step, why, word string) {
	if word == tr.lastWord {
		return
	}
	tr.traces = append(tr.traces, Trace{step, why, word})
	tr.lastWord = word
}

func (tr *Tracer) TraceWord(step, why, word string) {
	if tr == nil {
		return
	}
	tr.trace(step, why, word)
}

func (tr *Tracer) TraceSubwords(step, why string, subwords []Subword) {
	if tr == nil {
		return
	}
	word := NewSubwordsBuilder(subwords).String()
	tr.trace(step, why, word)
}