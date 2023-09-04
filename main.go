package main

import (
	"fmt"
	"log"

	"net/http"
)

func main() {
	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := ParseConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Start listening on port %s\nPress Ctrl+C to stop\n", cfg.Web.Port)
	if r := routers(cfg, cfgPath); r != nil {
		log.Fatal("Server exited:", http.ListenAndServe(fmt.Sprintf(":%s", cfg.Web.Port), r))
	}
}
