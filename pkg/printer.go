package gitlabsanitycli

import (
	"log"
	"os"
	"text/template"
)

// Printer interface for rendering output
type Printer interface {
	Print(rc <-chan content)
}

// TemplatePrinter implementation of Printer for templated output (see ../config/stdout.tpl)
type TemplatePrinter struct {
}

// Print get results from egress channels and prints to stdout by using a template
func (p TemplatePrinter) Print(rc <-chan content) {
	tpl, err := template.ParseGlob("config/*")
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for r := range rc {
			err := tpl.ExecuteTemplate(os.Stdout, "stdout.tpl", r)
			if err != nil {
				log.Panic(err)
			}
		}
	}()
}
