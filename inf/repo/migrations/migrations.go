package migrations

import (
	"embed"
	"sort"
	"strconv"
)

//go:embed *.sql
var files embed.FS

func ForVersion(version int) [][]byte {
	entries, err := files.ReadDir(`.`)
	if err != nil {
		panic(err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	if version == len(entries) {
		return nil
	}
	if version > len(entries) {
		panic(`invalid migration state, rollback required`)
	}

	var migrations [][]byte
	for _, name := range names[version:] {
		file, err := files.ReadFile(name)
		if err != nil {
			panic(err)
		}

		version++
		file = append(file, []byte(`UPDATE version SET version = `+strconv.Itoa(version)+`;`)...)

		migrations = append(migrations, file)
	}

	return migrations
}
