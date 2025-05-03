package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func getVideoAspectRatio(filePath string) (string, error) {
	args := fmt.Sprintf("-v error -select_streams v:0 -show_entries stream=display_aspect_ratio -of default=noprint_wrappers=1 %s", filePath)
	cmd := exec.Command("ffprobe", strings.Split(args, " ")...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get video aspect ratio: %w", err)
	}
	outputStr := string(output)
	if !strings.Contains(outputStr, "display_aspect_ratio=") {
		return "", fmt.Errorf("failed to parse aspect ratio from output: %s", outputStr)
	}

	aspectRatio := strings.TrimSpace(strings.Split(outputStr, "=")[1])
	switch aspectRatio {
	case "":
		return "", fmt.Errorf("aspect ratio is empty")
	case "16:9":
		return "landscape", nil
	case "9:16":
		return "portrait", nil
	default:
		return "other", nil
	}
}

func processVideoForFastStart(filePath string) (string, error) {
	outFilePath := fmt.Sprintf("%s.processing", filePath)

	args := fmt.Sprintf("-i %s -c copy -movflags faststart -f mp4 %s", filePath, outFilePath)
	cmd := exec.Command("ffmpeg", strings.Split(args, " ")...)
	_, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to process video file: %w", err)
	}

	return outFilePath, nil
}
