package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	PlaylistFile = "playlist.txt"
)

type Config struct {
	Ffmpeg  string `json:"ffmpeg"`
	FfProbe string `json:"ffprobe"`
	Stage   struct {
		Image  string `json:"image"`
		Width  string `json:"width"`
		Height string `json:"height"`
	}
	Input struct {
		Rectangle [4]string `json:"rectangle"`
		Episodes  []string  `json:"episodes"`
		Title     struct {
			Font   string `json:"font"`
			Prefix string `json:"prefix"`
			X      string `json:"x"`
			Y      string `json:"y"`
			Color  string `json:"color"`
			Size   string `json:"size"`
		} `json:"title"`
	} `json:"input"`
	Output struct {
		StreamURL string `json:"stream_url"`
	}
}

type Video struct {
	Name     string
	Duration float64
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	_ = os.Remove(PlaylistFile)

	configBytes, _ := ioutil.ReadFile("./config.json")
	var config Config
	_ = json.Unmarshal(configBytes, &config)

	videos := make([]*Video, 0)
	playlistLines := make([]string, 0)
	totalDuration := float64(0)
	// 1. Get Video metadata
	for _, v := range config.Input.Episodes {
		absoluePath, _ := filepath.Abs(v)
		playlistLines = append(playlistLines, "file "+absoluePath+"")
		out, err := exec.Command(
			config.FfProbe,
			"-v",
			"error",
			"-show_entries",
			"format=duration",
			"-of",
			"default=noprint_wrappers=1:nokey=1",
			"'"+v+"'").CombinedOutput()
		if err != nil {
			log.Fatalln(v, string(out))
		}
		duration, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
		segs := strings.Split(v, "/")
		filename := segs[len(segs)-1]
		for _, extesion := range []string{".flv", ".mp4"} {
			filename = strings.Replace(filename, extesion, "", 1)
		}
		videos = append(videos, &Video{
			Name:     filename,
			Duration: duration,
		})
		totalDuration += duration
	}
	// 2. Compose playlist and drawtext params
	drawtexts := make([]string, 0)
	startOffset := float64(0)
	loopOneTime := fmt.Sprintf("%.6f", totalDuration)
	for _, v := range videos {
		from := fmt.Sprintf("%.6f", startOffset)
		to := fmt.Sprintf("%.6f", startOffset+v.Duration)
		vDrawParam := fmt.Sprintf(
			"drawtext=fontfile='%s': text='%s%s':x=%s:y=%s:fontcolor=%s:fontsize=%s:box=1:boxcolor=0x00000099:enable='between(mod(t,%s),%s,%s)'",
			config.Input.Title.Font,
			config.Input.Title.Prefix,
			v.Name,
			config.Input.Title.X,
			config.Input.Title.Y,
			config.Input.Title.Color,
			config.Input.Title.Size,
			loopOneTime,
			from,
			to,
		)
		drawtexts = append(drawtexts, vDrawParam)
		startOffset += v.Duration
	}
	// 3. Write playlist file
	err := os.WriteFile(PlaylistFile, []byte(strings.Join(playlistLines, "\n")), 0660)
	if err != nil {
		log.Fatalln("Failed to save playlist file:", err)
	}
	// 4. Start streaming
	cmd := exec.Command(
		config.Ffmpeg,
		"-hide_banner",
		"-stream_loop",
		"-1",
		"-f",
		"concat",
		"-re",
		"-safe",
		"0",
		"-i",
		PlaylistFile,
		"-f",
		"flv",
		"-",
	)
	pushCmd := exec.Command(
		config.Ffmpeg,
		"-hide_banner",
		"-loglevel",
		"repeat+quiet",
		"-re",
		"-i",
		config.Stage.Image,
		"-i",
		"pipe:0",
		"-filter_complex",
		fmt.Sprintf(
			`[0:v] pad=%s:%s[bg];[1:v]scale=%s:%s[temp1];[bg][temp1] overlay=%s:%s[temp2];[temp2]%s`,
			config.Stage.Width,
			config.Stage.Height,
			config.Input.Rectangle[2],
			config.Input.Rectangle[3],
			config.Input.Rectangle[0],
			config.Input.Rectangle[1],
			strings.Join(drawtexts, ","),
		),
		"-f",
		"flv",
		"-c:a",
		"copy",
		"-c:v",
		"libx264",
		config.Output.StreamURL,
	)
	pushCmd.Stdin, _ = cmd.StdoutPipe()
	pushCmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	pushCmd.Stderr = os.Stderr
	Must(pushCmd.Start())
	Must(cmd.Run())
	Must(pushCmd.Wait())
}
