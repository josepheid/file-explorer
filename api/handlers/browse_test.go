package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// setupTestDirectory creates a temporary directory structure for testing
func setupTestDirectory(t *testing.T) (string, func()) {
	// Create temp directory
	rootDir, err := os.MkdirTemp("", "file-browser-test-*")
	if err != nil {
		t.Fatal(err)
	}

	// Create test directory structure
	dirs := []string{
		"empty",
		"dir1",
		"dir1/subdir",
	}

	files := []struct {
		path string
		size int
	}{
		{"dir1/file1.txt", 100},
		{"dir1/file2.txt", 200},
		{"dir1/subdir/file3.txt", 300},
	}

	// Create directories
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(rootDir, dir), 0755)
		if err != nil {
			os.RemoveAll(rootDir)
			t.Fatal(err)
		}
	}

	// Create files
	for _, file := range files {
		data := make([]byte, file.size)
		err := os.WriteFile(filepath.Join(rootDir, file.path), data, 0644)
		if err != nil {
			os.RemoveAll(rootDir)
			t.Fatal(err)
		}
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(rootDir)
	}

	return rootDir, cleanup
}

func TestBrowseHandler(t *testing.T) {
	// Setup test directory
	rootDir, cleanup := setupTestDirectory(t)
	defer cleanup()

	tests := []struct {
		name           string
		path           string
		method         string
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:           "Root Directory",
			path:           "/",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var response BrowseResponse
				err := json.Unmarshal(body, &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if len(response.Contents) != 2 { // dir1 and empty
					t.Errorf("Expected 2 items in root, got %d", len(response.Contents))
				}
			},
		},
		{
			name:           "Empty Directory",
			path:           "/empty",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var response BrowseResponse
				err := json.Unmarshal(body, &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if len(response.Contents) != 0 {
					t.Errorf("Expected empty directory, got %d items", len(response.Contents))
				}
			},
		},
		{
			name:           "Directory with Contents",
			path:           "/dir1",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var response BrowseResponse
				err := json.Unmarshal(body, &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if len(response.Contents) != 3 { // 2 files + 1 subdir
					t.Errorf("Expected 3 items, got %d", len(response.Contents))
				}
			},
		},
		{
			name:           "Non-existent Path",
			path:           "/nonexistent",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			validateBody:   nil,
		},
		{
			name:           "Invalid Method",
			path:           "/",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			validateBody:   nil,
		},
		{
			name:           "Path Traversal Attempt",
			path:           "/../../../etc",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			validateBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			handler := NewBrowseHandler(rootDir)
			req := httptest.NewRequest(tt.method, "/api/v1/browse?path="+tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Validate response body if needed
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
		})
	}
}

func TestBrowseHandlerContentValidation(t *testing.T) {
	rootDir, cleanup := setupTestDirectory(t)
	defer cleanup()

	handler := NewBrowseHandler(rootDir)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/browse?path=/dir1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var response BrowseResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Validate content type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Validate response structure
	if response.Type != "dir" {
		t.Errorf("Expected type 'dir', got %s", response.Type)
	}

	if response.Name != "dir1" {
		t.Errorf("Expected name 'dir1', got %s", response.Name)
	}

	// Check for specific files and directories
	fileMap := make(map[string]FileInfo)
	for _, item := range response.Contents {
		fileMap[item.Name] = item
	}

	// Check file1.txt
	if file, exists := fileMap["file1.txt"]; !exists {
		t.Error("file1.txt not found in response")
	} else {
		if file.Type != "file" {
			t.Errorf("Expected file1.txt type to be 'file', got %s", file.Type)
		}
		if file.Size != 100 {
			t.Errorf("Expected file1.txt size to be 100, got %d", file.Size)
		}
	}

	// Check subdir
	if dir, exists := fileMap["subdir"]; !exists {
		t.Error("subdir not found in response")
	} else {
		if dir.Type != "dir" {
			t.Errorf("Expected subdir type to be 'dir', got %s", dir.Type)
		}
	}
}

func TestPathTraversalAttempts(t *testing.T) {
	rootDir, cleanup := setupTestDirectory(t)
	defer cleanup()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "Simple Path Traversal",
			path:           "/../../../etc/passwd",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Encoded Path Traversal",
			path:           "/%2e%2e/%2e%2e/etc/passwd",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Double Slash Path Traversal",
			path:           "//etc/passwd",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unicode Path Traversal",
			path:           "/‥/‥/etc/passwd",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Valid Nested Path",
			path:           "/dir1/subdir",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewBrowseHandler(rootDir)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/browse?path="+tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
