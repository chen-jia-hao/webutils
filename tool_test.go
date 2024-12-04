package webutils

import (
	"github.com/mileusna/useragent"
	"testing"
)

func TestParseUA(t *testing.T) {
	u := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 TitansX/20.0.1.old KNB/1.0 iOS/17.4.1 meituangroup/com.meituan.imeituan/12.26.406 meituangroup/12.26.406 App/10110/12.26.406 iPhone/iPhoneSE2 WKWebView"
	userAgent := useragent.Parse(u)

	t.Log(userAgent)
	ua := ParseUA(u)

	t.Log(ua)
}
