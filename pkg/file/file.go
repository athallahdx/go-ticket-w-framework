package file

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	MaxFileSize  = 5 * 1024 * 1024 // 5MB
	MaxFileCount = 10
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

var uploadDirs = map[EntityType]string{
	EntityUser:      "uploads/users",
	EntityEvent:     "uploads/events",
	EntityOrganizer: "uploads/organizers",
}

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9_]`)

type EntityType string

const (
	EntityUser      EntityType = "user"
	EntityEvent     EntityType = "event"
	EntityOrganizer EntityType = "organizer"
)

// ----------------------
// helpers
// ----------------------

func sanitize(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = nonAlphanumeric.ReplaceAllString(name, "")
	if name == "" {
		name = "file"
	}
	return name
}

func IsAllowedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[ext]
}

func IsAllowedSize(header *multipart.FileHeader) bool {
	return header.Size <= MaxFileSize
}

func GetUploadDir(entity EntityType) string {
	dir, ok := uploadDirs[entity]
	if !ok {
		return "uploads/misc"
	}
	return dir
}

// ----------------------
// filename generation
// ----------------------

func GenerateFilename(entity EntityType, id int64, name string, original string) string {
	ext := strings.ToLower(filepath.Ext(original))
	base := sanitize(name)
	timestamp := time.Now().UnixNano()

	switch entity {
	case EntityUser:
		return fmt.Sprintf("user_%d_%s_%d%s", id, base, timestamp, ext)
	case EntityEvent:
		if base == "file" {
			return fmt.Sprintf("event_%d_%d%s", id, timestamp, ext)
		}
		return fmt.Sprintf("event_%d_%s_%d%s", id, base, timestamp, ext)
	case EntityOrganizer:
		return fmt.Sprintf("organizer_%d_%s_%d%s", id, base, timestamp, ext)
	default:
		return fmt.Sprintf("%s_%d_%s_%d%s", entity, id, base, timestamp, ext)
	}
}

// ----------------------
// save operations
// ----------------------

// SaveFile saves a single already-opened file to dstPath.
// Validates extension and size before writing.
func SaveFile(file multipart.File, header *multipart.FileHeader, dstPath string) error {
	if !IsAllowedExtension(header.Filename) {
		return errors.New("file type not allowed")
	}
	if !IsAllowedSize(header) {
		return errors.New("file size exceeds 5MB limit")
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return err
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return err
}

// saveOne opens, validates, and saves a single multipart file header.
// dstPath is derived automatically from EntityType.
func saveOne(header *multipart.FileHeader, entity EntityType, id int64, name string) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	filename := GenerateFilename(entity, id, name, header.Filename)
	dstPath := filepath.Join(GetUploadDir(entity), filename)

	if err := SaveFile(file, header, dstPath); err != nil {
		return "", err
	}

	return dstPath, nil
}

// SaveSingleFile saves one multipart file and returns the saved path.
// Use this for profile pictures, thumbnails, logos — single file uploads.
func SaveSingleFile(header *multipart.FileHeader, entity EntityType, id int64, name string) (string, error) {
	if header == nil {
		return "", errors.New("no file provided")
	}

	return saveOne(header, entity, id, name)
}

// SaveMultipleFiles saves multiple multipart files and returns their saved paths.
// If any file fails, already saved files are rolled back (deleted).
// Use this for event image galleries — multiple file uploads.
func SaveMultipleFiles(files []*multipart.FileHeader, entity EntityType, id int64, name string) ([]string, error) {
	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}
	if len(files) > MaxFileCount {
		return nil, fmt.Errorf("too many files, maximum allowed is %d", MaxFileCount)
	}

	var savedPaths []string

	for _, header := range files {
		path, err := saveOne(header, entity, id, name)
		if err != nil {
			// rollback — delete already saved files before returning
			for _, saved := range savedPaths {
				os.Remove(saved)
			}
			return nil, fmt.Errorf("failed on file %s: %w", header.Filename, err)
		}
		savedPaths = append(savedPaths, path)
	}

	return savedPaths, nil
}

// DeleteFile removes a file from disk.
// Call this when removing a profile picture or event image from the DB.
func DeleteFile(path string) error {
	if path == "" {
		return errors.New("file path is empty")
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not found")
		}
		return err
	}

	return nil
}

// ReplaceFile deletes the old file and saves the new one atomically.
// If saving the new file fails, old file is NOT restored — handle with care.
// Use this when updating a profile picture or organizer logo.
func ReplaceFile(oldPath string, newHeader *multipart.FileHeader, entity EntityType, id int64, name string) (string, error) {
	newPath, err := saveOne(newHeader, entity, id, name)
	if err != nil {
		return "", err
	}

	// only delete old file after new one is safely saved
	if oldPath != "" {
		os.Remove(oldPath) // non-fatal — old file might already be gone
	}

	return newPath, nil
}
