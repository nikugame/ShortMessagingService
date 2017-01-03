// Copyright 2016 Nikugame. All Rights Reserved.

package ShortMessagingService

import (
	"fmt"
	"net/url"
	"reflect"

	"io/ioutil"
	"net/http"

	"strings"

	"github.com/axgle/mahonia"
)

const (
	//XIAOKEY configure key by xiao
	XIAOKEY = "xiao"
)

//Xiao implements ShortMessagingService to choose xiao channel
type Xiao struct {
}

//Parse Create a mew Message by xiao channel
func (xiao *Xiao) Parse(filename string) (Message, error) {
	conf, err := LoadConfigure(filename)
	// fmt.Println(conf)
	if err != nil {
		return nil, err
	}
	xsms := &XiaoShortMesssagingService{}
	if confmap, found := conf[XIAOKEY]; found {
		xiaoElem := reflect.ValueOf(xsms).Elem()
		xiaoType := xiaoElem.Type()
		for i := 0; i < xiaoElem.NumField(); i++ {
			name := xiaoType.Field(i).Name
			if value, found := confmap[strings.ToLower(name)]; found {
				xiaoElem.FieldByName(name).SetString(value)
			} else {
				return nil, fmt.Errorf("can't find xiao configure: %s", name)
			}
		}
		return xsms, nil
	}
	return nil, fmt.Errorf("configure not found")
}

//XiaoShortMesssagingService a ShortMessagingService represents the xiao channel
type XiaoShortMesssagingService struct {
	UID string
	PWD string
	URL string
	CID string
}

//Send send short message
func (xsms *XiaoShortMesssagingService) Send(phone string, message string) (map[string]bool, error) {
	result := make(map[string]bool)
	result[phone] = false
	u, err := url.Parse(xsms.URL)
	if err != nil {
		return result, err
	}
	q := u.Query()
	q.Set("uid", xsms.UID)
	q.Set("auth", xsms.Auth())
	q.Set("expid", "0")
	q.Set("encode", "utf-8")
	decode := mahonia.NewEncoder("utf-8")
	q.Set("msg", decode.ConvertString(message))
	q.Set("mobile", phone)

	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if body[0] == '0' {
		result[phone] = true
	}
	return result, nil
}

//Auth xiao auth fucntion
func (xsms *XiaoShortMesssagingService) Auth() string {
	str := fmt.Sprintf("%s%s", xsms.CID, xsms.PWD)
	return MD5(str)
}

func init() {
	Register("xiao", &Xiao{})
}
