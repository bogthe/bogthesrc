package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bogthe/bogthesrc/app"
)

type subcmd struct {
	name        string
	description string
	run         func(args []string)
}

var cmds = []subcmd{
	{"serve", "Start a web server", serveCmd},
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `bogthesrc is web news and link server. (copy of thesrc)
Usage
	bogthesrc [options] command [arg...]

The commands are:
`)

		for _, c := range cmds {
			fmt.Fprintf(os.Stderr, "	%-24s %s\n", c.name, c.description)
		}

		fmt.Fprintln(os.Stderr, `
Use "bogthesrc command -s" for information about a command.

The options are:
`)

		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}

	log.SetFlags(0)

	cmd := flag.Arg(0)
	for _, c := range cmds {
		if cmd == c.name {
			c.run(flag.Args()[1:])
			return
		}
	}

	fmt.Fprintln(os.Stderr, "Failed to pass a known command")
	os.Exit(1)
}

func serveCmd(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.String("port", ":5000", "Port to listen to")
	staticDir := fs.String("staticDir", "./static", "Static assets directory")
	templateDir := fs.String("templateDir", "./tmpl", "Templates directory")
	reloadTemplates := fs.Bool("reloadTemplates", true, "Reloads templates on every request")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage bogthesrc serve [options]

Starts a web server with the specified options

The options are:
`)

		fs.PrintDefaults()
		os.Exit(1)
	}

	fs.Parse(args)
	if fs.NArg() != 0 {
		fs.Usage()
	}

	app.ReloadTemplates = *reloadTemplates

	static, errStatic := filepath.Abs(*staticDir)
	tmpl, errTemplate := filepath.Abs(*templateDir)
	if errStatic != nil || errTemplate != nil {
		fs.Usage()
	}

	app.StaticDir = static
	app.TemplateDir = tmpl

	handler := http.NewServeMux()
	handler.Handle("/", app.Handler())

	log.Print("Listening on ", *port)
	err := http.ListenAndServe(*port, handler)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}
