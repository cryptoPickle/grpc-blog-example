package main

import (
	"fmt"
	"github.com/cryptoPickle/blog/Server/blog"
	"github.com/cryptoPickle/blog/contract"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
)

func main (){
	app := cli.NewApp()
	app.Name = "Blog Server"
	app.Version = "0.0.1"
	app.Action = start
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "service-port",
			Value: "50052",
		},
		cli.BoolFlag{
			Name: "enable-tls",
		},
		cli.StringFlag{
			Name: "ssl-cert-file",
			Value: "ssl/server.crt",
		},
		cli.StringFlag{
			Name: "ssl-key-file",
			Value: "ssl/server.pem",
		},
		cli.BoolFlag{
			Name: "gen-ssl-cert",
		},
		cli.StringFlag{
			Name: "tls-folder",
			Value: "ssl/",
		},

	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start (c *cli.Context) {
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatal( err )
	}

	opts := []grpc.ServerOption{}

	//Blocks until certs are created.
	if c.Bool("gen-ssl-cert") {
		GenSSL(c.String("tls-folder"))
	}

	if c.Bool("enable-tls") {
		opts = ConfigureSSL(c.String("ssl-cert-file"), c.String("ssl-key-file"))
	}
	s := grpc.NewServer(opts...)
	contract.RegisterBlogServiceServer(s, blog.NewBlogService())


	go func(){
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	log.Info("Server successfully started")
	listenForSIGINT()
}
func ConfigureSSL(certPath, keyPath string) []grpc.ServerOption {
	creds, sslErr := credentials.NewServerTLSFromFile(certPath, keyPath)
	if sslErr != nil {
		log.Fatal(sslErr)
	}
	return []grpc.ServerOption{grpc.Creds(creds)}
}


func listenForSIGINT() {
	c := make (chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(c, os.Interrupt)

	go func () {
		for  range c {
			log.Info("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			done <- true
		}
	}()

	if <- done {
		os.Exit(0)
	}
}

func GenSSL(p string) error {
	log.Info("Generating SSL Cert...")
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		err := filepath.Walk(p, CheckSSLFolder)
		if err != nil {
			log.Warn( err )
			log.Warn("Could not create certs, it exists already. \nRemove ssl directory to recreate...")
			return err
		}
	}
	cmd := exec.Command("make", "generate-ssl")
	if err := cmd.Run(); err != nil {
		log.Warn(err)
		return err
	}
	log.Info("SSL Certs Generated!")
	return nil
}


func CheckSSLFolder(path string, info os.FileInfo, err error) error {
	files := []string{"ca.crt", "ca.key", "server.crt", "server.csr", "server.ket", "server.pem"}

	for _, file := range files {
		if file == info.Name() {
			return errors.New(fmt.Sprintf("%s already exisits skipping...", file))
		}
	}

	return nil

}