package util

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrAborted = errors.New("aborted")

// Confirm asks for user confirmation using prompt. Accepted answers are 'Y', 'y', 'N', 'n'.
// Returns nil on 'Y'/'y' or ErrAborted on 'N'/'n'.
func Confirm(out io.Writer, in io.Reader, prompt string) error {
	_, err := fmt.Fprintf(out, "%s [y/n]\n", prompt)
	if err != nil {
		return err
	}
	for {
		var ans string
		_, err := fmt.Fscanln(in, &ans)
		if err != nil {
			return err
		}
		ans = strings.ToLower(ans)
		switch ans {
		case "y":
			return nil
		case "n":
			return ErrAborted
		default:
			_, err = fmt.Fprint(out, "Unrecognized input, try again. [y/n]\n")
			if err != nil {
				return err
			}
		}
	}
}

// Confirmf is a printf-like version of Confirm.
func Confirmf(out io.Writer, in io.Reader, format string, a ...any) error {
	return Confirm(out, in, fmt.Sprintf(format, a...))
}
