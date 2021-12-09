package middleware

import (
	"api-service/pkg/log"
	"api-service/util"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"time"

	"api-service/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/willf/pad"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var regex = regexp.MustCompile("(/.*/user|/login)")

// Logging is a middleware function that logs the request.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path
		if !regex.MatchString(path) {
			return
		}

		// Skip for the health check requests.
		if path == "/monitor/health" || path == "/monitor/ram" || path == "/monitor/cpu" || path == "/monitor/disk" {
			return
		}

		// Read the Body content：读取后会被置空，因此需要重新赋值。
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// The basic information.
		method := c.Request.Method
		ip := c.ClientIP()

		log.Debugf("New request come in, path: %s, Method: %s, body `%s`", path, method, string(bodyBytes))
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// Continue.
		c.Next()

		// Calculates the latency.
		latency := time.Now().UTC().Sub(start)

		// get code and message
		// 通过重定向 HTTP 的 Response 到指定的 IO 流截获。
		code, message := -1, ""
		var response util.Response
		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			log.Errorf(err, "response body can not unmarshal to model.Response struct, body: `%s`", blw.body.Bytes())
			code = errno.InternalServerError.Code
			message = err.Error()
		} else {
			code = response.Code
			message = response.Message
		}

		log.Infof("%-13s | %-12s | %s %s | {code: %d, message: %s}", latency, ip, pad.Right(method, 5, ""), path, code, message)
	}
}
