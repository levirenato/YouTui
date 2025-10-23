package ui

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func CheckYtDlpVersion() (currentVersion, latestVersion string, needsUpdate bool) {
	cmd := exec.Command("yt-dlp", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", "", false
	}

	currentVersion = strings.TrimSpace(string(output))
	
	re := regexp.MustCompile(`(\d{4})\.(\d{2})\.(\d{2})`)
	matches := re.FindStringSubmatch(currentVersion)
	if len(matches) < 4 {
		return currentVersion, "", false
	}

	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])
	
	versionDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	daysOld := int(time.Since(versionDate).Hours() / 24)
	
	needsUpdate = daysOld > 14
	
	return currentVersion, "", needsUpdate
}
