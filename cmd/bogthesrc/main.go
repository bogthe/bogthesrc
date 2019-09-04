package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/bogthe/bogthesrc"
	"github.com/bogthe/bogthesrc/api"
	"github.com/bogthe/bogthesrc/app"
	"github.com/bogthe/bogthesrc/datastore"
	"github.com/bogthe/bogthesrc/importer"
)

type subcmd struct {
	name        string
	description string
	run         func(args []string)
}

var cmds = []subcmd{
	{"serve", "Start a web server", serveCmd},
	{"post", "Create a new post", postCmd},
	{"import", "Import posts", importCmd},
	{"create-db", "Creates the DB", createDbCmd},
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

	datastore.Connect()

	handler := http.NewServeMux()
	handler.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	handler.Handle("/", app.Handler())

	envPort := os.Getenv("PORT")
	if envPort != "" {
		*port = fmt.Sprintf(":%s", envPort)
	}

	createDbCmd(args)
	go importCmd(args)
	err := http.ListenAndServe(*port, handler)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

func createDbCmd(args []string) {
	datastore.Connect()
	datastore.Drop()
	datastore.Create()
}

func postCmd(args []string) {
	client := bogthesrc.NewApiClient(nil)
	post := &bogthesrc.Post{
		Title: "First of it's kind",
		Body:  "I AM THE FIRST OF MY KIND",
		Link:  "https://monzo.com",
	}

	err := client.Posts.Create(post)
	if err != nil {
		log.Fatalf("Failed creating a post %v", err)
	}
}

func importCmd(args []string) {
	var wg sync.WaitGroup

	for _, fetcher := range importer.Fetchers {
		f := fetcher
		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Printf("Starting %s", f.Site())
			errs := importer.Import(f)
			for _, err := range errs {
				log.Println("Failed to import " + err.Error())
			}

			log.Printf("Finished %s", f.Site())
		}()
	}

	wg.Wait()
}
