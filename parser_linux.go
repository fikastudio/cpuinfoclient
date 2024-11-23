package cpuinfoclient

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ProcessorName() (string, error) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", fmt.Errorf("could not open /proc/cpuinfo: %w", err)
	}

	var name string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if key == "model name" {
			name = value
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return name, nil
}
