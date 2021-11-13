package model

type Stories struct {
	Id          int    `json:"id"`
	Author      string `json:"by"`
	Descendants int    `json:"descendants"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"tittle"`
	Type        string `json:"type"`
	Url         string `json:"url"`
}
