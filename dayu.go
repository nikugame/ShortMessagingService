// Copyright 2016 Nikugame. All Rights Reserved.

package ShortMessagingService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"public/tools"
	"reflect"
	"sort"
	"strings"
	"time"
)

//Dayu implements ShortMessagingService to choose dayu channel
type Dayu struct{}

//Parse Create a new Message by dayu channel
func (dayu *Dayu) Parse(filename string) (Message, error) {
	conf, err := LoadConfigure(filename)
	if err != nil {
		return nil, err
	}
	dsms := &DayuShortMessagingService{}
	if confmap, found := conf["dayu"]; found {
		xiaoElem := reflect.ValueOf(dsms).Elem()
		xiaoType := xiaoElem.Type()
		for i := 0; i < xiaoElem.NumField(); i++ {
			name := xiaoType.Field(i).Name
			if value, found := confmap[strings.ToLower(name)]; found {
				xiaoElem.FieldByName(name).SetString(value)
			} else {
				return nil, fmt.Errorf("can't find xiao configure: %s", name)
			}
		}
		return dsms, nil
	}
	return nil, fmt.Errorf("configure not found")
}

//DayuShortMessagingService a ShortMessagingService represents the dayu channel
type DayuShortMessagingService struct {
	Name     string
	Key      string
	Sign     string
	Template string
	Secert   string
	URL      string
}

//Send send short message
func (dsms *DayuShortMessagingService) Send(phone string, message string) (map[string]bool, error) {
	result := make(map[string]bool)
	result[phone] = false

	var param struct {
		Code    string `json:"code"`
		Product string `json:"product"`
	}

	param.Code = message
	param.Product = dsms.Name

	str, _ := json.Marshal(param)

	m := make(map[string]string)
	m["app_key"] = dsms.Key
	m["format"] = "json"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["v"] = "2.0"
	m["method"] = "alibaba.aliqin.fc.sms.num.send"
	m["sign_method"] = "md5"
	m["sms_type"] = "normal"
	m["sms_free_sign_name"] = dsms.Sign
	m["sms_param"] = string(str)
	m["rec_num"] = phone
	m["sms_template_code"] = dsms.Template
	m["sign"] = sign(m, dsms.Secert)

	v := url.Values{}
	for k, values := range m {
		v.Set(k, values)
	}
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := http.Client{}
	request, err := http.NewRequest("POST", dsms.URL, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	var alim map[string]interface{}
	if err = json.Unmarshal(b, &alim); err != nil {
		return nil, err
	}
	if _, found := alim["alibaba_aliqin_fc_sms_num_send_response"]; found {
		var res struct {
			Key struct {
				Result struct {
					Code    string `json:"err_code"`
					Model   string `json:"model"`
					Success bool   `json:"success"`
				} `json:"result"`
				ID string `json:"request_id"`
			} `json:"alibaba_aliqin_fc_sms_num_send_response"`
		}
		if err := json.Unmarshal(b, &res); err != nil {
			return nil, err
		}

		switch res.Key.Result.Code {
		case "0":
			result[phone] = true
			return result, nil
		default:
			return nil, fmt.Errorf(res.Key.Result.Code)
		}
	}
	if _, found := alim["error_response"]; found {
		var res struct {
			Key struct {
				Code       int    `json:"code"`
				Message    string `json:"msg"`
				SubCode    string `json:"sub_code"`
				SubMessage string `json:"sub_msg"`
			} `json:"error_response"`
		}
		if err := json.Unmarshal(b, &res); err != nil {
			return nil, err
		}
		switch res.Key.SubCode {
		case "isv.OUT_OF_SERVICE":
			log.Println("短信业务停机，要充值！")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.PRODUCT_UNSUBSCRIBE":
			log.Println("产品服务未开通，参数配置错误")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.ACCOUNT_NOT_EXISTS":
			log.Println("账户信息不存在，参数配置错误")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.ACCOUNT_ABNORMAL":
			log.Println("账户信息异常，联系大于客服")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.SMS_TEMPLATE_ILLEGAL":
			log.Println("模板不合法， 参数配置错误")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.SMS_SIGNATURE_ILLEGAL":
			log.Println("签名不合法，参数配置错误")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.MOBILE_NUMBER_ILLEGAL":
			log.Println("手机号码格式错误，检查是否存在漏洞，被非法调用 ")
			return nil, fmt.Errorf("Bad Request, Must phone number!")
		case "isv.MOBILE_COUNT_OVER_LIMIT":
			log.Println("手机号码数量超过限制, 不能超过200个号码")
			return nil, fmt.Errorf("Too many phone number!")
		case "isv.TEMPLATE_MISSING_PARAMETERS":
			log.Println("短信模板变量缺少参数，配置错误")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.INVALID_PARAMETERS":
			log.Println("参数异常，配置错误")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.BUSINESS_LIMIT_CONTROL":
			log.Printf("%s : 超过1分钟1条 或者 1分钟7条 1天50条限制 \n", phone)
			return nil, fmt.Errorf("Send too often, Please try 1 min later")
		case "isv.INVALID_JSON_PARAM":
			log.Println("JSON参数不合法，代码被改了？")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isp.SYSTEM_ERROR":
			log.Println("大于挂了！改用其他通道发送")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.BLACK_KEY_CONTROL_LIMIT":
			log.Println("模板变量中存在黑名单关键字，大于后台配置错")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.PARAM_NOT_SUPPORT_URL":
			log.Println("不支持url为变量，配置错误")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.PARAM_LENGTH_LIMIT":
			log.Println("变量长度受限，配置错误")
			return nil, fmt.Errorf("System ERROR, Try later!")
		case "isv.AMOUNT_NOT_ENOUGH":
			log.Println("余额不足，要充值！")
			return nil, fmt.Errorf("System ERROR, Try later!")
		}
	}
	return result, nil
}

//sign func(map[string]string, string) string
func sign(m map[string]string, secret string) string {

	list := []string{}
	for k := range m {
		list = append(list, k)
	}
	sort.Strings(list)

	param := ""
	for i := 0; i < len(list); i++ {
		param = tools.StringJoin("", param, list[i], m[list[i]])
	}
	str := tools.StringJoin("", secret, param, secret)
	return strings.ToUpper(tools.BuilderMD5(str))
}

func init() {
	Register("dayu", &Dayu{})
}
