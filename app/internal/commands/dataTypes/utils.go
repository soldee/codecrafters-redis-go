package dataTypes

import "fmt"

const SEP string = "\r\n"

func GetUntilSeparator(raw *[]byte) ([]byte, int, error) {
	for i := 0; i < len(*raw); i++ {
		if (*raw)[i] == byte('\n') {
			if (*raw)[i-1] == byte('\r') {
				return (*raw)[0 : i-1], i + 1, nil
			}
		}
	}
	return nil, len(*raw), fmt.Errorf("hit EOF before getting separator")
}

func CheckSeparator(raw *[]byte) error {
	if len(*raw) < len(SEP) {
		return fmt.Errorf("expected separator '%v' at this position", []byte(SEP))
	}

	s := string((*raw)[:len(SEP)])
	if s != SEP {
		return fmt.Errorf("expected separator '%v' at this position but got %s", []byte(SEP), s)
	}

	*raw = (*raw)[len(SEP):]
	return nil
}
