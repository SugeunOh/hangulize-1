/*
Package hangulize provides an automatic transcriber into Hangul for non-Korean
words.

Originally, Hangulize was implemented in Python in 2010.  This implementation
is a reboot with attractive feature improvements.

Original implementation: https://github.com/sublee/hangulize

Brian Jongseong Park proposed the seed idea of Hangulize on his Blog.

Post by Brian: http://iceager.egloos.com/2610028

*/
package hangulize

import (
	"fmt"
	"strings"
)

// Hangulize transcribes a non-Korean word into Hangul, the Korean alphabet:
//
//    Hangulize("ita", "gloria")
//    // Output: "글로리아"
//
func Hangulize(lang string, word string) string {
	spec, ok := LoadSpec(lang)
	if !ok {
		// spec not found
		return word
	}

	h := NewHangulizer(spec)
	return h.Hangulize(word)
}

// -----------------------------------------------------------------------------

// Hangulizer ...
type Hangulizer struct {
	spec *Spec
}

// NewHangulizer ...
func NewHangulizer(spec *Spec) *Hangulizer {
	return &Hangulizer{spec}
}

// Hangulize transcribes a loanword into Hangul.
func (h *Hangulizer) Hangulize(word string) string {
	return h.HangulizeTrace(word, nil)
}

// HangulizeTrace transcribes a loanword into Hangul.  During
// transcribing, it sends internal traces to the given channel.
func (h *Hangulizer) HangulizeTrace(word string, ch chan<- Trace) string {
	if ch != nil {
		defer close(ch)
	}
	trace(ch, word, "", "input")

	// mask := h.newMask(word)
	// fmt.Println(word, mask)

	word = h.normalize(word, ch)

	chunks := []Chunk{Chunk{word, 0}}

	chunks = Rewrite(chunks, h.spec.rewrite)
	fmt.Println("rewrite", chunks)

	chunks = Replace(chunks, h.spec.transcribe)
	fmt.Println("transcribe", chunks)

	for i, chunk := range chunks {
		if chunk.age != 0 {
			chunk.word = AssembleJamo(chunk.word, ch)
			chunks[i] = chunk
		}
	}

	fmt.Println("jamo", chunks)

	word = NewChunkBuilder(chunks).String()

	fmt.Println("join", word)

	word = h.cleanUp(word)
	return word
}

// -----------------------------------------------------------------------------

func (h *Hangulizer) newMask(word string) []bool {
	text := []rune(word)
	return make([]bool, len(text))
}

// -----------------------------------------------------------------------------

func (h *Hangulizer) normalize(word string, ch chan<- Trace) string {
	// TODO(sublee): Language-specific normalizer
	orig := word

	except := make([]string, 0)

	args := make([]string, 0)
	for to, froms := range h.spec.normalize {
		for _, from := range froms {
			args = append(args, from, to)
		}

		except = append(except, to)
	}
	rep := strings.NewReplacer(args...)
	word = rep.Replace(word)

	word = Normalize(word, RomanNormalizer{}, except)

	word = strings.ToLower(word)

	trace(ch, word, orig, "roman")
	return word
}

func (h *Hangulizer) cleanUp(word string) string {
	markers := h.spec.Config.Markers

	args := make([]string, len(markers)*2)
	for i, marker := range markers {
		args[i*2] = string(marker)
		args[i*2+1] = ""
	}
	rep := strings.NewReplacer(args...)

	var buf strings.Builder

	for _, ch := range word {
		rep.WriteString(&buf, string(ch))
	}

	return buf.String()
}
