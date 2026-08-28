package manifest

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var errIsDir = errors.New("is directory")

func hashMd5File(path string, d os.DirEntry, root string, err error) (relPath, hash string, _ error) {
	if err != nil {
		return "", "", err
	}

	if d.IsDir() {
		return "", "", errIsDir
	}

	relPath, err = filepath.Rel(root, path)

	if err != nil {
		return "", "", err
	}

	relPath = filepath.ToSlash(relPath)

	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", "", err
	}
	return relPath, hex.EncodeToString(h.Sum(nil)), nil
}
