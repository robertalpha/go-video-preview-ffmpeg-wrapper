package go_video_preview_ffmpeg_wrapper

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

func CreatePreviewDefaults(sourcePath string, destinationPath string, parts int, partDurationSeconds float64) error {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return getError("could not find ffmpeg, please make sure it's installed", "")
	}

	secs, err := getDurationInSeconds(sourcePath)
	if err != nil {
		return err
	}

	parameters := ""
	if float64(parts+2)*partDurationSeconds > float64(secs) {
		// preview command for full length of the input video
		parameters = formatFfmpegCommandParametersWithout(sourcePath, destinationPath)
	} else {
		parameters = formatFfmpegCommandParametersWithFilter(sourcePath, destinationPath, parts, partDurationSeconds, secs)

	}
	workaround := fmt.Sprintf("%s %s", ffmpegPath, parameters)
	cmd := exec.Command("sh", "-c", workaround)
	if err := cmd.Run(); err != nil {
		return getError("Error: %v", err)
	}

	return nil
}

func getDurationInSeconds(sourcePath string) (int, error) {
	probePath, err := exec.LookPath("ffprobe")
	if err != nil {
		return 0, getError("could not find ffprobe, please make sure it's installed", "")
	}

	cmd := exec.Command(probePath, "-i", sourcePath, "-show_entries", "format=duration", "-v", "quiet")
	output, err := cmd.Output()
	if err != nil {
		return 0, getError("Error: %v", err)
	}
	if len(output) == 0 {
		return 0, getError("could not determine clip length", "")
	}
	lines := getLinesFromBytes(output)
	if len(lines) != 3 {
		return 0, getError("could not determine clip length", "")
	}

	length := lines[1]

	secs, err := matchSeconds(length)
	if err != nil {
		return 0, getError("Error: %v", err)
	}
	return secs, nil
}

func matchSeconds(duration string) (int, error) {
	r := regexp.MustCompile(`(?P<Pre>duration)=(?P<Seconds>\d+)\.(?P<Rest>\d+)`)
	submatch := r.FindStringSubmatch(duration)
	if len(submatch) != 4 {
		return 0, fmt.Errorf("could not determine clip length")
	}
	seconds, err := strconv.Atoi(submatch[2])
	if err != nil {
		return 0, fmt.Errorf("could not determine clip length for [%s]", duration)
	}
	return seconds, nil
}

func checkFiles(sourcePath string, destinationPath string) error {

	if _, err := os.Stat(sourcePath); errors.Is(err, os.ErrNotExist) {
		return getError("source file does not exist", destinationPath)
	}

	dstDir := path.Dir(destinationPath)
	if stat, err := os.Stat(dstDir); err == nil && !stat.IsDir() {
		return getError("destination directory does not exist", dstDir)
	}
	if _, err := os.Stat(destinationPath); err == nil {
		return getError("destination file already exists", destinationPath)
	}

	return nil
}

func getError(decription string, param interface{}) error {
	return fmt.Errorf("[go-video-preview-ffmpeg-wrapper] Error: %s [%v]", decription, param)
}

func CheckSystem() error {
	if err := checkCmdExists("ffmpeg"); err != nil {
		return getError("ffmpeg not found on system.. please install ", "ffprobe")
	}
	if err := checkCmdExists("ffprobe"); err != nil {
		return getError("ffprobe not found on system.. please install ", "ffmpeg/ffprobe")
	}
	return nil
}

func checkCmdExists(cmd string) error {
	_, err := exec.LookPath(cmd)
	return err
}

func getLinesFromBytes(input []byte) []string {
	bytesReader := bytes.NewReader(input)
	bufReader := bufio.NewReader(bytesReader)
	var output []string
	for {
		line, _, err := bufReader.ReadLine()

		if err == io.EOF {
			break
		}
		output = append(output, string(line))
	}
	return output
}

func formatFfmpegCommandParametersWithFilter(input string, output string, parts int, partSeconds float64, totalLenght int) string {
	return fmt.Sprintf("-i %s %s %s %s", input, formatFilter(parts, partSeconds, totalLenght), getOtherParams(), output)
}
func formatFfmpegCommandParametersWithout(input string, output string) string {
	return fmt.Sprintf("-i %s -vf scale=320:-2 %s %s", input, getOtherParams(), output)
}

func getOtherParams() string {
	return "-c:v libvpx -qmin 0 -qmax 25 -crf 23 -an -threads 0 -hide_banner -loglevel error -y"
}

func formatFilter(partsCount int, partSeconds float64, totalLength int) string {
	filterheader := fmt.Sprintf("\"[0]split=%v", partsCount)
	filterfooter := ""
	var lines string
	parts := getParts(partsCount, partSeconds, float64(totalLength))
	for idx, part := range parts {
		num := idx + 1
		filterheader = fmt.Sprintf("%s[v%v]", filterheader, num)
		start := part.startTime
		end := part.endTime
		lines = lines + fmt.Sprintf("[v%v]trim=%v:%v,setpts=PTS-STARTPTS[v%vt]; ", num, start, end, num)
		filterfooter = filterfooter + fmt.Sprintf("[v%vt]", num)
	}
	filterheader = filterheader + fmt.Sprintf("; %s", lines)
	filterfooter = fmt.Sprintf("%vconcat=n=%v:v=1:a=0[vc]; ", filterfooter, partsCount)
	filtercloser := "[vc]scale=320:-2[v]\" "

	return fmt.Sprintf("-filter_complex %s%s%s -map \"[v]\" ", filterheader, filterfooter, filtercloser)
}

type videoPart struct {
	startTime float64
	endTime   float64
}

func getParts(parts int, partLength float64, totalLength float64) (output []videoPart) {
	div := totalLength / float64(parts+1)
	for a := 1; a <= parts; a++ {
		start := (div * float64(a)) - (partLength / 2)
		end := (div * float64(a)) + (partLength / 2)
		output = append(output, videoPart{
			startTime: start,
			endTime:   end,
		})
	}
	return
}
