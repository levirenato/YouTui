// Package config
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type PlayerState struct {
	SearchTerm        string  `json:"search_term"`
	SearchResults     []Track `json:"search_results"`
	Playlist          []Track `json:"playlist"`
	CurrentTrackIdx   int     `json:"current_track_idx"`
	PlaylistMode      int     `json:"playlist_mode"` // 0=Normal, 1=RepeatOne, 2=RepeatAll, 3=Shuffle
	PlayMode          int     `json:"play_mode"`     // 0=Audio, 1=Video
	SearchScrollIdx   int     `json:"search_scroll_idx"`
	PlaylistScrollIdx int     `json:"playlist_scroll_idx"`
	SearchPage        int     `json:"search_page"`
	LastSaved         string  `json:"last_saved"`
}

type Track struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail"`
	Duration    string `json:"duration"`
	PublishedAt string `json:"published_at"`
	Description string `json:"description"`
}

func GetStatePath() string {
	if xdg := os.Getenv("XDG_STATE_HOME"); xdg != "" {
		return filepath.Join(xdg, "youtui-player", "state.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "state", "youtui-player", "state.json")
}

func LoadState() (*PlayerState, error) {
	statePath := GetStatePath()

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return &PlayerState{
			CurrentTrackIdx:   -1,
			PlaylistMode:      0,
			PlayMode:          0,
			SearchScrollIdx:   0,
			PlaylistScrollIdx: 0,
			SearchResults:     []Track{},
			Playlist:          []Track{},
		}, nil
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, err
	}

	var state PlayerState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func SaveState(state *PlayerState) error {
	statePath := GetStatePath()
	dir := filepath.Dir(statePath)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	state.LastSaved = getCurrentTimestamp()

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, data, 0o644)
}

func ClearState() error {
	statePath := GetStatePath()
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(statePath)
}

func getCurrentTimestamp() string {
	return filepath.Base(time.Now().Format("2006-01-02 15:04:05")) // Placeholder - u can use time.Now().Format()
}
