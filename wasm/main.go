package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"syscall/js" // wasm利用で必要なパッケージ
)

func main() {
	c := make(chan struct{}, 0)                // チャネル呼び出しはお作法
	js.Global().Set("getIp", js.FuncOf(GetIp)) // JS側で呼び出すための関数登録
	<-c
}

func GetIp(_ js.Value, _ []js.Value) interface{} {
	go func() { // HTTPリクエストを送信する場合は、goroutine化する必要がある
		c := http.Client{Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return nil, fmt.Errorf("あえてエラーにします")
			},
		}}

		resp, err := c.Get("https://api.ipify.org/")
		if err != nil {
			appendHTMLBody(fmt.Sprintf("http get: %s", err))
			return
		}
		defer resp.Body.Close()

		b := bytes.NewBuffer(nil)
		if _, err = io.Copy(b, resp.Body); err != nil {
			appendHTMLBody(fmt.Sprintf("read body: %s", err))
			return
		}
		appendHTMLBody(b.String())
	}()
	return "OK"
}

func appendHTMLBody(s string) {
	var document = js.Global().Get("document")
	var p = document.Call("createElement", "p")
	p.Set("textContent", s)
	document.Get("body").Call("appendChild", p)
}
