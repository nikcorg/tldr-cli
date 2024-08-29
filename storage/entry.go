package storage

import (
	"strings"
)

// Entry represents a single resource
type Entry struct {
	RelatedURLs []string `yaml:"related_urls"`
	SourceURL   string   `yaml:"source_url"`
	Tags        []string
	Title       string
	URL         string
	Unread      bool
	deleted     bool
}

// Contains returns `irue` if any string fields contains `needle`
func (e *Entry) Contains(needle string) bool {
	if strings.Contains(strings.ToLower(e.Title), needle) || strings.Contains(e.URL, needle) || strings.Contains(e.SourceURL, needle) {
		return true
	}

	for _, tag := range e.Tags {
		if strings.Contains(tag, needle) {
			return true
		}
	}

	return false
}

// SetDeleted sets the deleted state of the entry
func (e *Entry) SetDeleted(v bool) {
	e.deleted = v
}

// Deleted returns the deleted state of the entry
func (e *Entry) Deleted() bool {
	return e.deleted
}
