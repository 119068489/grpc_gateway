package easygo

import (
	"net/http"

	"github.com/astaxie/beego/logs"
)

func ServerRun(addr string, handler http.Handler, msg ...string) {
	msg = append(msg, "")
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	logs.Info("%s start to listen %s", msg[0], addr)

	err := httpServer.ListenAndServe()
	PanicError(err)
}
