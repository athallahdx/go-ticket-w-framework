package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Generate valid bcrypt hash for "qwerty123"
	hash, err := bcrypt.GenerateFromPassword([]byte("qwerty123"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	hashStr := string(hash)

	// Read the seed file
	content, err := ioutil.ReadFile("seeders/001_seed_users.sql")
	if err != nil {
		panic(err)
	}

	// The dummy hash I used previously
	oldHash := "$2a$10$rN91o7bL/A4I5fB8VwW5OuJ9K4QZ9pL/mP3I5jN/2T7dG8VqR1P3a"

	// Replace it everywhere
	newContent := strings.ReplaceAll(string(content), oldHash, hashStr)

	// Write it back
	err = ioutil.WriteFile("seeders/001_seed_users.sql", []byte(newContent), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully replaced all password hashes with valid hash for 'qwerty123':", hashStr)
}
