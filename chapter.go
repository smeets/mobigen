package main

type Chapter struct {
	Chapter ChapterInfo `json:"chapterInfo"`
}

type ChapterInfo struct {
	Id string `json:"chapterId"`
	Title string `json:"chapterName"`
	Index int `json:"chapterIndex"`
	Content string `json:"content"`
}
