package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"text/template"
)

var flagVersion bool
var fileName string

func main() {
	rootCmd := flag.NewFlagSet("Root", flag.ContinueOnError)
	rootCmd.BoolVar(&flagVersion, "v", false, "print version")
	rootCmd.BoolVar(&flagVersion, "version", false, "print version")
	addCmd := flag.NewFlagSet("Add", flag.ContinueOnError)
	addCmd.StringVar(&fileName, "name", time.Now().Format("2006-01-02")+".md", "File name")

	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}
	if flagVersion {
		fmt.Println("v0.0.1")
	}

	// Handle sub commands
	var err error
	args := rootCmd.Args()
	if len(args) > 0 {
		switch args[0] {
		case "add":
			_ = addCmd.Parse(args[1:])
			err = handleAddCmd(fileName)
			fmt.Println(fileName)
		default:
			fmt.Printf("Unknown command: %v\n", args[1:])
			os.Exit(2)
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}

func handleAddCmd(filename string) error {
	filePath := fmt.Sprintf("./templates/report.md.tmpl")
	byteTmpl, _ := ioutil.ReadFile(filePath)
	stringTmpl := string(byteTmpl)

	tmpl := template.Must(template.New("report").Parse(stringTmpl))
	// Todayを差し込む
	reportFile, _ := os.Create(fileName + ".md")
	reportMeta := struct {
		Today string
	}{
		Today: time.Now().Format("2006-01-02"),
	}
	// text/templateとhtml/templateで挙動が違うので注意
	_ = tmpl.Execute(reportFile, reportMeta)

	return nil
}
