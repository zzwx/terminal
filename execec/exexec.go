package execec

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/zzwx/terminal"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/mattn/go-colorable"
)

func Execute(ctx context.Context, alias string, fullCmd string) {
	fields := strings.Fields(fullCmd)
	if len(fields) > 0 {
		ExecuteFields(ctx, alias, fields[0], fields[1:]...)
	}
}

func ExecuteFields(ctx context.Context, alias string, path string, args ...string) {
	cmd := exec.CommandContext(ctx, path, args...)
	//cmd.Dir =
	//cmd.Env =
	Cmd(alias, cmd)
}

type dataWrap struct {
	base      string
	delimiter string
	bytes     []byte
}

var chStdOut = make(chan *dataWrap)
var chStdErr = make(chan *dataWrap)
var Done = make(chan bool)
var veryFirstLine int32

func init() {
	atomic.StoreInt32(&veryFirstLine, 1)
	tStdOut := terminal.NewTerminal(os.Stdout)
	tStdErr := terminal.NewTerminal(os.Stderr)
	go func() {
		for {
			select {
			case d := <-chStdOut:
				outputToTerminal(tStdOut, d)
			case d := <-chStdErr:
				outputToTerminal(tStdErr, d)
			case <-Done: // TODO: Find where to use Done to stop the loop
				close(chStdOut)
				close(chStdErr)
				return
			}
		}
	}()
}

var lastLastIsNewLine = false
var lastBase = ""
var lastDelimiter = ""

func outputToTerminal(term *terminal.Terminal, data *dataWrap) {
	if lastBase != "" && (!lastLastIsNewLine && (lastBase != data.base || lastDelimiter != data.delimiter)) {
		_, _ = term.Write([]byte{'\n'})
		//term.Write([]byte(data.base + data.delimiter))
	}
	betweenRs := bytes.Split(data.bytes, []byte{'\r'})
	for _, r := range betweenRs {
		_, _ = term.Write([]byte{'\r'})
		_, _ = term.Write([]byte(data.base + data.delimiter))
		_, _ = term.Write(r)
	}
	lastBase = data.base
	lastDelimiter = data.delimiter
	lnl := false
	if len(data.bytes) > 0 {
		lnl = data.bytes[len(data.bytes)-1] == '\n'
	}
	lastLastIsNewLine = lnl
}

type rgb struct{ r, g, b int }

var fgColors = make(map[int]rgb)
var fgColorsMu sync.Mutex
var randSource = rand.New(rand.NewSource(50))

func allocateFgColor(proc int) rgb {
	fgColorsMu.Lock()
	defer fgColorsMu.Unlock()
	if c, ok := fgColors[proc]; ok {
		return c
	}
	c := rgb{50 + randSource.Intn(205), 50 + randSource.Intn(205), 50 + randSource.Intn(205)}
	fgColors[proc] = c
	return c
}

var maxWidth = 0
var maxWidthMu sync.Mutex

// reMaxWidth(0) will shortcut to omit mutex and
// simply return the latest value
func reMaxWidth(possibleNewMaxWidth int) int {
	if possibleNewMaxWidth == 0 {
		return maxWidth
	}
	maxWidthMu.Lock()
	defer maxWidthMu.Unlock()
	if possibleNewMaxWidth > maxWidth {
		maxWidth = possibleNewMaxWidth
	}
	return maxWidth
}

func merge(v ...[]string) []string {
	kv := make(map[string]string)
	for _, vv := range v {
		for _, v := range vv {
			sp := strings.Split(v, "=")
			if len(sp) == 2 {
				kv[sp[0]] = sp[1]
			}
		}
	}
	var result = make([]string, 0, len(kv))
	for s, s2 := range kv {
		result = append(result, s+"="+s2)
	}
	return result
}

func Cmd(alias string, cmd *exec.Cmd) error {
	cmd.Env = merge(os.Environ(), cmd.Env, []string{"FORCE_TERMINAL_SEQUENCES=1"}) // Tick our Terminal to use sequences.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("can't create stdout pipe for %v %v: %w", cmd.Path, cmd.Args, err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("can't create stderr pipe for %v %v: %w", cmd.Path, cmd.Args, err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("can't start %v %v: %w", cmd.Path, cmd.Args, err)
	}
	var base string
	if alias != "" {
		base = alias
	} else {
		base = filepath.Base(cmd.Path)
		base = strings.TrimSuffix(base, ".exe")
		base = strings.TrimSuffix(base, ".cmd")
	}

	base += ":" + strconv.Itoa(cmd.Process.Pid)

	var wg sync.WaitGroup
	wg.Add(3)
	fg := allocateFgColor(cmd.Process.Pid)

	go func() {
		defer wg.Done()
		ioCopy(terminal.FgRGB(fg.r, fg.g, fg.b)+base+terminal.Reset, " | ", chStdOut, stdout)
	}()

	go func() {
		defer wg.Done()
		ioCopy(terminal.FgRGB(fg.r, fg.g, fg.b)+base+terminal.Reset, terminal.FgRed+" | "+terminal.Reset, chStdErr, stderr)
	}()

	waitErr := make(chan error)
	go func() {
		defer wg.Done()
		waitErr <- cmd.Wait()
	}()

	err = <-waitErr
	//err = cmd.Wait()
	if err != nil {
		//output, _ := cmd.CombinedOutput()
		//fmt.Println("error: " + err.Error())
		ioCopy(terminal.FgRGB(fg.r, fg.g, fg.b)+base+terminal.Reset, terminal.FgRed+" | "+terminal.Reset, chStdErr,
			strings.NewReader(fmt.Sprintf("%v\n%v\n", strings.Join(cmd.Args, " "), err)))
	}
	wg.Wait()
	if err != nil {
		return err
	}
	if cmd.ProcessState.ExitCode() != 0 {
		return errors.New("exit status != 0")
	}
	return nil
}

func ioCopy(base string, delimiter string, dst chan *dataWrap, r io.Reader) {
	baseLength := baseLengthReMax(base)
	buf := bufio.NewReader(r)
	// We'll be reading from io.Read until either \n found or a specific amount of time has passed after last
	// flush to "dst".
	type readRuneResult struct {
		rn   rune
		size int
	}
	runeCh := make(chan readRuneResult)
	exit := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			rn, size, err := buf.ReadRune()
			if size > 0 {
				runeCh <- readRuneResult{rn, size}
			}
			if err != nil {
				if err == io.EOF {
					//runeCh <- readRuneResult{'E', 1}
					//runeCh <- readRuneResult{'O', 1}
					//runeCh <- readRuneResult{'F', 1}
					//runeCh <- readRuneResult{'\n', 1}

					exit <- true
				}
				return
			}
		}
	}()
	go func() {
		var b bytes.Buffer
		defer wg.Done()
		for {
			t := time.After(time.Millisecond * 4150)
			select {
			case v := <-runeCh:
				if v.size > 0 {
					b.WriteRune(v.rn)
					if v.rn == '\n' {
						str := b.String()
						b.Reset()
						outputAccumulated(dst, base+strings.Repeat(" ", reMaxWidth(0)-baseLength), delimiter, str)
					} else if v.rn == '\r' {
						l, _ := b.ReadString('\r')
						if len(l) > 1 {
							// We only care if there was something before \r,
							// We'll issue everything minus the last \r, which we'll restore back till next time around.
							outputAccumulated(dst, base+strings.Repeat(" ", reMaxWidth(0)-baseLength), delimiter, l[:len(l)-1])
						}
						// Restore consumed \r
						b.Write([]byte{'\r'})
					}
				}
			case <-t:
				// Speed up slower output that hasn't yet ended with \n
				str := b.String()
				b.Reset()
				outputAccumulated(dst, base+strings.Repeat(" ", reMaxWidth(0)-baseLength), delimiter, str)
			case <-exit:
				// Output anything that left over as it might have never ended with \n
				str := b.String()
				b.Reset()
				outputAccumulated(dst, base+strings.Repeat(" ", reMaxWidth(0)-baseLength), delimiter, str)
				return
			}
		}
	}()
	wg.Wait()
}

func baseLengthReMax(base string) int {
	var ba bytes.Buffer
	w := colorable.NewNonColorable(&ba)
	_, _ = w.Write([]byte(base))
	baseLength := utf8.RuneCount(ba.Bytes()) // TODO: Use a terminal string length rather than rune count
	reMaxWidth(baseLength)                   // Trigger only
	return baseLength
}

func outputAccumulated(dst chan *dataWrap, base string, delimiter string, str string) {
	if len(str) > 0 {
		dst <- &dataWrap{base, delimiter, []byte(str)}
	}
}
