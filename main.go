package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <audio1> <audio2>")
		return
	}

	audio1Path := os.Args[1]
	audio2Path := os.Args[2]

	audio1, format1, err := loadAudio(audio1Path)
	if err != nil {
		fmt.Println("Error loading audio1:", err)
		return
	}
	defer audio1.Close()
	fmt.Println("Audio1 loaded successfully")

	audio2, format2, err := loadAudio(audio2Path)
	if err != nil {
		fmt.Println("Error loading audio2:", err)
		return
	}
	defer audio2.Close()
	fmt.Println("Audio2 loaded successfully")

	if format1.SampleRate != format2.SampleRate {
		fmt.Println("Sample rates do not match")
		return
	}

	err = mergeAudiosAndSave(audio1, audio2, format1.SampleRate, 4*time.Second)
	if err != nil {
		fmt.Println("Error merging audios:", err)
		return
	}

	fmt.Println("Audios merged successfully")
}

func loadAudio(path string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, err
	}

	var stream beep.StreamSeekCloser
	var format beep.Format

	if len(path) > 4 && path[len(path)-4:] == ".mp3" {
		stream, format, err = mp3.Decode(f)
	} else if len(path) > 4 && path[len(path)-4:] == ".wav" {
		stream, format, err = wav.Decode(f)
	} else {
		err = fmt.Errorf("unsupported file format")
	}

	if err != nil {
		return nil, beep.Format{}, err
	}

	return stream, format, nil
}

func mergeAudiosAndSave(audio1, audio2 beep.StreamSeekCloser, sampleRate beep.SampleRate, chunkDuration time.Duration) error {
	outputFile, err := os.Create("output.wav")
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	//chunkSize := int(sampleRate.N(chunkDuration))
	//buf1 := make([][2]float64, chunkSize)
	//buf2 := make([][2]float64, chunkSize)

	fmt.Println("Starting merge")

	return nil
}

/* ----------------------------------- */

func cut(streamer1, streamer2 beep.StreamSeekCloser, format1, format2 beep.Format) {
	// Calculate the durations of both audio files
	duration1 := streamer1.Len() / int(format1.SampleRate)
	duration2 := streamer2.Len() / int(format2.SampleRate)

	// Decide which audio file to trim and by how much
	var trimDuration time.Duration
	if duration1 > duration2 {
		trimDuration = time.Duration(duration1-duration2) * time.Second
		streamer1.Seek(int(-trimDuration.Nanoseconds() / int64(format1.SampleRate)))
	} else if duration2 > duration1 {
		trimDuration = time.Duration(duration2-duration1) * time.Second
		streamer2.Seek(int(-trimDuration.Nanoseconds() / int64(format2.SampleRate)))
	}

	// Choose which streamer to save
	var streamerToSave beep.StreamSeekCloser
	var formatToSave beep.Format
	if duration1 > duration2 {
		streamerToSave = streamer1
		formatToSave = format1
	} else {
		streamerToSave = streamer2
		formatToSave = format2
	}

	// Save the streamer as a .wav file
	saveAudio("output.wav", streamerToSave, formatToSave)
}

func saveAudio(path string, streamer beep.Streamer, format beep.Format) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = wav.Encode(f, streamer, format)
	if err != nil {
		return err
	}

	return nil
}
