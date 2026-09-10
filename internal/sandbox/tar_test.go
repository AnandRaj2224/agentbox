package sandbox

import (
	"archive/tar"
	"io"
	"testing"
)

func TestCreateTarArchive(t *testing.T) {
	reader, err := CreateTarArchive("main.py", "print('hello')")
	if err != nil {
		t.Fatalf("CreateTarArchive() error = %v", err)
	}

	tr := tar.NewReader(reader)

	header, err := tr.Next()
	if err != nil {
		t.Fatalf("failed to read tar header: %v", err)
	}

	if header.Name != "main.py" {
		t.Errorf("header.Name = %q, want %q", header.Name, "main.py")
	}

	if header.Size != int64(len("print('hello')")) {
		t.Errorf("header.Size = %d, want %d", header.Size, len("print('hello')"))
	}

	content, err := io.ReadAll(tr)
	if err != nil {
		t.Fatalf("failed to read file content: %v", err)
	}

	if string(content) != "print('hello')" {
		t.Errorf("content = %q, want %q", string(content), "print('hello')")
	}
}
