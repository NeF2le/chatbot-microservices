package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(baseURL string) gin.HandlerFunc {
	target, err := url.Parse(baseURL)
	if err != nil {
		panic(fmt.Sprintf("invalid base URL %q: %s", baseURL, err))
	}

	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(target)

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host
			req.URL.Path = c.Param("proxyPath")
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
