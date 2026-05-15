package internal

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRepository(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "resolv.conf")

	err := os.WriteFile(path, []byte(""), 0644)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewConfRepository(path)

	dns := []string{"8.8.8.8", "8.8.4.4", "0.0.0.0"}
	for i, dn := range dns {
		if err := repo.Add(context.Background(), dn); err != nil {
			t.Fatal("expected no err", dn, err)
		}

		list, err := repo.List(context.Background())
		if err != nil {
			t.Fatal("expected no err", dn, err)
		}

		if len(list) != 1+i {
			t.Fatal("expected", i+1, "dns", dn, len(list))
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	content := ""
	for _, dn := range dns {
		content += "nameserver " + dn + "\n"
	}

	if string(data) != content {
		t.Fatal("expected", content, "got", string(data))
	}

	if err := repo.Delete(context.Background(), dns[2]); err != nil {
		t.Fatal("expected no err on delete exists dns", dns[2], err)
	}

	if err := repo.Add(context.Background(), dns[1]); !errors.Is(err, ErrAlreadyExists) {
		t.Fatal("expected err on add exists dns", err)
	}

	list, err := repo.List(context.Background())
	if err != nil {
		t.Fatal("expected no err", err)
	}

	if len(list) != 2 {
		t.Fatal("expected 2 dns")
	}

	if list[0] != dns[0] || list[1] != dns[1] {
		t.Fatal("expected", dns[0], dns[1], "got", list[0], list[1])
	}
}
