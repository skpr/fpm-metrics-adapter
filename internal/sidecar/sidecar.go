package sidecar

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heptio/workgroup"

	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
)

// Current status cached from the last successful query.
var currentStatus fpm.Status

// Start the sidecar server.
func Start(srvAddr, srvPath, fpmAddr string, refresh time.Duration) error {
	var wg workgroup.Group

	wg.Add(func(stop <-chan struct{}) error {
		l, err := net.Listen("tcp", srvAddr)
		if err != nil {
			return err
		}

		go func() {
			<-stop
			l.Close()
		}()

		router := gin.Default()

		router.GET(srvPath, func(c *gin.Context) {
			c.JSON(http.StatusOK, currentStatus)
		})

		return http.Serve(l, router)

	})

	wg.Add(func(stop <-chan struct{}) error {
		ticker := time.NewTicker(refresh)

		for {
			select {
			case <-stop:
				return nil
			case <-ticker.C:
				s, err := fpm.QueryStatus(fpmAddr)
				if err != nil {
					log.Println(err)
					continue
				}

				currentStatus = s
			}
		}
	})

	return wg.Run()
}
