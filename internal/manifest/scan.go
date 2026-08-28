package manifest

import (
	"errors"
	"os"
	"path/filepath"
)

func GenerateHashMap(root string) (*Manifest, error) {
	m := New()

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		relPath, hash, err := hashMd5File(path, d, root, walkErr)
		if err != nil {
			if errors.Is(err, errIsDir) {
				return nil
			}
			return err
		}
		m.Add(relPath, hash)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}
