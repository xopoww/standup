package auth

import (
	"fmt"
	"strconv"
)

const base = 10

func idToString(id int64) string {
	return strconv.FormatInt(id, base)
}

func idFromString(s string) (id int64, err error) {
	id, err = strconv.ParseInt(s, base, 64)
	if err != nil {
		return 0, fmt.Errorf("id from string: %w", err)
	}
	return id, nil
}
