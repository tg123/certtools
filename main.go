//go:build windows
// +build windows

package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/tg123/certstore"
	"github.com/urfave/cli/v2"
)

var mainver string = "(devel)"

func version() string {

	var v = mainver

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return v
	}

	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			v = fmt.Sprintf("%v, %v", v, s.Value[:9])
		case "vcs.time":
			v = fmt.Sprintf("%v, %v", v, s.Value)
		}
	}

	v = fmt.Sprintf("%v, %v", v, bi.GoVersion)

	return v
}

func openstore(c *cli.Context) (certstore.Store, error) {

	name := c.String("store-name")
	loc := c.String("store-location")

	switch loc {
	case "local-machine":
		return certstore.OpenStoreWindows(name, certstore.StoreLocationLocalMachine)
	case "current-user":
		return certstore.OpenStoreWindows(name, certstore.StoreLocationCurrentUser)
	default:
		return nil, fmt.Errorf("unsupported store %v", loc)
	}
}

func main() {

	app := &cli.App{
		Usage:   "The missing cert management tools for windows nano",
		Version: version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "store-name",
				Value: "MY",
			},
			&cli.GenericFlag{
				Name: "store-location",
				Value: &EnumValue{
					Enum:    []string{"current-user", "local-machine"},
					Default: "current-user",
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "import",
				Usage:     "import a pfx to store",
				ArgsUsage: "<path/to/pfx>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "path to pfx",
					},
					&cli.StringFlag{
						Name:    "password",
						Aliases: []string{"p"},
						Usage:   "password to pfx",
					},
				},
				Action: func(c *cli.Context) error {
					pfx := c.Args().First()

					if pfx == "" {
						pfx = c.String("file")
					}
					if pfx == "" {
						return cli.ShowSubcommandHelp(c)
					}

					b, err := ioutil.ReadFile(pfx)
					if err != nil {
						return err
					}

					store, err := openstore(c)
					if err != nil {
						return err
					}

					return store.Import(b, c.String("password"))
				},
			},
			{
				Name:  "ls",
				Usage: "list certificates in store",
				Action: func(c *cli.Context) error {

					store, err := openstore(c)
					if err != nil {
						return err
					}

					certs, err := store.Identities()
					if err != nil {
						return err
					}

					for _, cert := range certs {
						p, err := cert.Certificate()
						if err != nil {
							log.Printf("error : %v", err)
							continue
						}

						fmt.Printf("%x %v\n", sha1.Sum(p.Raw), p.Subject.CommonName)
					}

					return nil
				},
			},
			// {
			// 	Name:    "export",
			// 	Usage:   "export certificate to pfx from store",
			// 	Action: func(c *cli.Context) error {
			// 		return nil
			// 	},
			// },
			{
				Name:      "rm",
				Usage:     "remove a certificate from store",
				ArgsUsage: "<thumbprint>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "thumbprint",
						Aliases:  []string{"t"},
						Usage:    "thumbprint of the certificate to be deleted",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					thumb := c.Args().First()

					if thumb == "" {
						thumb = c.String("thumbprint")
					}
					if thumb == "" {
						return cli.ShowSubcommandHelp(c)
					}

					thumb = strings.ToLower(thumb)

					store, err := openstore(c)
					if err != nil {
						return err
					}

					certs, err := store.Identities()
					if err != nil {
						return err
					}

					for _, cert := range certs {
						p, err := cert.Certificate()
						if err != nil {
							log.Printf("error : %v", err)
							continue
						}

						if strings.ToLower(fmt.Sprintf("%x", sha1.Sum(p.Raw))) == thumb {
							return cert.Delete()
						}

					}

					log.Printf("%v not found", thumb)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
