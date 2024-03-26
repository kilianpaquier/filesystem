// Code copied from https://github.com/andreyvit/diff; DO NOT EDIT.
// FIXME optimize code one day ...

package tests

import (
	"bytes"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffsToPatch transforms an input slice of diffs into a human-readable string with diffs together.
func DiffsToPatch(diffs []diffmatchpatch.Diff) string {
	patch := patchBuilder{Output: make([]string, 0, len(diffs))}
	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		patch.addChars(lines[0], diff.Type)

		for _, line := range lines[1:] {
			patch.addLine(diff.Type)
			patch.addChars(line, diff.Type)
		}
	}
	patch.flush()
	return strings.Join(patch.Output, "\n")
}

type patchBuilder struct {
	NewLineBuffer bytes.Buffer
	NewLines      []string

	OldLineBuffer bytes.Buffer
	OldLines      []string

	Output []string
}

func (p *patchBuilder) addLine(diff diffmatchpatch.Operation) {
	old := p.OldLineBuffer.String()
	new := p.NewLineBuffer.String()

	if diff == diffmatchpatch.DiffEqual && old == new {
		p.flushChunk()
		p.Output = append(p.Output, "  "+new)
		p.OldLineBuffer.Reset()
		p.NewLineBuffer.Reset()
		return
	}

	if diff == diffmatchpatch.DiffDelete || diff == diffmatchpatch.DiffEqual {
		p.OldLines = append(p.OldLines, "- "+old)
		p.OldLineBuffer.Reset()
	}

	if diff == diffmatchpatch.DiffInsert || diff == diffmatchpatch.DiffEqual {
		p.NewLines = append(p.NewLines, "+ "+new)
		p.NewLineBuffer.Reset()
	}
}

func (p *patchBuilder) addChars(line string, diff diffmatchpatch.Operation) {
	switch diff {
	case diffmatchpatch.DiffEqual:
		p.NewLineBuffer.WriteString(line)
		p.OldLineBuffer.WriteString(line)
	case diffmatchpatch.DiffDelete:
		p.OldLineBuffer.WriteString(line)
	case diffmatchpatch.DiffInsert:
		p.NewLineBuffer.WriteString(line)
	}
}

func (p *patchBuilder) flush() {
	switch {
	case p.OldLineBuffer.Len() > 0 && p.NewLineBuffer.Len() > 0:
		p.addLine(diffmatchpatch.DiffEqual)
	case p.OldLineBuffer.Len() > 0:
		p.addLine(diffmatchpatch.DiffDelete)
	case p.NewLineBuffer.Len() > 0:
		p.addLine(diffmatchpatch.DiffInsert)
	}
	p.flushChunk()
}

func (p *patchBuilder) flushChunk() {
	p.Output = append(p.Output, p.OldLines...)
	p.OldLines = nil

	p.Output = append(p.Output, p.NewLines...)
	p.NewLines = nil
}
