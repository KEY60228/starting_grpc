package game

type Color int

const (
	Empty Color = iota // 誰も打っていない
	Black              // 黒
	White              // 白
	Wall               // 壁
	None               // なんでもない
)

func ColorToStr(c Color) string {
	switch c {
	case Black:
		return "○"
	case White:
		return "◉"
	case Empty:
		return " "
	}
	return ""
}

func OpponentColor(me Color) Color {
	switch me {
	case Black:
		return White
	case White:
		return Black
	}

	panic("invalid state")
}
