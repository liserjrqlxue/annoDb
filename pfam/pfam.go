package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/liserjrqlxue/simple-util"
	"log"
	"os"
	"regexp"
	"strings"
)

// flag
var (
	input = flag.String(
		"input",
		"",
		"input Pfam-A.full.ncbi.gz",
	)
	output = flag.String(
		"output",
		*input+".HomoSapiens.txt",
		"out put Pfam domain db",
	)
)

// regexp
var (
	isAC          = regexp.MustCompile(`^#=GF\s+AC\s+(\S+)`)
	isDE          = regexp.MustCompile(`^#-GF\s+DE\s+(.+)`)
	isGS          = regexp.MustCompile(`^#=GS`)
	isHomoSapiens = regexp.MustCompile(`\[Homo sapiens\]`)
	isProtainPos  = regexp.MustCompile(`^#=GS\s+(\S+)/(\d+)-(\d+)\s+DE`)
)

func main() {
	flag.Parse()
	if *input == "" {
		flag.Usage()
		os.Exit(1)
	}
	file, err := os.Open(*input)
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(file)

	gr, err := gzip.NewReader(file)
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(gr)

	out, err := os.Create(*output)
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(out)
	fmt.Fprintln(out, strings.Join([]string{"#Protain", "Start", "End", "Accession", "Definition"}, "\t"))

	scanner := bufio.NewScanner(gr)

	var protain, start, end, accession, definition string

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case isAC.MatchString(line):
			matchs := isAC.FindStringSubmatch(line)
			if matchs != nil {
				accession = matchs[1]
			}
		case isDE.MatchString(line):
			matchs := isDE.FindStringSubmatch(line)
			if matchs != nil {
				definition = matchs[1]
			}
		case isGS.MatchString(line) && isHomoSapiens.MatchString(line):
			matchs := isProtainPos.FindStringSubmatch(line)
			if matchs != nil && len(matchs) == 4 {
				protain = matchs[1]
				start = matchs[2]
				end = matchs[3]
				fmt.Fprintln(out, strings.Join([]string{protain, start, end, accession, definition}, "\t"))
			} else {
				log.Fatalf("can not parser:[%s]\tmatchs:[%v]\n", line, matchs)
			}
		default:

		}
	}
	log.Println(scanner.Err())
}
