package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 800
	screenHeight = 600
	gameTime     = 60 // Game time in seconds
)

// Bee represents a bee in the game
type Bee struct {
	x, y        float64
	speedX      float64
	speedY      float64
	width       int
	height      int
	isHornets   bool
	isHighSpeed bool
	visible     bool
}

// Game represents the game state
type Game struct {
	bees            []Bee
	score           int
	hornetsClicked  int
	gameOver        bool
	gameStarted     bool
	startTime       time.Time
	remainingTime   int
	lightningEffect bool
	lightningTimer  int
	face            font.Face
}

// Initialize a new game
func NewGame() *Game {
	g := &Game{
		bees:           make([]Bee, 0),
		score:          0,
		hornetsClicked: 0,
		gameOver:       false,
		gameStarted:    false,
		remainingTime:  gameTime,
		face:           basicfont.Face7x13,
	}

	// Load bee images
	loadBeeImages()

	return g
}

var (
	beeImage     *ebiten.Image
	hornetsImage *ebiten.Image
	forestImage  *ebiten.Image
)

// Load bee images from local files
func loadBeeImages() {
	// Load bee image
	beeFile, err := ebitenutil.OpenFile("image/bee.png")
	if err != nil {
		log.Fatalf("Failed to open bee image: %v", err)
	}
	defer beeFile.Close()

	img, _, err := image.Decode(beeFile)
	if err != nil {
		log.Fatalf("Failed to decode bee image: %v", err)
	}
	beeImage = ebiten.NewImageFromImage(img)

	// Load hornets image
	hornetsFile, err := ebitenutil.OpenFile("image/hornet.png")
	if err != nil {
		log.Fatalf("Failed to open hornets image: %v", err)
	}
	defer hornetsFile.Close()

	img, _, err = image.Decode(hornetsFile)
	if err != nil {
		log.Fatalf("Failed to decode hornets image: %v", err)
	}
	hornetsImage = ebiten.NewImageFromImage(img)

	// Load forest background
	forestFile, err := ebitenutil.OpenFile("image/forest.jpg")
	if err != nil {
		log.Fatalf("Failed to open forest image: %v", err)
	}
	defer forestFile.Close()

	img, _, err = image.Decode(forestFile)
	if err != nil {
		log.Fatalf("Failed to decode forest image: %v", err)
	}
	forestImage = ebiten.NewImageFromImage(img)
}

// Draw forest background
func drawForestBackground(screen *ebiten.Image) {
	// Draw the forest background with proper scaling to fit the screen
	op := &ebiten.DrawImageOptions{}

	// Scale the background to fit the screen
	bgWidth, bgHeight := forestImage.Size()
	scaleX := float64(screenWidth) / float64(bgWidth)
	scaleY := float64(screenHeight) / float64(bgHeight)

	op.GeoM.Scale(scaleX, scaleY)
	screen.DrawImage(forestImage, op)

	// Add a slight green tint overlay to enhance forest feel
	overlayImage := ebiten.NewImage(screenWidth, screenHeight)
	overlayImage.Fill(color.RGBA{0, 100, 0, 40}) // Semi-transparent green
	screen.DrawImage(overlayImage, nil)
}

// Add a new bee to the game
func (g *Game) addBee() {
	isHornets := rand.Float64() < 0.2   // 20% chance to be a hornet
	isHighSpeed := rand.Float64() < 0.1 // 10% chance to be high speed

	var img *ebiten.Image
	if isHornets {
		img = hornetsImage
	} else {
		img = beeImage
	}

	width, height := img.Size()
	width = width / 2 // Make the hitbox smaller than the actual image
	height = height / 2

	speedBase := 2.0
	if isHighSpeed {
		speedBase = 5.0
	}

	bee := Bee{
		x:           float64(rand.Intn(screenWidth - width)),
		y:           float64(rand.Intn(screenHeight - height)),
		speedX:      (rand.Float64()*2 - 1) * speedBase,
		speedY:      (rand.Float64()*2 - 1) * speedBase,
		width:       width,
		height:      height,
		isHornets:   isHornets,
		isHighSpeed: isHighSpeed,
		visible:     true,
	}

	g.bees = append(g.bees, bee)
}

// Update the game state
func (g *Game) Update() error {
	if !g.gameStarted {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.gameStarted = true
			g.startTime = time.Now()
		}
		return nil
	}

	if g.gameOver {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// Reset game
			g.bees = make([]Bee, 0)
			g.score = 0
			g.hornetsClicked = 0
			g.gameOver = false
			g.startTime = time.Now()
			g.remainingTime = gameTime
			g.lightningEffect = false
		}
		return nil
	}

	// Update remaining time
	elapsed := time.Since(g.startTime)
	g.remainingTime = gameTime - int(elapsed.Seconds())
	if g.remainingTime <= 0 {
		g.gameOver = true
		g.remainingTime = 0
		return nil
	}

	// Add new bees randomly
	if rand.Float64() < 0.05 && len(g.bees) < 10 {
		g.addBee()
	}

	// Update lightning effect timer
	if g.lightningEffect {
		g.lightningTimer--
		if g.lightningTimer <= 0 {
			g.lightningEffect = false
		}
	}

	// Update bee positions
	for i := range g.bees {
		if !g.bees[i].visible {
			continue
		}

		g.bees[i].x += g.bees[i].speedX
		g.bees[i].y += g.bees[i].speedY

		// Bounce off walls
		if g.bees[i].x <= 0 || g.bees[i].x >= float64(screenWidth-g.bees[i].width) {
			g.bees[i].speedX = -g.bees[i].speedX
		}
		if g.bees[i].y <= 0 || g.bees[i].y >= float64(screenHeight-g.bees[i].height) {
			g.bees[i].speedY = -g.bees[i].speedY
		}
	}

	// Check for mouse clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for i := range g.bees {
			if !g.bees[i].visible {
				continue
			}

			// Check if click is within bee bounds
			if float64(x) >= g.bees[i].x && float64(x) <= g.bees[i].x+float64(g.bees[i].width) &&
				float64(y) >= g.bees[i].y && float64(y) <= g.bees[i].y+float64(g.bees[i].height) {

				g.bees[i].visible = false

				if g.bees[i].isHornets {
					g.hornetsClicked++
					g.lightningEffect = true
					g.lightningTimer = 30 // Show lightning for 30 frames
					if g.hornetsClicked >= 3 {
						g.gameOver = true
					}
				} else {
					// Add score based on bee type
					if g.bees[i].isHighSpeed {
						g.score += 3 // High-speed bees worth more points
					} else {
						g.score++
					}
				}
				break
			}
		}
	}

	// Remove invisible bees
	newBees := make([]Bee, 0)
	for _, bee := range g.bees {
		if bee.visible {
			newBees = append(newBees, bee)
		}
	}
	g.bees = newBees

	return nil
}

// Draw the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw forest background
	drawForestBackground(screen)

	if !g.gameStarted {
		// Draw start screen
		msg := "Click to start the Bee Catching Game!"
		x := (screenWidth - len(msg)*7) / 2
		y := screenHeight / 2

		// Draw text with shadow for better visibility against forest background
		text.Draw(screen, msg, g.face, x+1, y+1, color.Black)
		text.Draw(screen, msg, g.face, x, y, color.White)
		return
	}

	if g.gameOver {
		// Draw game over screen
		msg := fmt.Sprintf("Game Over! Your score: %d", g.score)
		x := (screenWidth - len(msg)*7) / 2
		y := screenHeight / 2

		// Draw text with shadow for better visibility
		text.Draw(screen, msg, g.face, x+1, y+1, color.Black)
		text.Draw(screen, msg, g.face, x, y, color.White)

		msg = "Click to play again"
		x = (screenWidth - len(msg)*7) / 2
		y += 30
		text.Draw(screen, msg, g.face, x+1, y+1, color.Black)
		text.Draw(screen, msg, g.face, x, y, color.White)
		return
	}

	// Draw bees
	for _, bee := range g.bees {
		if !bee.visible {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(bee.x, bee.y)

		if bee.isHornets {
			screen.DrawImage(hornetsImage, op)
		} else {
			screen.DrawImage(beeImage, op)
		}
	}

	// Draw lightning effect
	if g.lightningEffect {
		for i := 0; i < 10; i++ {
			x1 := rand.Intn(screenWidth)
			y1 := 0
			x2 := rand.Intn(screenWidth)
			y2 := screenHeight
			ebitenutil.DrawLine(screen, float64(x1), float64(y1), float64(x2), float64(y2), color.RGBA{255, 255, 0, 192})
		}
	}

	// Draw score and time with shadow for better visibility against forest background
	scoreText := fmt.Sprintf("Score: %d", g.score)
	text.Draw(screen, scoreText, g.face, 11, 21, color.Black)
	text.Draw(screen, scoreText, g.face, 10, 20, color.White)

	timeText := fmt.Sprintf("Time: %d", g.remainingTime)
	text.Draw(screen, timeText, g.face, screenWidth-99, 21, color.Black)
	text.Draw(screen, timeText, g.face, screenWidth-100, 20, color.White)

	hornetsText := fmt.Sprintf("Hornets: %d/3", g.hornetsClicked)
	text.Draw(screen, hornetsText, g.face, 11, 41, color.Black)
	text.Draw(screen, hornetsText, g.face, 10, 40, color.White)
}

// Layout implements ebiten.Game's Layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Set up the game
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bee Catching Game")

	// ヘッドレスモードを有効にする
	ebiten.SetScreenClearedEveryFrame(false)

	// Create and run the game
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
