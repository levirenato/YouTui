// Package search
package search

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Video struct {
	Title         string `json:"title"`
	VideoID       string `json:"videoId"`
	Author        string `json:"author"`
	LengthSeconds int    `json:"lengthSeconds"`
	ViewCount     int64  `json:"viewCount"`
}

type Result struct {
	Title       string
	Author      string
	Duration    string
	URL         string
	Thumbnail   string
	PublishedAt string
	Description string
}

func humanDuration(sec int) string {
	if sec <= 0 {
		return ""
	}
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

type ytdlpItem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Uploader    string  `json:"uploader"`
	Duration    float64 `json:"duration"`
	WebpageURL  string  `json:"webpage_url"`
	URL         string  `json:"url"`
	Thumbnails  []any   `json:"thumbnails"`
	Description string  `json:"description"`
	UploadDate  string  `json:"upload_date"`
}

func SearchVideos(ctx context.Context, q string, limit int) ([]Result, error) {
	t := getTexts()
	if strings.TrimSpace(q) == "" {
		return nil, errors.New(t.EmptyQuery)
	}

	N := limit
	if N <= 0 {
		N = 30
	}
	if N > 50 {
		N = 50
	}

	query := fmt.Sprintf("ytsearch%d:%s", N, q)

	args := []string{
		"-j",
		"--no-warnings",
		"--flat-playlist",
		query,
	}

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("%s", t.YtDlpNotFound)
		}
		return nil, fmt.Errorf("%s: %w", t.YtDlpStartFailed, err)
	}

	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	results := make([]Result, 0, N)
	for sc.Scan() {
		line := sc.Text()
		var it ytdlpItem
		if err := json.Unmarshal([]byte(line), &it); err != nil {
			continue
		}
		if it.ID == "" && it.WebpageURL == "" && it.Title == "" {
			continue
		}

		url := it.WebpageURL
		if url == "" && it.ID != "" {
			url = "https://www.youtube.com/watch?v=" + it.ID
		}

		dur := ""
		if it.Duration > 0 {
			dur = humanDuration(int(it.Duration))
		}

		thumb := ""
		if it.ID != "" {
			thumb = fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", it.ID)
		}

		publishedAt := t.UnknownDate
		if len(it.UploadDate) == 8 {
			year := it.UploadDate[0:4]
			month := it.UploadDate[4:6]
			day := it.UploadDate[6:8]
			publishedAt = fmt.Sprintf("%s/%s/%s", day, month, year)
		}

		description := it.Description
		if description == "" {
			description = t.NoDescription
		}

		results = append(results, Result{
			Title:       it.Title,
			Author:      it.Uploader,
			Duration:    dur,
			URL:         url,
			Thumbnail:   thumb,
			PublishedAt: publishedAt,
			Description: description,
		})
		if limit > 0 && len(results) >= limit {
			break
		}
	}

	_ = stderr

	waitErr := cmd.Wait()
	if waitErr != nil && len(results) == 0 {
		return nil, fmt.Errorf("yt-dlp erro: %w", waitErr)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("%s: %q", t.NoResultsFor, q)
	}
	return results, nil
}

func GetVideoDetails(ctx context.Context, url string) (*Result, error) {
	t := getTexts()

	if url == "" {
		return nil, errors.New(t.EmptyQuery)
	}

	args := []string{
		"-j",
		"--no-warnings",
		"--skip-download",
		url,
	}

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp erro: %w", err)
	}

	var it ytdlpItem
	if err := json.Unmarshal(output, &it); err != nil {
		return nil, fmt.Errorf("parse erro: %w", err)
	}

	dur := ""
	if it.Duration > 0 {
		dur = humanDuration(int(it.Duration))
	}

	thumb := ""
	if it.ID != "" {
		thumb = fmt.Sprintf("https://i.ytimg.com/vi/%s/hqdefault.jpg", it.ID)
	}

	publishedAt := t.UnknownDate
	if len(it.UploadDate) == 8 {
		year := it.UploadDate[0:4]
		month := it.UploadDate[4:6]
		day := it.UploadDate[6:8]
		publishedAt = fmt.Sprintf("%s/%s/%s", day, month, year)
	}

	description := it.Description
	if description == "" {
		description = t.NoDescription
	}

	url = it.WebpageURL
	if url == "" && it.ID != "" {
		url = "https://www.youtube.com/watch?v=" + it.ID
	}

	return &Result{
		Title:       it.Title,
		Author:      it.Uploader,
		Duration:    dur,
		URL:         url,
		Thumbnail:   thumb,
		PublishedAt: publishedAt,
		Description: description,
	}, nil
}
