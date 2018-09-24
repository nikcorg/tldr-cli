package config

// Config is a configuration
type Config struct {
	// Home is the path to TLDR_HOME
	Home string
	// Home is the path to TLDR_HOME/archive
	Archive string
	// Format is a Dateformat string used to construct new archive filenames
	Format string
}
