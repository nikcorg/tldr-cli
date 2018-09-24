package archive

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Load loads an archive file returning an Archive
func Load(file string) Archive {
	data, ioerr := ioutil.ReadFile(file)

	if ioerr != nil {
		panic(ioerr)
	}

	var raw rawArchive

	yamlerr := yaml.Unmarshal(data, &raw)

	if yamlerr != nil {
		panic(yamlerr)
	}

	archive := Archive{
		Title:   raw.Title,
		Entries: []Entry{},
	}

	for _, rawEntry := range raw.Entries {
		relateds := make([]RelatedEntry, 0, len(rawEntry.Related))

		for _, rawRelated := range rawEntry.Related {
			relateds = append(relateds, mapRelatedEntry(rawRelated))
		}

		archive.Entries = append(archive.Entries, Entry{
			Title:   rawEntry.Title,
			Unread:  rawEntry.Unread,
			URL:     rawEntry.URL,
			Source:  rawEntry.Source,
			Related: relateds,
		})
	}

	return archive
}
