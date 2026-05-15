package internal

import (
	"context"
	"os"
	"strings"
	"sync"
)

type ConfRepository struct {
	mu   sync.Mutex
	path string
}

func NewConfRepository(path string) *ConfRepository {
	return &ConfRepository{
		path: path,
	}
}

func (c *ConfRepository) readLines() ([]string, error) {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return nil, err
	}
	content := strings.TrimRight(string(data), "\n")
	if content == "" {
		return []string{}, nil
	}

	return strings.Split(content, "\n"), nil
}

func (c *ConfRepository) writeLines(lines []string) error {
	info, err := os.Stat(c.path)
	if err != nil {
		return err
	}

	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}

	return os.WriteFile(c.path, []byte(content), info.Mode())
}

func (c *ConfRepository) Add(ctx context.Context, newDns string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	lines, err := c.readLines()
	if err != nil {
		return err
	}

	for _, line := range lines {
		dnsLine, ok := parseLine(line)
		if !ok {
			continue
		}
		if dnsLine == newDns {
			return ErrAlreadyExists
		}
	}

	lines = append(lines, "nameserver "+newDns)
	return c.writeLines(lines)
}

func (c *ConfRepository) Delete(ctx context.Context, dnsToDelete string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	lines, err := c.readLines()
	if err != nil {
		return err
	}

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		dnsLine, ok := parseLine(line)
		if ok && dnsLine == dnsToDelete {
			continue
		}
		result = append(result, line)
	}

	if len(lines) == len(result) {
		return ErrNotFound
	}

	return c.writeLines(result)
}

func (c *ConfRepository) List(ctx context.Context) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	lines, err := c.readLines()
	if err != nil {
		return nil, err
	}
	dnsList := make([]string, 0, len(lines))
	for _, line := range lines {
		dnsLine, ok := parseLine(line)
		if ok {
			dnsList = append(dnsList, dnsLine)
		}
	}
	return dnsList, nil
}

func parseLine(line string) (string, bool) {
	line = strings.TrimSpace(line)
	fields := strings.Fields(line)
	if len(fields) < 2 || fields[0] != "nameserver" {
		return "", false
	}

	return fields[1], true
}
