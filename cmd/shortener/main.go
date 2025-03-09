package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sheinsviatoslav/shortener/internal/cert"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/routes"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	config.Init()

	log.Println("listen on", *config.ServerAddr)
	if *config.EnableHTTPS == "true" {
		if err := cert.CreateTLSCertificate(); err != nil {
			log.Fatal(err)
		}
		log.Fatal(http.ListenAndServeTLS(*config.ServerAddr, cert.CertificateFileName, cert.KeyFileName, routes.MainRouter()))
	} else {
		log.Fatal(http.ListenAndServe(*config.ServerAddr, routes.MainRouter()))
	}

}
