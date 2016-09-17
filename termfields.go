//Package termfields creates updateable form fields at specified locations in the console.
package termfields

import (
	"fmt"

	tb "github.com/nsf/termbox-go"
)

var boxRunesMap map[boxStyle][]rune

type (
	boxStyle uint16
	moveDir  uint16
)

// Flags to style a box border around a field
const (
	boxStyleClear boxStyle = iota
	BoxStyleNone
	BoxStyleASCII
	BoxStyleUnicode
)

// Flags to move a field in a specified direction
const (
	FieldMoveLeft moveDir = iota
	FieldMoveRight
	FieldMoveUp
	FieldMoveDown
)

// Field is the identifier for a specific form field on the screen.
type Field struct {
	field
}

type field struct {
	x, y   int
	len    int
	border boxStyle
	text   string
}

func init() {
	boxRunesMap = map[boxStyle][]rune{
		boxStyleClear:   {' ', ' ', ' ', ' ', ' ', ' '},
		BoxStyleNone:    {0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		BoxStyleASCII:   {'+', '+', '+', '+', '-', '|'},
		BoxStyleUnicode: {0x250c, 0x2510, 0x2514, 0x2518, 0x2500, 0x2502},
	}
}

// Init Initializes termfields library. This function should be called before any other functions.
// After successful initialization, the library must be finalized using 'Close' function.
//
// Example usage:
//      err := termfields.Init()
//      if err != nil {
//              panic(err)
//      }
//      defer termfields.Close()
func Init() error {
	return tb.Init()
}

// Close Finalizes termbox library, should be called after successful initialization
// when termbox's functionality isn't required anymore.
func Close() {
	tb.SetCursor(0, 0)
	tb.Close()
}

func (f *field) Move(dir moveDir) {
	border := f.border
	f.DrawBox(boxStyleClear)
	switch {
	case dir == FieldMoveLeft:
		f.x--
	case dir == FieldMoveRight:
		f.x++
	case dir == FieldMoveUp:
		f.y--
	case dir == FieldMoveDown:
		f.y++
	}
	f.DrawBox(border)
	f.Update(f.text)
}

// NewField creates a new field at location y,x of lenth len with contents text.
func NewField(y, x, len int, text string) (*Field, error) {
	f := field{
		x:   x,
		y:   y,
		len: len,
	}
	err := f.Update(text)
	if err != nil {
		return nil, err
	}
	return &Field{f}, nil
}

func (f *field) DrawBox(boxType boxStyle) error {
	if !tb.IsInit {
		return fmt.Errorf("Term not Initialized")
	}
	if _, ok := boxRunesMap[boxType]; !ok {
		return fmt.Errorf("Unknown Box Style")
	}

	//Draw Corners
	tb.SetCell(f.x-1, f.y-1, boxRunesMap[boxType][0], tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(f.x+f.len+1, f.y-1, boxRunesMap[boxType][1], tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(f.x-1, f.y+1, boxRunesMap[boxType][2], tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(f.x+f.len+1, f.y+1, boxRunesMap[boxType][3], tb.ColorDefault, tb.ColorDefault)
	//Draw Sides
	tb.SetCell(f.x-1, f.y, boxRunesMap[boxType][5], tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(f.x+f.len+1, f.y, boxRunesMap[boxType][5], tb.ColorDefault, tb.ColorDefault)
	//Draw Top
	for i := 0; i < f.len+1; i++ {
		tb.SetCell(f.x+i, f.y-1, boxRunesMap[boxType][4], tb.ColorDefault, tb.ColorDefault)
		tb.SetCell(f.x+i, f.y+1, boxRunesMap[boxType][4], tb.ColorDefault, tb.ColorDefault)
	}
	tb.Flush()
	f.border = boxType
	return nil
}

func (f *field) Update(s string) error {
	if !tb.IsInit {
		return fmt.Errorf("Term not Initialized")
	}
	for i, c := range s {
		tb.SetCell(f.x+i, f.y, c, tb.ColorDefault, tb.ColorDefault)
	}
	tb.Flush()
	f.text = s
	return nil
}
