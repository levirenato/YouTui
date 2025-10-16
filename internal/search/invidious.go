package search

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Video struct {
	Title         string `json:"title"`
	VideoId       string `json:"videoId"`
	Author        string `json:"author"`
	LengthSeconds int    `json:"lengthSeconds"`
	ViewCount     int64  `json:"viewCount"`
}

type Result struct {
	Title    string
	Author   string
	Duration string
	URL      string
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

// payload que o yt-dlp emite com -j (--dump-json)
type ytdlpItem struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Uploader   string  `json:"uploader"`
	Duration   float64 `json:"duration"`    // pode vir ausente/0 se usar --flat-playlist
	WebpageURL string  `json:"webpage_url"` // normalmente presente
	URL        string  `json:"url"`         // nem sempre útil aqui
	Thumbnails []any   `json:"thumbnails"`  // ignorado
	// … há muitos outros campos; mantemos o necessário
}

// SearchVideos: executa `yt-dlp -j ytsearchN:<q>` e parseia o NDJSON.
// Evita depender de instâncias Invidious e funciona offline (exceto rede pro YouTube).
func SearchVideos(ctx context.Context, q string, limit int) ([]Result, error) {
	if strings.TrimSpace(q) == "" {
		return nil, errors.New("consulta vazia")
	}

	// saneia limite
	N := limit
	if N <= 0 {
		N = 30
	}
	if N > 50 {
		N = 50 // manter rápido e gentil
	}

	// Monta a “URL de busca” especial do yt-dlp
	query := fmt.Sprintf("ytsearch%d:%s", N, q)

	// Preferimos --flat-playlist para ser rápido (1 request) — porém muitas vezes não traz duração.
	// Se você quiser duração precisa, remova --flat-playlist (mais lento, chama por item).
	args := []string{
		"-j", // JSON por item (NDJSON)
		"--no-warnings",
		"--flat-playlist", // rápido; pode omitir pra tentar obter duração
		query,
	}

	// Executa como subprocesso com timeout via ctx
	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	// Importante: não redirecionar Stderr pro terminal do TUI; mantemos silencioso.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	// Start
	if err := cmd.Start(); err != nil {
		// Dica comum: yt-dlp ausente
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("yt-dlp não encontrado no PATH. Instale com 'pipx install yt-dlp' ou 'pip install --user yt-dlp'")
		}
		return nil, fmt.Errorf("falha ao iniciar yt-dlp: %w", err)
	}

	// Leitor de linhas NDJSON
	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	results := make([]Result, 0, N)
	for sc.Scan() {
		line := sc.Text()
		var it ytdlpItem
		if err := json.Unmarshal([]byte(line), &it); err != nil {
			// tenta ignorar linhas quebradas e segue
			continue
		}
		if it.ID == "" && it.WebpageURL == "" && it.Title == "" {
			continue
		}

		// Monta URL final
		url := it.WebpageURL
		if url == "" && it.ID != "" {
			url = "https://www.youtube.com/watch?v=" + it.ID
		}

		// Duração: pode vir 0 se --flat-playlist (aceitamos vazio)
		dur := ""
		if it.Duration > 0 {
			dur = humanDuration(int(it.Duration))
		}

		results = append(results, Result{
			Title:    it.Title,
			Author:   it.Uploader,
			Duration: dur,
			URL:      url,
		})
		if limit > 0 && len(results) >= limit {
			break
		}
	}

	// Consome stderr (se houver) para evitar deadlock e enriquecer mensagens
	_ = stderr

	// Aguarda término
	waitErr := cmd.Wait()
	if waitErr != nil && len(results) == 0 {
		// Se falhou e não retornou nada, repassa erro
		return nil, fmt.Errorf("yt-dlp erro: %w", waitErr)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("nenhum resultado para: %q", q)
	}
	return results, nil
}

// (opcional) conversor seguro quando duration vier como string em algum cenário incomum
func atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
