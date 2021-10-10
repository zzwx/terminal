// Code generated by internal/generate.go; DO NOT EDIT.

package terminal

// MoveByX moves cursor position by x difference. Negative means left, positive -
// right. Never passes the edges.
func (t *Terminal) MoveByX(xDiff int) *Terminal {
	t.Print(MoveByX(xDiff))
	return t
}

// MoveByY moves cursor position by yDiff difference. Negative means up, positive
// down. Doesn't cause scrolling.
func (t *Terminal) MoveByY(yDiff int) *Terminal {
	t.Print(MoveByY(yDiff))
	return t
}

// MoveUpScroll ("Reverse Index") moves up maintaining x cursor position.
// Upon reaching the top of the screen it begins appending empty
// lines with the current background color.
func (t *Terminal) MoveUpScroll() *Terminal {
	t.Print(MoveUpScroll())
	return t
}

// MoveNextLineBy moves the cursor down by amount,
// to the first column, without scrolling.
func (t *Terminal) MoveNextLineBy(amount int) *Terminal {
	t.Print(MoveNextLineBy(amount))
	return t
}

// MovePreviousLineBy moves the cursor up by amount,
// to the first column, without scrolling.
func (t *Terminal) MovePreviousLineBy(amount int) *Terminal {
	t.Print(MovePreviousLineBy(amount))
	return t
}

// MoveToXY moves cursor to absolute x.y. Accepts numbers from (0,0) as top left
// corner.
func (t *Terminal) MoveToXY(x, y int) *Terminal {
	t.Print(MoveToXY(x, y))
	return t
}

// MoveTopLeft moves the cursor to absolute (0,0) corner of the screen, equals to MoveToXY(0,0).
func (t *Terminal) MoveTopLeft() *Terminal {
	t.Print(MoveTopLeft())
	return t
}

// MoveToX moves cursor to absolute x column, starting from 0 as left-most column.
func (t *Terminal) MoveToX(x int) *Terminal {
	t.Print(MoveToX(x))
	return t
}

// MoveToY moves cursor to absolute y row, starting from 0 as left-most column.
func (t *Terminal) MoveToY(y int) *Terminal {
	t.Print(MoveToY(y))
	return t
}

// SavePos issues terminal command to save cursor position for upcoming
// RestorePos.
func (t *Terminal) SavePos() *Terminal {
	t.Print(SavePos())
	return t
}

// RestorePos issues terminal command to restore cursor position saved previously
// using SavePos.
func (t *Terminal) RestorePos() *Terminal {
	t.Print(RestorePos())
	return t
}

// SetCursorVisible sets cursor visibility.
func (t *Terminal) SetCursorVisible(visible bool) *Terminal {
	t.Print(SetCursorVisible(visible))
	return t
}

// SetBlinking sets cursor blinking on / off.
func (t *Terminal) SetBlinking(on bool) *Terminal {
	t.Print(SetBlinking(on))
	return t
}

// SetBright sets bright / bold flag to foreground color.
func (t *Terminal) SetBright(on bool) *Terminal {
	t.Print(SetBright(on))
	return t
}

// SetUnderline sets font with underline.
func (t *Terminal) SetUnderline(on bool) *Terminal {
	t.Print(SetUnderline(on))
	return t
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
func (t *Terminal) SetScrollRegion(h, b int) *Terminal {
	t.Print(SetScrollRegion(h, b))
	return t
}

// ScrollBy will scroll from the current vertical cursor position. The text will
// go up for diff < 0 and down for diff > 0. Empty lines will be added to fill
// the gap with the current background color. The area affected can be controlled
// using SetScrollRegion.
func (t *Terminal) ScrollBy(yDiff int) *Terminal {
	t.Print(ScrollBy(yDiff))
	return t
}

// EraseRestOfLine clears from current cursor position (including) to the end of
// line without moving the cursor. Clearing happens with the current background
// color.
func (t *Terminal) EraseRestOfLine() *Terminal {
	t.Print(EraseRestOfLine())
	return t
}

// EraseRestOfScreen clears from current cursor position (including) to the right
// and down until the bottom right of the screen without moving the cursor.
// Clearing happens with the current background color.
func (t *Terminal) EraseRestOfScreen() *Terminal {
	t.Print(EraseRestOfScreen())
	return t
}

// EraseFrontOfLine erases from the beginning of the current line to and
// including current cursor position without moving the cursor. Clearing happens
// with the current background color.
func (t *Terminal) EraseFrontOfLine() *Terminal {
	t.Print(EraseFrontOfLine())
	return t
}

// EraseFrontOfScreen erases from the top left of the screen to and including
// current cursor position without moving the cursor. Clearing happens with the
// current background color.
func (t *Terminal) EraseFrontOfScreen() *Terminal {
	t.Print(EraseFrontOfScreen())
	return t
}

// EraseLine erases the whole current line without moving the cursor.
// Clearing happens with the current background color.
func (t *Terminal) EraseLine() *Terminal {
	t.Print(EraseLine())
	return t
}

// EraseScreen clears the while screen without moving the cursor.
// Clearing happens with the current background color.
func (t *Terminal) EraseScreen() *Terminal {
	t.Print(EraseScreen())
	return t
}

// StartAlternativeBuffer clears the screen, moves the cursor to (0,0) and allows
// to return to the original buffer using EndAlternativeBuffer. This allows for
// isolated modifications.
func (t *Terminal) StartAlternativeBuffer() *Terminal {
	t.Print(StartAlternativeBuffer())
	return t
}

// EndAlternativeBuffer returns back to the screen before StartAlternativeBuffer
// as it was left off. If no modifications were made before
// StartAlternativeBuffer to the cursor visibility and colors, sending Reset is
// unnecessary. If the program interrupts in the middle, it seems necessary
// to implicitly call EndAlternativeBuffer, otherwise the console will print
// out the prompt with the alternative buffer settings, and at least in case of cmd.exe
// continues typing with these settings.
func (t *Terminal) EndAlternativeBuffer() *Terminal {
	t.Print(EndAlternativeBuffer())
	return t
}

// ShiftRight moves current line by amount from the current column position.
// Spaces will be added to fill the gap, and anything going beyond the borders of
// viewport will be trimmed.
func (t *Terminal) ShiftRight(amount int) *Terminal {
	t.Print(ShiftRight(amount))
	return t
}

// EraseShiftLeft deletes amount of characters at the current cursor position,
// shifting in space character from the right edge of the viewport.
func (t *Terminal) EraseShiftLeft(amount int) *Terminal {
	t.Print(EraseShiftLeft(amount))
	return t
}

// Erase erases amount characters from the current cursor position without moving
// the cursor by overwriting characters with a space character and not wrapping
// after reaching the right screen border.
func (t *Terminal) Erase(amount int) *Terminal {
	t.Print(Erase(amount))
	return t
}

// ShiftDown shifts the current line down by amount, adding empty line(s) to fill
// the gap. Scrolling margins set with SetScrollRegion will be respected.
func (t *Terminal) ShiftDown(amount int) *Terminal {
	t.Print(ShiftDown(amount))
	return t
}

// DeleteLines deletes amount of lines from the buffer, starting with the row the
// cursor is on. Scrolling margins set with SetScrollRegion will be respected.
func (t *Terminal) DeleteLines(amount int) *Terminal {
	t.Print(DeleteLines(amount))
	return t
}

// Swap swaps foreground and background colors.
// This actually seems to swap the meaning of fg and bg and can be stacked.
// Output CancelSwap to return to normal.
func (t *Terminal) Swap() *Terminal {
	t.Print(Swap())
	return t
}

// CancelSwap returns foreground/background to normal after any proceeding Swap.
func (t *Terminal) CancelSwap() *Terminal {
	t.Print(CancelSwap())
	return t
}

func (t *Terminal) FgRGB(r, g, b int) *Terminal {
	t.Print(FgRGB(r, g, b))
	return t
}

func (t *Terminal) BgRGB(r, g, b int) *Terminal {
	t.Print(BgRGB(r, g, b))
	return t
}
