package main

import (
	"github.com/766b/mobi"
	"regexp"
	"strconv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"path"
)

var projectPath = flag.String("path", "~/", "base path to root")
var outName = flag.String("out", "book.mobi", "output filename")
var start = flag.Int("start", 0, "start chapter id")
var end = flag.Int("end", -1, "end chapter id")

func main() {
	flag.Parse()

	configPath := path.Join(*projectPath, "project.json")
	rawsPath := path.Join(*projectPath, "raw/")
	fmt.Printf("%s %s %d %d\n", *projectPath, configPath, *start, *end)

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln("error:", err)
	}

	var config Project
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalln("error:", err)
	}
	if *end == -1 {
		*end = len(config.TOC)
	}

	kindle, err := mobi.NewWriter(path.Join(*projectPath, *outName))
	if err != nil {
		log.Fatalln(err)
	}

	chapterSpan := " (" + strconv.Itoa(1+*start) + "-" + strconv.Itoa(*end) + ")"
	kindle.Title(config.Title + chapterSpan)
	kindle.Compression(mobi.CompressionNone)
	kindle.NewExthRecord(mobi.EXTH_DOCTYPE, "EBOK")
	kindle.NewExthRecord(mobi.EXTH_AUTHOR, config.Author)
	kindle.NewExthRecord(mobi.EXTH_TITLE, config.Author)
	kindle.NewExthRecord(mobi.EXTH_PUBLISHINGDATE, config.Author)
	kindle.NewExthRecord(mobi.EXTH_AUTHOR, config.Author)
	kindle.NewExthRecord(mobi.EXTH_AUTHOR, config.Author)
	kindle.NewExthRecord(mobi.EXTH_AUTHOR, config.Author)
	kindle.NewExthRecord(mobi.EXTH_AUTHOR, config.Author)

	for i := *start; i < *end; i++ {
		chapterFile := "chapter-" + strconv.Itoa(i) + ".json"
		chapterPath := path.Join(rawsPath, chapterFile)
		chapterData, err := ioutil.ReadFile(chapterPath)
		if err != nil {
			log.Fatalln("error:", err)
		}

		var blob Chapter
		err = json.Unmarshal(chapterData, &blob)
		if err != nil {
			log.Fatalln("error:", err)
		}

		title := "Chapter " + strconv.Itoa(i+1) + ": " + blob.Chapter.Title
		content := regexp.MustCompile("[\r\n]+").ReplaceAllString(blob.Chapter.Content, "<br>")
		kindle.NewChapter(title, []byte(content))
	}

	kindle.Write()
}