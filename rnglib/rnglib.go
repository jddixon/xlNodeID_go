package rnglib

import "os"
import "strings"

func Version() (string, string) {
	return "0.1.1", "2013-06-14"
}

// a crude attempt at properties
var _FILE_NAME_STARTERS = strings.Split(
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_", "")

func FILE_NAME_STARTERS() []string {
	return _FILE_NAME_STARTERS
}

var _FILE_NAME_CHARS = strings.Split(
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-.", "")

func FILE_NAME_CHARS() []string {
	return _FILE_NAME_CHARS
}

/* XXX this should be in another package */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
