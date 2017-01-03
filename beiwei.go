// Copyright 2016 Nikugame. All Rights Reserved.

package ShortMessagingService

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/axgle/mahonia"
)

const (
	//BEIWEIKEY configure key by xiao
	BEIWEIKEY = "beiwei"
)

//BeiWei implements ShortMessagingService to choose beiwei channel
type BeiWei struct {
}

//Parse Create a mew Message by beiwei channel
func (beiwei *BeiWei) Parse(filename string) (Message, error) {
	conf, err := LoadConfigure(filename)
	fmt.Println(conf)
	if err != nil {
		return nil, err
	}
	bsms := &BeiWeiShortMesssagingService{}
	if confmap, found := conf[BEIWEIKEY]; found {
		beiweiElem := reflect.ValueOf(bsms).Elem()
		beiweiType := beiweiElem.Type()
		for i := 0; i < beiweiElem.NumField(); i++ {
			name := beiweiType.Field(i).Name
			if value, found := confmap[strings.ToLower(name)]; found {
				beiweiElem.FieldByName(name).SetString(value)
			} else {
				return nil, fmt.Errorf("can't find xiao configure: %s", name)
			}
		}
		return bsms, nil
	}
	return nil, fmt.Errorf("configure not found")
}

//Send send short message
func (bsms *BeiWeiShortMesssagingService) Send(phone string, message string) (map[string]bool, error) {
	result := make(map[string]bool)
	result[phone] = false
	u, err := url.Parse(bsms.URL)
	if err != nil {
		return result, err
	}

	rrid := RandomString(8, NUMBER, LOWCHARACTER)

	q := u.Query()
	q.Set("sn", bsms.SN)
	q.Set("pwd", bsms.Auth())
	q.Set("mobile", phone)
	decode := mahonia.NewEncoder("GBK")
	q.Set("content", decode.ConvertString(message))
	q.Set("ext", bsms.EXT)
	q.Set("stime", "")
	q.Set("rrid", rrid)

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
	if strings.Contains(string(body), rrid) {
		result[phone] = true
	} else {
		return nil, fmt.Errorf("beiwei sms error %s", string(body))
	}
	return result, nil
}

//Auth beiwei auth function
func (bsms *BeiWeiShortMesssagingService) Auth() string {
	str := fmt.Sprintf("%s%s", bsms.SN, bsms.PWD)
	return strings.ToUpper(MD5(str))
}

//BeiWeiShortMesssagingService a ShortMessagingService represents the beiwei channel
type BeiWeiShortMesssagingService struct {
	SN  string
	PWD string
	URL string
	EXT string
}

func init() {
	Register("beiwei", &BeiWei{})
}
