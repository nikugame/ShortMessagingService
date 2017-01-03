// Copyright 2016 Nikugame. All Rights Reserved.

package ShortMessagingService

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

//LoadConfigure load configure file
func LoadConfigure(filename string) (map[string]map[string]string, error) {

	result := make(map[string]map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	defer file.Close()

	section := "default"
	m := make(map[string]string)
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
		if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			l := strings.TrimSpace(l)
			section = strings.ToLower(l[1 : len(l)-1])
			m = make(map[string]string)
			continue
		}
		if strings.HasPrefix(l, "#") || strings.HasPrefix(l, ";") {
			continue
		}
		parts := strings.SplitN(l, "=", 2)
		// fmt.Printf("configure value : %s = %s \n", parts[0], parts[1])
		name := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		if strings.HasPrefix(value, "\"") {
			// fmt.Printf("has \" : xx%sxx \n", value)
			value = value[1 : len(value)-1]
		}
		m[name] = value
		result[section] = m
	}

	return result, nil
}

//MD5 md5 function
func MD5(str string) string {
	ctx := md5.New()
	ctx.Write([]byte(str))
	s := ctx.Sum(nil)
	return hex.EncodeToString(s)
}

//RandomString build random string
func RandomString(length int, character ...string) string {
	random := strings.Join([]string(character), "")
	bytes := []byte(random)
	res := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		res = append(res, bytes[r.Intn(len(bytes))])
	}
	return string(res)
}

const (
	//NUMBER The pure digital
	NUMBER = "0123456789"
	//LOWCHARACTER lower-case letters
	LOWCHARACTER = "abcdefghijklmnopqrstuvwxyz"
	//UPCHARACTER capital letter
	UPCHARACTER = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)
