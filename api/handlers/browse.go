package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/goteleport-interview/fs4/api/internal/respond"
)

type BrowseHandler struct {
	rootDir string
}

func NewBrowseHandler(rootDir string) *BrowseHandler {
	return &BrowseHandler{rootDir: rootDir}
}

type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type BrowseResponse struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	Size     int64      `json:"size"`
	Contents []FileInfo `json:"contents"`
}

func (h *BrowseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		respond.WithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the path from query parameters
	requestPath := r.URL.Query().Get("path")
	if requestPath == "" {
		requestPath = "/"
	}

	// Clean and validate the path format
	cleanPath := filepath.Clean(requestPath)
	if !strings.HasPrefix(cleanPath, "/") {
		respond.WithError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	absRootDir, err := filepath.Abs(h.rootDir)
	if err != nil {
		respond.WithError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Construct and validate the full filesystem path
	fullPath := filepath.Join(absRootDir, cleanPath)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		respond.WithError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Validate that the requested path is within the root directory
	// This needs to happen before we check if the path exists
	if !isSubpath(absRootDir, absPath) {
		respond.WithError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Now that we've validated the path is within our root,
	// check if it exists and get file info
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			respond.WithError(w, "Path not found", http.StatusNotFound)
		} else {
			respond.WithError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Check if path is a directory
	if !info.IsDir() {
		respond.WithError(w, "Path is not a directory", http.StatusBadRequest)
		return
	}

	// Read directory contents
	dir, err := os.ReadDir(absPath)
	if err != nil {
		respond.WithError(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	// Build response
	contents := make([]FileInfo, 0, len(dir))
	var totalSize int64

	for _, entry := range dir {
		info, err := entry.Info()
		if err != nil {
			continue // Skip entries we can't read
		}

		fileType := "file"
		if info.IsDir() {
			fileType = "dir"
		}

		contents = append(contents, FileInfo{
			Name: info.Name(),
			Type: fileType,
			Size: info.Size(),
		})

		totalSize += info.Size()
	}

	response := BrowseResponse{
		Name:     filepath.Base(cleanPath),
		Type:     "dir",
		Size:     totalSize,
		Contents: contents,
	}

	respond.WithJSON(w, response, http.StatusOK)
}

// isSubpath checks if childPath is a subpath of parentPath
func isSubpath(parentPath, childPath string) bool {
	relativePath, err := filepath.Rel(parentPath, childPath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) && relativePath != ".."
}
