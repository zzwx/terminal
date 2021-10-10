// Package terminal provides console terminal codes for windows
package terminal

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"runtime"

	"github.com/mattn/go-colorable"
	"golang.org/x/term"
)

// TODO: Check out
// https://github.com/gosuri/uilive
// https://github.com/rthornton128/goncurses
// https://github.com/inancgumus/screen
// https://github.com/buger/goterm
// https://github.com/c-bata/go-prompt
// https://github.com/baulk/bulk/tree/master/progressbar
// https://github.com/moby/term
// https://github.com/containerd/console

// Windows Terminal Sequences
// https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
const (
	ESC = "\x1b"    // ESC consists of one 0x1b symbol.
	CSI = ESC + "[" // CSI stands for Control Sequence Introducer, used in the majority of sequences.

	// [X] Reset       = CSI + "0m" - see ResetAll(), Reset
	// [X] Bright      = CSI + "1m" - see SetBright()
	// [?] Dim         = CSI + "2m" // Dim doesn't work in Windows. Use SetBright instead.
	// [X] NoBright    = CSI + "22m" - see SetBright()
	// [X] Underline   = CSI + "4m"  // - see SetUnderline()
	// [X] NoUnderline = CSI + "24m" // - see SetUnderline()
	// [?] Blink       = CSI + "5m" // Blink doesn't seem to be supported in Windows.
	// [X] Swap        = CSI + "7m" - see Swap()
	// [?] Hidden      = CSI + "8m" - Doesn't seem to work under Windows
	// [X] CancelSwap  = CSI + "27m" - see CancelSwap()

	FgBlack   = CSI + "30m"
	FgRed     = CSI + "31m"
	FgGreen   = CSI + "32m"
	FgYellow  = CSI + "33m"
	FgBlue    = CSI + "34m"
	FgMagenta = CSI + "35m"
	FgCyan    = CSI + "36m"
	FgWhite   = CSI + "37m"

	// [X] 38 ; 2 ; <r> ; <g> ; <b>m	 | Set foreground color to RGB value specified in <r>, <g>, <b> parameters*
	// [X] 48 ; 2 ; <r> ; <g> ; <b>m   | Set background color to RGB value specified in <r>, <g>, <b> parameters*
	// [ ] 38m | Applies extended color value to the foreground (see details below)
	// [ ] 39m | Applies only the foreground portion of the defaults (see 0)

	BgBlack   = CSI + "40m"
	BgRed     = CSI + "41m"
	BgGreen   = CSI + "42m"
	BgYellow  = CSI + "43m"
	BgBlue    = CSI + "44m"
	BgMagenta = CSI + "45m"
	BgCyan    = CSI + "46m"
	BgWhite   = CSI + "47m"

	// [ ] 48m // Applies extended color value to the background (see details below)
	// [ ] 49m // Applies only the background portion of the defaults (see 0)

	FgHiBlack   = CSI + "90m"
	FgHiRed     = CSI + "91m"
	FgHiGreen   = CSI + "92m"
	FgHiYellow  = CSI + "93m"
	FgHiBlue    = CSI + "94m"
	FgHiMagenta = CSI + "95m"
	FgHiCyan    = CSI + "96m"
	FgHiWhite   = CSI + "97m"

	BgHiBlack   = CSI + "100m"
	BgHiRed     = CSI + "101m"
	BgHiGreen   = CSI + "102m"
	BgHiYellow  = CSI + "103m"
	BgHiBlue    = CSI + "104m"
	BgHiMagenta = CSI + "105m"
	BgHiCyan    = CSI + "106m"
	BgHiWhite   = CSI + "107m"

	// Reset returns text formatting attributes to the default state
	// prior to modification, which includes Bold, Underline, Negative
	// and color choice.
	//
	// Reset is equal to Terminal.Reset, just shorter to write in string concatenations.
	Reset = CSI + "0m"
)

var fgColors []string
var bgColors []string

func init() {
	fgColors = []string{
		FgBlack,
		FgRed,
		FgGreen,
		FgYellow,
		FgBlue,
		FgMagenta,
		FgCyan,
		FgWhite,
		FgHiBlack,
		FgHiRed,
		FgHiGreen,
		FgHiYellow,
		FgHiBlue,
		FgHiMagenta,
		FgHiCyan,
		FgHiWhite,
	}
	bgColors = []string{
		BgBlack,
		BgRed,
		BgGreen,
		BgYellow,
		BgBlue,
		BgMagenta,
		BgCyan,
		BgWhite,
		BgHiBlack,
		BgHiRed,
		BgHiGreen,
		BgHiYellow,
		BgHiBlue,
		BgHiMagenta,
		BgHiCyan,
		BgHiWhite,
	}
}

// FgColors returns all supported named colors
func FgColors() []string {
	return fgColors
}

// BgColors returns all supported named colors
func BgColors() []string {
	return bgColors
}

/*
┌┬─┐
├┼─┤
││ │
└┴─┘
*/

// Terminal is ready to use with os.Stdout if not initialized using
// NewTerminal with the output file description.
//
// Special sequences are available through terminal package constants
// and functions. terminal package functions may be used as well if desired.
//
// It's safe to output special sequences even if destination
// appears to be a non-terminal, as special sequences will be automatically
// discarded. Since no special meaning will be applied, like cursor moving,
// checking for !Terminal.IsTerminal could be utilized to provide a different
// formatting.
type Terminal struct {
	f      *os.File
	out    io.Writer
	once   sync.Once
	raw    *term.State
	isTerm bool
}

func (t *Terminal) Write(p []byte) (n int, err error) {
	return t.Print(string(p)) // TODO: Thoroughly test whether incomplete utf bytes could cause an issue
}

// NewTerminal returns a new Terminal instance attached
// to specified file, typically os.Stdout.
//
// NewTerminal(os.Stdout) is equal to simply using a Terminal
// instance.
func NewTerminal(f *os.File) *Terminal {
	if f == nil {
		panic("destination f can't be nil")
	}
	t := &Terminal{f: f}
	runtime.SetFinalizer(t, syncFileFinalizer)
	return t
}

func syncFileFinalizer(t *Terminal) {
	if t != nil && t.f != nil {
		_ = t.f.Sync()
		_ = t.f.Close()
		runtime.SetFinalizer(t, nil)
	}
}

func (t *Terminal) OverrideOut(out io.Writer) {
	t.out = out
}

// SetRaw puts the terminal connection into raw mode or back.
func (t *Terminal) SetRaw(raw bool) {
	t.init()
	if !t.IsTerminal() {
		return
	}
	if t.raw != nil {
		return
	}
	if raw {
		if t.raw == nil { // Already raw
			r, err := term.MakeRaw(int(t.f.Fd()))
			if err != nil {
				log.Fatalf("can't make terminal raw due to: %v", err)
			}
			t.raw = r
		}
	} else {
		if t.raw != nil { // Not been set to raw
			err := term.Restore(int(t.f.Fd()), t.raw)
			if err != nil {
				log.Fatalf("can't restore terminal from raw state due to: %v", err)
			}
			t.raw = nil
		}
	}
}

func (t *Terminal) init() {
	t.once.Do(func() {
		if t.out == nil {
			if t.f == nil {
				t.f = os.Stdout
			}
			//if isatty.IsTerminal(t.f.Fd()) {
			if IsTerminal(int(t.f.Fd())) {
				if runtime.GOOS == "windows" /*&& os.Getenv("TERM") != ""*/ {
					// TODO: Find a way to say if this windows version already supports terminal commands
					err := EnableVirtualTerminalProcessing(t.f, true)
					if err != nil {
						log.Fatalf("Can't enable virtual terminal processing due to %v", err)
					}
					t.out = t.f
					t.isTerm = true
				} else {
					t.out = colorable.NewColorable(t.f)
				}
			}
			if t.out == nil {
				if os.Getenv("FORCE_TERMINAL_SEQUENCES") == "1" {
					t.out = colorable.NewColorable(t.f)
				} else {
					t.out = colorable.NewNonColorable(t.f)
				}
			}
		}
	})
}

func (t *Terminal) Printf(format string, a ...interface{}) (n int, err error) {
	t.init()
	return fmt.Fprintf(t.out, format, a...)
}

func (t *Terminal) Println(a ...interface{}) (n int, err error) {
	t.init()
	return fmt.Fprintln(t.out, a...)
}

func (t *Terminal) Print(a ...interface{}) (n int, err error) {
	t.init()
	return fmt.Fprint(t.out, a...)
}

// SameLinePrintf erases last line and outputs format with provided variables.
// For a non-terminal case it uses "\r", otherwise MoveToX(0).
//
// Appending "\n" will retain this line in place during the next SameLineTerminalReport call.
func (t *Terminal) SameLinePrintf(format string, a ...interface{}) {
	t.init()
	t.MoveToX(0)
	if !t.IsTerminal() {
		t.Printf("\r")
	}
	p := fmt.Sprintf(format, a...)
	// To force line erasing at every line we have to split the output
	s := strings.Split(p, "\n")
	for i, s_ := range s {
		if i > 0 {
			t.Printf("\n")
		}
		t.Print(s_)
		t.EraseRestOfLine()
	}
}

// IsTerminal returns true if during initialization the output
// has been recognized as terminal. This may be used to make quick decisions
// as to how to output the result.
func (t *Terminal) IsTerminal() bool {
	t.init()
	return t.isTerm
}

// GetSize reports current (width, height) of the viewport\
// or 80,24 if it can't retrieve it
func (t *Terminal) GetSize() (w, h int) {
	t.init()
	w, h, err := term.GetSize(int(t.f.Fd()))
	if err != nil {
		w, h = 80, 24
	}
	return
}

// SetTitle sets console title which will be automatically restored
// by Windows when program exists.
func (t *Terminal) SetTitle(title string) {
	t.init()
	if t.IsTerminal() {
		// ESC ] 0 ; <string> BEL
		t.Printf(ESC + "]0;" + strings.ReplaceAll(title, "\x07", "") + "\x07")
	}
}

// Reset returns text formatting attributes to the default state
// prior to modification, which includes Bold, Underline, Negative
// and color choice.
//
// For quick concatenations use terminal.Reset.
func (t *Terminal) Reset() {
	t.Printf(Reset)
}
