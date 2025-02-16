package application

import "os"

type ThemeName string

const (
	WhiteTheme ThemeName = "white"
	BlackTheme ThemeName = "black"
)

type ThemeObject struct {
	Text       string `json:"text"`
	TextInvert string `json:"text-invert"`
	TextSmooth string `json:"text-smooth"`

	BodyBackground string `json:"body-background"`

	WindowBackground string `json:"window-background"`
	WindowElement    string `json:"window-element"`

	ErrorErr    string `json:"error-error"`
	ErrorAccent string `json:"error-accent"`
	Link        string `json:"link"`

	BarColorOneMin string `json:"bar-1-min"`
	BarColorOneMax string `json:"bar-1-max"`

	BarColorTwoMin string `json:"bar-2-min"`
	BarColorTwoMax string `json:"bar-2-max"`

	BarColorThreeMin string `json:"bar-3-min"`
	BarColorThreeMax string `json:"bar-3-max"`
}

var ThemeStorage = map[ThemeName]ThemeObject{

	BlackTheme: ThemeObject{
		Text:       "rgb(0, 0, 0)",
		TextInvert: "rgb(216, 216, 216)",
		TextSmooth: "rgb(160, 160, 160)",

		BodyBackground: "rgb(220, 220, 220)",

		WindowBackground: "rgb(220, 220, 220)",
		WindowElement:    "rgb(0, 0, 0)",

		ErrorErr:    "rgb(156, 15, 15)",
		ErrorAccent: "rgb(197, 138, 12)",
		Link:        "rgb(7, 102, 180)",

		BarColorOneMin: "purple",
		BarColorOneMax: "rgb(220, 0, 220)",

		BarColorTwoMin: "green",
		BarColorTwoMax: "rgb(0, 220, 0)",

		BarColorThreeMin: "rgb(0, 125, 125)",
		BarColorThreeMax: "rgb(0, 220, 220)",
	},

	WhiteTheme: ThemeObject{
		Text:       "rgb(223, 223, 223)",
		TextInvert: "rgb(0, 0, 0)",
		TextSmooth: "rgb(160, 160, 160)",

		BodyBackground: "rgb(0, 0, 0)",

		WindowBackground: "rgb(0, 0, 0)",
		WindowElement:    "rgb(220, 220, 220)",

		ErrorErr:    "rgb(156, 15, 15)",
		ErrorAccent: "rgb(197, 138, 12)",
		Link:        "rgb(7, 102, 180)",

		BarColorOneMin: "purple",
		BarColorOneMax: "rgb(220, 0, 220)",

		BarColorTwoMin: "green",
		BarColorTwoMax: "rgb(0, 220, 0)",

		BarColorThreeMin: "rgb(0, 125, 125)",
		BarColorThreeMax: "rgb(0, 220, 220)",
	},
}

func FileBG() []byte {
	buff, err := os.ReadFile("./login-bg.jpg")
	if err != nil {
		return nil
	}
	return buff
}
