package config

import (
	"testing"

	"github.com/micro-editor/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestSimpleStringToStyle(t *testing.T) {
	s := StringToStyle("lightblue,magenta")

	fg, bg, _ := s.Decompose()

	assert.Equal(t, tcell.ColorBlue, fg)
	assert.Equal(t, tcell.ColorPurple, bg)
}

func TestAttributeStringToStyle(t *testing.T) {
	s := StringToStyle("bold cyan,brightcyan")

	fg, bg, attr := s.Decompose()

	assert.Equal(t, tcell.ColorTeal, fg)
	assert.Equal(t, tcell.ColorAqua, bg)
	assert.NotEqual(t, 0, attr&tcell.AttrBold)
}

func TestMultiAttributesStringToStyle(t *testing.T) {
	s := StringToStyle("bold italic underline cyan,brightcyan")

	fg, bg, attr := s.Decompose()

	assert.Equal(t, tcell.ColorTeal, fg)
	assert.Equal(t, tcell.ColorAqua, bg)
	assert.NotEqual(t, 0, attr&tcell.AttrBold)
	assert.NotEqual(t, 0, attr&tcell.AttrItalic)
	assert.NotEqual(t, 0, attr&tcell.AttrUnderline)
}

func TestColor256StringToStyle(t *testing.T) {
	s := StringToStyle("128,60")

	fg, bg, _ := s.Decompose()

	assert.Equal(t, tcell.Color128, fg)
	assert.Equal(t, tcell.Color60, bg)
}

func TestColorHexStringToStyle(t *testing.T) {
	s := StringToStyle("#deadbe,#ef1234")

	fg, bg, _ := s.Decompose()

	assert.Equal(t, tcell.NewRGBColor(222, 173, 190), fg)
	assert.Equal(t, tcell.NewRGBColor(239, 18, 52), bg)
}

func TestHardcodedColorscheme(t *testing.T) {
	// Test that our hardcoded colorscheme initializes correctly
	err := InitColorscheme()
	assert.Nil(t, err)

	// Test that essential colorscheme entries exist
	assert.NotNil(t, Colorscheme["default"])
	assert.NotNil(t, Colorscheme["comment"])
	assert.NotNil(t, Colorscheme["constant"])
	assert.NotNil(t, Colorscheme["statusline"])

	// Test that the comment color is gray as expected
	fg, _, _ := Colorscheme["comment"].Decompose()
	assert.Equal(t, tcell.ColorGray, fg)
}
