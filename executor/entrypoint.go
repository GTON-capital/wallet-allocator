package executor

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
	"github.com/valyala/fasthttp"
)

type InputConfig struct {
	PKPath string
	Port int
}

func Run() {
	inputCfg := &InputConfig{}

	app := &cli.App{
		Name:  "Wallet alloc for solana",
		Usage: "Wallet alloc for solana",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Value:       8081,
				Usage:       "port for allocator",
				Destination: &inputCfg.Port,
			},
			&cli.StringFlag{
				Name:        "keypair",
				Value:       "",
				Usage:       "allocator private key path",
				Destination: &inputCfg.PKPath,
			},
		},
		Action: func(c *cli.Context) error {
			allocator, err := NewAllocator(inputCfg.PKPath)
			if err != nil {
				debug.PrintStack()
				return err
			}

			fasthttp.ListenAndServe(fmt.Sprintf(":%v", inputCfg.Port), allocator.RequestHandle)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
