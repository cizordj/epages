package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func GenerateFileHashMap(root string) (map[string]string, error) {
	hashMap := make(map[string]string)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		relPath = filepath.ToSlash(relPath)

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			return err
		}

		hashMap[hex.EncodeToString(hash.Sum(nil))] = relPath

		return nil
	})

	if err != nil {
		return nil, err
	}

	return hashMap, nil
}
