package cmdline

import "github.com/urfave/cli/v2"

var Flags []cli.Flag
func init() {
	Flags = []cli.Flag {
		&cli.StringFlag{
			Name: "address",
			Aliases: []string{"a"},
			Value: "",
			Usage: "Set the listen address",
		},
		&cli.StringFlag{
			Name: "config-file",
			Aliases: []string{"cf"},
			Value: "",
			Usage: "Set the config file",
		},
		&cli.StringFlag{
			Name: "fuzz",
			Aliases: []string{"f"},
			Value:"",
			Usage: "Set the fuzz module",
		},
		&cli.StringFlag{
			Name:"fuzz-file",
			Aliases: []string{"ff"},
			Value:"",
			Usage: "Set the fuzz module file",
		},
	}
}