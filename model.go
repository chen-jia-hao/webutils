package webutils

import (
    "fmt"
    "github.com/go-resty/resty/v2"
    "github.com/wxpusher/wxpusher-sdk-go"
    "github.com/wxpusher/wxpusher-sdk-go/model"
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

    CronExpr string
}

// 推送接口
type Pusher interface {
    pushMarkdown(title string, content string) (string, error)
}

type WxPusher struct {
    AppToken string
    UId      string
}

func (w *WxPusher) PushMarkdown(title string, content string) (string, error) {
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

func (p *PushPlus) PushMarkdown(title string, content string) (string, error) {
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
