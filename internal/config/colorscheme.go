package config

import (
	"strconv"
	"strings"

	"github.com/micro-editor/tcell/v2"
)

// DefStyle is Micro's default style
var DefStyle tcell.Style = tcell.StyleDefault

// Colorscheme is the current colorscheme - hardcoded for micromini
var Colorscheme map[string]tcell.Style

// GetColor takes in a syntax group and returns the colorscheme's style for that group
func GetColor(color string) tcell.Style {
	st := DefStyle
	if color == "" {
		return st
	}
	groups := strings.Split(color, ".")
	if len(groups) > 1 {
		curGroup := ""
		for i, g := range groups {
			if i != 0 {
				curGroup += "."
			}
			curGroup += g
			if style, ok := Colorscheme[curGroup]; ok {
				st = style
			}
		}
	} else if style, ok := Colorscheme[color]; ok {
		st = style
	} else {
		st = StringToStyle(color)
	}

	return st
}

// InitColorscheme initializes the hardcoded default dark colorscheme for micromini
func InitColorscheme() error {
	Colorscheme = make(map[string]tcell.Style)
	DefStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	// Hardcoded dark theme colors - simplified for micromini
	Colorscheme["default"] = DefStyle
	Colorscheme["comment"] = DefStyle.Foreground(tcell.ColorGray)
	Colorscheme["comment.line"] = DefStyle.Foreground(tcell.ColorGray)
	Colorscheme["comment.block"] = DefStyle.Foreground(tcell.ColorGray)
	Colorscheme["constant"] = DefStyle.Foreground(tcell.ColorRed)
	Colorscheme["constant.bool"] = DefStyle.Foreground(tcell.ColorRed)
	Colorscheme["constant.number"] = DefStyle.Foreground(tcell.ColorRed)
	Colorscheme["constant.string"] = DefStyle.Foreground(tcell.ColorYellow)
	Colorscheme["identifier"] = DefStyle.Foreground(tcell.ColorWhite)
	Colorscheme["identifier.function"] = DefStyle.Foreground(tcell.ColorBlue)
	Colorscheme["identifier.class"] = DefStyle.Foreground(tcell.ColorBlue)
	Colorscheme["statement"] = DefStyle.Foreground(tcell.ColorGreen)
	Colorscheme["preproc"] = DefStyle.Foreground(tcell.ColorPurple)
	Colorscheme["type"] = DefStyle.Foreground(tcell.ColorTeal)
	Colorscheme["special"] = DefStyle.Foreground(tcell.ColorPurple)
	Colorscheme["underlined"] = DefStyle.Underline(true)
	Colorscheme["error"] = DefStyle.Foreground(tcell.ColorRed).Background(tcell.ColorWhite)
	Colorscheme["todo"] = DefStyle.Foreground(tcell.ColorYellow).Bold(true)
	Colorscheme["statusline"] = DefStyle.Reverse(true)
	Colorscheme["tabbar"] = DefStyle.Reverse(true)
	Colorscheme["indent-char"] = DefStyle.Foreground(tcell.ColorGray)
	Colorscheme["line-number"] = DefStyle.Foreground(tcell.ColorGray)
	Colorscheme["current-line-number"] = DefStyle.Foreground(tcell.ColorWhite).Bold(true)
	Colorscheme["diff-added"] = DefStyle.Foreground(tcell.ColorGreen)
	Colorscheme["diff-modified"] = DefStyle.Foreground(tcell.ColorYellow)
	Colorscheme["diff-deleted"] = DefStyle.Foreground(tcell.ColorRed)
	Colorscheme["gutter-error"] = DefStyle.Foreground(tcell.ColorRed)
	Colorscheme["gutter-warning"] = DefStyle.Foreground(tcell.ColorYellow)
	Colorscheme["cursor-line"] = DefStyle.Background(tcell.ColorNavy)
	Colorscheme["color-column"] = DefStyle.Background(tcell.ColorNavy)
	Colorscheme["ignore"] = DefStyle.Foreground(tcell.ColorGray)
	Colorscheme["scrollbar"] = DefStyle.Foreground(tcell.ColorWhite).Background(tcell.ColorGray)
	Colorscheme["divider"] = DefStyle.Foreground(tcell.ColorGray)

	return nil
}

// StringToStyle returns a style from a string
// The strings must be in the format "extra foregroundcolor,backgroundcolor"
// The 'extra' can be bold, reverse, italic or underline
func StringToStyle(str string) tcell.Style {
	var fg, bg string
	spaceSplit := strings.Split(str, " ")
	split := strings.Split(spaceSplit[len(spaceSplit)-1], ",")
	if len(split) > 1 {
		fg, bg = split[0], split[1]
	} else {
		fg = split[0]
	}
	fg = strings.TrimSpace(fg)
	bg = strings.TrimSpace(bg)

	var fgColor, bgColor tcell.Color
	var ok bool
	if fg == "" || fg == "default" {
		fgColor, _, _ = DefStyle.Decompose()
	} else {
		fgColor, ok = StringToColor(fg)
		if !ok {
			fgColor, _, _ = DefStyle.Decompose()
		}
	}
	if bg == "" || bg == "default" {
		_, bgColor, _ = DefStyle.Decompose()
	} else {
		bgColor, ok = StringToColor(bg)
		if !ok {
			_, bgColor, _ = DefStyle.Decompose()
		}
	}

	style := DefStyle.Foreground(fgColor).Background(bgColor)
	if strings.Contains(str, "bold") {
		style = style.Bold(true)
	}
	if strings.Contains(str, "italic") {
		style = style.Italic(true)
	}
	if strings.Contains(str, "reverse") {
		style = style.Reverse(true)
	}
	if strings.Contains(str, "underline") {
		style = style.Underline(true)
	}
	return style
}

// StringToColor returns a tcell color from a string representation of a color
// We accept either bright... or light... to mean the brighter version of a color
func StringToColor(str string) (tcell.Color, bool) {
	switch str {
	case "black":
		return tcell.ColorBlack, true
	case "red":
		return tcell.ColorMaroon, true
	case "green":
		return tcell.ColorGreen, true
	case "yellow":
		return tcell.ColorOlive, true
	case "blue":
		return tcell.ColorNavy, true
	case "magenta":
		return tcell.ColorPurple, true
	case "cyan":
		return tcell.ColorTeal, true
	case "white":
		return tcell.ColorSilver, true
	case "brightblack", "lightblack":
		return tcell.ColorGray, true
	case "brightred", "lightred":
		return tcell.ColorRed, true
	case "brightgreen", "lightgreen":
		return tcell.ColorLime, true
	case "brightyellow", "lightyellow":
		return tcell.ColorYellow, true
	case "brightblue", "lightblue":
		return tcell.ColorBlue, true
	case "brightmagenta", "lightmagenta":
		return tcell.ColorFuchsia, true
	case "brightcyan", "lightcyan":
		return tcell.ColorAqua, true
	case "brightwhite", "lightwhite":
		return tcell.ColorWhite, true
	case "default":
		return tcell.ColorDefault, true
	default:
		// Check if this is a 256 color
		if num, err := strconv.Atoi(str); err == nil {
			return GetColor256(num), true
		}
		// Check if this is a truecolor hex value
		if len(str) == 7 && str[0] == '#' {
			return tcell.GetColor(str), true
		}
		return tcell.ColorDefault, false
	}
}

// GetColor256 returns the tcell color for a number between 0 and 255
func GetColor256(color int) tcell.Color {
	if color == 0 {
		return tcell.ColorDefault
	}
	return tcell.PaletteColor(color)
}