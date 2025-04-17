package cpuinfoclient

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func scanFor(searchKey string) (string, error) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", fmt.Errorf("could not open /proc/cpuinfo: %w", err)
	}

	var textValue string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if key == searchKey {
			textValue = value
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return textValue, nil
}

func ProcessorName() (string, error) {
	return scanFor("model name")
}

func NumCores() (uint64, error) {
	textValue, err := scanFor("cpu cores")
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(textValue, 10, 64)
}
