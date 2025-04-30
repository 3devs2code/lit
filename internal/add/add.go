package add

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Add locates a file by name, creates a blob, compresses it,
// generates its SHA-1 hash, and stores it in .lit/objects/
func Add(filename string) (string, error) {
	// 1. Search for the file recursively from current directory
	rootPath, _ := os.Getwd()
	var foundPath string

	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == filename {
			foundPath = path
			return filepath.SkipDir
		}
		return nil
	})

	if foundPath == "" {
		return "", fmt.Errorf("file '%s' not found", filename)
	}

	// 2. Read the file contents
	data, err := os.ReadFile(foundPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// 3. Create the Git-style blob (e.g., "blob 42\x00<content>")
	blob, err := blobCreation(string(data))
	if err != nil {
		return "", fmt.Errorf("error in blob creation: %v", err)
	}

	// 4. Generate SHA-1 hash of the blob
	sha, err := shaCreation(blob)
	if err != nil {
		return "", fmt.Errorf("error in SHA-1 creation: %v", err)
	}

	// 5. Compress the blob (zlib format)
	compressedFile, err := compression(blob)
	if err != nil {
		return "", fmt.Errorf("compression failed: %v", err)
	}

	// 6. Save the compressed blob to .lit/objects/<first-2-sha>/<rest>
	err = createFile(compressedFile, sha)
	if err != nil {
		return "", fmt.Errorf("failed to save blob: %v", err)
	}

	fmt.Println("✅ Blob saved at:", filepath.Join(".lit", "objects", sha[:2], sha[2:]))
	return sha, nil
}

// createFile saves compressed blob using Git-style sha path
func createFile(compressedFile []byte, sha string) error {
	dir := filepath.Join(".lit", "objects", sha[:2])
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(dir, sha[2:])
	return os.WriteFile(fullPath, compressedFile, 0644)
}

// compression compresses blob string using zlib
func compression(blob string) ([]byte, error) {
	var buffer bytes.Buffer
	w := zlib.NewWriter(&buffer)
	_, err := w.Write([]byte(blob))
	if err != nil {
		return nil, err
	}
	w.Close()
	return buffer.Bytes(), nil
}

// deCompression is for testing: decompresses compressed file and prints content
func deCompression(compressedFile []byte) (string, error) {
	buffer := bytes.NewBuffer(compressedFile)
	r, err := zlib.NewReader(buffer)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// shaCreation generates SHA-1 from blob string
func shaCreation(blob string) (string, error) {
	sha := sha1.New()
	sha.Write([]byte(blob))
	return fmt.Sprintf("%x", sha.Sum(nil)), nil
}

// blobCreation adds Git-style header to file content
func blobCreation(fileData string) (string, error) {
	header := fmt.Sprintf("blob %d\x00", len(fileData))
	return header + fileData, nil
}
