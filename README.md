# sms
golang send message tools

install 

    go get -u github.com/axgle/mahonia 
    go get -u github.com/nikugame/ShortMessagingService

use

    package main

    import "github.com/nikugame/ShortMessagingService"
    import "fmt"

    func main() {
        sms, err ：= ShortMessagingService.NewShortMessagingService("xiao", "conf/xiao.ini")
        if err != nil {
            //doing some thing
        }
        status, err := sms.Send("18001901021", "this is a test message")
        if err != nil {
            //doing some thing
        }
        for phone, stat := range status {
            fmt.Printf("%s send state %b", phone, stat)
        }
    }
