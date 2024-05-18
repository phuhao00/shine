package conf

var (
	LenStackBuf = 4096

	// log
	LogLevel string
	LogPath  string
	LogFlag  int

	// console
	ConsolePort   int
	ConsolePrompt string = "shine# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int
)
