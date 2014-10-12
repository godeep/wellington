package sprite_sass

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/seateam/color"
)

type lError struct {
	Pos     int
	Message string
}

type lErrors struct {
	Errors []lError
	Pos    int
}

func (ctx *Context) ErrorTokenizer(src string) lErrors {
	errors := []lError{}
	r := strings.NewReplacer(":", " ", ",", " ")
	src = r.Replace(src)
	scanner := bufio.NewScanner(strings.NewReader(src))
	scanner.Split(bufio.ScanWords)
	var (
		line int64
		str  string
	)
	for scanner.Scan() {
		var err error

		if scanner.Text() != "Backtrace" && scanner.Text() != "stdin" {
			str += scanner.Text() + " "
		} else {
			if line == 0 && str == "" {
			} else {
				le := lError{int(line - 1), strings.TrimSpace(str)}
				errors = append(errors, le)
				str = ""
				line = 0
			}
		}

		if scanner.Text() == "stdin" {
			if scanner.Scan() {
				line, err = strconv.ParseInt(scanner.Text(), 10, 16)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	errors = append(errors, lError{int(line - 1), strings.TrimSpace(str)})
	ctx.errors = lErrors{
		Pos:    int(line),
		Errors: errors,
	}
	return ctx.errors
}

// Error reads the original libsass error and creates helpful debuggin
// information for debuggin that error.
func (ctx *Context) ProcessSassError(err string) {
	// Attempt to find the source error
	split := strings.Split(err, ":")
	if len(split) == 0 {
		return
	}
	eObj := ctx.ErrorTokenizer(err)
	pos := eObj.Pos
	// Decrement for $rel line
	pos = pos - 1
	lines := bytes.Split(ctx.Parser.Output, []byte("\n"))
	// Line number is off by one from libsass
	// Find previous lines to maximum available
	errLines := "" //"error in " + ctx.Parser.LookupFile(pos)
	red := color.NewStyle(color.BlackPaint, color.RedPaint).Brush()
	first := pos - 7
	if first < 0 {
		first = 0
	}
	last := pos + 7
	if last > len(lines) {
		last = len(lines)
	}
	for i := first; i < last; i++ {
		// translate 0 index to 1 index
		str := fmt.Sprintf("\n%3d: %s", i+1, lines[i])
		if i == pos-1 {
			str = red(str)
		}
		errLines += str
	}
	ctx.errorString = err + "\n" + errLines
}

func (ctx *Context) Error() string {
	return ctx.errorString
}
