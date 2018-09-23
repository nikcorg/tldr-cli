package types

import (
	"reflect"
)

// Archive is a collection of entries
type Archive struct {
	Title   string  `yaml:"title"`
	Entries []Entry `yaml:"entries"`
}

// Entry is a single entry in an archive
type Entry struct {
	Title   string        `yaml:"title"`
	URL     string        `yaml:"url"`
	Unread  bool          `yaml:"unread"`
	Source  string        `yaml:"source,omitempty"`
	Related []interface{} `yaml:"related,omitempty"`
}

// RelatedEntry is a placeholder type for reading archives
type RelatedEntry interface {
	IsSimple() bool
	IsConcise() bool
}

// ConciseEntry is a concise entry
type ConciseEntry struct {
	Title string
	URL   string
}

// IsSimple beep boop
func (x ConciseEntry) IsSimple() bool {
	return false
}

// IsConcise beep boop
func (x ConciseEntry) IsConcise() bool {
	return true
}

// SimpleEntry is a simple entry
type SimpleEntry struct {
	URL string
}

// IsSimple beep boop
func (x SimpleEntry) IsSimple() bool {
	return true
}

// IsConcise beep boop
func (x SimpleEntry) IsConcise() bool {
	return false
}

// MapRelatedEntry maps from interface{} to ConciseEntry and SimpleEntry
func MapRelatedEntry(related interface{}) RelatedEntry {
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
