package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	arcadeFont          font.Face
	space, rocket, rock *ebiten.Image
)

type mode int

const (
	GameMode mode = iota
	GameOverMode
	NewGameMode
)

type Game struct {
	background    *ebiten.Image
	height, width int

	mode

	player struct {
		image        *ebiten.Image
		xPos, yPos   float64
		speed        float64
		angleRadians float64
	}

	obstacle struct {
		image      *ebiten.Image
		xPos, yPos float64
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.mode == GameOverMode {
		text.Draw(screen, "Game Over!", arcadeFont, g.width/2, g.height/2, color.RGBA{R: 255})
		return
	}
	// draw background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(g.background, op)

	// draw player
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(g.player.angleRadians)
	op.GeoM.Translate(g.player.xPos, g.player.yPos)
	screen.DrawImage(g.player.image, op)

	// draw obstacle
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.obstacle.xPos, g.obstacle.yPos)
	screen.DrawImage(g.obstacle.image, op)
}

func (g *Game) movePlayer() {
	if g.player.xPos+float64(g.player.image.Bounds().Dx()) > float64(g.background.Bounds().Dx()) {
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.yPos -= g.player.speed
		g.player.angleRadians = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.yPos += g.player.speed
		g.player.angleRadians = math.Pi
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.xPos -= g.player.speed
		g.player.angleRadians = 270 * math.Pi / 180
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.xPos += g.player.speed
		g.player.angleRadians = 90 * math.Pi / 180
	}
}

func (g *Game) controlSpeed() {
	speed := g.player.speed
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		speed += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		speed -= 1
	}

	if speed > 20 {
		speed = 20
	}
	if speed < 1 {
		speed = 1
	}
	if speed != g.player.speed {
		g.player.speed = speed
		fmt.Printf("Speed: %v\n", g.player.speed)
	}
}

func (g *Game) hit() bool {
	if g.player.xPos > g.obstacle.xPos &&
		g.player.xPos < g.obstacle.xPos+float64(g.obstacle.image.Bounds().Dx()) &&
		g.player.yPos > g.obstacle.yPos &&
		g.player.yPos < g.obstacle.yPos+float64(g.obstacle.image.Bounds().Dy()) {
		return true
	}
	return false
}

func (g *Game) gameOverMenuControls() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.mode = NewGameMode
	}
}

// Update ...
func (g *Game) Update(screen *ebiten.Image) error {
	switch g.mode {
	case GameMode:
		g.controlSpeed()
		g.movePlayer()
		if g.hit() {
			g.mode = GameOverMode
		}
	case GameOverMode:
		g.gameOverMenuControls()
	case NewGameMode:
		g = NewGame()
		g.mode = GameMode
	}
	return nil
}

func NewGame() *Game {
	myGame := Game{
		mode:       GameMode,
		background: space,
		width:      640,
		height:     480,
	}

	h, w := myGame.Layout(0, 0)
	myGame.player.xPos, myGame.player.yPos = float64(h/2), float64(w/2)
	myGame.player.speed = 12
	myGame.player.image = rocket

	myGame.obstacle.xPos, myGame.obstacle.yPos = float64(h/4), float64(w/4)
	myGame.obstacle.image = rock
	return &myGame
}

func main() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	space, _, err = ebitenutil.NewImageFromFile("assets/background.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	rocket, _, err = ebitenutil.NewImageFromFile("assets/rocket.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	rock, _, err = ebitenutil.NewImageFromFile("assets/rock.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
