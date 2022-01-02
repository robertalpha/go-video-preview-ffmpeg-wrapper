package go_video_preview_ffmpeg_wrapper

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
)

type testCase struct {
	inputVideoPath              string
	expectedOutputLengthSeconds int
	parts                       int
	partDurationSeconds         float64
}

func TestCmd_HappyFlows(t *testing.T) {
	testCases := []testCase{
		{
			inputVideoPath:              "testdata/testvid.mp4",
			parts:                       1,
			partDurationSeconds:         1,
			expectedOutputLengthSeconds: 1,
		},
		{
			inputVideoPath:              "testdata/testvid.mp4",
			parts:                       2,
			partDurationSeconds:         1,
			expectedOutputLengthSeconds: 2,
		},
		{
			inputVideoPath:              "testdata/testvid.mp4",
			parts:                       1,
			partDurationSeconds:         2,
			expectedOutputLengthSeconds: 2,
		},
		{
			inputVideoPath:              "testdata/testvid.mp4",
			parts:                       2,
			partDurationSeconds:         .5,
			expectedOutputLengthSeconds: 1,
		},
		{
			inputVideoPath:              "testdata/shortvid.avi", // 2 seconds video
			parts:                       4,
			partDurationSeconds:         4,
			expectedOutputLengthSeconds: 2, // parts * partDuration exceeds video length, so output will be the full length
		},
	}

	for _, test := range testCases {
		inputSeconds, errInput := getDurationInSeconds(test.inputVideoPath)
		if errInput != nil {
			t.Fatalf("Test failed, could not get duration of clip: %v", errInput)
		}

		dir := t.TempDir()
		tmpOutputFile := filepath.Join(dir, "output.webm")

		err := CreatePreviewDefaults(test.inputVideoPath, tmpOutputFile, test.parts, test.partDurationSeconds)
		if err != nil {
			t.Fatalf("Test failed, could not create test clip due to error: %v", err)
		}

		seconds, err := getDurationInSeconds(tmpOutputFile)
		if err != nil {
			t.Fatalf("Test failed, could not get duration of clip: %v", err)
		}
		if seconds != test.expectedOutputLengthSeconds {
			t.Fatalf("Generated clip expected to have duration %vs, but instead is [%vs]", test.expectedOutputLengthSeconds, seconds)
		}

		log.Printf("Source: %v , preview: %v", inputSeconds, seconds)
	}
}

func TestCheckSystem(t *testing.T) {
	err := CheckSystem()
	if err != nil {
		t.Fatalf("System checks did not pass due to error: %v", err)
	}
}

func TestCheckFiles_happy(t *testing.T) {
	err := checkFiles("./main.go", "./main")
	if err != nil {
		t.Fatalf("System checks did not pass due to error: %v", err)
	}
}

func TestCheckFiles_no_source(t *testing.T) {
	err := checkFiles("./some-non-existing-file.txt", "./main")
	if err == nil {
		t.Fatalf("System checks did not throw an error for missing source...")
	}
}

func TestCheckFiles_no_destination(t *testing.T) {
	err := checkFiles("./main.go", "./main.go")
	if err == nil {
		t.Fatalf("System checks did not throw an error for existing destination...")
	}
}

func TestCheckFiles_getSeconds_happy(t *testing.T) {
	seconds := 6132412
	secs, err := matchSeconds(fmt.Sprintf("duration=%v.01243905712597", seconds))
	if err != nil {
		t.Fatalf("could not determine duration")
	}
	if secs != seconds {
		t.Fatalf("wrong duration found")
	}
}

func TestGetParts(t *testing.T) {
	parts := getParts(3, 1, 5)
	if len(parts) != 3 {
		t.Fatalf("parts wrong size")
	}
	if parts[0].startTime != 0.75 {
		t.Fatalf("start time wrong")
	}
	if parts[0].endTime != 1.75 {
		t.Fatalf("end time wrong")
	}
	if parts[1].startTime != 2 {
		t.Fatalf("start time wrong")
	}
	if parts[1].endTime != 3 {
		t.Fatalf("end time wrong")
	}
	if parts[2].startTime != 3.25 {
		t.Fatalf("start time wrong")
	}
	if parts[2].endTime != 4.25 {
		t.Fatalf("end time wrong")
	}
}

func TestGetParts_whole(t *testing.T) {
	parts := getParts(5, 1, 6)
	if len(parts) != 5 {
		t.Fatalf("parts wrong size")
	}
	if parts[0].startTime != 0.5 {
		t.Fatalf("start time wrong")
	}
	if parts[0].endTime != 1.5 {
		t.Fatalf("end time wrong")
	}
	if parts[1].startTime != 1.5 {
		t.Fatalf("start time wrong")
	}
	if parts[1].endTime != 2.5 {
		t.Fatalf("end time wrong")
	}
	if parts[2].startTime != 2.5 {
		t.Fatalf("start time wrong")
	}
	if parts[2].endTime != 3.5 {
		t.Fatalf("end time wrong")
	}
	if parts[3].startTime != 3.5 {
		t.Fatalf("start time wrong")
	}
	if parts[3].endTime != 4.5 {
		t.Fatalf("end time wrong")
	}
	if parts[4].startTime != 4.5 {
		t.Fatalf("start time wrong")
	}
	if parts[4].endTime != 5.5 {
		t.Fatalf("end time wrong")
	}
}
