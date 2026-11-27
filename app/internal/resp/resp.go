package resp

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse reads one command. A zero byte count with no error means the command
// is incomplete.
func Parse(s string) ([]string, int, error) {
	if len(s) == 0 {
		return nil, 0, nil
	}

	count, consumed := 1, 0
	if s[0] == '*' {
		var err error
		count, consumed, err = parseLength(s)
		if err != nil || consumed == 0 {
			return nil, 0, err
		}
	}

	var args []string
	for range count {
		remaining := s[consumed:]
		if len(remaining) == 0 {
			return nil, 0, nil
		}
		if remaining[0] != '$' {
			return nil, 0, fmt.Errorf("expected a RESP bulk string")
		}

		length, headerSize, err := parseLength(remaining)
		if err != nil || headerSize == 0 {
			return nil, 0, err
		}
		body := remaining[headerSize:]
		if length > len(body)-2 {
			return nil, 0, nil
		}
		if body[length:length+2] != "\r\n" {
			return nil, 0, fmt.Errorf("missing CRLF after RESP bulk string")
		}

		args = append(args, body[:length])
		consumed += headerSize + length + 2
	}
	return args, consumed, nil
}

func parseLength(s string) (int, int, error) {
	end := strings.Index(s, "\r\n")
	if end < 0 {
		return 0, 0, nil
	}

	length, err := strconv.Atoi(s[1:end])
	if err != nil || length < 0 {
		return 0, 0, fmt.Errorf("invalid RESP length %q", s[1:end])
	}
	return length, end + 2, nil
}


func BulkStr(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}
