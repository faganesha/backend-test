package controller

import (
	"backend-test/config"
	"backend-test/connection"
	"backend-test/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func getIdTopStories() []int {
	db := connection.Connection()
	defer db.Close()

	url := fmt.Sprintf("%s%s", config.HackerNewsApiBasePath, config.HackerNewsApiTopStoryPAth)
	response, err := http.Get(url)
	idsTopStory := []int{}

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(responseData, &idsTopStory)
	return idsTopStory
}

func InsertTopStories(res http.ResponseWriter, req *http.Request) {
	db := connection.Connection()
	defer db.Close()

	idsStories := getIdTopStories()

	deleteQuery := "DELETE FROM `stories` WHERE id IN"
	insertQuery := "INSERT INTO stories (id) VALUES "

	var (
		insert, delete string
	)

	for i := 0; i < len(idsStories); i++ {
		insert += fmt.Sprintf("(%d), ", idsStories[i])
		delete += fmt.Sprintf("%d, ", idsStories[i])
	}

	insert = strings.TrimSuffix(insert, ", ")
	delete = strings.TrimSuffix(delete, ", ")

	fullQueryInsert := fmt.Sprintf("%s%s", insertQuery, insert)
	fullQueryDelete := fmt.Sprintf("%s (%s)", deleteQuery, delete)

	db.Exec(fullQueryDelete)
	responseDbSelect, err := db.Exec(fullQueryInsert)
	if err != nil {
		panic(err.Error())
	}

	lastId, err := responseDbSelect.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Last id: %d\n", lastId)
}

func InsertTopStoryDetail(res http.ResponseWriter, req *http.Request) {
	db := connection.Connection()
	defer db.Close()

	idsStories := getIdTopStories()

	for i := 0; i < len(idsStories); i++ {
		HackerNewsApiStoryPAth := fmt.Sprintf("/v0/item/%d.json?print=pretty", idsStories[i])
		urlDetail := fmt.Sprintf("%s%s", config.HackerNewsApiBasePath, HackerNewsApiStoryPAth)
		response, err := http.Get(urlDetail)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			log.Fatal(err)
		}

		var detailStory model.Stories
		json.Unmarshal(responseData, &detailStory)

		kids, _ := json.Marshal(detailStory.Kids)
		kidsStr := strings.Trim(string(kids), "[]")

		updateQueries := fmt.Sprintf(
			`UPDATE stories
			SET
			author='%s',
			kids='%s',
			descendants=%d,
			score=%d,
			time=%d,
			tittle='%s',
			type='%s',
			url='%s'
			WHERE id = %d`,

			detailStory.Author,
			kidsStr,
			detailStory.Descendants,
			detailStory.Score,
			detailStory.Time,
			detailStory.Title,
			detailStory.Type,
			detailStory.Url,
			detailStory.Id)

		_, err = db.Exec(updateQueries)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(updateQueries)
	}
}

func GetTopStories(res http.ResponseWriter, req *http.Request) {

	url := fmt.Sprintf("%s%s", config.HackerNewsApiBasePath, config.HackerNewsApiTopStoryPAth)
	response, _ := http.Get(url)
	responseData, _ := ioutil.ReadAll(response.Body)
	var responseJson []int
	json.Unmarshal(responseData, &responseJson)
	fmt.Println(string(responseData))
	if err := json.NewEncoder(res).Encode(responseJson); err != nil {
		fmt.Println(err.Error(), nil)
	}
}

func GetDetailStories(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json?print=pretty", ps.ByName("story"))
	response, _ := http.Get(url)
	responseData, _ := ioutil.ReadAll(response.Body)
	var detailStory model.Stories
	json.Unmarshal(responseData, &detailStory)

	if err := json.NewEncoder(res).Encode(detailStory); err != nil {
		fmt.Println(err.Error(), nil)
	}

	fmt.Print(url)
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
