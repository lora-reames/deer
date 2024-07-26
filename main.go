package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	// "github.com/hajimehoshi/ebiten/v2/audio"
	// "github.com/hajimehoshi/ebiten/v2/audio/wav"
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type deer struct {
	waiting    bool
	x          int
	y          int
	distance   int
	count      int
	min        int
	max        int
	state      int
	sprite     string
	lastSprite string
	img        *ebiten.Image
}

type Config struct {
	Speed int     `cfg:"speed" cfgDefault:"2" cfgHelper:"The speed of the deer."`
	Scale float64 `cfg:"scale" cfgDefault:"2.0" cfgHelper:"The scale of the deer."`
	// Quiet            bool    `cfg:"quiet" cfgDefault:"false" cfgHelper:"Disable sound."`
	MousePassthrough bool `cfg:"mousepassthrough" cfgDefault:"false" cfgHelper:"Enable mouse passthrough."`
}

const (
	width  = 32
	height = 32
)

var (
	// loaded  = false
	// mSprite map[string]*ebiten.Image
	//go:embed deer.png
	Deer_png       []byte
	spriteSheet, _ = LoadSpriteSheet()

	monitorWidth, monitorHeight = ebiten.Monitor().Size()

	cfg = &Config{}
)

// SpriteSheet represents a collection of sprite images.
type SpriteSheet map[string]*ebiten.Image

// {
// 	StandingFrontTaillUp       *ebiten.Image
// 	StandingFrontTailDown      *ebiten.Image
// 	StandingFrontLegUpTailUp   *ebiten.Image
// 	StandingFrontLegUpTailDown *ebiten.Image
// 	StandingBackTaillUp        *ebiten.Image
// 	StandingBackTailDown       *ebiten.Image
// 	StandingBackLegUpTailUp    *ebiten.Image
// 	StandingBackLegUpTailDown  *ebiten.Image
// 	StandingRightTaillUp       *ebiten.Image
// 	StandingRightTailDown      *ebiten.Image
// 	StandingRightLegUpTailUp   *ebiten.Image
// 	StandingRightLegUpTailDown *ebiten.Image
// 	StandingLeftTaillUp        *ebiten.Image
// 	StandingLeftTailDown       *ebiten.Image
// 	StandingLeftLegUpTailUp    *ebiten.Image
// 	StandingLeftLegUpTailDown  *ebiten.Image
// 	WalkingForward1            *ebiten.Image
// 	WalkingForward2            *ebiten.Image
// 	WalkingForward3            *ebiten.Image
// 	WalkingForward4            *ebiten.Image
// 	WalkingAway1               *ebiten.Image
// 	WalkingAway2               *ebiten.Image
// 	WalkingAway3               *ebiten.Image
// 	WalkingAway4               *ebiten.Image
// 	WalkingRight1              *ebiten.Image
// 	WalkingRight2              *ebiten.Image
// 	WalkingRight3              *ebiten.Image
// 	WalkingRight4              *ebiten.Image
// 	WalkingLeft1               *ebiten.Image
// 	WalkingLeft2               *ebiten.Image
// 	WalkingLeft3               *ebiten.Image
// 	WalkingLeft4               *ebiten.Image
// 	JumpingLeft1               *ebiten.Image
// 	JumpingLeft2               *ebiten.Image
// 	JumpingLeft3               *ebiten.Image
// 	JumpingLeft4               *ebiten.Image
// 	JumpingRight1              *ebiten.Image
// 	JumpingRight2              *ebiten.Image
// 	JumpingRight3              *ebiten.Image
// 	JumpingRight4              *ebiten.Image
// 	JumpingForward1            *ebiten.Image
// 	JumpingForward2            *ebiten.Image
// 	JumpingForward3            *ebiten.Image
// 	JumpingForward4            *ebiten.Image
// 	JumpingAway1               *ebiten.Image
// 	JumpingAway2               *ebiten.Image
// 	JumpingAway3               *ebiten.Image
// 	JumpingAway4               *ebiten.Image
// }

// LoadSpriteSheet loads the embedded SpriteSheet.
func LoadSpriteSheet() (SpriteSheet, error) {
	var tileSize = 32

	img, _, err := image.Decode(bytes.NewReader(Deer_png))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*tileSize, (y+1)*tileSize, (x+1)*tileSize, y*tileSize)).(*ebiten.Image)
	}

	// Populate SpriteSheet.

	s := SpriteSheet{
		"StandingFrontTailUp": spriteAt(0, 0),

		"StandingFrontTaillUp":       spriteAt(0, 0),
		"StandingFrontTailDown":      spriteAt(1, 0),
		"StandingFrontLegUpTailDown": spriteAt(2, 0),
		"StandingFrontLegUpTailUp":   spriteAt(3, 0),

		"StandingBackTaillUp":       spriteAt(0, 1),
		"StandingBackTailDown":      spriteAt(1, 1),
		"StandingBackLegUpTailUp":   spriteAt(2, 1),
		"StandingBackLegUpTailDown": spriteAt(3, 1),

		"StandingRightTaillUp":       spriteAt(0, 3),
		"StandingRightTailDown":      spriteAt(1, 3),
		"StandingRightLegUpTailUp":   spriteAt(2, 3),
		"StandingRightLegUpTailDown": spriteAt(3, 3),

		"StandingLeftTaillUp":       spriteAt(0, 2),
		"StandingLeftTailDown":      spriteAt(1, 2),
		"StandingLeftLegUpTailUp":   spriteAt(2, 2),
		"StandingLeftLegUpTailDown": spriteAt(3, 2),

		"WalkingForward1": spriteAt(0, 4),
		"WalkingForward2": spriteAt(1, 4),
		"WalkingForward3": spriteAt(2, 4),
		"WalkingForward4": spriteAt(3, 4),

		"WalkingAway1": spriteAt(0, 5),
		"WalkingAway2": spriteAt(1, 5),
		"WalkingAway3": spriteAt(2, 5),
		"WalkingAway4": spriteAt(3, 5),

		"WalkingRight1": spriteAt(0, 6),
		"WalkingRight2": spriteAt(1, 6),
		"WalkingRight3": spriteAt(2, 6),
		"WalkingRight4": spriteAt(3, 6),

		"WalkingLeft1": spriteAt(0, 7),
		"WalkingLeft2": spriteAt(1, 7),
		"WalkingLeft3": spriteAt(2, 7),
		"WalkingLeft4": spriteAt(3, 7),

		"JumpingLeft1": spriteAt(4, 7),
		"JumpingLeft2": spriteAt(5, 7),
		"JumpingLeft3": spriteAt(6, 7),
		"JumpingLeft4": spriteAt(7, 7),

		"JumpingRight1": spriteAt(4, 6),
		"JumpingRight2": spriteAt(5, 6),
		"JumpingRight3": spriteAt(6, 6),
		"JumpingRight4": spriteAt(7, 6),

		"JumpingForward1": spriteAt(4, 4),
		"JumpingForward2": spriteAt(5, 4),
		"JumpingForward3": spriteAt(6, 4),
		"JumpingForward4": spriteAt(7, 4),

		"JumpingAway1": spriteAt(4, 5),
		"JumpingAway2": spriteAt(5, 5),
		"JumpingAway3": spriteAt(6, 5),
		"JumpingAway4": spriteAt(7, 5),
	}

	return s, nil
}

func (m *deer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func (m *deer) Update() error {
	m.count++
	// Prevents deer from being stuck on the side of the screen
	// or randomly traveling to another monitor
	m.x = max(0, min(m.x, monitorWidth))
	m.y = max(0, min(m.y, monitorHeight))
	ebiten.SetWindowPosition(m.x, m.y)

	mx, my := ebiten.CursorPosition()
	x := mx - (height / 2)
	y := my - (width / 2)

	dy, dx := y, x
	if dy < 0 {
		dy = -dy
	}
	if dx < 0 {
		dx = -dx
	}

	m.distance = dx + dy
	if m.distance < width || m.waiting {
		// m.stayIdle()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			m.waiting = !m.waiting
		}
		return nil
	}

	m.catchCursor(x, y)
	return nil
}

// func (m *deer) stayIdle() {
// 	// idle state
// 	switch m.state {
// 	case 0:
// 		m.state = 1
// 		fallthrough

// 	case 1, 2, 3:
// 		m.sprite = "awake"

// 	case 4, 5, 6:
// 		m.sprite = "scratch"

// 	case 7, 8, 9:
// 		m.sprite = "wash"

// 	case 10, 11, 12:
// 		m.min = 32
// 		m.max = 64
// 		m.sprite = "yawn"

// 	default:
// 		m.sprite = "sleep"
// 	}
// }

func (m *deer) catchCursor(x, y int) {
	m.state = 0
	m.min = 8
	m.max = 16
	tr := 0.0
	// get mouse direction
	r := math.Atan2(float64(y), float64(x))
	if r <= 0 {
		tr = 360
	}

	a := (r / math.Pi * 180) + tr

	switch {
	case a <= 292.5 && a > 247.5: // up
		m.y -= cfg.Speed
	case a <= 337.5 && a > 292.5: // up right
		m.x += cfg.Speed
		m.y -= cfg.Speed
	case a <= 22.5 || a > 337.5: // right
		m.x += cfg.Speed
	case a <= 67.5 && a > 22.5: // down right
		m.x += cfg.Speed
		m.y += cfg.Speed
	case a <= 112.5 && a > 67.5: // down
		m.y += cfg.Speed
	case a <= 157.5 && a > 112.5: // down left
		m.x -= cfg.Speed
		m.y += cfg.Speed
	case a <= 202.5 && a > 157.5: // left
		m.x -= cfg.Speed
	case a <= 247.5 && a > 202.5: // up left
		m.x -= cfg.Speed
		m.y -= cfg.Speed
	}

	switch {
	case a < 292 && a > 247:
		m.sprite = "WalkingAway1"
		// m.img = spriteSheet.WalkingAway1
	case a < 337 && a > 292:
		m.sprite = "JumpingRight1"
		// m.img = spriteSheet.JumpingRight1
	case a < 22 || a > 337:
		m.sprite = "WalkingRight1"
		// m.img = spriteSheet.WalkingRight1
	case a < 67 && a > 22:
		m.sprite = "JumpingRight3"
		// m.img = spriteSheet.JumpingRight3
	case a < 112 && a > 67:
		m.sprite = "WalkingForward1"
		// m.img = spriteSheet.WalkingForward1
	case a < 157 && a > 112:
		m.sprite = "JumpingLeft3"
		// m.img = spriteSheet.JumpingLeft3
	case a < 202 && a > 157:
		m.sprite = "WalkingLeft1"
		// m.img = spriteSheet.WalkingLeft1
	case a < 247 && a > 202:
		m.sprite = "JumpingLeft1"
		// m.img = spriteSheet.JumpingLeft1
	}
}

func (m *deer) Draw(screen *ebiten.Image) {
	// var sprite string
	// switch {
	// case m.sprite == "awake":
	// 	sprite = m.sprite
	// case m.count < m.min:
	// 	sprite = m.sprite + "1"
	// default:
	// 	sprite = m.sprite + "2"
	// }

	m.img = spriteSheet[m.sprite]
	// m.img = spriteSheet.WalkingForward1

	// if m.count > m.max {
	// 	m.count = 0

	// 	if m.state > 0 {
	// 		m.state++
	// 		// switch m.state {
	// 		// case 13:
	// 		// 	playSound(mSound["sleep"])
	// 		// }
	// 	}
	// }

	if m.lastSprite == m.sprite {
		return
	}

	m.lastSprite = m.sprite

	screen.Clear()

	screen.DrawImage(m.img, nil)
}

func main() {
	config.PrefixEnv = "DEER"
	config.File = "deer.ini"
	config.Parse(cfg)

	d := &deer{
		x:      monitorWidth / 2,
		y:      monitorHeight / 2,
		sprite: "WalkingForward1",
		min:    8,
		max:    16,
	}

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetTPS(50)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowMousePassthrough(cfg.MousePassthrough)
	ebiten.SetWindowSize(int(float64(width)*cfg.Scale), int(float64(height)*cfg.Scale))
	ebiten.SetWindowTitle("Deer")

	err := ebiten.RunGameWithOptions(d, &ebiten.RunGameOptions{
		InitUnfocused:     true,
		ScreenTransparent: true,
		SkipTaskbar:       true,
		X11ClassName:      "Deer",
		X11InstanceName:   "Deer",
	})
	if err != nil {
		log.Fatal(err)
	}
}
