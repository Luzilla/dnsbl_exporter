package resolvconf

import (
	"bufio"
	"os"
	"strings"
)

func GetServers(path string) (servers []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "nameserver ") {
			continue
		}

		server, ok := strings.CutPrefix(line, "nameserver ")
		if !ok {
			continue
		}

		servers = append(servers, server)
	}

	err = scanner.Err()
	return
}
