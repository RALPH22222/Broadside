package ui

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"log"

	"github.com/RALPH22222/Broadside/game"

	"unicode/utf8"

	"bytes"
	"image/png"

	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 768
)

// Color palette
var (
	NavyBlue     = color.RGBA{0x00, 0x1f, 0x3f, 0xff} // #001f3f
	GunmetalGray = color.RGBA{0x2f, 0x4f, 0x4f, 0xff} // #2f4f4f
	OceanTeal    = color.RGBA{0x1e, 0x65, 0x6d, 0xff} // #1e656d
	SmokeWhite   = color.RGBA{0xf0, 0xf0, 0xf0, 0xff} // #f0f0f0
	AlertRed     = color.RGBA{0xff, 0x3b, 0x3b, 0xff} // #ff3b3b
	VictoryGold  = color.RGBA{0xff, 0xd7, 0x00, 0xff} // #ffd700
)

// QuizQuestion represents a quiz question and its answers
type QuizQuestion struct {
	Question string
	Options  []string
	Answer   int // index of correct answer
}

// Level constants
const (
	LevelEasy = iota
	LevelMedium
	LevelHard
	LevelExtreme
)

var levelNames = []string{"Easy", "Medium", "Hard", "Extreme"}

// Game represents the main game state
type Game struct {
	gameFont    font.Face
	state       GameState
	menuRects   []image.Rectangle // clickable menu option areas
	hoveredMenu int               // -1 if none

	// Quiz game state
	quizQuestions []QuizQuestion
	currentQ      int
	score         int
	answerRects   []image.Rectangle // clickable answer areas
	selectedAns   int               // -1 if none
	showFeedback  bool
	feedbackTime  time.Time
	feedbackRight bool

	// Combat system
	level              int
	playerHP           int
	playerMaxHP        int
	playerShields      int
	playerMaxShields   int
	bonusActive        bool
	bonusAnswered      bool
	enemyHP            int
	enemyMaxHP         int
	bonusQDone         bool
	mainQDone          int
	bonusQIndex        int
	combatOver         bool
	rank               string
	scorePercent       int
	quiz               *game.Quiz
	subjects           []string
	selectedSubject    string
	selectedDifficulty string

	// Animation state (removed ship/fire animation fields)

	playerName      string
	enteredName     string
	nameInputActive bool

	leaderboardEntries []game.LeaderboardEntry
	leaderboardFetched bool

	// Add a smaller font for confirm button
	confirmFont font.Face

	logoImg *ebiten.Image

	// Ship images
	playerShipImg *ebiten.Image
	enemyShipImg  *ebiten.Image
	// Add for slicing enemy ship
	playerShipFrames [4]*ebiten.Image
	enemyShipFrames  [4]*ebiten.Image

	// Add userID to Game struct
	userID int64

	// Add this field to track previous mouse state
	prevMousePressed bool

	// Add this field to Game struct
	unlockedDifficulties map[string]bool // key: difficulty, value: unlocked

	// Add new state for star modal
	showStarModal   bool
	starCount       float64           // 0, 1, 1.5, 2, 2.5, 3
	continueRects   []image.Rectangle // clickable areas for modal buttons
	starModalResult int               // 0: undecided, 1: continue, 2: exit

	answeredSubjects map[string]bool // key: subject, value: answered for current difficulty

	// Timer for question answering
	questionTimer    time.Time
	questionDuration time.Duration
	timerActive      bool

	// Star modal animation state
	starModalStartTime time.Time
	starAnimationDone  bool
	starPopIndex       int
	starShineAngle     float64

	// --- Pending answer system ---
	pendingAnswer bool
	// In Game struct, add:
	// lastNavClick time.Time

	// In Game struct, add:
	compassImg *ebiten.Image

	// In Game struct, add:
	starStripImg *ebiten.Image

	// In Game struct, add fire animation state and images
	fireImgPlayer *ebiten.Image
	fireImgEnemy  *ebiten.Image
	showFire      bool
	fireStartTime time.Time
	fireType      int // 0: none, 1: player fire, 2: enemy fire
}

// GameState represents the current state of the game UI
type GameState int

const (
	StateTitle GameState = iota
	StateMenu
	StateSelectDifficulty
	StateSelectSubject
	StatePlaying
	StateGameOver
	StateNameEntry
	StateHowToPlay
	StateLeaderboard
	StateStarModal
)

var whiteImg *ebiten.Image

func ensureWhiteImg() {
	if whiteImg == nil {
		whiteImg = ebiten.NewImage(1, 1)
		whiteImg.Fill(color.White)
	}
}

func colorToRGBA(c color.Color) color.RGBA {
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < ScreenHeight; y++ {
		frac := float64(y) / float64(ScreenHeight)
		col := color.RGBA{
			R: uint8(20 + 60*frac),
			G: uint8(40 + 80*frac),
			B: uint8(100 + 100*frac),
			A: 255,
		}
		for x := 0; x < ScreenWidth; x++ {
			screen.Set(x, y, col)
		}
	}
	// Only draw sun, clouds, and seagulls if not in leaderboard state
	if g.state != StateLeaderboard {
		// Draw a small realistic sun at the top right edge
		drawSun(screen, ScreenWidth-80, 80, 38)
		// Animated scanline overlay
		t := float64(time.Now().UnixNano()%1_000_000_000) / 1_000_000_000
		for y := 0; y < ScreenHeight; y += 4 {
			alpha := uint8(24 + 16*math.Sin(2*math.Pi*(float64(y)/32)+t*2*math.Pi))
			for x := 0; x < ScreenWidth; x++ {
				screen.Set(x, y, color.RGBA{0, 0, 0, alpha})
			}
		}
		// --- Static realistic clouds ---
		drawRealisticCloud(screen, 150, 80, 1.2)
		drawRealisticCloud(screen, 450, 60, 0.9)
		drawRealisticCloud(screen, 750, 100, 1.5)
		drawRealisticCloud(screen, 250, 120, 0.8)
		drawRealisticCloud(screen, 650, 90, 1.1)
		// --- Smoother animated seagulls (left to right) ---
		for i := 0; i < 5; i++ {
			phase := float64(i) * 0.7
			gullX := math.Mod(t*ScreenWidth*(0.12+0.04*float64(i))+float64(i)*220, float64(ScreenWidth+120)) - 60
			gullY := 120 + float64(i)*22 + 18*math.Sin(t*1.2+phase)
			wing := 0.7 + 0.5*math.Sin(t*3+phase)
			drawSeagull(screen, gullX, gullY, 0.9+0.13*float64(i), wing)
		}
		// --- Smoother animated seagulls (right to left, mirrored) ---
		for i := 0; i < 5; i++ {
			phase := float64(i) * 0.9
			gullX := ScreenWidth - (math.Mod(t*ScreenWidth*(0.09+0.03*float64(i))+float64(i)*180, float64(ScreenWidth+120)) - 60)
			gullY := 180 + float64(i)*18 + 14*math.Sin(t*1.5+phase)
			wing := -(0.7 + 0.5*math.Sin(t*2.7+phase))
			drawSeagull(screen, gullX, gullY, 0.8+0.11*float64(i), wing)
		}
	}

	switch g.state {
	case StateTitle:
		g.drawTitle(screen)
	case StateMenu:
		g.drawMenu(screen)
	case StateSelectDifficulty:
		g.drawSelectDifficulty(screen)
	case StateSelectSubject:
		g.drawSelectSubject(screen)
	case StateHowToPlay:
		g.drawHowToPlayOverlay(screen)
	case StateLeaderboard:
		g.drawLeaderboardOverlay(screen)
	case StateGameOver:
		g.drawGameOver(screen)
	case StateNameEntry:
		g.drawNameEntry(screen)
	case StateStarModal:
		g.drawStarModal(screen)
	case StatePlaying:
		g.drawPlaying(screen)
	}

	// Show feedback prominently in the center of the screen when active
	if g.showFeedback {
		var msg string
		var col color.Color
		if g.selectedAns == -1 {
			msg = "Time's up!"
			col = AlertRed // Use AlertRed for time's up (same as incorrect)
		} else if g.feedbackRight {
			msg = "Correct!"
			col = VictoryGold // Use VictoryGold for correct
		} else {
			msg = "Incorrect!"
			col = AlertRed
		}

		feedbackFontSize := 32 // Smaller font
		bounds, _ := font.BoundString(g.gameFont, msg)
		msgWidth := (bounds.Max.X - bounds.Min.X).Ceil()
		msgHeight := (bounds.Max.Y - bounds.Min.Y).Ceil()

		bgPadding := 12 // Smaller padding
		bgWidth := msgWidth + bgPadding*2
		bgHeight := msgHeight + bgPadding*2
		bgX := (ScreenWidth - bgWidth) / 2
		bgY := (ScreenHeight-bgHeight)/2 + 200 // Move further downward

		for r := 0; r < 3; r++ {
			alpha := uint8(60 - r*20)
			glowCol := color.RGBA{col.(color.RGBA).R, col.(color.RGBA).G, col.(color.RGBA).B, alpha}
			vector.StrokeRect(screen, float32(bgX-r*2), float32(bgY-r*2), float32(bgWidth+2*r*2), float32(bgHeight+2*r*2), 6, glowCol, true)
		}

		vector.DrawFilledRect(screen, float32(bgX), float32(bgY), float32(bgWidth), float32(bgHeight), GunmetalGray, true)
		vector.StrokeRect(screen, float32(bgX), float32(bgY), float32(bgWidth), float32(bgHeight), 3, col, true)

		textX := bgX + bgPadding
		textY := bgY + bgPadding + msgHeight
		drawWrappedTextWithShadow(screen, msg, g.gameFont, textX, textY, bgWidth, feedbackFontSize, col)
	}
}

// Helper: wrap text to fit within a max width (in pixels)
func wrapText(face font.Face, textStr string, maxWidth int) []string {
	var lines []string
	words := []rune(textStr)
	start := 0
	for start < len(words) {
		end := start
		lastSpace := -1
		for end < len(words) {
			ch := words[end]
			if ch == '\n' {
				break
			}
			testStr := string(words[start : end+1])
			bounds, _ := font.BoundString(face, testStr)
			width := (bounds.Max.X - bounds.Min.X).Ceil()
			if width > maxWidth {
				break
			}
			if ch == ' ' {
				lastSpace = end
			}
			end++
		}
		if end < len(words) && words[end] == '\n' {
			lines = append(lines, string(words[start:end]))
			start = end + 1
			continue
		}
		bounds, _ := font.BoundString(face, string(words[start:end]))
		width := (bounds.Max.X - bounds.Min.X).Ceil()
		if end > start && width > maxWidth && lastSpace > start {
			lines = append(lines, string(words[start:lastSpace]))
			start = lastSpace + 1
		} else {
			lines = append(lines, string(words[start:end]))
			start = end
		}
	}
	return lines
}

// Helper: draw wrapped text with shadow
func drawWrappedTextWithShadow(screen *ebiten.Image, str string, face font.Face, x, y, maxWidth, lineHeight int, col color.Color) {
	lines := wrapText(face, str, maxWidth)
	for i, line := range lines {
		shadow := color.RGBA{0, 0, 0, 180}
		text.Draw(screen, line, face, x+2, y+2+i*lineHeight, shadow)
		text.Draw(screen, line, face, x, y+i*lineHeight, col)
	}
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	// --- Premium Title Screen Redesign with Animation & Effects ---
	t := float64(time.Now().UnixNano()) / 1e9
	centerX := ScreenWidth / 2

	// 1. Animated Logo Bobbing
	logoH := 220
	logoW := 220
	bob := int(18 * math.Sin(t*1.2))
	logoY := ScreenHeight/7 + bob
	logoX := centerX - logoW/2
	if g.logoImg != nil {
		op := &ebiten.DrawImageOptions{}
		sx := float64(logoW) / float64(g.logoImg.Bounds().Dx())
		sy := float64(logoH) / float64(g.logoImg.Bounds().Dy())
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(float64(logoX), float64(logoY))
		screen.DrawImage(g.logoImg, op)
	}

	// 2. Title Text with Pulse and Gold Gradient/Glow
	title := "BROADSIDE"
	fontSize := 64
	fontFace := g.gameFont
	bounds, _ := font.BoundString(fontFace, title)
	titleW := (bounds.Max.X - bounds.Min.X).Ceil()
	pulse := 1.0 + 0.04*math.Sin(t*2.1)
	titleX := float64(centerX) - float64(titleW)*pulse/2
	titleY := float64(logoY + logoH + 32)

	// Draw glow and gold gradient with pulse (draw to temp image, then scale)
	tempW := int(float64(titleW) * 1.2)
	tempH := fontSize + 32
	tempImg := ebiten.NewImage(tempW, tempH)
	// Glow
	for r := 8; r > 0; r -= 2 {
		glowCol := color.RGBA{255, 215, 0, uint8(30 + 10*r)}
		text.Draw(tempImg, title, fontFace, tempW/2-titleW/2+r, tempH/2+fontSize/2+r, glowCol)
		text.Draw(tempImg, title, fontFace, tempW/2-titleW/2-r, tempH/2+fontSize/2-r, glowCol)
	}
	// Gold gradient effect
	for i := 0; i < fontSize; i += 2 {
		frac := float64(i) / float64(fontSize)
		col := color.RGBA{
			R: uint8(255 - 40*frac),
			G: uint8(215 + 30*frac),
			B: uint8(0 + 80*frac),
			A: 255,
		}
		text.Draw(tempImg, title, fontFace, tempW/2-titleW/2, tempH/2+fontSize/2-i/4, col)
	}
	// Draw scaled temp image for pulse
	titleOp := &ebiten.DrawImageOptions{}
	titleOp.GeoM.Scale(pulse, pulse)
	titleOp.GeoM.Translate(titleX-float64(tempW)*pulse/2+float64(tempW)/2, titleY-float64(tempH)*pulse/2+float64(tempH)/2)
	screen.DrawImage(tempImg, titleOp)

	// 3. Sparkle Particles
	nSparkles := 18
	for i := 0; i < nSparkles; i++ {
		angle := t*0.7 + float64(i)*2*math.Pi/float64(nSparkles)
		radius := float64(titleW)/2 + 40 + 18*math.Sin(t*1.5+float64(i))
		sx := float64(centerX) + radius*math.Cos(angle)
		sy := titleY + 32 + 18*math.Sin(angle*2+t)
		size := 2.5 + 2.5*math.Abs(math.Sin(t*2+float64(i)))
		alpha := uint8(120 + 80*math.Sin(t*2+float64(i)))
		vector.DrawFilledCircle(screen, float32(sx), float32(sy), float32(size), color.RGBA{255, 215, 0, alpha}, true)
	}

	// 4. Subtitle with Fade-in
	sub := "Press SPACE to start"
	fadeIn := math.Min(1, math.Max(0, (t-0.7)/1.2))
	bounds, _ = font.BoundString(fontFace, sub)
	subW := (bounds.Max.X - bounds.Min.X).Ceil()
	subH := (bounds.Max.Y - bounds.Min.Y).Ceil()
	subX := centerX - subW/2
	subY := int(titleY) + fontSize + 48
	bgPadX := 32
	bgPadY := 18
	bgCol := color.RGBA{30, 65, 109, uint8(180 * fadeIn)}
	vector.DrawFilledRect(screen, float32(subX-bgPadX), float32(subY-bgPadY), float32(subW+2*bgPadX), float32(subH+2*bgPadY), bgCol, true)
	vector.StrokeRect(screen, float32(subX-bgPadX), float32(subY-bgPadY), float32(subW+2*bgPadX), float32(subH+2*bgPadY), 3, VictoryGold, true)
	col := color.RGBA{SmokeWhite.R, SmokeWhite.G, SmokeWhite.B, uint8(255 * fadeIn)}
	drawWrappedTextWithShadow(screen, sub, fontFace, subX, subY+subH, subW, 36, col)

	// 5. Improved Shine Sweep (curved)
	shineW := titleW / 3
	shineFrac := math.Mod(t/2, 1.0)
	shineX := titleX + float64(titleW+shineW)*shineFrac - float64(shineW)
	shineY := titleY - float64(fontSize)/2 + 12*math.Sin(shineFrac*math.Pi)
	for dx := 0; dx < shineW; dx++ {
		alpha := uint8(80 * (1 - math.Abs(float64(dx-shineW/2))/float64(shineW/2)))
		if alpha > 0 {
			vector.DrawFilledRect(screen, float32(shineX+float64(dx)), float32(shineY), 2, float32(fontSize+12), color.RGBA{255, 255, 255, alpha}, true)
		}
	}

	// --- Soft vignette overlay for depth ---
	for i := 0; i < 80; i++ {
		alpha := uint8(120 * float64(i) / 80)
		vector.StrokeRect(screen, float32(i), float32(i), float32(ScreenWidth-2*i), float32(ScreenHeight-2*i), 2, color.RGBA{0, 0, 0, alpha}, true)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	msg := "Main Menu"
	drawWrappedTextWithShadow(screen, msg, g.gameFont, ScreenWidth/10, ScreenHeight/8, ScreenWidth*8/10, 36, VictoryGold)

	options := []string{
		"Play",
		"How to Play",
		"Leaderboard",
		"Exit",
	}
	g.menuRects = g.menuRects[:0]
	menuW := ScreenWidth * 5 / 12
	menuX := (ScreenWidth - menuW) / 2
	menuH := ScreenHeight / 14
	startY := ScreenHeight/2 - (len(options)*menuH+(len(options)-1)*ScreenHeight/48)/2
	for i, opt := range options {
		btnY := startY + i*(menuH+ScreenHeight/48)
		w := menuW
		h := menuH
		// Neon glow effect
		glowColor := OceanTeal
		if g.hoveredMenu == i {
			glowColor = VictoryGold
		}
		for r := 0; r < 3; r++ {
			alpha := uint8(60 - r*20)
			vector.StrokeRect(screen, float32(menuX-r*3), float32(btnY-r*3), float32(w+2*r*3), float32(h+2*r*3), 6, color.RGBA{glowColor.R, glowColor.G, glowColor.B, alpha}, true)
		}
		// Button background (metallic gradient)
		for dy := 0; dy < h; dy++ {
			frac := float64(dy) / float64(h)
			c := color.RGBA{
				R: uint8(47 + 40*frac),
				G: uint8(79 + 40*frac),
				B: uint8(79 + 60*frac),
				A: 255,
			}
			for dx := 0; dx < w; dx++ {
				screen.Set(menuX+dx, btnY+dy, c)
			}
		}
		// Button border
		vector.StrokeRect(screen, float32(menuX), float32(btnY), float32(w), float32(h), 4, glowColor, true)
		// Button text
		textCol := SmokeWhite
		if g.hoveredMenu == i {
			textCol = NavyBlue
		}
		bounds, _ := font.BoundString(g.gameFont, opt)
		width := (bounds.Max.X - bounds.Min.X).Ceil()
		strX := menuX + (w-width)/2
		strY := btnY + h/2 + 12
		drawWrappedTextWithShadow(screen, opt, g.gameFont, strX, strY, width, 36, textCol)
		// Save clickable area
		rect := image.Rect(menuX, btnY, menuX+w, btnY+h)
		g.menuRects = append(g.menuRects, rect)
	}
}

func (g *Game) drawSelectDifficulty(screen *ebiten.Image) {
	msg := "Select Difficulty"
	drawWrappedTextWithShadow(screen, msg, g.gameFont, ScreenWidth/10, ScreenHeight/8, ScreenWidth*8/10, 36, VictoryGold)

	difficulties := []string{"Easy", "Medium", "Hard", "Extreme"}
	g.menuRects = g.menuRects[:0]
	menuW := ScreenWidth * 5 / 12
	menuX := (ScreenWidth - menuW) / 2
	menuH := ScreenHeight / 14
	startY := ScreenHeight/2 - (len(difficulties)*menuH+(len(difficulties)-1)*ScreenHeight/48)/2
	for i, diff := range difficulties {
		btnY := startY + i*(menuH+ScreenHeight/48)
		w := menuW
		h := menuH
		locked := false
		if diff != "Easy" && !g.unlockedDifficulties[diff] {
			locked = true
		}
		glowColor := OceanTeal
		if g.hoveredMenu == i && !locked {
			glowColor = VictoryGold
		}
		for r := 0; r < 3; r++ {
			alpha := uint8(60 - r*20)
			vector.StrokeRect(screen, float32(menuX-r*3), float32(btnY-r*3), float32(w+2*r*3), float32(h+2*r*3), 6, color.RGBA{glowColor.R, glowColor.G, glowColor.B, alpha}, true)
		}
		for dy := 0; dy < h; dy++ {
			frac := float64(dy) / float64(h)
			c := color.RGBA{
				R: uint8(47 + 40*frac),
				G: uint8(79 + 40*frac),
				B: uint8(79 + 60*frac),
				A: 255,
			}
			if locked {
				c = GunmetalGray
			}
			for dx := 0; dx < w; dx++ {
				screen.Set(menuX+dx, btnY+dy, c)
			}
		}
		vector.StrokeRect(screen, float32(menuX), float32(btnY), float32(w), float32(h), 4, glowColor, true)
		textCol := SmokeWhite
		if locked {
			textCol = GunmetalGray
		} else if g.hoveredMenu == i {
			textCol = NavyBlue
		}
		bounds, _ := font.BoundString(g.gameFont, diff)
		width := (bounds.Max.X - bounds.Min.X).Ceil()
		strX := menuX + (w-width)/2
		strY := btnY + h/2 + 12
		drawWrappedTextWithShadow(screen, diff, g.gameFont, strX, strY, width, 36, textCol)
		rect := image.Rect(menuX, btnY, menuX+w, btnY+h)
		g.menuRects = append(g.menuRects, rect)
	}
}

func (g *Game) drawSelectSubject(screen *ebiten.Image) {
	msg := "Select Subject"
	drawWrappedTextWithShadow(screen, msg, g.gameFont, ScreenWidth/10, ScreenHeight/8, ScreenWidth*8/10, 36, VictoryGold)

	if len(g.subjects) == 0 {
		g.subjects = g.quiz.ListSubjects()
	}
	g.menuRects = g.menuRects[:0]
	menuW := ScreenWidth * 5 / 12
	menuX := (ScreenWidth - menuW) / 2
	menuH := ScreenHeight / 14
	startY := ScreenHeight/2 - (len(g.subjects)*menuH+(len(g.subjects)-1)*ScreenHeight/48)/2
	for i, subj := range g.subjects {
		btnY := startY + i*(menuH+ScreenHeight/48)
		w := menuW
		h := menuH
		glowColor := OceanTeal
		if g.hoveredMenu == i {
			glowColor = VictoryGold
		}
		for r := 0; r < 3; r++ {
			alpha := uint8(60 - r*20)
			vector.StrokeRect(screen, float32(menuX-r*3), float32(btnY-r*3), float32(w+2*r*3), float32(h+2*r*3), 6, color.RGBA{glowColor.R, glowColor.G, glowColor.B, alpha}, true)
		}
		for dy := 0; dy < h; dy++ {
			frac := float64(dy) / float64(h)
			c := color.RGBA{
				R: uint8(47 + 40*frac),
				G: uint8(79 + 40*frac),
				B: uint8(79 + 60*frac),
				A: 255,
			}
			for dx := 0; dx < w; dx++ {
				screen.Set(menuX+dx, btnY+dy, c)
			}
		}
		vector.StrokeRect(screen, float32(menuX), float32(btnY), float32(w), float32(h), 4, glowColor, true)
		textCol := SmokeWhite
		if g.hoveredMenu == i {
			textCol = NavyBlue
		}
		bounds, _ := font.BoundString(g.gameFont, subj)
		width := (bounds.Max.X - bounds.Min.X).Ceil()
		strX := menuX + (w-width)/2
		strY := btnY + h/2 + 12
		drawWrappedTextWithShadow(screen, subj, g.gameFont, strX, strY, width, 36, textCol)
		rect := image.Rect(menuX, btnY, menuX+w, btnY+h)
		g.menuRects = append(g.menuRects, rect)
	}
}

// Enhance overlays
func (g *Game) drawHowToPlayOverlay(screen *ebiten.Image) {
	w, h := 600, 340
	x := (ScreenWidth - w) / 2
	y := (ScreenHeight - h) / 2
	// Neon border
	for r := 0; r < 3; r++ {
		alpha := uint8(60 - r*20)
		vector.StrokeRect(screen, float32(x-r*3), float32(y-r*3), float32(w+2*r*3), float32(h+2*r*3), 6, color.RGBA{OceanTeal.R, OceanTeal.G, OceanTeal.B, alpha}, true)
	}
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), GunmetalGray, true)
	vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 4, VictoryGold, true)
	msg := "How to Play:\n- Choose the correct answer within 10 seconds to kill enemies\n- Use a safety shield to survive wrong answers up to 3 times\n- Defeat the enemies!!\n- ESC to return to menu\n(Press ESC or click to close)"
	textPaddingX := 20
	drawWrappedTextWithShadow(screen, msg, g.gameFont, x+textPaddingX, y+80, w-(textPaddingX*2), 36, SmokeWhite)
}

func (g *Game) drawLeaderboardOverlay(screen *ebiten.Image) {
	w, h := ScreenWidth, ScreenHeight
	fontFace := g.confirmFont
	if fontFace == nil {
		fontFace = g.gameFont
	}

	// --- Background vignette ---
	for i := 0; i < 80; i++ {
		alpha := uint8(120 * float64(i) / 80)
		vector.StrokeRect(screen, float32(i), float32(i), float32(w-2*i), float32(h-2*i), 2, color.RGBA{0, 0, 0, alpha}, true)
	}

	// --- Leaderboard dimensions ---
	leaderW := w * 2 / 5
	lineHeight := 40
	cellPadY := 6 // Vertical padding for cell backgrounds
	maxRows := 11 // Includes header row (row 0) + 10 data rows (rows 1-10)
	leaderH := lineHeight*maxRows + 160
	leaderX := w / 8
	leaderY := h/2 - leaderH/2

	// --- Logo ---
	logoW, logoH := 220, 220
	logoX := leaderX + leaderW + 80
	logoY := leaderY + leaderH/2 - logoH/2
	if g.logoImg != nil {
		logoOp := &ebiten.DrawImageOptions{}
		sx := float64(logoW) / float64(g.logoImg.Bounds().Dx())
		sy := float64(logoH) / float64(g.logoImg.Bounds().Dy())
		logoOp.GeoM.Scale(sx, sy)
		logoOp.GeoM.Translate(float64(logoX), float64(logoY))
		screen.DrawImage(g.logoImg, logoOp)
	}

	// --- Table background ---
	shadowCol := color.RGBA{0, 0, 0, 80}
	vector.DrawFilledRect(screen, float32(leaderX+8), float32(leaderY+8), float32(leaderW), float32(leaderH), shadowCol, true)
	vector.DrawFilledRect(screen, float32(leaderX), float32(leaderY), float32(leaderW), float32(leaderH), GunmetalGray, true)
	vector.StrokeRect(screen, float32(leaderX), float32(leaderY), float32(leaderW), float32(leaderH), 5, VictoryGold, true)

	// --- Title ---
	premiumTitle := "LEADERBOARD"
	fontSize := 96
	titleBounds, _ := font.BoundString(fontFace, premiumTitle)
	titleW := (titleBounds.Max.X - titleBounds.Min.X).Ceil()
	titleX := leaderX + (leaderW-titleW)/2
	titleY := leaderY + 10
	text.Draw(screen, premiumTitle, fontFace, titleX, titleY+fontSize, SmokeWhite)

	// --- Columns ---
	headers := []string{"Name", "Score"}
	const tableInnerPadX = 24

	// Calculate the total usable width for columns within the table's inner padding
	effectiveTableWidth := leaderW - 2*tableInnerPadX

	// Distribute the effective width to columns
	col1Width := effectiveTableWidth * 2 / 3
	col2Width := effectiveTableWidth - col1Width

	colWidths := []int{col1Width, col2Width}

	// colX[j] will now be the left edge of the background rectangle for column j
	colX := []int{
		leaderX + tableInnerPadX,             // Left edge of the first column
		leaderX + tableInnerPadX + col1Width, // Left edge of the second column
	}

	tableY := leaderY + 120 // Top Y for the text baseline of the first row (header)

	// --- FIX --- Define the exact Y coordinates for the top and bottom of the entire grid area
	// These will be used for both row backgrounds and horizontal/vertical lines.
	tableGridTopY := float32(tableY - cellPadY)                         // Top of the header row's background
	tableGridBottomY := float32(tableY + maxRows*lineHeight + cellPadY) // Bottom of the last data row's background

	// --- Rows (including header) ---
	for i := 0; i < maxRows; i++ {
		rowY := tableY + i*lineHeight     // Y position for text baseline
		rowBgY := rowY - cellPadY         // Y position for background rectangle top
		rowBgH := lineHeight + 2*cellPadY // Height of background rectangle

		for j := 0; j < len(colX); j++ {
			cellLeftX := colX[j]      // Left edge of cell background
			cellWidth := colWidths[j] // Width of cell background

			var bgColor color.RGBA
			if i == 0 {
				bgColor = color.RGBA{60, 60, 60, 240} // Header background (dark gray)
			} else if (i % 2) == 1 {
				bgColor = color.RGBA{45, 45, 45, 180} // Slightly lighter dark gray
			} else {
				bgColor = color.RGBA{35, 35, 35, 150} // Slightly darker dark gray
			}

			//Draw cell background using the precise column boundaries
			vector.DrawFilledRect(screen, float32(cellLeftX), float32(rowBgY), float32(cellWidth), float32(rowBgH), bgColor, true)
		}

		var rowVals []string
		if i == 0 {
			rowVals = headers
		} else if i-1 < len(g.leaderboardEntries) {
			entry := g.leaderboardEntries[i-1]
			rowVals = []string{strings.TrimSpace(entry.PlayerName), itoa(entry.Score)}
		} else {
			rowVals = []string{"", ""} // Empty rows if not enough entries
		}

		for j, val := range rowVals {
			bounds, _ := font.BoundString(fontFace, val)
			textW := (bounds.Max.X - bounds.Min.X).Ceil()
			textH := (bounds.Max.Y - bounds.Min.Y).Ceil()
			cellLeftX := colX[j]
			cellWidth := colWidths[j]
			textX := cellLeftX + (cellWidth-textW)/2
			textY := rowY + (lineHeight+textH)/2 - 2 // Adjust for font baseline
			text.Draw(screen, val, fontFace, textX, textY, SmokeWhite)
		}
	}

	// --- Grid lines ---
	// --- FIX --- Define the exact X coordinates for the left and right of the entire grid area
	// These must match the boundaries used for drawing cell backgrounds.
	tableGridLeftX := float32(leaderX + tableInnerPadX)
	tableGridRightX := float32(leaderX + tableInnerPadX + col1Width + col2Width) // Right edge of the last column

	// Vertical lines
	// --- FIX --- Draw each vertical line exactly at the column boundaries
	vector.StrokeLine(screen,
		tableGridLeftX, tableGridTopY,
		tableGridLeftX, tableGridBottomY,
		2, VictoryGold, true) // Leftmost vertical line

	vector.StrokeLine(screen,
		float32(colX[1]), tableGridTopY, // Line between column 1 and 2
		float32(colX[1]), tableGridBottomY,
		2, VictoryGold, true)

	vector.StrokeLine(screen,
		tableGridRightX, tableGridTopY,
		tableGridRightX, tableGridBottomY,
		2, VictoryGold, true) // Rightmost vertical line (drawn only once)

	// Horizontal lines
	for i := 0; i <= maxRows; i++ {
		// --- FIX --- Calculate lineY to align with the top/bottom of row backgrounds
		lineY := float32(tableY + i*lineHeight - cellPadY)
		if i == maxRows { // For the very last line, ensure it's at the bottom of the last row's background
			lineY = tableGridBottomY
		}
		vector.StrokeLine(screen,
			tableGridLeftX, lineY,
			tableGridRightX, lineY, // --- FIX --- Ensure horizontal lines span the exact table width
			2, VictoryGold, true)
	}

	// --- Footer ---
	msg2 := "(Press ESC or click to close)"
	msg2Bounds, _ := font.BoundString(fontFace, msg2)
	msg2Width := (msg2Bounds.Max.X - msg2Bounds.Min.X).Ceil()
	drawWrappedTextWithShadow(screen, msg2, fontFace, (w-msg2Width)/2, h-48, w-48, 32, SmokeWhite)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Draw HP, shields, enemy HP
	barY := 40
	drawWrappedTextWithShadow(screen, "Your HP: "+itoa(g.playerHP)+"/"+itoa(g.playerMaxHP), g.gameFont, 40, barY, ScreenWidth-200, 36, VictoryGold)
	drawWrappedTextWithShadow(screen, "Shields: "+itoa(g.playerShields)+"/"+itoa(g.playerMaxShields), g.gameFont, 40, barY+40, ScreenWidth-200, 36, OceanTeal)
	drawWrappedTextWithShadow(screen, "Enemy HP: "+itoa(g.enemyHP)+"/"+itoa(g.enemyMaxHP), g.gameFont, ScreenWidth-340, barY, ScreenWidth-200, 36, AlertRed)
	drawWrappedTextWithShadow(screen, "Level: "+levelNames[g.level], g.gameFont, ScreenWidth-340, barY+40, ScreenWidth-200, 36, SmokeWhite)

	// Draw timer if active
	if g.timerActive && !g.showFeedback {
		remaining := g.questionDuration - time.Since(g.questionTimer)
		if remaining > 0 {
			seconds := int(remaining.Seconds()) + 1
			timerColor := SmokeWhite
			if seconds <= 3 {
				timerColor = AlertRed
			} else if seconds <= 5 {
				timerColor = VictoryGold
			}
			timerText := "Time: " + itoa(seconds) + "s"
			drawWrappedTextWithShadow(screen, timerText, g.gameFont, 40, barY+80, ScreenWidth-200, 36, timerColor)
		} else {
			// Time's up - handle timeout
			if g.currentQ < len(g.quizQuestions) {
				isBonusQ := g.currentQ == 0
				questionsLeft := len(g.quizQuestions) - g.currentQ - 1

				// Process timeout as incorrect answer
				cs := game.ProcessAnswer(
					g.level, g.playerHP, g.playerShields, g.enemyHP, g.score, g.mainQDone, g.bonusActive, g.bonusAnswered,
					g.bonusQIndex, g.currentQ, isBonusQ, false, questionsLeft,
				)

				// Update game state
				g.playerHP = cs.PlayerHP
				g.playerShields = cs.PlayerShields
				g.enemyHP = cs.EnemyHP
				g.score = cs.Score
				g.mainQDone = cs.MainQDone
				g.bonusActive = cs.BonusActive
				g.bonusAnswered = cs.BonusAnswered
				g.combatOver = cs.CombatOver

				// Show timeout feedback first
				g.showFeedback = true
				g.feedbackTime = time.Now()
				g.feedbackRight = false
				g.selectedAns = -1 // No answer selected
				g.timerActive = false
			}
		}
	}

	// Only show question box and options
	if g.currentQ < len(g.quizQuestions) {
		q := g.quizQuestions[g.currentQ]
		optionFontSize := 32
		optionFontFace := g.gameFont
		btnGap := 20
		// --- Calculate question lines and height ---
		maxQuestionLines := 2
		questionFontSizeNormal := 48
		questionFontSizeSmall := 32
		questionFontSizeToUse := questionFontSizeNormal
		questionLineHeight := questionFontSizeToUse + 10
		questionW := ScreenWidth * 7 / 10
		qLines := wrapText(g.gameFont, q.Question, questionW)
		if len(qLines) > maxQuestionLines {
			questionFontSizeToUse = questionFontSizeSmall
			questionLineHeight = questionFontSizeToUse + 8
			qLines = wrapText(g.gameFont, q.Question, questionW)
		}
		questionHeight := len(qLines) * questionLineHeight
		// --- Calculate options height ---
		optionHeights := make([]int, len(q.Options))
		optionLines := make([][]string, len(q.Options))
		totalOptionsHeight := 0
		btnW := questionW - 48
		for i, opt := range q.Options {
			optLines := wrapText(optionFontFace, opt, btnW-24)
			btnH := len(optLines)*(optionFontSize+6) + 12
			optionHeights[i] = btnH
			optionLines[i] = optLines
			totalOptionsHeight += btnH
			if i < len(q.Options)-1 {
				totalOptionsHeight += btnGap
			}
		}
		// --- Calculate total box size and position ---
		padding := 32
		boxW := questionW + padding*2
		boxH := questionHeight + totalOptionsHeight + padding*3 + 16
		boxX := (ScreenWidth - boxW) / 2
		boxY := (ScreenHeight - boxH) / 2
		// --- Draw the main rectangle ---
		vector.DrawFilledRect(screen, float32(boxX), float32(boxY), float32(boxW), float32(boxH), GunmetalGray, true)
		vector.StrokeRect(screen, float32(boxX), float32(boxY), float32(boxW), float32(boxH), 4, VictoryGold, true)
		// --- Draw question lines ---
		questionX := boxX + padding
		questionY := boxY + padding
		for i, line := range qLines {
			lineY := questionY + i*questionLineHeight
			shadow := color.RGBA{0, 0, 0, 180}
			text.Draw(screen, line, g.gameFont, questionX+2, lineY+2, shadow)
			text.Draw(screen, line, g.gameFont, questionX, lineY, SmokeWhite)
		}
		// --- Draw options as buttons inside the rectangle ---
		optionsStartY := questionY + questionHeight + padding
		btnX := questionX
		btnY := optionsStartY
		g.answerRects = g.answerRects[:0]

		for i, optLines := range optionLines {
			btnH := optionHeights[i]
			bgCol := GunmetalGray
			if g.selectedAns == i {
				bgCol = OceanTeal
			}
			vector.DrawFilledRect(screen, float32(btnX), float32(btnY), float32(btnW), float32(btnH), bgCol, true)
			vector.StrokeRect(screen, float32(btnX), float32(btnY), float32(btnW), float32(btnH), 3, VictoryGold, true)
			textCol := SmokeWhite
			if g.selectedAns == i {
				textCol = VictoryGold
			}
			// Vertically center the text block in the button
			textBlockHeight := len(optLines) * (optionFontSize + 6)
			startTextY := btnY + (btnH-textBlockHeight)/2 + 25
			for j, line := range optLines {
				bounds, _ := font.BoundString(optionFontFace, line)
				optWidth := (bounds.Max.X - bounds.Min.X).Ceil()
				optX := btnX + (btnW-optWidth)/2
				optY := startTextY + j*(optionFontSize+6)
				shadow := color.RGBA{0, 0, 0, 180}
				text.Draw(screen, line, optionFontFace, optX+2, optY+2, shadow)
				text.Draw(screen, line, optionFontFace, optX, optY, textCol)
			}
			// Save clickable area
			rect := image.Rect(btnX, btnY, btnX+btnW, btnY+btnH)
			g.answerRects = append(g.answerRects, rect)
			btnY += btnH + btnGap

		}
	} else if g.currentQ >= len(g.quizQuestions) && !g.combatOver {
		// Show a placeholder when all questions are answered but combat isn't over yet
		placeholderMsg := "All questions answered!"
		placeholderCol := VictoryGold
		// Calculate placeholder box size
		questionW := ScreenWidth * 7 / 10
		questionHeight := 48 // Single line height
		totalOptionsHeight := 48
		// Calculate total box size and position
		padding := 32
		boxW := questionW + padding*2
		boxH := questionHeight + totalOptionsHeight + padding*3 + 16
		boxX := (ScreenWidth - boxW) / 2
		boxY := (ScreenHeight - boxH) / 2
		// Draw the main rectangle
		vector.DrawFilledRect(screen, float32(boxX), float32(boxY), float32(boxW), float32(boxH), GunmetalGray, true)
		vector.StrokeRect(screen, float32(boxX), float32(boxY), float32(boxW), float32(boxH), 4, VictoryGold, true)
		// Draw placeholder message
		questionX := boxX + padding
		questionY := boxY + padding
		bounds, _ := font.BoundString(g.gameFont, placeholderMsg)
		msgWidth := (bounds.Max.X - bounds.Min.X).Ceil()
		msgX := questionX + (questionW-msgWidth)/2
		drawWrappedTextWithShadow(screen, placeholderMsg, g.gameFont, msgX, questionY+36, questionW, 48, placeholderCol)
	}

	// Draw ships below the question box
	g.drawShips(screen)

	if g.combatOver {
		// Handle unanswered questions based on victory/defeat conditions
		if g.currentQ < len(g.quizQuestions) {
			unanswered := len(g.quizQuestions) - g.currentQ
			if unanswered > 0 {
				// If enemy HP is 0 (victory), award points for unanswered questions
				if g.enemyHP == 0 {
					// Award 10 points per unanswered question (as specified)
					g.score += unanswered * 10
					g.mainQDone += unanswered
				}
				// If player HP is 0 (defeat), no points for unanswered questions
				g.currentQ = len(g.quizQuestions) // Mark all as answered
			}
		}
		// Show star modal immediately
		if !g.showStarModal {
			g.showFeedback = false
			g.selectedAns = -1
			g.showStarModal = true
			g.state = StateStarModal
			// Calculate star count (updated thresholds)
			percent := g.scorePercent
			var stars float64
			switch {
			case g.playerHP == 0 || percent < 50:
				stars = 0
			case percent == 100:
				stars = 3
			case percent >= 90:
				stars = 2.5
			case percent >= 80:
				stars = 2
			case percent >= 70:
				stars = 1.5
			case percent >= 51:
				stars = 1
			default:
				stars = 0
			}
			g.starCount = stars
			// Reset animation state
			g.starModalStartTime = time.Now()
			g.starAnimationDone = false
			g.starPopIndex = 0
			g.starShineAngle = 0.0
		}
		g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		return
	}

	// After answer, show fire if needed
	if g.showFeedback && !g.showFire && g.selectedAns != -1 && g.currentQ < len(g.quizQuestions) {
		isBonusQ := g.currentQ == 0
		if !isBonusQ {
			if g.feedbackRight {
				g.showFire = true
				g.fireStartTime = time.Now()
				g.fireType = 1 // player fire
			} else {
				g.showFire = true
				g.fireStartTime = time.Now()
				g.fireType = 2 // enemy fire
			}
		}
	}
	// Block question advance until fire is done
	if g.showFire {
		if time.Since(g.fireStartTime) >= time.Second {
			g.showFire = false
			g.fireType = 0
			// After fire, hide feedback and advance question
			g.showFeedback = false
			g.selectedAns = -1
			g.currentQ++
			if g.currentQ < len(g.quizQuestions) {
				g.timerActive = true
				g.questionTimer = time.Now()
			}
		} else {
			// Draw ships with fire overlay
			g.drawShips(screen)
			return
		}
	}
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	msg := "Game Over! (Press ESC to return to menu)"
	x := (ScreenWidth - len(msg)*14) / 2
	y := ScreenHeight / 2
	drawWrappedTextWithShadow(screen, msg, g.gameFont, x, y, ScreenWidth-200, 36, AlertRed)
	scoreMsg := "Final Score: " + itoa(g.score)
	drawWrappedTextWithShadow(screen, scoreMsg, g.gameFont, x, y+60, ScreenWidth-200, 36, VictoryGold)
	// Show rank
	if g.rank != "" {
		rankMsg := "Rank: " + g.rank + " (" + itoa(g.scorePercent) + "% )"
		drawWrappedTextWithShadow(screen, rankMsg, g.gameFont, x, y+120, ScreenWidth-200, 36, SmokeWhite)
	}
}

// Helper to convert int to string (no strconv for minimalism)
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var digits [20]byte
	i := len(digits)
	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		digits[i] = '-'
	}
	return string(digits[i:])
}

// NewGame creates a new Game instance and initializes the font and state.
func NewGame() *Game {
	g := &Game{
		state:                StateNameEntry,
		menuRects:            nil,
		hoveredMenu:          -1,
		quiz:                 game.NewQuiz(),
		unlockedDifficulties: make(map[string]bool),
		answeredSubjects:     make(map[string]bool),
		questionDuration:     10 * time.Second, // 10 second timer
		timerActive:          false,
	}
	g.initFont()
	g.initLogo()
	g.initShips()
	g.initCompass()
	g.initStarStrip()
	g.initFireAnimation() // <-- add this
	return g
}

// Add fire animation image loader
func (g *Game) initFireAnimation() {
	// Player fire
	firePath := "assets/fire_animation_whole.png"
	data, err := os.ReadFile(firePath)
	if err == nil {
		img, err := png.Decode(bytes.NewReader(data))
		if err == nil {
			g.fireImgPlayer = ebiten.NewImageFromImage(img)
		}
	}
	// Enemy fire
	fireEnemyPath := "assets/fire_animation_whole_enemy.png"
	data, err = os.ReadFile(fireEnemyPath)
	if err == nil {
		img, err := png.Decode(bytes.NewReader(data))
		if err == nil {
			g.fireImgEnemy = ebiten.NewImageFromImage(img)
		}
	}
}

// Add a stub for initFont if missing
func (g *Game) initFont() {
	fontBytes, err := os.ReadFile("ui/PressStart2P.ttf")
	if err != nil {
		log.Printf("failed to load font: %v", err)
		return
	}
	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Printf("failed to parse font: %v", err)
		return
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Printf("failed to create font face: %v", err)
		return
	}
	g.gameFont = face
	cf, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err == nil {
		g.confirmFont = cf
	}
}

// Layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// Update implements ebiten.Game interface
func (g *Game) Update() error {
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	mouseJustPressed := mousePressed && !g.prevMousePressed
	fmt.Println("Update called, state:", g.state)
	// Handle SPACE on title screen to go to menu
	if g.state == StateTitle {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.state = StateMenu
		}
		g.prevMousePressed = mousePressed
		return nil
	}
	// Only handle menu input in menu state
	if g.state == StateMenu {
		x, y := ebiten.CursorPosition()
		g.hoveredMenu = -1
		for i, rect := range g.menuRects {
			if x >= rect.Min.X && x < rect.Max.X && y >= rect.Min.Y && y < rect.Max.Y {
				g.hoveredMenu = i
				break
			}
		}
		if g.hoveredMenu != -1 && mouseJustPressed {
			switch g.hoveredMenu {
			case 0: // Play
				g.state = StateSelectDifficulty
			case 1: // How to Play
				g.state = StateHowToPlay
			case 2: // Leaderboard
				g.leaderboardFetched = false // <-- ensure leaderboard always refreshes
				g.state = StateLeaderboard
			case 3: // Exit
				os.Exit(0)
			}
		}
		g.prevMousePressed = mousePressed
		return nil
	}
	// Difficulty selection
	if g.state == StateSelectDifficulty {
		x, y := ebiten.CursorPosition()
		g.hoveredMenu = -1
		difficulties := []string{"Easy", "Medium", "Hard", "Extreme"}
		for i, rect := range g.menuRects {
			locked := false
			if difficulties[i] != "Easy" && !g.unlockedDifficulties[difficulties[i]] {
				locked = true
			}
			if locked {
				continue // skip locked difficulties for hover/click
			}
			if x >= rect.Min.X && x < rect.Max.X && y >= rect.Min.Y && y < rect.Max.Y {
				g.hoveredMenu = i
				break
			}
		}
		if g.hoveredMenu != -1 && mouseJustPressed {
			g.selectedDifficulty = difficulties[g.hoveredMenu]
			g.answeredSubjects = make(map[string]bool) // Reset for new difficulty
			g.state = StateSelectSubject
		}
		g.prevMousePressed = mousePressed
		return nil
	}
	// Subject selection
	if g.state == StateSelectSubject {
		x, y := ebiten.CursorPosition()
		g.hoveredMenu = -1
		if len(g.subjects) == 0 {
			g.subjects = g.quiz.ListSubjects()
		}
		for i, rect := range g.menuRects {
			if x >= rect.Min.X && x < rect.Max.X && y >= rect.Min.Y && y < rect.Max.Y {
				g.hoveredMenu = i
				break
			}
		}
		if g.hoveredMenu != -1 && mouseJustPressed {
			g.selectedSubject = g.subjects[g.hoveredMenu]
			g.startCombatWithSubjectAndDifficulty(g.selectedSubject, g.selectedDifficulty)
			g.state = StatePlaying
		}
		g.prevMousePressed = mousePressed
		return nil
	}
	// When entering StateLeaderboard, fetch leaderboard entries from the database if not already fetched
	if g.state == StateLeaderboard && !g.leaderboardFetched {
		fmt.Println("About to fetch leaderboard entries")
		entries, err := game.GetTopLeaderboard(10)
		fmt.Println("After fetch call")
		if err != nil {
			fmt.Println("Leaderboard DB error:", err)
		} else {
			fmt.Printf("Fetched %d leaderboard entries\n", len(entries))
			for i, e := range entries {
				fmt.Printf("Fetched Entry %d: %+v\n", i, e)
				g.leaderboardEntries = entries
			}
		}
		g.leaderboardFetched = true
	}
	if g.state == StateHowToPlay || g.state == StateLeaderboard {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) || mouseJustPressed {
			g.state = StateMenu
		}
		g.prevMousePressed = mousePressed
		return nil
	}
	// Handle quiz combat
	if g.state == StatePlaying {
		// Handle answer option clicks
		if !g.showFeedback && g.currentQ < len(g.quizQuestions) {
			x, y := ebiten.CursorPosition()
			for i, rect := range g.answerRects {
				if x >= rect.Min.X && x < rect.Max.X && y >= rect.Min.Y && y < rect.Max.Y && mouseJustPressed {
					g.selectedAns = i
					// Process the answer
					q := g.quizQuestions[g.currentQ]
					isCorrect := i == q.Answer
					isBonusQ := g.currentQ == 0 // First question is bonus
					questionsLeft := len(g.quizQuestions) - g.currentQ - 1

					// Use game logic to process answer
					cs := game.ProcessAnswer(
						g.level, g.playerHP, g.playerShields, g.enemyHP, g.score, g.mainQDone, g.bonusActive, g.bonusAnswered,
						g.bonusQIndex, g.currentQ, isBonusQ, isCorrect, questionsLeft,
					)

					// Update game state
					g.playerHP = cs.PlayerHP
					g.playerShields = cs.PlayerShields
					g.enemyHP = cs.EnemyHP
					g.score = cs.Score
					g.mainQDone = cs.MainQDone
					g.bonusActive = cs.BonusActive
					g.bonusAnswered = cs.BonusAnswered
					g.combatOver = cs.CombatOver

					// Show feedback first
					g.showFeedback = true
					g.feedbackTime = time.Now()
					g.feedbackRight = isCorrect
					g.timerActive = false
					break
				}
			}
		}

		// Hide feedback after 2 seconds and move to next question
		if g.showFeedback && time.Since(g.feedbackTime) > 2*time.Second {
			g.showFeedback = false
			g.selectedAns = -1

			// Move to next question
			g.currentQ++

			// Start timer for next question if not at the end
			if g.currentQ < len(g.quizQuestions) {
				g.timerActive = true
				g.questionTimer = time.Now()
			}
		}

		if g.combatOver {
			// Handle unanswered questions based on victory/defeat conditions
			if g.currentQ < len(g.quizQuestions) {
				unanswered := len(g.quizQuestions) - g.currentQ
				if unanswered > 0 {
					// If enemy HP is 0 (victory), award points for unanswered questions
					if g.enemyHP == 0 {
						// Award 10 points per unanswered question (as specified)
						g.score += unanswered * 10
						g.mainQDone += unanswered
					}
					// If player HP is 0 (defeat), no points for unanswered questions
					g.currentQ = len(g.quizQuestions) // Mark all as answered
				}
			}
			if !g.showStarModal {
				// Save score to leaderboard when star modal appears
				if g.userID > 0 && g.score > 0 {
					err := game.InsertLeaderboard(
						g.userID,
						g.score,
						g.mainQDone,
						0,
						float64(g.scorePercent)/100.0, // accuracy
						boolToFloat(g.bonusAnswered),  // bonus success
					)
					if err != nil {
						log.Printf("failed to save leaderboard: %v", err)
					}
				}

				g.showFeedback = false
				g.selectedAns = -1
				g.showStarModal = true
				g.state = StateStarModal
				_, percent := game.CalcRankAndPercent(g.score, g.mainQDone, g.bonusAnswered, g.enemyHP, g.playerHP, g.level)
				var stars float64
				switch {
				case g.playerHP == 0 || percent < 50:
					stars = 0
				case percent == 100:
					stars = 3
				case percent >= 90:
					stars = 2.5
				case percent >= 80:
					stars = 2
				case percent >= 70:
					stars = 1.5
				case percent >= 51:
					stars = 1
				default:
					stars = 0
				}
				g.starCount = stars
				g.starModalStartTime = time.Now()
				g.starAnimationDone = false
				g.starPopIndex = 0
				g.starShineAngle = 0.0
			}
			g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
			return nil
		}
		if g.currentQ >= len(g.quizQuestions) {
			g.combatOver = true
			g.rank, g.scorePercent = game.CalcRankAndPercent(g.score, g.mainQDone, g.bonusAnswered, g.enemyHP, g.playerHP, g.level)
			g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
			return nil
		}
		g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		return nil
	}
	// Handle ESC in Game Over to return to menu
	if g.state == StateGameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			// Save to leaderboard if userID is set, score > 0, and not defeated
			if g.userID > 0 && g.score > 0 && g.rank != "Defeated" {
				err := game.InsertLeaderboard(
					g.userID,
					g.score,
					g.mainQDone,
					0,
					float64(g.scorePercent)/100.0, // accuracy
					boolToFloat(g.bonusAnswered),  // bonus success
				)
				if err != nil {
					log.Printf("failed to save leaderboard: %v", err)
				}
			}
			g.state = StateNameEntry
		}
		g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		return nil
	}
	// Handle name entry
	if g.state == StateNameEntry {
		// Keyboard input
		if !g.nameInputActive {
			g.nameInputActive = true
		}
		for _, r := range ebiten.InputChars() {
			if r == '\n' || r == '\r' {
				continue
			}
			if len(g.enteredName) < 16 {
				g.enteredName += string(r)
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(g.enteredName) > 0 {
			// Remove last rune
			r, size := utf8.DecodeLastRuneInString(g.enteredName)
			if r != utf8.RuneError {
				g.enteredName = g.enteredName[:len(g.enteredName)-size]
			}
		}
		if (ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyKPEnter)) && g.enteredName != "" {
			g.playerName = g.enteredName
			userID, err := game.InsertUser(g.playerName)
			if err != nil {
				log.Printf("failed to save user: %v", err)
			} else {
				g.userID = userID
			}
			g.state = StateMenu
			g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
			return nil
		}
		// Mouse click on confirm button
		w, h := 600, 200
		x := (ScreenWidth - w) / 2
		y := (ScreenHeight - h) / 2
		btnW := 180
		btnH := 48
		btnX := x + w - btnW - 40
		btnY := y + 110 + 60
		mx, my := ebiten.CursorPosition()
		if mx >= btnX && mx < btnX+btnW && my >= btnY && my < btnY+btnH && mouseJustPressed && g.enteredName != "" {
			g.playerName = g.enteredName
			userID, err := game.InsertUser(g.playerName)
			if err != nil {
				log.Printf("failed to save user: %v", err)
			} else {
				g.userID = userID
			}
			g.state = StateMenu
			g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
			return nil
		}
		g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		return nil
	}
	// Handle star modal
	if g.state == StateStarModal {
		x, y := ebiten.CursorPosition()
		if g.starModalResult == 0 {
			for i, rect := range g.continueRects {
				if x >= rect.Min.X && x < rect.Max.X && y >= rect.Min.Y && y < rect.Max.Y && mouseJustPressed {
					g.starModalResult = i + 1 // 1: continue/retry, 2: exit
				}
			}
			if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyKPEnter) {
				g.starModalResult = 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyEscape) {
				g.starModalResult = 2
			}
		} else {
			if g.starModalResult == 1 {
				if g.starCount < 1 {
					// Retry: reset all state and restart the same subject/difficulty
					g.startCombatWithSubjectAndDifficulty(g.selectedSubject, g.selectedDifficulty)
					g.showStarModal = false
					g.starModalResult = 0
					g.state = StatePlaying
					return nil
				} else {
					// Continue: mark subject as answered for this difficulty
					if g.answeredSubjects == nil {
						g.answeredSubjects = make(map[string]bool)
					}
					g.answeredSubjects[g.selectedSubject] = true
					// Check if all subjects are answered for this difficulty
					allSubjects := g.quiz.ListSubjects()
					allAnswered := true
					for _, subj := range allSubjects {
						if !g.answeredSubjects[subj] {
							allAnswered = false
							break
						}
					}
					if allAnswered {
						// Unlock next difficulty if not already unlocked
						difficulties := []string{"Easy", "Medium", "Hard", "Extreme"}
						curIdx := 0
						for i, d := range difficulties {
							if d == g.selectedDifficulty {
								curIdx = i
								break
							}
						}
						if curIdx+1 < len(difficulties) {
							nextDiff := difficulties[curIdx+1]
							g.unlockedDifficulties[nextDiff] = true
						}
						// Reset answeredSubjects for next difficulty
						g.answeredSubjects = make(map[string]bool)
						g.state = StateSelectDifficulty
						g.starModalResult = 0
						return nil
					}
					// Otherwise, go to subject selection for more subjects
					g.state = StateSelectSubject
					g.starModalResult = 0
				}
			} else if g.starModalResult == 2 {
				// Exit: terminate program
				os.Exit(0)
			}
		}
		g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		return nil
	}
	g.prevMousePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	return nil
}

// startCombatWithSubjectAndDifficulty sets up the UI state for a new combat session.
//
// Purpose: This function is called when the player selects a subject and difficulty.
// It prepares the UI state for a new combat round, but delegates all game logic and
// combat state initialization to game.InitCombatState (in the game logic package).
//
// How the logic works:
// - Calls game.InitCombatState(difficulty) to get all initial combat values (HP, shields, etc.).
// - Sets up the quiz questions for the selected subject and difficulty.
// - The UI only manages display and input; all combat rules and state are handled in the game package.
//
// This separation keeps UI and game logic clean, maintainable, and testable.
func (g *Game) startCombatWithSubjectAndDifficulty(subject, difficulty string) {
	// --- UI should NOT handle combat state logic directly. ---
	// All combat state initialization is now handled by game.InitCombatState.
	// This keeps UI and game logic cleanly separated for maintainability.
	combatState := game.InitCombatState(difficulty)
	g.level = combatState.Level
	g.playerMaxHP = combatState.PlayerMaxHP
	g.playerHP = combatState.PlayerHP
	g.playerMaxShields = combatState.PlayerMaxShields
	g.playerShields = combatState.PlayerShields
	g.enemyMaxHP = combatState.EnemyMaxHP
	g.enemyHP = combatState.EnemyHP
	g.bonusActive = combatState.BonusActive
	g.bonusAnswered = combatState.BonusAnswered
	g.bonusQDone = combatState.BonusQDone
	g.mainQDone = combatState.MainQDone
	g.combatOver = combatState.CombatOver
	g.rank = ""
	g.scorePercent = 0

	// Start timer for first question
	g.timerActive = true
	g.questionTimer = time.Now()
	// Filter questions by subject and difficulty
	var filtered []game.Question
	for _, q := range g.quiz.Questions {
		if q.Subject == subject && q.Difficulty == difficulty {
			filtered = append(filtered, q)
		}
	}
	// Shuffle questions
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(filtered), func(i, j int) { filtered[i], filtered[j] = filtered[j], filtered[i] })
	// Pick enough questions for this level
	mainCount := game.GetMainQuestionsCount(combatState.Level)
	if len(filtered) < mainCount+1 {
		mainCount = len(filtered) - 1
	}
	g.quizQuestions = make([]QuizQuestion, 0, mainCount+1)
	if len(filtered) > 0 {
		// First is bonus - initialize with no animation/damage flags
		g.quizQuestions = append(g.quizQuestions, QuizQuestion{
			Question: "[BONUS] " + filtered[0].Text,
			Options:  filtered[0].Choices,
			Answer:   indexOf(filtered[0].Answer, filtered[0].Choices),
		})
		// Ensure no animations or damage for bonus question
		g.pendingAnswer = false
		// The rest are main questions
		for i := 1; i <= mainCount && i < len(filtered); i++ {
			g.quizQuestions = append(g.quizQuestions, QuizQuestion{
				Question: filtered[i].Text,
				Options:  filtered[i].Choices,
				Answer:   indexOf(filtered[i].Answer, filtered[i].Choices),
			})
		}
	}
	g.bonusQIndex = 0
	g.currentQ = 0
	g.selectedAns = -1
	g.showFeedback = false
	g.answerRects = nil
	g.showStarModal = false
}

// indexOf returns the index of ans in choices, or 0 if not found
func indexOf(ans string, choices []string) int {
	for i, c := range choices {
		if c == ans {
			return i
		}
	}
	return 0
}

func (g *Game) drawNameEntry(screen *ebiten.Image) {
	w, h := 600, 200
	x := (ScreenWidth - w) / 2
	y := (ScreenHeight - h) / 2
	// Draw logo at the top center above the box
	logoH := 240
	logoW := 240
	logoY := y - logoH - 36 // reduce offset so logo is closer to the box
	logoX := x + (w-logoW)/2
	if g.logoImg != nil {
		op := &ebiten.DrawImageOptions{}
		sx := float64(logoW) / float64(g.logoImg.Bounds().Dx())
		sy := float64(logoH) / float64(g.logoImg.Bounds().Dy())
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(float64(logoX), float64(logoY))
		screen.DrawImage(g.logoImg, op)
	}
	// Draw BROADSIDE below the logo
	broadsideY := logoY + logoH + 10
	bounds, _ := font.BoundString(g.gameFont, "BROADSIDE")
	width := (bounds.Max.X - bounds.Min.X).Ceil()
	drawWrappedTextWithShadow(screen, "BROADSIDE", g.gameFont, x+(w-width)/2, broadsideY, w-80, 36, VictoryGold)
	// Standard rectangular input box
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), GunmetalGray, true)
	vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 4, VictoryGold, true)
	msg := "Enter your name:"
	drawWrappedTextWithShadow(screen, msg, g.gameFont, x+40, y+60, w-80, 36, VictoryGold)
	// --- Standard input box ---
	inputBoxY := y + 110
	inputBoxH := 48
	inputBoxW := w - 120
	inputBoxX := x + 60
	vector.DrawFilledRect(screen, float32(inputBoxX), float32(inputBoxY), float32(inputBoxW), float32(inputBoxH), SmokeWhite, true)
	vector.StrokeRect(screen, float32(inputBoxX), float32(inputBoxY), float32(inputBoxW), float32(inputBoxH), 2, NavyBlue, true)
	// Draw name inside input box
	nameToShow := g.enteredName
	if g.nameInputActive && (time.Now().UnixNano()/500_000_000)%2 == 0 {
		nameToShow += "|"
	}
	maxNameWidth := inputBoxW - 24
	nameRunes := []rune(nameToShow)
	truncName := ""
	for i := range nameRunes {
		bounds, _ := font.BoundString(g.gameFont, truncName+string(nameRunes[i]))
		width := (bounds.Max.X - bounds.Min.X).Ceil()
		if width > maxNameWidth {
			truncName += "..."
			break
		}
		truncName += string(nameRunes[i])
	}
	nameBounds, _ := font.BoundString(g.gameFont, truncName)
	nameWidth := (nameBounds.Max.X - nameBounds.Min.X).Ceil()
	nameX := inputBoxX + (inputBoxW-nameWidth)/2
	nameY := inputBoxY + inputBoxH/2 + 12
	drawWrappedTextWithShadow(screen, truncName, g.gameFont, nameX, nameY, maxNameWidth, 36, NavyBlue)
	// --- Standard confirm button ---
	btnW := 180
	btnH := 48
	btnX := x + w - btnW - 40
	btnY := y + 110 + 60
	vector.DrawFilledRect(screen, float32(btnX), float32(btnY), float32(btnW), float32(btnH), OceanTeal, true)
	vector.StrokeRect(screen, float32(btnX), float32(btnY), float32(btnW), float32(btnH), 2, VictoryGold, true)
	btnText := "Confirm"
	cf := g.confirmFont
	if cf == nil {
		cf = g.gameFont // fallback
	}
	bounds, _ = font.BoundString(cf, btnText)
	width = (bounds.Max.X - bounds.Min.X).Ceil()
	if width > btnW-24 {
		btnText = btnText[:4] + "..."
	}
	cfMetrics := cf.Metrics()
	cfAscent := cfMetrics.Ascent.Round()
	cfDescent := cfMetrics.Descent.Round()
	cfHeight := cfAscent + cfDescent
	cfY := btnY + (btnH+cfHeight)/2 - cfDescent
	drawWrappedTextWithShadow(screen, btnText, cf, btnX+(btnW-width)/2, cfY, btnW-24, 24, SmokeWhite)
}

func (g *Game) initLogo() {
	logoPath := "ui/logo.png"
	data, err := os.ReadFile(logoPath)
	if err != nil {
		g.logoImg = nil
		return
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		g.logoImg = nil
		return
	}
	g.logoImg = ebiten.NewImageFromImage(img)
}

func (g *Game) initShips() {
	// Load player ship
	playerShipPath := "assets/ship.png"
	data, err := os.ReadFile(playerShipPath)
	if err != nil {
		g.playerShipImg = nil
		log.Printf("failed to load player ship: %v", err)
		for i := 0; i < 4; i++ {
			g.playerShipFrames[i] = nil
		}
	} else {
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			g.playerShipImg = nil
			log.Printf("failed to decode player ship: %v", err)
			for i := 0; i < 4; i++ {
				g.playerShipFrames[i] = nil
			}
		} else {
			g.playerShipImg = ebiten.NewImageFromImage(img)
			imgBounds := g.playerShipImg.Bounds()
			if imgBounds.Dx() < 256 || imgBounds.Dy() < 64 {
				log.Printf("player ship image is too small: got %dx%d, need at least 256x64", imgBounds.Dx(), imgBounds.Dy())
				for i := 0; i < 4; i++ {
					g.playerShipFrames[i] = nil
				}
			} else {
				for i := 0; i < 4; i++ {
					frame := g.playerShipImg.SubImage(image.Rect(i*64, 0, (i+1)*64, 64)).(*ebiten.Image)
					g.playerShipFrames[i] = frame
				}
			}
		}
	}

	// Load enemy ship sprite sheet
	enemyShipPath := "assets/enemy ship.png"
	data, err = os.ReadFile(enemyShipPath)
	if err != nil {
		g.enemyShipImg = nil
		log.Printf("failed to load enemy ship: %v", err)
		for i := 0; i < 4; i++ {
			g.enemyShipFrames[i] = nil
		}
		return
	} else {
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			g.enemyShipImg = nil
			log.Printf("failed to decode enemy ship: %v", err)
			for i := 0; i < 4; i++ {
				g.enemyShipFrames[i] = nil
			}
			return
		} else {
			g.enemyShipImg = ebiten.NewImageFromImage(img)
			imgBounds := g.enemyShipImg.Bounds()
			if imgBounds.Dx() < 256 || imgBounds.Dy() < 64 {
				log.Printf("enemy ship image is too small: got %dx%d, need at least 256x64", imgBounds.Dx(), imgBounds.Dy())
				for i := 0; i < 4; i++ {
					g.enemyShipFrames[i] = nil
				}
				return
			}
			// Slice into 4 frames (64x64 each)
			for i := 0; i < 4; i++ {
				frame := g.enemyShipImg.SubImage(image.Rect(i*64, 0, (i+1)*64, 64)).(*ebiten.Image)
				g.enemyShipFrames[i] = frame
			}
		}
	}
}

func (g *Game) drawShips(screen *ebiten.Image) {
	// Calculate ship positions
	shipY := ScreenHeight - 60 // Position ships closer to the bottom

	// Player ship on the left (frame by frame, safe nil check)
	playerFrameIdx := 3
	if g.playerHP > 0 && g.playerMaxHP > 0 {
		percent := float64(g.playerHP) / float64(g.playerMaxHP)
		switch {
		case percent >= 0.76:
			playerFrameIdx = 0 // 100–76%
		case percent >= 0.51:
			playerFrameIdx = 1 // 75–51%
		case percent > 0:
			playerFrameIdx = 2 // 50–1%
		default:
			playerFrameIdx = 3 // 0%
		}
	} else {
		playerFrameIdx = 3
	}
	playerFrame := g.playerShipFrames[playerFrameIdx]
	if playerFrame != nil {
		playerShipW := 200
		playerShipH := 140
		playerShipX := 60
		playerShipY := shipY - playerShipH

		op := &ebiten.DrawImageOptions{}
		sx := float64(playerShipW) / 64.0
		sy := float64(playerShipH) / 64.0
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(float64(playerShipX), float64(playerShipY))
		screen.DrawImage(playerFrame, op)
	}

	// Enemy ship on the right (use correct frame, safe nil check)
	frameIdx := 0
	if g.enemyHP > 0 && g.enemyMaxHP > 0 {
		percent := float64(g.enemyHP) / float64(g.enemyMaxHP)
		switch {
		case percent >= 0.76:
			frameIdx = 3 // 100–76%
		case percent >= 0.51:
			frameIdx = 2 // 75–51%
		case percent > 0:
			frameIdx = 1 // 50–1%
		default:
			frameIdx = 0 // 0%
		}
	} else {
		frameIdx = 0
	}
	enemyFrame := g.enemyShipFrames[frameIdx]
	if enemyFrame != nil {
		enemyShipW := 200
		enemyShipH := 140
		enemyShipX := ScreenWidth - enemyShipW - 60
		enemyShipY := shipY - enemyShipH

		op := &ebiten.DrawImageOptions{}
		sx := float64(enemyShipW) / 64.0
		sy := float64(enemyShipH) / 64.0
		op.GeoM.Scale(sx, sy)
		op.GeoM.Translate(float64(enemyShipX), float64(enemyShipY))
		screen.DrawImage(enemyFrame, op)
	}

	// Draw fire overlay if active
	if g.showFire && g.fireType != 0 {
		if g.fireType == 1 && g.fireImgPlayer != nil {
			// Player fire: right edge of player ship to left edge of enemy ship
			playerShipW := 200
			playerShipH := 140
			playerShipX := 60
			playerShipY := ScreenHeight - 60 - playerShipH
			enemyShipX := ScreenWidth - 200 - 60
			fireW := enemyShipX - (playerShipX + playerShipW)
			fireH := 80
			fireY := playerShipY + playerShipH/2 - fireH/2
			if fireW > 0 {
				op := &ebiten.DrawImageOptions{}
				sx := float64(fireW) / float64(g.fireImgPlayer.Bounds().Dx())
				sy := float64(fireH) / float64(g.fireImgPlayer.Bounds().Dy())
				op.GeoM.Scale(sx, sy)
				op.GeoM.Translate(float64(playerShipX+playerShipW), float64(fireY))
				screen.DrawImage(g.fireImgPlayer, op)
			}
		}
		if g.fireType == 2 && g.fireImgEnemy != nil {
			// Enemy fire: left edge of enemy ship to right edge of player ship
			playerShipX := 60
			playerShipW := 200
			playerShipY := ScreenHeight - 60 - 140
			enemyShipX := ScreenWidth - 200 - 60
			fireW := enemyShipX - (playerShipX + playerShipW)
			fireH := 80
			fireY := playerShipY + 140/2 - fireH/2
			if fireW > 0 {
				op := &ebiten.DrawImageOptions{}
				sx := float64(fireW) / float64(g.fireImgEnemy.Bounds().Dx())
				sy := float64(fireH) / float64(g.fireImgEnemy.Bounds().Dy())
				op.GeoM.Scale(sx, sy)
				op.GeoM.Translate(float64(playerShipX+playerShipW), float64(fireY))
				screen.DrawImage(g.fireImgEnemy, op)
			}
		}
	}
}

// Add helper to convert bool to float64
func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

func (g *Game) drawStarModal(screen *ebiten.Image) {
	w, h := 700, 380 // Modal size unchanged
	x := (ScreenWidth - w) / 2
	y := (ScreenHeight - h) / 2
	// Border color: red if failed, gold otherwise
	borderColor := VictoryGold
	if g.starCount < 1 {
		borderColor = AlertRed
	}
	vector.StrokeRect(screen, float32(x-4), float32(y-4), float32(w+8), float32(h+8), 4, borderColor, true)
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), GunmetalGray, true)
	vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 2, borderColor, true)

	title := "Test Results"
	drawWrappedTextWithShadow(screen, title, g.confirmFont, x+32, y+48, w-64, 28, VictoryGold)

	// --- Draw star strip ---
	starStripY := y + 70 // Move star further upward
	if g.starStripImg != nil {
		percent := g.scorePercent
		frame := 5 // default: empty
		if g.enemyHP == 0 {
			frame = 0
		} else if percent == 100 {
			frame = 0
		} else if percent >= 90 {
			frame = 1
		} else if percent >= 80 {
			frame = 2
		} else if percent >= 70 {
			frame = 3
		} else if percent >= 1 {
			frame = 4
		} else if percent == 0 {
			frame = 5
		}
		imgW := g.starStripImg.Bounds().Dx()
		imgH := g.starStripImg.Bounds().Dy()
		frameH := imgH / 6
		starX := x + w/2 - imgW/2
		srcRect := image.Rect(0, frame*frameH, imgW, (frame+1)*frameH)
		frameImg := g.starStripImg.SubImage(srcRect).(*ebiten.Image)
		starW := 120
		starH := 120
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(starW)/float64(imgW), float64(starH)/float64(frameH))
		op.GeoM.Translate(float64(starX), float64(starStripY))
		screen.DrawImage(frameImg, op)
	}

	// Add more space below the star strip
	scoreY := starStripY + 120 + 18 // 18px extra space

	// Get total score from all subjects
	totalScore := 0
	if g.userID > 0 {
		if total, err := game.GetTotalScore(g.userID); err == nil {
			totalScore = total
		}
	}

	// Show current subject score (total points from this quiz)
	subjectScoreText := "Score (this subject): " + itoa(g.score) + " pts"
	drawWrappedTextWithShadow(screen, subjectScoreText, g.confirmFont, x+32, scoreY, w-64, 24, VictoryGold)

	// Show total score from all subjects
	totalScoreText := "Total Score (all subjects): " + itoa(totalScore) + " pts"
	totalScoreY := scoreY + 32 // More space between score and total
	drawWrappedTextWithShadow(screen, totalScoreText, g.confirmFont, x+32, totalScoreY, w-64, 24, OceanTeal)

	rank := g.getRankText()
	rankY := totalScoreY + 36 // More space between total and rank
	drawWrappedTextWithShadow(screen, rank, g.confirmFont, x+32, rankY, w-64, 24, g.getRankColor())

	// Feedback
	meaning := g.getMeaningText()
	meaningY := rankY + 36 // More space between rank and feedback
	drawWrappedTextWithShadow(screen, meaning, g.confirmFont, x+32, meaningY, w-64, 20, SmokeWhite)

	// Buttons (keep at bottom, add more space above)
	btnW := 160
	btnH := 44
	// Add extra space between feedback and buttons
	extraBtnSpace := 80                   // Set to 32px
	btnY := meaningY + extraBtnSpace + 20 // 20 is line height for feedback
	// Ensure buttons don't go below modal (clamp if needed)
	if btnY > y+h-52 {
		btnY = y + h - 52
	}
	btnX1 := x + 120
	btnX2 := x + w - btnW - 80
	vector.DrawFilledRect(screen, float32(btnX1), float32(btnY), float32(btnW), float32(btnH), OceanTeal, true)
	vector.StrokeRect(screen, float32(btnX1), float32(btnY), float32(btnW), float32(btnH), 2, VictoryGold, true)
	vector.DrawFilledRect(screen, float32(btnX2), float32(btnY), float32(btnW), float32(btnH), AlertRed, true)
	vector.StrokeRect(screen, float32(btnX2), float32(btnY), float32(btnW), float32(btnH), 2, VictoryGold, true)
	if g.starCount < 1 {
		drawWrappedTextWithShadow(screen, "Retry", g.confirmFont, btnX1+32, btnY+30, btnW-36, 20, SmokeWhite)
	} else {
		drawWrappedTextWithShadow(screen, "Continue", g.confirmFont, btnX1+32, btnY+30, btnW-36, 20, SmokeWhite)
	}
	drawWrappedTextWithShadow(screen, "Exit", g.confirmFont, btnX2+56, btnY+30, btnW-36, 20, SmokeWhite)
	g.continueRects = []image.Rectangle{
		image.Rect(btnX1, btnY, btnX1+btnW, btnY+btnH),
		image.Rect(btnX2, btnY, btnX2+btnW, btnY+btnH),
	}
}

// Add rank/meaning helpers:
func (g *Game) getRankText() string {
	switch {
	case g.starCount == 3:
		return "Rank: S+ (Perfect)"
	case g.starCount == 2.5:
		return "Rank: Gold (Excellent)"
	case g.starCount == 2:
		return "Rank: Silver (Good)"
	case g.starCount == 1.5:
		return "Rank: Bronze (Fair)"
	case g.starCount == 1:
		return "Rank: Pass (Needs Work)"
	default:
		return "Rank: Fail (Try Again)"
	}
}
func (g *Game) getRankColor() color.Color {
	switch {
	case g.starCount == 3:
		return VictoryGold
	case g.starCount >= 2:
		return VictoryGold
	case g.starCount >= 1:
		return OceanTeal
	default:
		return AlertRed
	}
}
func (g *Game) getMeaningText() string {
	switch {
	case g.starCount == 3:
		return "Outstanding performance! You've mastered this subject!"
	case g.starCount == 2.5:
		return "Excellent work! You're very close to perfection!"
	case g.starCount == 2:
		return "Good job! You have a solid understanding."
	case g.starCount == 1.5:
		return "Not bad! A bit more practice will help."
	case g.starCount == 1:
		return "You passed, but there's room for improvement."
	default:
		return "Don't give up! Review the material and try again."
	}
}

// --- Drawing helpers for new visuals ---
// Draw a realistic cloud with shadows
func drawRealisticCloud(screen *ebiten.Image, x, y, scale float64) {
	// Cloud shadow (darker, slightly offset)
	shadowColor := color.RGBA{180, 180, 180, 200}
	shadowOffset := 3.0
	for i := 0; i < 8; i++ {
		angle := float64(i) * math.Pi / 4
		cx := x + math.Cos(angle)*35*scale + shadowOffset
		cy := y + math.Sin(angle)*15*scale + shadowOffset
		radius := 20*scale + float64(i%3)*5*scale
		vector.DrawFilledCircle(screen, float32(cx), float32(cy), float32(radius), shadowColor, true)
	}
	// Main cloud (white)
	cloudColor := SmokeWhite
	for i := 0; i < 8; i++ {
		angle := float64(i) * math.Pi / 4
		cx := x + math.Cos(angle)*35*scale
		cy := y + math.Sin(angle)*15*scale
		radius := 20*scale + float64(i%3)*5*scale
		vector.DrawFilledCircle(screen, float32(cx), float32(cy), float32(radius), cloudColor, true)
	}
	// Cloud highlights (brighter white)
	highlightColor := color.RGBA{255, 255, 255, 255}
	for i := 0; i < 4; i++ {
		angle := float64(i) * math.Pi / 2
		cx := x + math.Cos(angle)*25*scale
		cy := y + math.Sin(angle)*10*scale - 5*scale
		radius := 15*scale + float64(i%2)*3*scale
		vector.DrawFilledCircle(screen, float32(cx), float32(cy), float32(radius), highlightColor, true)
	}
}

// Draw a simple seagull (two wings)
func drawSeagull(screen *ebiten.Image, x, y, scale, wingAngle float64) {
	c := SmokeWhite
	w := 24 * scale
	h := 8 * scale
	// Left wing
	vector.StrokeLine(screen, float32(x), float32(y), float32(x-w*math.Cos(wingAngle)), float32(y-h*math.Sin(wingAngle)), 3, c, true)
	// Right wing
	vector.StrokeLine(screen, float32(x), float32(y), float32(x+w*math.Cos(wingAngle)), float32(y-h*math.Sin(wingAngle)), 3, c, true)
}

// Draw a stylized sun with sharp rays and a glossy highlight
func drawSun(screen *ebiten.Image, x, y, r int) {
	// Draw rays
	rayColor := color.RGBA{255, 180, 40, 180}
	nRays := 24
	outerR := float32(r) * 1.45
	innerR := float32(r) * 1.08
	for i := 0; i < nRays; i++ {
		angle := float64(i) * 2 * math.Pi / float64(nRays)
		nextAngle := float64(i+1) * 2 * math.Pi / float64(nRays)
		x0 := float32(x) + innerR*float32(math.Cos(angle))
		y0 := float32(y) + innerR*float32(math.Sin(angle))
		x1 := float32(x) + outerR*float32(math.Cos((angle+nextAngle)/2))
		y1 := float32(y) + outerR*float32(math.Sin((angle+nextAngle)/2))
		x2 := float32(x) + innerR*float32(math.Cos(nextAngle))
		y2 := float32(y) + innerR*float32(math.Sin(nextAngle))
		drawFilledTriangle(screen, x0, y0, x1, y1, x2, y2, rayColor)
	}
	// Sun core (gradient effect with two circles)
	coreColor := color.RGBA{255, 210, 60, 255}
	vector.DrawFilledCircle(screen, float32(x), float32(y), float32(r), coreColor, true)
	vector.DrawFilledCircle(screen, float32(x), float32(y), float32(r*7/10), color.RGBA{255, 240, 100, 255}, true)
	// Glossy highlight (ellipse simulated with circles)
	highlightColor := color.RGBA{255, 255, 255, 180}
	hw := float32(r) * 0.55
	hh := float32(r) * 0.22
	hx := float32(x) - float32(r)*0.25
	hy := float32(y) - float32(r)*0.32
	for i := 0; i < 10; i++ {
		frac := float32(i) / 10.0
		cx := hx + hw*frac
		cy := hy - (hh * (frac - 0.5) * (frac - 0.5) * 4)
		vector.DrawFilledCircle(screen, cx, cy, float32(r)/13, highlightColor, true)
	}
}

// Draws a filled triangle with the given color
func drawFilledTriangle(screen *ebiten.Image, x0, y0, x1, y1, x2, y2 float32, col color.Color) {
	c := colorToRGBA(col)
	vs := []ebiten.Vertex{
		{DstX: x0, DstY: y0, ColorR: float32(c.R) / 255, ColorG: float32(c.G) / 255, ColorB: float32(c.B) / 255, ColorA: float32(c.A) / 255},
		{DstX: x1, DstY: y1, ColorR: float32(c.R) / 255, ColorG: float32(c.G) / 255, ColorB: float32(c.B) / 255, ColorA: float32(c.A) / 255},
		{DstX: x2, DstY: y2, ColorR: float32(c.R) / 255, ColorG: float32(c.G) / 255, ColorB: float32(c.B) / 255, ColorA: float32(c.A) / 255},
	}
	is := []uint16{0, 1, 2}
	ensureWhiteImg()
	screen.DrawTriangles(vs, is, whiteImg, nil)
}

func (g *Game) initCompass() {
	compassPath := "assets/compass.png"
	data, err := os.ReadFile(compassPath)
	if err != nil {
		g.compassImg = nil
		return
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		g.compassImg = nil
		return
	}
	g.compassImg = ebiten.NewImageFromImage(img)
}

func (g *Game) initStarStrip() {
	starPath := "assets/star.png"
	data, err := os.ReadFile(starPath)
	if err != nil {
		g.starStripImg = nil
		return
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		g.starStripImg = nil
		return
	}
	g.starStripImg = ebiten.NewImageFromImage(img)
}
