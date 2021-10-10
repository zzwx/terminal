package tests

import (
	"fmt"
	"github.com/zzwx/terminal"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	sshT "golang.org/x/crypto/ssh/terminal"
)

func mainChat() {
	if err := chat(); err != nil {
		log.Fatal(err)
	}
}

func chat() error {
	f := os.Stdout
	if !sshT.IsTerminal(int(f.Fd())) {
		return fmt.Errorf("stdin/stdout should be terminal")
	}
	oldState, err := sshT.MakeRaw(int(f.Fd()))
	if err != nil {
		return err
	}
	defer sshT.Restore(int(f.Fd()), oldState)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := sshT.NewTerminal(screen, "> ")
	//term.SetPrompt(string(term.Escape.Red) + "> " + string(term.Escape.Reset))
	rePrefix := string(term.Escape.Cyan) + "Human says:" + string(term.Escape.Reset)
	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if line == "" {
			continue
		}
		fmt.Fprintln(term, rePrefix, line)
	}
}

func TestTerminal2(t_ *testing.T) {
	t := terminal.NewTerminal(os.Stdout)
	w, h := t.GetSize()
	t.StartAlternativeBuffer().
		SetCursorVisible(false)
	t.Printf(terminal.BgWhite + terminal.FgRed + "Test")
	for i := 0; i < h; i++ {
		//t.Printf(terminal.MoveByY(-1))
		//t.Printf(terminal.ESC + "M")
		t.MoveUpScroll()
		//t.Printf(terminal.ESC + "C" + ".")
		//time.Sleep(time.Duration(100) * time.Millisecond)
	}
	//┌┬─┐
	//├┼─┤
	//││ │
	//└┴─┘
	t.MoveToXY(0, 0)
	t.Printf("┌")
	for x := 1; x < w-1; x++ {
		t.Printf("─")
	}
	for y := 1; y < h-1; y++ {
		t.MoveToY(y).MoveToX(0)
		t.Printf("│")
		t.MoveToX(w - 1)
		t.Printf("│")
	}
	t.MoveToXY(w-1, 0)
	t.Printf("┐")
	t.MoveToXY(w-1, h-1)
	t.Printf("┘")
	t.MoveToXY(0, h-1)
	t.Printf("└")
	t.MoveToX(1)
	for x := 1; x < w-1; x++ {
		t.Printf("─")
	}
	for b := 0; b < len(terminal.BgColors()); b++ {
		for f := 0; f < len(terminal.FgColors()); f++ {
			t.SetBright((f+b)%2 == 0)
			t.Printf(terminal.BgColors()[b])
			t.Printf(terminal.FgColors()[f])
			t.MoveToXY(b*2+1, f+1)
			t.Printf("XX")
		}
	}
	t.Reset()

	t.SetRaw(true)

	var buffer [1]byte
	_, _ = os.Stdin.Read(buffer[:])

	t.SetRaw(false)

	t.EndAlternativeBuffer()
	t.Reset()
}

func TestTerminal(t_ *testing.T) {
	var t terminal.Terminal

	t.StartAlternativeBuffer()

	w, h := t.GetSize()

	t.SetCursorVisible(false)

	t.Printf("Start (%d,%d)", w, h)
	end := fmt.Sprintf(terminal.SetBright(true)+"End (%d,%d)"+terminal.SetBright(false), w, h)
	t.MoveToY(h - 1)

	//  01234     w = 5
	// |     |
	// |   xx| l = 2. 5-2= 3
	t.MoveToX(w - len(end))
	t.SavePos()
	t.Printf(end)
	time.Sleep(time.Second)

	t.RestorePos()

	t.MoveByY(-1)
	t.Printf("Nice!")
	t.MoveByX(0)
	time.Sleep(time.Second)

	t.EraseRestOfScreen()
	time.Sleep(time.Second)
	t.Printf("And?")

	t.MoveToXY(w-1, 0)
	t.Printf("x")

	time.Sleep(time.Second * 5)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			t.MoveToXY(x, y)
			if x == 0 || y == 0 || x == w-1 || y == h-1 {
				t.Printf(terminal.BgBlue + "|")
			} else {
				v := strconv.Itoa(y)
				if (x+y)%2 == 0 {
					t.Printf(terminal.BgBlack + v[len(v)-1:])
				} else {
					t.Printf(terminal.BgBlack + " ")
				}
			}
		}
	}
	t.SetScrollRegion(1, h-2)
	t.MoveToY(h / 2)
	t.MoveToX(5)

	time.Sleep(time.Second)

	for i := 0; i < w; i++ {
		t.EraseShiftLeft(1)
		time.Sleep(time.Millisecond * 20)
	}

	for i := 0; i < 10; i++ {
		t.Printf(terminal.ScrollBy(-1))
		time.Sleep(time.Millisecond * 20)
	}

	for i := 0; i < h; i++ {
		if i%2 == 0 {
			t.SetBright(true)
		} else {
			t.SetBright(false)
		}
		t.Swap()
		t.Printf("%d", i)

		time.Sleep(time.Millisecond * 20)
		t.Printf("\n")
	}
	t.CancelSwap()
	t.SetUnderline(false)
	t.Printf(terminal.BgGreen + terminal.FgBlue)
	for i := 0; i < h; i++ {
		t.ScrollBy(1)
		time.Sleep(time.Millisecond * 20)
	}
	t.Reset()
	t.MoveToXY(w/2, h/2)

	t.EraseRestOfScreen()
	time.Sleep(time.Second * 2)

	//var c [1]byte
	//os.Stdin.Read(c[:])
	t.EndAlternativeBuffer()
	t.SetCursorVisible(true)
	t.Reset()
}
