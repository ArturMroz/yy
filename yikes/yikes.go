package yikes

import (
	"fmt"
	"strings"
)

type YYError struct {
	Msg    string // description of error
	Offset int    // error occurred after reading Offset bytes
}

func (e *YYError) Error() string { return e.Msg }

func PrettyError(src []byte, offset int, errMsg string) string {
	if offset < 0 {
		return "error: " + errMsg
	}

	line := 1
	col := 0
	lastNewLineStart := 0

	i := 0
	for ; i < offset; i++ {
		if src[i] == '\n' {
			line++
			col = 0
			lastNewLineStart = i
		} else {
			col++
		}
	}

	if lastNewLineStart > 0 {
		lastNewLineStart++
	}

	lastNewLineEnd := i
	for lastNewLineEnd < len(src) && src[lastNewLineEnd] != '\n' {
		lastNewLineEnd++
	}

	// fmt.Println("line", line, "col", col, "offset", offset, "start", lastNewLineStart, "end", lastNewLineEnd, "curToken", tok)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("error: %s\n", errMsg))
	b.WriteString(fmt.Sprintf("%3d | %s\n", line, src[lastNewLineStart:lastNewLineEnd]))
	b.WriteString(fmt.Sprintf("      %s^", strings.Repeat(" ", col)))

	return b.String()
}
