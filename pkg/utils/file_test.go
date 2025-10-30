package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDir(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "utils-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "create new directory",
			path:    filepath.Join(tmpDir, "newdir"),
			wantErr: false,
		},
		{
			name:    "create nested directory",
			path:    filepath.Join(tmpDir, "parent", "child", "grandchild"),
			wantErr: false,
		},
		{
			name:    "existing directory",
			path:    tmpDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify directory exists
			if err == nil {
				if _, statErr := os.Stat(tt.path); os.IsNotExist(statErr) {
					t.Errorf("Directory was not created: %s", tt.path)
				}
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "utils-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		content string
		wantErr bool
	}{
		{
			name:    "write simple file",
			path:    filepath.Join(tmpDir, "test.txt"),
			content: "Hello, World!",
			wantErr: false,
		},
		{
			name:    "write file in nested directory",
			path:    filepath.Join(tmpDir, "nested", "dir", "file.txt"),
			content: "Nested content",
			wantErr: false,
		},
		{
			name:    "overwrite existing file",
			path:    filepath.Join(tmpDir, "existing.txt"),
			content: "New content",
			wantErr: false,
		},
		{
			name:    "write empty file",
			path:    filepath.Join(tmpDir, "empty.txt"),
			content: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteFile(tt.path, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Verify file exists and has correct content
				data, readErr := os.ReadFile(tt.path)
				if readErr != nil {
					t.Errorf("Failed to read written file: %v", readErr)
					return
				}

				if string(data) != tt.content {
					t.Errorf("File content = %v, want %v", string(data), tt.content)
				}
			}
		})
	}
}

func TestWriteExecutableFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "utils-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		content string
		wantErr bool
	}{
		{
			name:    "write executable script",
			path:    filepath.Join(tmpDir, "script.sh"),
			content: "#!/bin/bash\necho 'Hello'",
			wantErr: false,
		},
		{
			name:    "write executable in nested directory",
			path:    filepath.Join(tmpDir, "scripts", "run.sh"),
			content: "#!/bin/bash\necho 'Run'",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteExecutableFile(tt.path, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteExecutableFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Verify file exists
				info, statErr := os.Stat(tt.path)
				if statErr != nil {
					t.Errorf("Failed to stat written file: %v", statErr)
					return
				}

				// Verify file is executable
				mode := info.Mode()
				if mode&0111 == 0 {
					t.Errorf("File is not executable: %s (mode: %v)", tt.path, mode)
				}

				// Verify content
				data, readErr := os.ReadFile(tt.path)
				if readErr != nil {
					t.Errorf("Failed to read written file: %v", readErr)
					return
				}

				if string(data) != tt.content {
					t.Errorf("File content = %v, want %v", string(data), tt.content)
				}
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "utils-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "exists.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file",
			path: testFile,
			want: true,
		},
		{
			name: "non-existing file",
			path: filepath.Join(tmpDir, "notexists.txt"),
			want: false,
		},
		{
			name: "directory",
			path: tmpDir,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FileExists(tt.path)
			if got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
