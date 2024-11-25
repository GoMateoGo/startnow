package core

import (
	"fmt"
	"log"
	"net/http"
	"second_hand_mall/internal/global"
	"second_hand_mall/router"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	Router := router.InitRouter()

	address := fmt.Sprintf("%s:%s", global.GVAL_CONFIG.Server.Host, global.GVAL_CONFIG.Server.Port)

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
