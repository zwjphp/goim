package middleware

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"net/http"
	"goim/util/vlog"
)

type ResponseWithRecorder struct{
	http.ResponseWriter
	statusCode int
	body bytes.Buffer
}


// 记录日志
func AccessLogging(f http.Handler) http.Handler {
	// 创建一个新的handler包装http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logEntry := vlog.AccessLog.WithFields(logrus.Fields{
			"ip":           r.RemoteAddr,
			"method":       r.Method,
			"path":         r.RequestURI,
			"query":        r.URL.RawQuery,
			"request_body": r.PostForm.Encode(),
		})
		wc := &ResponseWithRecorder{
			ResponseWriter: w,
			statusCode: http.StatusOK,
		}
		f.ServeHTTP(wc, r)
		defer logEntry.WithFields(logrus.Fields{
			"status":         wc.statusCode,
			"response_body":  wc.body.String(),
		}).Info()
	})
}