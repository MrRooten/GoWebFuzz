package cmdline

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gowebfuzz/library/utils"
	"os"
)

func Execute() {
	cmd := &cli.App{
		//Global flag
		Flags: []cli.Flag {

		},
		//Comands
		Commands: []*cli.Command{
			{
				Name: "webproxy",
				Usage: "Use web proxy to detect target website",
				Action: func(context *cli.Context) error {
					fmt.Println(context.String("config-file"))
					return nil
				},
				Flags: Flags,
			},
		},
	}

	err := cmd.Run(os.Args)
	logger := &utils.Logger{}
	if err != nil {
		logger.Error(err.Error())
		return
	}
}