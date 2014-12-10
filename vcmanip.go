package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	cardExp  *regexp.Regexp = regexp.MustCompile("(?sU)BEGIN:VCARD.*END:VCARD")
	nameProp *regexp.Regexp = regexp.MustCompile("FN:(.*)")
)

func main() {
	split := flag.Bool("s", false, "Split a monolithic vCard.")
	merge := flag.Bool("m", false, "Merge a directory of vCards.")
	input := flag.String("i", "", "vCard file or directory.")
	output := flag.String("o", "", "Output directory.")
	flag.Parse()

	if *input == "" {
		log.Fatal("Didn't specify an input.")
	}

	if *split {
		splitCard(*input, *output)
	} else if *merge {
		mergeCards(*input, *output)
	} else {
		log.Fatal("Didn't specify split or merge flag.")
	}
}

func mergeCards(directory, outputFile string) {
	vCardPaths := []string{}
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, _ error) error {
			if info.IsDir() {
				return nil
			}

			vCardPaths = append(vCardPaths, path)
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	if outputFile == "" {
		outputFile = "monolith.vcf"
	}

	mergeFile, err := os.OpenFile(outputFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer mergeFile.Close()

	for _, p := range vCardPaths {
		fBytes, err := ioutil.ReadFile(p)
		if err != nil {
			log.Fatal(err)
		}

		_, err = mergeFile.Write(fBytes)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func splitCard(filename, outputDir string) {
	fBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	vCards := [][]byte{}
	vcFileNames := []string{}

	for _, match := range cardExp.FindAll(fBytes, -1) {
		vCards = append(vCards, match)
	}

	if outputDir == "" {
		var clean string
		for _, match := range nameProp.FindAllSubmatch(fBytes, -1) {
			clean = string(match[1])
			clean = clean[:len(clean)-1]
			vcFileNames = append(vcFileNames, fmt.Sprintf("%s.vcf", clean))
		}
	} else {
		var clean string
		for _, match := range nameProp.FindAllSubmatch(fBytes, -1) {
			clean = string(match[1])
			clean = clean[:len(clean)-1]
			clean = filepath.Join(outputDir, fmt.Sprintf("%s.vcf", clean))
			vcFileNames = append(vcFileNames, clean)
		}
	}

	if len(vCards) != len(vcFileNames) {
		log.Fatal("Parsing error")
	}

	for i, name := range vcFileNames {
		if err := ioutil.WriteFile(name, vCards[i], 0644); err != nil {
			log.Fatal(err)
		}
	}
}
