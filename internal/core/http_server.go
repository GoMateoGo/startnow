package core

import (
	"fmt"
	"log"
	"net/http"
	"second_hand_mall/internal/global"
	"second_hand_mall/router"
	"sync"
	"time"
)

type http_server interface {
	ListenAndServe() error
}

func RunHttpServer(wg *sync.WaitGroup) {
	defer wg.Done()
	Router := router.InitRouter()

	address := fmt.Sprintf("%s:%s", global.GVAL_CONFIG.HttpServer.Host, global.GVAL_CONFIG.HttpServer.Port)

	s := &http.Server{
		Addr:         address,
		Handler:      Router,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		//MaxHeaderBytes: 1 << 20, 1M
	}
	log.Printf("server starting... bind address:%s", address)
	global.GVAL_LOG.Error(s.ListenAndServe().Error())
}
