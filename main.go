package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"strings"

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
	state      string
	sprite     string
	lastSprite string
	lastSpriteCount int
	img        *ebiten.Image
}

type Config struct {
	Speed int     `cfg:"speed" cfgDefault:"4" cfgHelper:"The speed of the deer."`
	Scale float64 `cfg:"scale" cfgDefault:"3.0" cfgHelper:"The scale of the deer."`
	// Quiet            bool    `cfg:"quiet" cfgDefault:"false" cfgHelper:"Disable sound."`
	MousePassthrough bool `cfg:"mousepassthrough" cfgDefault:"false" cfgHelper:"Enable mouse passthrough."`
}

const (
	width  = 32
	height = 32
)

var (
	//go:embed deer.png
	Deer_png       []byte
	spriteSheet, _ = LoadSpriteSheet()
	monitorWidth, monitorHeight = ebiten.Monitor().Size()
	cfg = &Config{}
)

// SpriteSheet represents a collection of sprite images.
type SpriteSheet map[string]*ebiten.Image

// LoadSpriteSheet loads the embedded SpriteSheet.
func LoadSpriteSheet() (SpriteSheet, error) {
	var tileSize = 32
	var halfTileSize = 16

	img, _, err := image.Decode(bytes.NewReader(Deer_png))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y, xSize,ySize int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*xSize, (y+1)*ySize, (x+1)*xSize, y*ySize)).(*ebiten.Image)
	}

	// Populate SpriteSheet.

	s := SpriteSheet{
		"StandingFrontTailUp":       spriteAt(0, 0, tileSize, tileSize),
		"StandingFrontTailDown":      spriteAt(1, 0, tileSize, tileSize),
		"StandingFrontLegUpTailDown": spriteAt(2, 0, tileSize, tileSize),
		"StandingFrontLegUpTailUp":   spriteAt(3, 0, tileSize, tileSize),

		"StandingBackTailUp":       spriteAt(0, 1, tileSize, tileSize),
		"StandingBackTailDown":      spriteAt(1, 1, tileSize, tileSize),
		"StandingBackLegUpTailUp":   spriteAt(2, 1, tileSize, tileSize),
		"StandingBackLegUpTailDown": spriteAt(3, 1, tileSize, tileSize),

		"StandingRightTailUp":       spriteAt(0, 3, tileSize, tileSize),
		"StandingRightTailDown":      spriteAt(1, 3, tileSize, tileSize),
		"StandingRightLegUpTailUp":   spriteAt(2, 3, tileSize, tileSize),
		"StandingRightLegUpTailDown": spriteAt(3, 3, tileSize, tileSize),

		"StandingLeftTailUp":       spriteAt(0, 2, tileSize, tileSize),
		"StandingLeftTailDown":      spriteAt(1, 2, tileSize, tileSize),
		"StandingLeftLegUpTailUp":   spriteAt(2, 2, tileSize, tileSize),
		"StandingLeftLegUpTailDown": spriteAt(3, 2, tileSize, tileSize),

		"WalkingForward1": spriteAt(0, 4, tileSize, tileSize),
		"WalkingForward2": spriteAt(1, 4, tileSize, tileSize),
		"WalkingForward3": spriteAt(2, 4, tileSize, tileSize),
		"WalkingForward4": spriteAt(3, 4, tileSize, tileSize),

		"WalkingAway1": spriteAt(0, 5, tileSize, tileSize),
		"WalkingAway2": spriteAt(1, 5, tileSize, tileSize),
		"WalkingAway3": spriteAt(2, 5, tileSize, tileSize),
		"WalkingAway4": spriteAt(3, 5, tileSize, tileSize),

		"WalkingRight1": spriteAt(0, 6, tileSize, tileSize),
		"WalkingRight2": spriteAt(1, 6, tileSize, tileSize),
		"WalkingRight3": spriteAt(2, 6, tileSize, tileSize),
		"WalkingRight4": spriteAt(3, 6, tileSize, tileSize),

		"WalkingLeft1": spriteAt(0, 7, tileSize, tileSize),
		"WalkingLeft2": spriteAt(1, 7, tileSize, tileSize),
		"WalkingLeft3": spriteAt(2, 7, tileSize, tileSize),
		"WalkingLeft4": spriteAt(3, 7, tileSize, tileSize),

		"JumpingLeft1": spriteAt(4, 7, tileSize, tileSize),
		"JumpingLeft2": spriteAt(5, 7, tileSize, tileSize),
		"JumpingLeft3": spriteAt(6, 7, tileSize, tileSize),
		"JumpingLeft4": spriteAt(7, 7, tileSize, tileSize),

		"JumpingRight1": spriteAt(4, 6, tileSize, tileSize),
		"JumpingRight2": spriteAt(5, 6, tileSize, tileSize),
		"JumpingRight3": spriteAt(6, 6, tileSize, tileSize),
		"JumpingRight4": spriteAt(7, 6, tileSize, tileSize),

		"JumpingForward1": spriteAt(4, 4, tileSize, tileSize),
		"JumpingForward2": spriteAt(5, 4, tileSize, tileSize),
		"JumpingForward3": spriteAt(6, 4, tileSize, tileSize),
		"JumpingForward4": spriteAt(7, 4, tileSize, tileSize),

		"JumpingAway1": spriteAt(4, 5, tileSize, tileSize),
		"JumpingAway2": spriteAt(5, 5, tileSize, tileSize),
		"JumpingAway3": spriteAt(6, 5, tileSize, tileSize),
		"JumpingAway4": spriteAt(7, 5, tileSize, tileSize),

		"HeadLookingForward1": spriteAt(4,0,tileSize, halfTileSize),
		"HeadLookingForward2": spriteAt(4,1,tileSize, halfTileSize),
		"HeadLookingForward3": spriteAt(4,2,tileSize, halfTileSize),
		"HeadLookingBack": spriteAt(5,0,tileSize, halfTileSize),
		"HeadDownForward": spriteAt(8,0,tileSize, tileSize),
		"HeadLookingRight": spriteAt(6,0, tileSize, halfTileSize),
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
		m.stayIdle()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			m.waiting = !m.waiting
		}
		return nil
	}

	m.catchCursor(x, y)
	return nil
}

func (m *deer) stayIdle() {
	if (m.state == "idle") {
		if m.count - m.lastSpriteCount > rand.Intn(200) + 50 {
			var random = rand.Intn(15)
			switch {
			case 0 < random && random < 2:
				m.sprite = "StandingFrontTailUp+HeadLookingForward1"
			case 2 < random && random < 5:
				m.sprite = "StandingFrontTailDown+HeadLookingForward2"
			case random > 5 && random < 7:
				m.sprite = "StandingFrontTailUp+HeadLookingForward2"
			case random > 7 && random < 9:
				m.sprite = "StandingFrontLegUpTailUp+HeadLookingRight"
			case random == 9:
				m.sprite = "StandingFrontTailDown+HeadDownForward"
			case random >= 10:
				m.sprite = "StandingFrontTailUp+HeadLookingBack"
			}
		}

	} else {
		m.sprite = "StandingFrontTailDown+HeadLookingForward1"
		m.state = "idle"
	}
}

func (m *deer) catchCursor(x, y int) {
	m.state = "catchCursor"


	// get mouse direction
	rad := math.Atan2(float64(y), float64(x))
	turn := 0.0
	if rad <= 0 {
		turn = 360
	}

	angle := (rad / math.Pi * 180) + turn

	switch {
	case angle <= 292.5 && angle > 247.5: // up
		m.y -= cfg.Speed
	case angle <= 337.5 && angle > 292.5: // up right
		m.x += cfg.Speed
		m.y -= cfg.Speed
	case angle <= 22.5 || angle > 337.5: // right
		m.x += cfg.Speed
	case angle <= 67.5 && angle > 22.5: // down right
		m.x += cfg.Speed
		m.y += cfg.Speed
	case angle <= 112.5 && angle > 67.5: // down
		m.y += cfg.Speed
	case angle <= 157.5 && angle > 112.5: // down left
		m.x -= cfg.Speed
		m.y += cfg.Speed
	case angle <= 202.5 && angle > 157.5: // left
		m.x -= cfg.Speed
	case angle <= 247.5 && angle > 202.5: // up left
		m.x -= cfg.Speed
		m.y -= cfg.Speed
	}


	if (m.lastSpriteMinDiff()) {
		switch {
		case angle < 292 && angle > 247:
			switch m.lastSprite {
				case "WalkingAway1":
					m.sprite = "WalkingAway2"
				case "WalkingAway2":
					m.sprite = "WalkingAway3"
				case "WalkingAway3":
					m.sprite = "WalkingAway4"
				default:
					m.sprite = "WalkingAway1"
				}
		case angle < 337 && angle > 292:
			switch m.lastSprite {
				case "JumpingRight1":
					m.sprite = "JumpingRight2"
				case "JumpingRight2":
					m.sprite = "JumpingRight3"
				case "JumpingRight3":
					m.sprite = "JumpingRight4"
				default:
					m.sprite = "JumpingRight1"
				}
		case angle < 22 || angle > 337:
			switch m.lastSprite {
				case "WalkingRight1":
					m.sprite = "WalkingRight2"
				case "WalkingRight2":
					m.sprite = "WalkingRight3"
				case "WalkingRight3":
					m.sprite = "WalkingRight4"
				default:
					m.sprite = "WalkingRight1"
				}
		case angle < 67 && angle > 22:
			switch m.lastSprite {
				case "JumpingRight1":
					m.sprite = "JumpingRight2"
				case "JumpingRight2":
					m.sprite = "JumpingRight3"
				case "JumpingRight3":
					m.sprite = "JumpingRight4"
				default:
					m.sprite = "JumpingRight1"
				}
			
		case angle < 112 && angle > 67:
			switch m.lastSprite {
				case "WalkingForward1":
					m.sprite = "WalkingForward2"
				case "WalkingForward2":
					m.sprite = "WalkingForward3"
				case "WalkingForward3":
					m.sprite = "WalkingForward4"
				default:
					m.sprite = "WalkingForward1"
				}
		case angle < 157 && angle > 112:
			switch m.lastSprite {
				case "JumpingLeft1":
					m.sprite = "JumpingLeft2"
				case "JumpingLeft2":
					m.sprite = "JumpingLeft3"
				case "JumpingLeft3":
					m.sprite = "JumpingLeft4"
				default:
					m.sprite = "JumpingLeft1"
				}
		case angle < 202 && angle > 157:
			switch m.lastSprite {
				case "WalkingLeft1":
					m.sprite = "WalkingLeft2"
				case "WalkingLeft2":
					m.sprite = "WalkingLeft3"
				case "WalkingLeft3":
					m.sprite = "WalkingLeft4"
				default:
					m.sprite = "WalkingLeft1"
				}
		case angle < 247 && angle > 202:
			switch m.lastSprite {
				case "JumpingLeft1":
					m.sprite = "JumpingLeft2"
				case "JumpingLeft2":
					m.sprite = "JumpingLeft3"
				case "JumpingLeft3":
					m.sprite = "JumpingLeft4"
				default:
					m.sprite = "JumpingLeft1"
				}
		}
	}

}

func (m *deer) lastSpriteMinDiff() bool {
	var lastSpriteDiff = m.count - m.lastSpriteCount

	var lastSpriteMinDiff = lastSpriteDiff > 5
	return lastSpriteMinDiff
}

func (m *deer) Draw(screen *ebiten.Image) {
	screen.Clear()
	split := strings.Split(m.sprite, "+")
	if (len(split) > 1) {
		layer1 := split[0]
		layer2 := split[1]
		screen.DrawImage(spriteSheet[layer1], nil)
		screen.DrawImage(spriteSheet[layer2], nil)
	} else {
		m.img = spriteSheet[m.sprite]
		screen.DrawImage(m.img, nil)
	}


	if m.lastSprite == m.sprite {
		return
	}

	m.lastSprite = m.sprite
	m.lastSpriteCount = m.count
}

func main() {
	config.PrefixEnv = "DEER"
	config.File = "deer.ini"
	config.Parse(cfg)

	d := &deer{
		x:      monitorWidth / 2,
		y:      monitorHeight / 2,
		sprite: "WalkingForward1",
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
