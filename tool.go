package webutils

import (
    "errors"
    "fmt"
    "github.com/go-resty/resty/v2"
    "html"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

func WaitToSecond(targetSecond int, aheadMs int) time.Duration {
    now := time.Now()

    waitDuration := time.Duration((targetSecond-now.Second())*1000_000_000 - now.Nanosecond() - aheadMs*1000_000)
    return waitDuration
}

func GetWaiMaiSeverTime() (timestamp int64, err error) {
    client := resty.New()
    client.SetTimeout(time.Second * 20)
    timeUrl := "https://promotion.waimai.meituan.com/lottery/limitcouponcomponent/getTime"

    headers := map[string]string{
        "Sec-Fetch-Site":  "same-site",
        "Host":            "promotion.waimai.meituan.com",
        "Sec-Fetch-Dest":  "empty",
        "Accept-Encoding": "gzip, deflate, br",
        "Sec-Fetch-Mode":  "cors",
        "Accept":          "application/json, text/plain, */*",
        "Connection":      "keep-alive",
        "Accept-Language": "zh-CN,zh-Hans;q=0.9",
        "Origin":          "https://market.waimai.meituan.com",
        "Referer":         "https://market.waimai.meituan.com/",
        "User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.41(0x1800292e) NetType/WIFI Language/zh_CN miniProgram/wx2c348cf579062e56",
    }

    var timeResp TimeResponse

    response, err := client.R().
        EnableTrace().
        SetHeaders(headers).
        SetResult(&timeResp).
        Get(timeUrl)

    if err != nil {
        return
    }

    if response.StatusCode() != 200 {
        err = errors.New(strconv.Itoa(response.StatusCode()))
        return
    }

    ti := response.Request.TraceInfo()

    timestamp = timeResp.Data + ti.ResponseTime.Milliseconds()
    return
}

func GetNextMinuteMs(aheadMs int) (time.Duration, time.Duration, error) {
    timestamp, err := GetWaiMaiSeverTime()
    duration := WaitToSecond(60, aheadMs)
    if err != nil {
        return 0, 0, fmt.Errorf("获取服务器时间失败，%s", err.Error())
    }

    gap := time.Duration(time.Now().UnixMilli()-timestamp) * time.Millisecond
    duration += gap

    return duration, gap, nil
}

func CookieToMap(cookie string) map[string]string {
    result := make(map[string]string)
    for _, s := range strings.Split(cookie, ";") {
        pair := strings.SplitN(s, "=", 2)
        if len(pair) == 2 {
            result[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
        }
    }
    return result
}

func RandomEmoji() string {
    // http://apps.timwhitlock.info/emoji/tables/unicode
    emoji := [][]int{
        // Emoticons icons
        {128513, 128591},
        // Transport and map symbols
        {128640, 128704},
    }
    r := emoji[rand.Int()%len(emoji)]
    min := r[0]
    max := r[1]
    n := rand.Intn(max-min+1) + min
    return html.UnescapeString("&#" + strconv.Itoa(n) + ";")
}
