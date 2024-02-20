package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		os.Exit(1)
	}

	inputFilePath := os.Args[1]
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	inputData, err := io.ReadAll(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	wavFiles := extractWavFiles(inputData)
	for i, wavData := range wavFiles {
		outputFilePath := filepath.Join(filepath.Dir(inputFilePath), fmt.Sprintf("%s_%d.wav", filepath.Base(inputFilePath), i+1))
		err := os.WriteFile(outputFilePath, wavData, 0644)
		if err != nil {
			fmt.Printf("Error writing WAV file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("WAV file extracted: %s\n", outputFilePath)
	}
}

func extractWavFiles(data []byte) [][]byte {
	var wavFiles [][]byte
	var wavEnd int
	riffPattern := []byte("RIFF")
	wavePattern := []byte("WAVE")

	for {
		riffIndex := bytes.Index(data, riffPattern)
		if riffIndex == -1 {
			break
		}

		waveIndex := bytes.Index(data[riffIndex+4:], wavePattern)
		if waveIndex == -1 {
			break
		}

		// Find the next "RIFF" pattern to determine the end of the current WAV file
		nextRiffIndex := bytes.Index(data[riffIndex+4:], riffPattern)
		if nextRiffIndex == -1 {
			// If there is no next "RIFF" pattern, the current WAV file extends to the end of the data
			wavEnd = len(data)
		} else {
			// If there is a next "RIFF" pattern, the current WAV file ends just before it
			wavEnd = riffIndex + nextRiffIndex
		}

		// Extract the WAV file
		wavFiles = append(wavFiles, data[riffIndex:wavEnd])

		// Move to the next chunk after the current WAV file
		data = data[wavEnd:]
	}

	return wavFiles
}
