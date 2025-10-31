// Package search
package search

import "sync/atomic"

type Texts struct {
	EmptyQuery       string
	YtDlpNotFound    string
	YtDlpStartFailed string
	YtDlpError       string
	NoResultsFor     string

	UnknownDate   string
	NoDescription string
}

var texts atomic.Value

func init() {
	texts.Store(Texts{
		EmptyQuery:       "Consulta vazia",
		YtDlpNotFound:    "yt-dlp não encontrado no PATH. Instale com 'pipx install yt-dlp' ou 'pip install --user yt-dlp'",
		YtDlpStartFailed: "Falha ao iniciar yt-dlp",
		YtDlpError:       "Erro do yt-dlp",
		NoResultsFor:     "Nenhum resultado para: %q",

		UnknownDate:   "Data desconhecida",
		NoDescription: "Sem descrição disponível",
	})
}

func setTexts(t Texts) { texts.Store(t) }
func getTexts() Texts  { return texts.Load().(Texts) }

func SetTexts(t Texts) { setTexts(t) }
