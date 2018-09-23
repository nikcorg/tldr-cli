package archive

import (
	"reflect"
)

type rawArchive struct {
	Entries []rawEntry `yaml:"entries"`
	Title   string     `yaml:"title"`
}

type rawEntry struct {
	Related []interface{} `yaml:"related,omitempty"`
	Source  string        `yaml:"source,omitempty"`
	Title   string        `yaml:"title"`
	URL     string        `yaml:"url"`
	Unread  bool          `yaml:"unread"`
}

// Archive is a collection of entries
type Archive struct {
	Entries []Entry `yaml:"entries"`
	Title   string  `yaml:"title"`
}

// Entry is a single entry in an archive
type Entry struct {
	Related []RelatedEntry `yaml:"related,omitempty"`
	Source  string         `yaml:"source,omitempty"`
	Title   string         `yaml:"title"`
	URL     string         `yaml:"url"`
	Unread  bool           `yaml:"unread"`
}

// RelatedEntry is a placeholder type for reading archives
type RelatedEntry interface {
	IsConcise() bool
	IsSimple() bool
}

// ConciseEntry is a concise entry
type ConciseEntry struct {
	Title string
	URL   string
}

// IsSimple type discriminator for ConciseEntry
func (x ConciseEntry) IsSimple() bool {
	return false
}

// IsConcise type discriminator for ConciseEntry
func (x ConciseEntry) IsConcise() bool {
	return true
}

// SimpleEntry is a plain URL entry
type SimpleEntry struct {
	URL string
}

// IsSimple type discriminator for SimpleEntry
func (x SimpleEntry) IsSimple() bool {
	return true
}

// IsConcise type discriminator for SimpleEntry
func (x SimpleEntry) IsConcise() bool {
	return false
}

func mapRelatedEntry(related interface{}) RelatedEntry {
	if reflect.TypeOf(related).Kind() == reflect.Map {
		asMap := related.(map[interface{}]interface{})
		return ConciseEntry{
			Title: asMap["title"].(string),
			URL:   asMap["url"].(string),
		}
	}

	return SimpleEntry{
		URL: related.(string),
	}
}
