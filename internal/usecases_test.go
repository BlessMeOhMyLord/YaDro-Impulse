package internal

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func createFile(t *testing.T, content string) string {
	dir := t.TempDir()
	path := filepath.Join(dir, "resolv.conf")

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func TestUsecasesDefault(t *testing.T) {
	path := createFile(t, "")
	repo := NewConfRepository(path)
	service := NewService(repo)

	dns, err := service.Delete(context.Background(), "something")

	if !errors.Is(err, ErrIsIncorrect) {
		t.Fatalf("expected ErrIsIncorrect got %v", err)
	}

	dns, err = service.Add(context.Background(), " 8.8.8.8 \t\n")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if dns != "8.8.8.8" {
		t.Fatalf("expected normalized dns 8.8.8.8, got %s", dns)
	}

	list, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 dns, got %d: %v", len(list), list)
	}

	if list[0] != "8.8.8.8" {
		t.Fatalf("expected 8.8.8.8, got %s", list[0])
	}

	deletedDNS, err := service.Delete(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if deletedDNS != "8.8.8.8" {
		t.Fatalf("expected deleted dns 8.8.8.8, got %s", deletedDNS)
	}

	list, err = service.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(list) != 0 {
		t.Fatalf("expected empty list, got %v", list)
	}
}

func TestUsecasesExistst(t *testing.T) {
	path := createFile(t, "nameserver 8.8.8.8\n")
	repo := NewConfRepository(path)
	service := NewService(repo)

	_, err := service.Add(context.Background(), "8.8.8.8")

	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestServiceDeleteNotFound(t *testing.T) {
	path := createFile(t, "nameserver 1.1.1.1\n")

	repo := NewConfRepository(path)
	service := NewService(repo)

	dns, err := service.Delete(context.Background(), "8.8.8.8")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if dns != "" {
		t.Fatalf("expected empty dns on error, got %s", dns)
	}
}
