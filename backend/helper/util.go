package helper

import (
	"strings"
)


func RemoveFirstIfMatch(s, match string) string {
    if len(s) > 0 && strings.HasPrefix(s, match) {
        return s[1:]
    }
    return s
}
