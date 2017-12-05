package main

import (
	"github.com/766b/mobi"
	"regexp"
	"strings"
	"strconv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"path"
)

var projectPath = flag.String("path", "~/", "base path to root")
var folderName = flag.String("folder", "raw", "relative path to chapter folder")
var outName = flag.String("out", "book.mobi", "output filename")
var start = flag.Int("start", 0, "start chapter id")
var end = flag.Int("end", -1, "end chapter id")

func main() {
	flag.Parse()

	configPath := path.Join(*projectPath, "project.json")
	rawsPath := path.Join(*projectPath, *folderName)
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
	kindle.Compression(mobi.CompressionNone)
	kindle.NewExthRecord(mobi.EXTH_DOCTYPE, "EBOK")
	kindle.NewExthRecord(mobi.EXTH_AUTHOR, config.Author)
	kindle.NewExthRecord(mobi.EXTH_TITLE, config.Title + chapterSpan)
	kindle.Title(config.Title + chapterSpan)

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

		title := blob.Chapter.Title
		if !strings.Contains(title, "Chapter") {
			title = "Chapter " + strconv.Itoa(i+1) + " – " + title
		}
		rows := regexp.MustCompile("[\r\n]+").Split(blob.Chapter.Content, -1)
		content := "<p>" + strings.Join(rows[:], "</p><p>") + "</p>"
		kindle.NewChapter(title, []byte(content))
	}

	kindle.Write()
}