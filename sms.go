// Copyright 2016 Nikugame. All Rights Reserved.

package ShortMessagingService

import "fmt"

//Message defines how to send sms
type Message interface {
	Send(phone []string, message string) (map[string]bool, error)
}

var channels = make(map[string]ShortMessagingService)

//ShortMessagingService is the adapter interface for choose sms channel
type ShortMessagingService interface {
	Parse(filename, filetype string) (Message, error)
}

//NewShortMessagingService channel is xiao/beiwei/dayu
//filetype only support ini
func NewShortMessagingService(channel, filename, filetype string) (Message, error) {
	if filetype == "" {
		filetype = "ini"
	}
	if filename == "" {
		filename = "conf/settings.conf"
	}
	if service, ok := channels[channel]; ok {
		return service.Parse(filename, filetype)
	}
	return nil, fmt.Errorf("unknown sms channel: %q", channel)
}

//Register make a message channel available by the channel name
func Register(name string, sms ShortMessagingService) {
	if sms == nil {
		panic("[SMS]: Register channel is nil")
	}
	if _, ok := channels[name]; ok {
		panic("[SMS]: Register called teice for channel " + name)
	}
	channels[name] = sms
}
