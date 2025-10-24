package ui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//
//go:embed all:webui
var embeddedFS embed.FS

// webuiFS holds the filesystem for the contents of the 'webui' directory.
var webuiFS fs.FS

func init() {
	var err error
	// We create a sub-filesystem that starts at the 'webui' directory,
	// effectively stripping the 'webui/' prefix from file paths.
	webuiFS, err = fs.Sub(embeddedFS, "webui")
	if err != nil {
		log.Fatalf("failed to create webui sub-filesystem: %v", err)
	}

	// --- DEBUGGING: List all embedded files ---
	// This code will run on startup and print every file path found in the
	// embedded webuiFS. If you don't see your index.html and asset files
	// listed in the console, it means the 'webui' directory is in the
	// wrong place or has the wrong name.
	log.Println("--- Checking Embedded Filesystem ---")
	err = fs.WalkDir(webuiFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		log.Printf("Found embedded file: %s\n", path)
		return nil
	})
	if err != nil {
		log.Printf("Error walking embedded files: %v", err)
	}
	log.Println("------------------------------------")
}

// DistFS returns the filesystem for the entire webui build output.
func DistFS() http.FileSystem {
	return http.FS(webuiFS)
}

// AssetsFS returns a filesystem specifically for the 'assets' subdirectory.
func AssetsFS() http.FileSystem {
	sub, err := fs.Sub(webuiFS, "assets")
	if err != nil {
		// This should never happen if the Vite build is correct.
		log.Fatalf("failed to create assets sub-filesystem: %v", err)
	}
	return http.FS(sub)
}
