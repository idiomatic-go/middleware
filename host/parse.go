package host

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	//token = byte(' ')
	//eol       = byte('\n')
	comment   = "//"
	delimiter = ":"
)

func ParseMap(buf []byte) (map[string]string, error) {
	m := make(map[string]string)
	if len(buf) == 0 {
		return m, nil
	}
	r := bytes.NewReader(buf)
	reader := bufio.NewReader(r)
	var line string
	var err error
	for {
		line, err = reader.ReadString('\n')
		k, v, err0 := ParseLine(line)
		if err0 != nil {
			return m, err0
		}
		if len(k) > 0 {
			m[k] = v
		}
		if err == io.EOF {
			break
		} else {
			if err != nil {
				break
			}
		}
	}
	return m, nil
}

func ParseLine(line string) (string, string, error) {
	if len(line) == 0 {
		return "", "", nil
	}
	line = strings.TrimLeft(line, " ")
	if len(line) == 0 || strings.HasPrefix(line, comment) {
		return "", "", nil
	}
	i := strings.Index(line, delimiter)
	if i == -1 {
		return "", "", fmt.Errorf("invalid argument : line does not contain the ':' delimeter : [%v]", line)
	}
	key := line[:i]
	val := line[i+1:]
	return strings.TrimSpace(key), strings.TrimLeft(val, " "), nil
}
