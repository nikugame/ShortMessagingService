// Copyright 2016 Nikugame. All Rights Reserved.

package ShortMessagingService

//Xiao implements ShortMessagingService to choose xiao channel
type Xiao struct {
}

//Parse Create a mew Message by xiao channel
func (xiao *Xiao) Parse(filename, filetype string) (Message, error) {
	return &XiaoShortMesssagingService{}, nil
}

//XiaoShortMesssagingService a ShortMessagingService represents the xiao channel
type XiaoShortMesssagingService struct {
	uid    string
	auth   string
	expid  string
	encode string
}

//Send send short message
func (xsms *XiaoShortMesssagingService) Send(phone []string, message string) (map[string]bool, error) {
	result := make(map[string]bool)
	return result, nil
}
