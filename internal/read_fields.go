package internal

import (
	"bufio"
	"encoding/json"
	"io"
)

func StreamTestOutput(reader io.Reader, f func(TestOutputLine) error) error {
	r := bufio.NewReader(reader)

	for {
		packet, err := ReadField(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
		err = f(packet)
		if err != nil {
			return err
		}
	}
}

func ReadField(reader *bufio.Reader) (TestOutputLine, error) {
	var field TestOutputLine
	line, err := reader.ReadString('\n')
	if err != nil {
		return field, err
	}
	err = json.Unmarshal([]byte(line), &field)
	return field, err
}
