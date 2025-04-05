package main

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sheinsviatoslav/shortener/internal/cert"
	"github.com/sheinsviatoslav/shortener/internal/config"
	"github.com/sheinsviatoslav/shortener/internal/grpcserv"
	"github.com/sheinsviatoslav/shortener/internal/routes"
	pb "github.com/sheinsviatoslav/shortener/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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

	go func() {
		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatal(err)
		}
		s := grpc.NewServer()
		pb.RegisterUrlsServer(s, &grpcserv.UrlsServer{})

		fmt.Println("Сервер gRPC начал работу")
		if err := s.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	var srv = http.Server{
		Addr:    *config.ServerAddr,
		Handler: routes.MainRouter(),
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigs
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("Server Shutdown error: %v", err)
		}
	}()

	log.Println("listen on", *config.ServerAddr)
	if *config.EnableHTTPS == "true" {
		if err := cert.CreateTLSCertificate(); err != nil {
			log.Fatal(err)
		}
		if err := srv.ListenAndServeTLS(cert.CertificateFileName, cert.KeyFileName); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTPS server error: %v", err)
		}
	} else {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}
}
