# sms
golang send message tools

#### 2017-01-13

add support dayu to send verification code

configure

    [dayu]
    url = "xxxx"x
    name = "xxxxx"x
    key = "xxxx"
    secert = "xxxx"
    sign = "xxxxx"
    template = "xxxxx"


#### 2016

support

    xiao/beiwei

install 

    go get -u github.com/axgle/mahonia 
    go get -u github.com/nikugame/ShortMessagingService

use

    package main

    import "github.com/nikugame/ShortMessagingService"
    import "fmt"

    func main() {
        sms, err ：= ShortMessagingService.NewShortMessagingService("xiao", "conf/xiao.ini")
        //sms, err ：= ShortMessagingService.NewShortMessagingService("beiwei", "conf/beiwei.ini")
        if err != nil {
            //doing some thing
        }
        status, err := sms.Send("180****1021", "this is a test message")
        if err != nil {
            //doing some thing
        }
        for phone, stat := status {
            fmt.Printf("%s send state %b", phone, stat)
        }
    }


configure 

    [xiao]
    uid = "000000"
    cid = "xxxx"
    pwd = "xxxxxx"
    url = "http://xiao.url"
    [beiwei]
    url = "http://beiwei.url"
    sn = "xxxxxxxxxxx"
    pwd = "xxxxx"
    ext = "1"

