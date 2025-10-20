package ui

// Language representa um idioma
type Language string

const (
	LanguagePT Language = "pt"
	LanguageEN Language = "en"
)

// Strings contém todas as traduções
type Strings struct {
	// Menus e títulos
	Search          string
	Results         string
	Playlist        string
	Player          string
	Help            string
	Config          string
	
	// Ações
	Play            string
	Pause           string
	Stop            string
	Next            string
	Previous        string
	Add             string
	Remove          string
	
	// Modos
	Audio           string
	Video           string
	Shuffle         string
	RepeatOne       string
	RepeatAll       string
	NoRepeat        string
	
	// Mensagens do Player
	NoTrackPlaying      string
	Playing             string
	PlayingWithoutPlaylist string
	Paused              string
	Stopped             string
	PlaylistFinished    string
	PlaybackFinished    string
	SkippingTo          string
	EnteringPlaylist    string
	AlreadyLastSong     string
	AlreadyFirstSong    string
	NothingPlaying      string
	
	// Mensagens de Playlist
	AddedToPlaylist     string
	RemovedFromPlaylist string
	ItemMoved           string
	PlaylistEmpty       string
	
	// Mensagens de Busca
	Searching           string
	FoundResults        string
	SearchError         string
	NextPage            string
	PrevPage            string
	AlreadyLastPage     string
	AlreadyFirstPage    string
	
	// Mensagens de Erro
	MpvError            string
	StateError          string
	Error               string
	
	// Mensagens de Modo
	ModeChanged         string
	ThemeComingSoon     string
	LanguageChanged     string
	
	// Configuração
	Language        string
	Theme           string
	Close           string
	
	// Comandos
	TypeToSearch    string
	NavigateLists   string
	ShowHelp        string
	Quit            string
	
	// Help
	HelpTitle       string
	HelpNavigation  string
	HelpSearch      string
	HelpResults     string
	HelpPlaylist    string
	HelpPlayer      string
	HelpGlobal      string
	HelpIcons       string
}

// translations guarda todas as traduções
var translations = map[Language]Strings{
	LanguagePT: {
		// Menus e títulos
		Search:          "Busca",
		Results:         "Resultados",
		Playlist:        "Playlist",
		Player:          "Player",
		Help:            "Ajuda",
		Config:          "Configurações",
		
		// Ações
		Play:            "Tocar",
		Pause:           "Pausar",
		Stop:            "Parar",
		Next:            "Próxima",
		Previous:        "Anterior",
		Add:             "Adicionar",
		Remove:          "Remover",
		
		// Modos
		Audio:           "Áudio",
		Video:           "Vídeo",
		Shuffle:         "Aleatório",
		RepeatOne:       "Repetir Uma",
		RepeatAll:       "Repetir Todas",
		NoRepeat:        "Sem Repetição",
		
		// Mensagens do Player
		NoTrackPlaying:         "Nenhuma faixa tocando",
		Playing:                "Tocando",
		PlayingWithoutPlaylist: "Tocando: %s (sem playlist)",
		Paused:                 "Pausado",
		Stopped:                "Parado",
		PlaylistFinished:       "Playlist finalizada",
		PlaybackFinished:       "Reprodução finalizada",
		SkippingTo:             "Pulando para: %d/%d - %s",
		EnteringPlaylist:       "Entrando na playlist...",
		AlreadyLastSong:        "Já está na última música",
		AlreadyFirstSong:       "Já está na primeira música",
		NothingPlaying:         "Nada tocando. Inicie a playlist primeiro.",
		
		// Mensagens de Playlist
		AddedToPlaylist:     "Adicionado: %s",
		RemovedFromPlaylist: "Removido da playlist",
		ItemMoved:           "Item movido",
		PlaylistEmpty:       "Playlist vazia",
		
		// Mensagens de Busca
		Searching:        "Buscando...",
		FoundResults:     "Encontrados %d resultados (Página %d/%d)",
		SearchError:      "Erro: %v",
		NextPage:         "Página %d/%d",
		PrevPage:         "Página %d/%d",
		AlreadyLastPage:  "Já está na última página",
		AlreadyFirstPage: "Já está na primeira página",
		
		// Mensagens de Erro
		MpvError:   "Erro mpv: %v",
		StateError: "Estado: isPlaying=%v socket=%s",
		Error:      "Erro: %v | %s",
		
		// Mensagens de Modo
		ModeChanged:     "Modo: %s",
		ThemeComingSoon: "Tema: Em breve!",
		LanguageChanged: "Idioma alterado para: %s",
		
		// Configuração
		Language:        "Idioma",
		Theme:           "Tema",
		Close:           "Fechar",
		
		// Comandos
		TypeToSearch:    "Digite para buscar",
		NavigateLists:   "Navegar nas listas",
		ShowHelp:        "Mostrar ajuda",
		Quit:            "Sair",
		
		// Help
		HelpTitle:       "ATALHOS DO YOUTUI",
		HelpNavigation:  "NAVEGAÇÃO",
		HelpSearch:      "BUSCA",
		HelpResults:     "RESULTADOS",
		HelpPlaylist:    "PLAYLIST (quando focada)",
		HelpPlayer:      "PLAYER (quando focado)",
		HelpGlobal:      "CONTROLES GLOBAIS",
		HelpIcons:       "ÍCONES DA PLAYLIST",
	},
	
	LanguageEN: {
		// Menus e títulos
		Search:          "Search",
		Results:         "Results",
		Playlist:        "Playlist",
		Player:          "Player",
		Help:            "Help",
		Config:          "Settings",
		
		// Ações
		Play:            "Play",
		Pause:           "Pause",
		Stop:            "Stop",
		Next:            "Next",
		Previous:        "Previous",
		Add:             "Add",
		Remove:          "Remove",
		
		// Modos
		Audio:           "Audio",
		Video:           "Video",
		Shuffle:         "Shuffle",
		RepeatOne:       "Repeat One",
		RepeatAll:       "Repeat All",
		NoRepeat:        "No Repeat",
		
		// Mensagens do Player
		NoTrackPlaying:         "No track playing",
		Playing:                "Playing",
		PlayingWithoutPlaylist: "Playing: %s (no playlist)",
		Paused:                 "Paused",
		Stopped:                "Stopped",
		PlaylistFinished:       "Playlist finished",
		PlaybackFinished:       "Playback finished",
		SkippingTo:             "Skipping to: %d/%d - %s",
		EnteringPlaylist:       "Entering playlist...",
		AlreadyLastSong:        "Already at last song",
		AlreadyFirstSong:       "Already at first song",
		NothingPlaying:         "Nothing playing. Start playlist first.",
		
		// Mensagens de Playlist
		AddedToPlaylist:     "Added: %s",
		RemovedFromPlaylist: "Removed from playlist",
		ItemMoved:           "Item moved",
		PlaylistEmpty:       "Playlist empty",
		
		// Mensagens de Busca
		Searching:        "Searching...",
		FoundResults:     "Found %d results (Page %d/%d)",
		SearchError:      "Error: %v",
		NextPage:         "Page %d/%d",
		PrevPage:         "Page %d/%d",
		AlreadyLastPage:  "Already at last page",
		AlreadyFirstPage: "Already at first page",
		
		// Mensagens de Erro
		MpvError:   "mpv error: %v",
		StateError: "State: isPlaying=%v socket=%s",
		Error:      "Error: %v | %s",
		
		// Mensagens de Modo
		ModeChanged:     "Mode: %s",
		ThemeComingSoon: "Theme: Coming soon!",
		LanguageChanged: "Language changed to: %s",
		
		// Configuração
		Language:        "Language",
		Theme:           "Theme",
		Close:           "Close",
		
		// Comandos
		TypeToSearch:    "Type to search",
		NavigateLists:   "Navigate lists",
		ShowHelp:        "Show help",
		Quit:            "Quit",
		
		// Help
		HelpTitle:       "YOUTUI SHORTCUTS",
		HelpNavigation:  "NAVIGATION",
		HelpSearch:      "SEARCH",
		HelpResults:     "RESULTS",
		HelpPlaylist:    "PLAYLIST (when focused)",
		HelpPlayer:      "PLAYER (when focused)",
		HelpGlobal:      "GLOBAL CONTROLS",
		HelpIcons:       "PLAYLIST ICONS",
	},
}

// GetStrings retorna as traduções para um idioma
func GetStrings(lang Language) Strings {
	if s, ok := translations[lang]; ok {
		return s
	}
	// Fallback para português
	return translations[LanguagePT]
}

// GetLanguageName retorna o nome do idioma para exibição
func GetLanguageName(lang Language) string {
	switch lang {
	case LanguagePT:
		return "Português"
	case LanguageEN:
		return "English"
	default:
		return "Português"
	}
}

// GetAllLanguages retorna todos os idiomas disponíveis
func GetAllLanguages() []Language {
	return []Language{LanguagePT, LanguageEN}
}
