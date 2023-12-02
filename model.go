package webutils

import (
    "fmt"
    "github.com/go-resty/resty/v2"
    "github.com/wxpusher/wxpusher-sdk-go"
    "github.com/wxpusher/wxpusher-sdk-go/model"
    "strings"
    "time"
)

type TimeResponse struct {
    Code    int    `json:"code"`
    SubCode int    `json:"subcode"`
    Msg     string `json:"msg"`
    Data    int64  `json:"data"`
}

type Coupon struct {
    Id           string  `json:"id"`
    MerchantName string  `json:"merchantName"`
    Money        float64 `json:"money"`
}

type Task struct {
    CronExpr string
    Coupon   Coupon
}

type GlobalOption struct {
    // 发包总数
    Count int
    // 发包间隔
    Interval int

    // 提前毫秒数，微调使用
    AheadMs int

    // 推送方式
    PushType int

    // fetch url
    ApiUrl string
}

// 推送接口
type pusher interface {
    pushMarkdown(title string, content string) (string, error)
}

type WxPusher struct {
    AppToken string
    UId      string
}

func (w *WxPusher) pushMarkdown(title string, content string) (string, error) {
    if w.AppToken == "" || w.UId == "" {
        return "", fmt.Errorf("token参数配置有误 [appToken: %s, uId: %s]", w.AppToken, w.UId)
    }

    msg := model.NewMessage(w.AppToken).
        SetContentType(3).
        SetSummary(title).
        SetContent(content).
        AddUId(w.UId)
    msgArr, err := wxpusher.SendMessage(msg)
    if err != nil {
        return "", err
    }
    return msgArr[0].Status, nil
}

type PushPlus struct {
    Token string
}

func (p *PushPlus) pushMarkdown(title string, content string) (string, error) {
    if p.Token == "" {
        return "", fmt.Errorf("token参数配置有误 [token: %s]", p.Token)
    }

    pushPlusUrl := "https://www.pushplus.plus/api/send"

    client := resty.New()

    body := map[string]string{
        "token":    p.Token,
        "title":    title,
        "content":  content,
        "template": "markdown",
        "channel":  "wechat",
    }

    response, err := client.R().
        SetBody(body).
        Post(pushPlusUrl)

    if err != nil {
        return "", err
    }
    return response.String(), nil
}

type CubeResp struct {
    Code int        `json:"code"`
    Msg  string     `json:"msg"`
    Data []CubeItem `json:"data"`
}

type CubeItem struct {
    Code               string      `json:"code"`
    Msg                interface{} `json:"msg"`
    PlayWaySecret      string      `json:"playWaySecret"`
    PlayWayId          int         `json:"playWayId"`
    RestCount          int         `json:"restCount"`
    TotalCount         int         `json:"totalCount"`
    PlayWayStartTime   int64       `json:"playWayStartTime"`
    Title              string      `json:"title"`
    BusinessCouponType int         `json:"businessCouponType"`
    Money              float64     `json:"money"`
    MinConsume         float64     `json:"minConsume"`
    Discount           float64     `json:"discount"`
    MaxMinus           float64     `json:"maxMinus"`
    Extend             struct {
        CouponJumpUrl   string `json:"couponJumpUrl"`
        StudentDiscount string `json:"studentDiscount"`
    } `json:"extend"`
    CouponIcon              interface{} `json:"couponIcon"`
    BgConfig                string      `json:"bgConfig"`
    ChannelUrlKey           string      `json:"channelUrlKey"`
    PlayWayEndTime          int64       `json:"playWayEndTime"`
    CouponType              int         `json:"couponType"`
    PlayWayTitle            string      `json:"playWayTitle"`
    CouponValidityStartTime int64       `json:"couponValidityStartTime"`
    CouponValidityEndTime   int64       `json:"couponValidityEndTime"`
}

func (ci *CubeItem) showInfoHuman() string {
    var sb strings.Builder

    startTime := time.UnixMilli(ci.PlayWayStartTime).Format("2006-01-02 15:04:05")
    endTime := time.UnixMilli(ci.PlayWayEndTime).Format("2006-01-02 15:04:05")
    sb.WriteString(fmt.Sprintf("[%s->%s] %s %s ", startTime, endTime, ci.PlayWaySecret, ci.Code))
    sb.WriteString(fmt.Sprintf("%.1f-%.1f (%d/%d) ", ci.MinConsume, ci.Money, ci.RestCount, ci.TotalCount))
    sb.WriteString(ci.Title)
    sb.WriteString(ci.PlayWayTitle)

    return sb.String()

}

func (ci *CubeItem) toVo() CubeItemVo {
    var civ CubeItemVo
    civ.Title = ci.Title
    civ.Code = ci.Code
    civ.PlayWaySecret = ci.PlayWaySecret
    civ.PlayWayId = ci.PlayWayId
    civ.RestCount = ci.RestCount
    civ.TotalCount = ci.TotalCount
    civ.PlayWayStartTime = time.UnixMilli(ci.PlayWayStartTime)
    civ.PlayWayEndTime = time.UnixMilli(ci.PlayWayEndTime)
    civ.Money = ci.Money
    civ.MinConsume = ci.MinConsume
    civ.CouponValidityStartTime = time.UnixMilli(ci.CouponValidityStartTime)
    civ.CouponValidityEndTime = time.UnixMilli(ci.CouponValidityEndTime)

    return civ
}

type CubeItemVo struct {
    Title                   string    `json:"title"`
    Code                    string    `json:"code"`
    PlayWaySecret           string    `json:"playWaySecret"`
    PlayWayId               int       `json:"playWayId"`
    RestCount               int       `json:"restCount"`
    TotalCount              int       `json:"totalCount"`
    PlayWayStartTime        time.Time `json:"playWayStartTime"`
    PlayWayEndTime          time.Time `json:"playWayEndTime"`
    Money                   float64   `json:"money"`
    MinConsume              float64   `json:"minConsume"`
    CouponValidityStartTime time.Time `json:"couponValidityStartTime"`
    CouponValidityEndTime   time.Time `json:"couponValidityEndTime"`
}

type ActivityData struct {
    ActivityCode string `json:"activityCode"`
    ActivityId   string `json:"activityId"`
}
