package ui

type Language string

const (
	LanguagePT Language = "pt"
	LanguageEN Language = "en"
)

type Strings struct {
	Search   string
	Results  string
	Playlist string
	Player   string
	Help     string
	Config   string

	Play     string
	Pause    string
	Stop     string
	Next     string
	Previous string
	Add      string
	Remove   string

	Audio     string
	Video     string
	Shuffle   string
	RepeatOne string
	RepeatAll string
	NoRepeat  string

	NoTrackPlaying         string
	Playing                string
	PlayingWithoutPlaylist string
	Paused                 string
	Stopped                string
	PlaylistFinished       string
	PlaybackFinished       string
	SkippingTo             string
	EnteringPlaylist       string
	AlreadyLastSong        string
	AlreadyFirstSong       string
	NothingPlaying         string

	AddedToPlaylist     string
	RemovedFromPlaylist string
	ItemMoved           string
	PlaylistEmpty       string

	Searching        string
	FoundResults     string
	SearchError      string
	NextPage         string
	PrevPage         string
	AlreadyLastPage  string
	AlreadyFirstPage string

	MpvError   string
	StateError string
	Error      string

	ModeChanged     string
	ThemeComingSoon string
	ThemeChanged    string
	LanguageChanged string

	Language string
	Theme    string
	Close    string

	TypeToSearch  string
	NavigateLists string
	ShowHelp      string
	Quit          string

	HelpTitle      string
	HelpNavigation string
	HelpSearch     string
	HelpResults    string
	HelpPlaylist   string
	HelpPlayer     string
	HelpGlobal     string
	HelpIcons      string

	NoTitle          string
	Unknown          string
	Channel          string
	Duration         string
	PressEnterToPlay string
	UnknownDate      string
	NoDescription    string
	Page             string

	CmdSearchBar   string
	CmdResultsBar  string
	CmdPlaylistBar string
	CmdPlayerBar   string
	CmdDefaultBar  string

	HelpNavigationText string
	HelpSearchText     string
	HelpResultsText    string
	HelpPlaylistText   string
	HelpPlayerText     string
	HelpGlobalText     string
	HelpIconsText      string

	ConfigText string
}

var translations = map[Language]Strings{
	LanguagePT: {
		Search:   "Busca",
		Results:  "Resultados",
		Playlist: "Playlist",
		Player:   "Player",
		Help:     "Ajuda",
		Config:   "Configurações",

		Play:     "Tocar",
		Pause:    "Pausar",
		Stop:     "Parar",
		Next:     "Próxima",
		Previous: "Anterior",
		Add:      "Adicionar",
		Remove:   "Remover",

		Audio:     "Áudio",
		Video:     "Vídeo",
		Shuffle:   "Aleatório",
		RepeatOne: "Repetir Uma",
		RepeatAll: "Repetir Todas",
		NoRepeat:  "Sem Repetição",

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

		AddedToPlaylist:     "Adicionado: %s",
		RemovedFromPlaylist: "Removido da playlist",
		ItemMoved:           "Item movido",
		PlaylistEmpty:       "Playlist vazia",

		Searching:        "Buscando...",
		FoundResults:     "Encontrados %d resultados (Página %d/%d)",
		SearchError:      "Erro: %v",
		NextPage:         "Página %d/%d",
		PrevPage:         "Página %d/%d",
		AlreadyLastPage:  "Já está na última página",
		AlreadyFirstPage: "Já está na primeira página",

		MpvError:   "Erro mpv: %v",
		StateError: "Estado: isPlaying=%v socket=%s",
		Error:      "Erro: %v | %s",

		ModeChanged:     "Modo: %s",
		ThemeComingSoon: "Tema: Em breve!",
		ThemeChanged:    "Tema alterado para: %s",
		LanguageChanged: "Idioma alterado para: %s",

		Language: "Idioma",
		Theme:    "Tema",
		Close:    "Fechar",

		TypeToSearch:  "Digite para buscar",
		NavigateLists: "Navegar nas listas",
		ShowHelp:      "Mostrar ajuda",
		Quit:          "Sair",

		HelpTitle:      "ATALHOS DO YOUTUI",
		HelpNavigation: "NAVEGAÇÃO",
		HelpSearch:     "BUSCA",
		HelpResults:    "RESULTADOS",
		HelpPlaylist:   "PLAYLIST (quando focada)",
		HelpPlayer:     "PLAYER (quando focado)",
		HelpGlobal:     "CONTROLES GLOBAIS",
		HelpIcons:      "ÍCONES DA PLAYLIST",

		NoTitle:          "Sem título",
		Unknown:          "Desconhecido",
		Channel:          "Canal",
		Duration:         "Duração",
		PressEnterToPlay: "Pressione Enter para tocar",
		UnknownDate:      "Data desconhecida",
		NoDescription:    "Sem descrição disponível",
		Page:             "Página",

		CmdSearchBar:   "Digite para buscar | [#89b4fa]Enter[-] Buscar | [#89b4fa]Tab[-] Próximo | [#f38ba8]Ctrl+Q[-] Sair | [#cba6f7]Ctrl+C[-] Config",
		CmdResultsBar:  "[#89b4fa]↑/↓[-] Nav | [#89b4fa]Enter[-] Play | [#a6e3a1]a[-] Add | [#cba6f7][ ][-] Pág | [#89b4fa]/[-] Buscar | [#f38ba8]Ctrl+Q[-] Sair | [#cba6f7]Ctrl+C[-] Config",
		CmdPlaylistBar: "[#89b4fa]↑/↓[-] Nav | [#89b4fa]Enter[-] Play | [#f38ba8]d[-] Del | [#cba6f7]J/K[-] Move | [#fab387]r[-] Repeat | [#94e2d5]h[-] Shuffle | [#f38ba8]Ctrl+Q[-] Sair | [#cba6f7]Ctrl+C[-] Config",
		CmdPlayerBar:   "[#a6e3a1]Space[-] Pause/Play | [#89dceb]n[-] Next | [#89dceb]p[-] Prev | [#f38ba8]s[-] Stop | [#cba6f7]m[-] Modo | [#f38ba8]Ctrl+Q[-] Sair | [#cba6f7]Ctrl+C[-] Config",
		CmdDefaultBar:  "[#89b4fa]Tab[-] Navegar entre painéis | [#f38ba8]Ctrl+Q[-] Sair | [#cba6f7]Ctrl+C[-] Config",

		HelpNavigationText: "  Tab       Alternar entre painéis (Busca → Resultados → Playlist → Player)\n  /         Focar na busca\n  ↑/↓       Navegar nas listas\n  ?         Mostrar esta ajuda",
		HelpSearchText:     "  Digite    Texto para buscar\n  Enter     Executar busca",
		HelpResultsText:    "  Enter     Tocar faixa diretamente (sem playlist)\n  a         Adicionar à playlist\n  [ ]       Navegar entre páginas (anterior/próxima)",
		HelpPlaylistText:   "  Enter     Tocar faixa da playlist\n  Space     Tocar playlist do início\n  d         Remover item\n  J         Mover item para baixo\n  K         Mover item para cima\n  r         Ciclar repetição (󰑗 → 󰑘 → 󰑖 → 󰑗)\n  h         Toggle shuffle ()",
		HelpPlayerText:     "  Space     Pause/Play\n  s         Stop\n  n         Próxima música\n  p         Música anterior",
		HelpGlobalText:     "  m         Alternar áudio/vídeo\n  Ctrl+Q    Sair da aplicação\n  Ctrl+C    Configurações",
		HelpIconsText:      "  󰑗 Sem Repetição  󰑘 Repetir Uma  󰑖 Repetir Todas   Aleatório",

		ConfigText: "⚙️  CONFIGURAÇÕES\n\nEscolha uma opção abaixo para configurar o YouTui.\nUse as setas ←/→ para navegar e Enter para selecionar.",
	},

	LanguageEN: {
		Search:   "Search",
		Results:  "Results",
		Playlist: "Playlist",
		Player:   "Player",
		Help:     "Help",
		Config:   "Settings",

		Play:     "Play",
		Pause:    "Pause",
		Stop:     "Stop",
		Next:     "Next",
		Previous: "Previous",
		Add:      "Add",
		Remove:   "Remove",

		Audio:     "Audio",
		Video:     "Video",
		Shuffle:   "Shuffle",
		RepeatOne: "Repeat One",
		RepeatAll: "Repeat All",
		NoRepeat:  "No Repeat",

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

		AddedToPlaylist:     "Added: %s",
		RemovedFromPlaylist: "Removed from playlist",
		ItemMoved:           "Item moved",
		PlaylistEmpty:       "Playlist empty",

		Searching:        "Searching...",
		FoundResults:     "Found %d results (Page %d/%d)",
		SearchError:      "Error: %v",
		NextPage:         "Page %d/%d",
		PrevPage:         "Page %d/%d",
		AlreadyLastPage:  "Already at last page",
		AlreadyFirstPage: "Already at first page",

		MpvError:   "mpv error: %v",
		StateError: "State: isPlaying=%v socket=%s",
		Error:      "Error: %v | %s",

		ModeChanged:     "Mode: %s",
		ThemeComingSoon: "Theme: Coming soon!",
		ThemeChanged:    "Theme changed to: %s",
		LanguageChanged: "Language changed to: %s",

		Language: "Language",
		Theme:    "Theme",
		Close:    "Close",

		TypeToSearch:  "Type to search",
		NavigateLists: "Navigate lists",
		ShowHelp:      "Show help",
		Quit:          "Quit",

		HelpTitle:      "YOUTUI SHORTCUTS",
		HelpNavigation: "NAVIGATION",
		HelpSearch:     "SEARCH",
		HelpResults:    "RESULTS",
		HelpPlaylist:   "PLAYLIST (when focused)",
		HelpPlayer:     "PLAYER (when focused)",
		HelpGlobal:     "GLOBAL CONTROLS",
		HelpIcons:      "PLAYLIST ICONS",

		NoTitle:          "No title",
		Unknown:          "Unknown",
		Channel:          "Channel",
		Duration:         "Duration",
		PressEnterToPlay: "Press Enter to play",
		UnknownDate:      "Unknown date",
		NoDescription:    "No description available",
		Page:             "Page",

		CmdSearchBar:   "Type to search | [#89b4fa]Enter[-] Search | [#89b4fa]Tab[-] Next | [#f38ba8]Ctrl+Q[-] Quit | [#cba6f7]Ctrl+C[-] Config",
		CmdResultsBar:  "[#89b4fa]↑/↓[-] Nav | [#89b4fa]Enter[-] Play | [#a6e3a1]a[-] Add | [#cba6f7][ ][-] Page | [#89b4fa]/[-] Search | [#f38ba8]Ctrl+Q[-] Quit | [#cba6f7]Ctrl+C[-] Config",
		CmdPlaylistBar: "[#89b4fa]↑/↓[-] Nav | [#89b4fa]Enter[-] Play | [#f38ba8]d[-] Del | [#cba6f7]J/K[-] Move | [#fab387]r[-] Repeat | [#94e2d5]h[-] Shuffle | [#f38ba8]Ctrl+Q[-] Quit | [#cba6f7]Ctrl+C[-] Config",
		CmdPlayerBar:   "[#a6e3a1]Space[-] Pause/Play | [#89dceb]n[-] Next | [#89dceb]p[-] Prev | [#f38ba8]s[-] Stop | [#cba6f7]m[-] Mode | [#f38ba8]Ctrl+Q[-] Quit | [#cba6f7]Ctrl+C[-] Config",
		CmdDefaultBar:  "[#89b4fa]Tab[-] Navigate panels | [#f38ba8]Ctrl+Q[-] Quit | [#cba6f7]Ctrl+C[-] Config",

		HelpNavigationText: "  Tab       Switch panels (Search → Results → Playlist → Player)\n  /         Focus search\n  ↑/↓       Navigate lists\n  ?         Show this help",
		HelpSearchText:     "  Type      Text to search\n  Enter     Execute search",
		HelpResultsText:    "  Enter     Play track directly (no playlist)\n  a         Add to playlist\n  [ ]       Navigate pages (previous/next)",
		HelpPlaylistText:   "  Enter     Play track from playlist\n  Space     Play playlist from start\n  d         Remove item\n  J         Move item down\n  K         Move item up\n  r         Cycle repeat (󰑗 → 󰑘 → 󰑖 → 󰑗)\n  h         Toggle shuffle ()",
		HelpPlayerText:     "  Space     Pause/Play\n  s         Stop\n  n         Next song\n  p         Previous song",
		HelpGlobalText:     "  m         Toggle audio/video\n  Ctrl+Q    Quit application\n  Ctrl+C    Settings",
		HelpIconsText:      "  󰑗 No Repeat  󰑘 Repeat One  󰑖 Repeat All   Shuffle",

		ConfigText: "⚙️  SETTINGS\n\nChoose an option below to configure YouTui.\nUse ←/→ arrows to navigate and Enter to select.",
	},
}

func GetStrings(lang Language) Strings {
	if s, ok := translations[lang]; ok {
		return s
	}
	return translations[LanguagePT]
}

func GetLanguageName(lang Language) string {
	switch lang {
	case LanguagePT:
		return "Português"
	case LanguageEN:
		return "English"
	default:
		return "English"
	}
}

func GetAllLanguages() []Language {
	return []Language{LanguagePT, LanguageEN}
}
