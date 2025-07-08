package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/RALPH22222/Broadside/game"
	"github.com/RALPH22222/Broadside/ui"
)

func playBackgroundMusic() {
	const sampleRate = 44100
	audioContext := audio.NewContext(sampleRate)
	f, err := os.Open("assets/bgm.wav")
	if err != nil {
		log.Println("failed to open bgm:", err)
		return
	}
	stream, err := wav.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		log.Println("failed to decode bgm:", err)
		return
	}
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	player, err := audioContext.NewPlayer(loop)
	if err != nil {
		log.Println("failed to create audio player:", err)
		return
	}
	player.SetVolume(0.5)
	player.Play()
}

func main() {
	// Initialize DB connection
	game.InitDB("root:@tcp(127.0.0.1:3306)/broadside")

	// Play background music
	go playBackgroundMusic()

	// Create and run the game
	ebiten.SetWindowSize(ui.ScreenWidth, ui.ScreenHeight)
	ebiten.SetWindowTitle("Broadside: Naval Quiz Battle")

	game := ui.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
