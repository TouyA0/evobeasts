package main

import (
	"os"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

const (
	musicAmbiance = "assets/musique_ambiance.mp3"
	musicCombat   = "assets/musique_combat.mp3"
	musicBoss     = "assets/musique_boss.mp3"
)

const baseSampleRate beep.SampleRate = 44100

var (
	audioMu         sync.Mutex
	speakerReady    bool
	currentStreamer beep.StreamSeekCloser
)

func initSpeaker() bool {
	if speakerReady {
		return true
	}
	if err := speaker.Init(baseSampleRate, baseSampleRate.N(time.Second/10)); err != nil {
		return false
	}
	speakerReady = true
	return true
}

func playMusic(path string) {
	audioMu.Lock()
	defer audioMu.Unlock()

	speaker.Clear()
	if currentStreamer != nil {
		currentStreamer.Close()
		currentStreamer = nil
	}

	if !initSpeaker() {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		return
	}

	looped := beep.Loop(-1, streamer)
	resampled := beep.Resample(4, format.SampleRate, baseSampleRate, looped)
	currentStreamer = streamer

	speaker.Play(resampled)
}

func stopMusic() {
	audioMu.Lock()
	defer audioMu.Unlock()

	speaker.Clear()
	if currentStreamer != nil {
		currentStreamer.Close()
		currentStreamer = nil
	}
}
