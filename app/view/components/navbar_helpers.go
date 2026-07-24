package components

import (
	"strings"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

func navActive(item models.NavItem, activePath string) bool {
	for _, sub := range item.SubItems {
		if navPathActive(sub.Path, activePath) {
			return true
		}
	}
	return false
}

func navPathActive(path, activePath string) bool {
	if path == "/" {
		return activePath == path
	}
	return activePath == path || strings.HasPrefix(activePath, strings.TrimSuffix(path, "/")+"/")
}
