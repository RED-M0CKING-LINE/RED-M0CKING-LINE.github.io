package middleware

import (
	"encoding/base64"
	"encoding/binary"
	"hash/crc32"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	crc32cTable = crc32.MakeTable(crc32.Castagnoli)
	assetHashes = map[string]string{}
)

func DumpAssetHashes() map[string]string { return assetHashes }

// Opens a file and generates a base64 representation fo its crc32c hash
func CRC32cBase64(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return "000000"
	}
	defer f.Close()

	h := crc32.New(crc32cTable)
	io.Copy(h, f)

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], h.Sum32())
	return base64.RawURLEncoding.EncodeToString(buf[:])
}

// generate a map of hashes for static assets
func InitAssetHashes(files map[string]string) {
	for urlPath, fsPath := range files {
		assetHashes[urlPath] = CRC32cBase64(fsPath)
	}
}

// walks the directory and initializes asset hashes for all files found
func InitAssetDirHashes(staticDir, urlPrefix string) error {
	return fs.WalkDir(os.DirFS(staticDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		assetHashes[urlPrefix+"/"+path] = CRC32cBase64(filepath.Join(staticDir, path))
		return nil
	})
}

// add the GET parameter to the URL based on the map
func AssetURL(urlPath string) string {
	if hash, ok := assetHashes[urlPath]; ok {
		return urlPath + "?c=" + hash
	}
	return urlPath
}

// FuncMap of asset URLs
func AssetURLFuncMap() template.FuncMap {
	return template.FuncMap{
		"assetURL": AssetURL,
	}
}
