package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/Sirupsen/logrus"
)

var (
	isFile bool
)

func init() {
	// parse flags
	flag.BoolVar(&isFile, "r", false, "read from file")
	flag.Parse()

	if len(flag.Args()) < 1 {
		usageAndExit(1, "must pass a version string or version file\nex. %s v0.1.0", os.Args[0])
	}
}

func main() {
	version := flag.Arg(0)
	if isFile {
		b, err := ioutil.ReadFile(version) // just pass the file name
		if err != nil {
			logrus.Fatal(err)
		}
		version = string(b)
		re := regexp.MustCompile(`\r?\n`)
		version = re.ReplaceAllString(version, "")
	}

	if len(flag.Args()) > 1 {
		for _, file := range flag.Args()[1:] {
			writeOut(version, file)
		}
		return
	}
	fmt.Printf("%s", version)
}

func writeOut(version, file string) {
	re := regexp.MustCompile(`v?\d*\.\d*\.\d*`)

	input, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Fatal(err)
	}

	output := re.ReplaceAllString(string(input), version)

	err = ioutil.WriteFile(file, []byte(output), 0666)
	if err != nil {
		log.Fatalln(err)
	}
}

func usageAndExit(exitCode int, message string, args ...interface{}) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message, args...)
		fmt.Fprint(os.Stderr, "\n\n")
	}
	fmt.Println("vupdate <version|version file> <file to update>...")
	flag.Usage()
	fmt.Fprintln(os.Stderr, "")
	os.Exit(exitCode)
}
