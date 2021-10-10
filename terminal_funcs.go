package terminal

import "strconv"

// Functions from this file are used to generate methods
// in Terminal

// MoveByX moves cursor position by x difference. Negative means left, positive -
// right. Never passes the edges.
func MoveByX(xDiff int) string {
	switch {
	//---BUGGY--- - doesn't work in "Windows Terminal" as of Version: 1.5.10411.0
	//case xDiff == -1:
	//	// ESC D | Cursor Backward (Left) by 1
	//	return ESC + "D"
	case xDiff < 0:
		// ESC [ <n> D | Cursor backward (Left) by <n>
		return CSI + strconv.Itoa(-xDiff) + "D"
	//---BUGGY--- - doesn't work in "Windows Terminal" as of Version: 1.5.10411.0
	//case xDiff == 1:
	//	// ESC C | Cursor forward by 1 (not passing the edge)
	//	return ESC + "C"
	case xDiff > 0:
		// ESC [ <n> C | Cursor forward (Right) by <n>
		return CSI + strconv.Itoa(xDiff) + "C"
	}
	return ""
}

// MoveByY moves cursor position by yDiff difference. Negative means up, positive
// down. Doesn't cause scrolling.
func MoveByY(yDiff int) string {
	switch {
	//---BUGGY--- - doesn't work in "Windows Terminal" as of Version: 1.5.10411.0
	//case yDiff == -1:
	//	// ESC A | Cursor Up by 1 (no scrolling)
	//	return ESC + "A"
	case yDiff < 0:
		// ESC [ <n> A | Cursor up by <n>
		return CSI + strconv.Itoa(-yDiff) + "A"
	//---BUGGY--- - doesn't work in "Windows Terminal" as of Version: 1.5.10411.0
	//case yDiff == 1:
	//	// ESC B | Cursor Down by 1 (no scrolling)
	//	return ESC + "B"
	case yDiff > 0:
		// ESC [ <n> B | Cursor down by <n>
		return CSI + strconv.Itoa(yDiff) + "B"
	}
	return ""
}

// MoveUpScroll ("Reverse Index") moves up maintaining x cursor position.
// Upon reaching the top of the screen it begins appending empty
// lines with the current background color.
func MoveUpScroll() string {
	// ESC M | Reverse Index – Performs the reverse operation of \n, moves cursor up
	// one line, maintains horizontal position, scrolls buffer if necessary*
	return ESC + "M"
}

// MoveNextLineBy moves the cursor down by amount,
// to the first column, without scrolling.
func MoveNextLineBy(amount int) string {
	if amount < 0 {
		return ""
	}
	// ESC [ <n> E | Cursor down <n> lines from current position
	return CSI + strconv.Itoa(amount) + "E"
}

// MovePreviousLineBy moves the cursor up by amount,
// to the first column, without scrolling.
func MovePreviousLineBy(amount int) string {
	if amount < 0 {
		return ""
	}
	// ESC [ <n> F | Cursor up <n> lines from current position
	return CSI + strconv.Itoa(amount) + "F"
}

// MoveToXY moves cursor to absolute x.y. Accepts numbers from (0,0) as top left
// corner.
func MoveToXY(x, y int) string {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	// Seem to be doing the same:
	// ESC [ <y> ; <x> H
	// ESC [ <y> ; <x> f
	return CSI + strconv.Itoa(y+1) + ";" + strconv.Itoa(x+1) + "H"
}

// MoveTopLeft moves the cursor to absolute (0,0) corner of the screen, equals to MoveToXY(0,0).
func MoveTopLeft() string {
	return CSI + "H"
}

// MoveToX moves cursor to absolute x column, starting from 0 as left-most column.
func MoveToX(x int) string {
	if x < 0 {
		return ""
	}
	// ESC [ <n> G | Cursor moves to <n>th position horizontally in the current line
	return CSI + strconv.Itoa(x+1) + "G"
}

// MoveToY moves cursor to absolute y row, starting from 0 as left-most column.
func MoveToY(y int) string {
	if y < 0 {
		return ""
	}
	// ESC [ <n> d | Cursor moves to the <n>th position vertically in the current
	// column
	return CSI + strconv.Itoa(y+1) + "d"
}

// SavePos issues terminal command to save cursor position for upcoming
// RestorePos.
func SavePos() string {
	// ESC 7 | Save Cursor Position in Memory
	// ESC [ s
	return CSI + "s"
}

// RestorePos issues terminal command to restore cursor position saved previously
// using SavePos.
func RestorePos() string {
	// ESC 8 | Restore Cursor Position from Memory
	// ESC [ u
	return CSI + "u"
}

// SetCursorVisible sets cursor visibility.
func SetCursorVisible(visible bool) string {
	if visible {
		// ESC [ ? 25 h
		return CSI + "?25h"
	}
	// ESC [ ? 25 l
	return CSI + "?25l"
}

// SetBlinking sets cursor blinking on / off.
func SetBlinking(on bool) string {
	if on {
		// ESC [ ? 12 h
		return CSI + "?12h"
	}
	// ESC [ ? 12 l
	return CSI + "?12l"
}

// SetBright sets bright / bold flag to foreground color.
func SetBright(on bool) string {
	if on {
		return CSI + "1m"
	}
	return CSI + "22m"
}

// SetUnderline sets font with underline.
func SetUnderline(on bool) string {
	if on {
		return CSI + "4m"
	}
	return CSI + "24m"
}

// SetScrollRegion sets the region for scrolling using ScrollBy, by specifying
// top and bottom fixed areas. Scrolling also happens if \n is printed at the
// last line of the scroll region or MoveUpScroll at the top of it, filling the
// gap of the opposite side of the scrolling region with an empty line with the
// current background color.
//
// h==0 means no region on top is set aside as fixed.
//
// h==1 means first row will be fixed.
//
// b==h-1 means now bottom region is going to be fixed during scrolling.
//
// b==h-2 means one last row will be fixed during scrolling.
func SetScrollRegion(h, b int) string {
	// ESC [ <t> ; <b> r
	return CSI + strconv.Itoa(h+1) + ";" + strconv.Itoa(b+1) + "r"
}

// ScrollBy will scroll from the current vertical cursor position. The text will
// go up for diff < 0 and down for diff > 0. Empty lines will be added to fill
// the gap with the current background color. The area affected can be controlled
// using SetScrollRegion.
func ScrollBy(yDiff int) string {
	if yDiff < 0 {
		// ESC [ <n> S | Scroll text up by <n>. Also known as pan down, new lines fill in
		// from the bottom of the screen
		return CSI + strconv.Itoa(-yDiff) + "S"
	} else if yDiff > 0 {
		//ESC [ <n> T Scroll down by <n>. Also known as pan up, new lines fill in from
		//the top of the screen
		return CSI + strconv.Itoa(yDiff) + "T"
	}
	return ""
}

// EraseRestOfLine clears from current cursor position (including) to the end of
// line without moving the cursor. Clearing happens with the current background
// color.
func EraseRestOfLine() string {
	return CSI + "K" // Equals to "0K"
}

// EraseRestOfScreen clears from current cursor position (including) to the right
// and down until the bottom right of the screen without moving the cursor.
// Clearing happens with the current background color.
func EraseRestOfScreen() string {
	return CSI + "J" // Equals to "0J"
}

// EraseFrontOfLine erases from the beginning of the current line to and
// including current cursor position without moving the cursor. Clearing happens
// with the current background color.
func EraseFrontOfLine() string {
	return CSI + "1K"
}

// EraseFrontOfScreen erases from the top left of the screen to and including
// current cursor position without moving the cursor. Clearing happens with the
// current background color.
func EraseFrontOfScreen() string {
	return CSI + "1J"
}

// EraseLine erases the whole current line without moving the cursor.
// Clearing happens with the current background color.
func EraseLine() string {
	return CSI + "2K"
}

// EraseScreen clears the while screen without moving the cursor.
// Clearing happens with the current background color.
func EraseScreen() string {
	return CSI + "2J"
}

// StartAlternativeBuffer clears the screen, moves the cursor to (0,0) and allows
// to return to the original buffer using EndAlternativeBuffer. This allows for
// isolated modifications.
func StartAlternativeBuffer() string {
	return CSI + "?1049h"
}

// EndAlternativeBuffer returns back to the screen before StartAlternativeBuffer
// as it was left off. If no modifications were made before
// StartAlternativeBuffer to the cursor visibility and colors, sending Reset is
// unnecessary. If the program interrupts in the middle, it seems necessary
// to implicitly call EndAlternativeBuffer, otherwise the console will print
// out the prompt with the alternative buffer settings, and at least in case of cmd.exe
// continues typing with these settings.
func EndAlternativeBuffer() string {
	return CSI + "?1049l"
}

// ShiftRight moves current line by amount from the current column position.
// Spaces will be added to fill the gap, and anything going beyond the borders of
// viewport will be trimmed.
func ShiftRight(amount int) string {
	if amount > 0 {
		// ESC [ <n> @ | Insert <n> spaces at the current cursor position, shifting all
		// existing text to the right. Text exiting the screen to the right is removed.
		return CSI + strconv.Itoa(amount) + "@"
	}
	return ""
}

// EraseShiftLeft deletes amount of characters at the current cursor position,
// shifting in space character from the right edge of the viewport.
func EraseShiftLeft(amount int) string {
	if amount > 0 {
		// ESC [ <n> P | Delete <n> characters at the current cursor position, shifting
		// in space characters from the right edge of the screen.
		return CSI + strconv.Itoa(amount) + "P"
	}
	return ""
}

// Erase erases amount characters from the current cursor position without moving
// the cursor by overwriting characters with a space character and not wrapping
// after reaching the right screen border.
func Erase(amount int) string {
	if amount > 0 {
		// ESC [ <n> X | Erase <n> characters from the current cursor position by
		// overwriting them with a space character.
		return CSI + strconv.Itoa(amount) + "X"
	}
	return ""
}

// ShiftDown shifts the current line down by amount, adding empty line(s) to fill
// the gap. Scrolling margins set with SetScrollRegion will be respected.
func ShiftDown(amount int) string {
	if amount > 0 {
		// ESC [ <n> L | Inserts <n> lines into the buffer at the cursor position. The
		// line the cursor is on, and lines below it, will be shifted downwards.
		return CSI + strconv.Itoa(amount) + "L"
	}
	return ""
}

// DeleteLines deletes amount of lines from the buffer, starting with the row the
// cursor is on. Scrolling margins set with SetScrollRegion will be respected.
func DeleteLines(amount int) string {
	if amount > 0 {
		// ESC [ <n> M | Deletes <n> lines from the buffer, starting with the row the
		// cursor is on.
		return CSI + strconv.Itoa(amount) + "M"
	}
	return ""
}

// Swap swaps foreground and background colors.
// This actually seems to swap the meaning of fg and bg and can be stacked.
// Output CancelSwap to return to normal.
func Swap() string {
	return CSI + "7m"
}

// CancelSwap returns foreground/background to normal after any proceeding Swap.
func CancelSwap() string {
	return CSI + "27m"
}

func FgRGB(r, g, b int) string {
	return CSI + "38;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m"
}

func BgRGB(r, g, b int) string {
	return CSI + "48;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m"
}
