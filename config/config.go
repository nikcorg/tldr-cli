package config

// Config is a configuration
type Config struct {
	Home        string // Home is the path to TLDR_HOME
	Archive     string // Home is the path to TLDR_HOME/archive
	Format      string // Format is a Dateformat string used to construct new archive filenames
	TitleFormat string // TitleFormat is a DateFormat string for daily archive titles
}
