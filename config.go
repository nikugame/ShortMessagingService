// Copyright 2016 Nikugame. All Rights Reserved.

package ShortMessagingService

import (
	"bufio"
	"io"
	"os"
	"strings"
)

//LoadDataByIni load configure by ini type file
func LoadDataByIni(filename string, object interface{}) (*interface{}, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if len(line) == 0 {
			continue
		}
		l := strings.TrimSpace(string(line))
	}

	return nil, nil
}
