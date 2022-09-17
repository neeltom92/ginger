// https://github.com/imtoori/gin-redis-ip-limiter
// https://le-gall.bzh/post/go/a-reverse-proxy-in-go-using-gin/

package main

import (
	"fmt"
	"github.com/Salvatore-Giordano/gin-redis-ip-limiter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
  "strconv"

)

func main() {

	/*
	   replace each of the variables according to your requirements

	   host: is the API endpoint to which you want to Proxy pass your requests to
	   requestCount: number of request thats allowed
	   secondsTime: time within which the requests will be throttled
	   redisHost: redis host, replace "127.0.0.1" with the actual redis host
     port: the port you want the application listen to
	   here , I throttle all requests if it cross more than "10 requests within 60 seconds"

	*/

	host := "https://vpic.nhtsa.dot.gov"
	requestCount := 10
	const secondsTime = 60
	redisHost := "127.0.0.1:6379"
	port := 8000

	r := gin.Default()

	r.Use(iplimiter.NewRateLimiterMiddleware(redis.NewClient(&redis.Options{
		// replace the "127.0.0.1" with the actual redis endpoint

		Addr:     redisHost,
		Password: "",
		DB:       1,
	}), "general", requestCount, secondsTime*time.Second))
	// 10 requests within 1 min

	r.Any("/*proxyPath", func(c *gin.Context) {

		remote, err := url.Parse(host)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		//Define the director func
		//This is a good place to log, for example
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Param("proxyPath")

			if req.URL.Path == "/api/vehicles/getallmanufacturers" {
				fmt.Println(req.URL.Path)
			}
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	})
  portmapping := strconv.Itoa(port)

	r.Run(":"+portmapping)
}
