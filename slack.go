package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func commandList() string {
	return `specify command: [help, recipe]`
}

func help() string {
	return `
	help: list commands
	recipe [category]: show recommended recipe
		if does not specify category: list categories.
		specify category: show recommended recipe.
	`
}

func messageHandling(msg string) string {
	log.Printf("[DEBUG] msg: %s", msg)
	word := strings.Split(strings.TrimSpace(msg), " ")
	if len(word) == 1 {
		return commandList()
	}
	switch word[1] {
	case "recipe":
		// TODO: 2patterの返却値を返す

		return "recipe"
	default:
		return help()
	}
}

type resCategory struct {
	Result categoryResult `json:"result"`
}

type categoryResult struct {
	Large  []categories `json:"large"`
	Medium []categories `json:"medium"`
	Small  []categories `json:"small"`
}

type categories struct {
	CategoryID   string `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	CatrgoryURL  string `json:"categoryUrl"`
}

func getCategories(categoryType string) map[string]string {
	resp, err := http.Get(fmt.Sprintf("https://app.rakuten.co.jp/services/api/Recipe/CategoryList/20170426?applicationId=%s&categoryType=%s", env.RakutenID, categoryType))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var categories resCategory
	err = json.Unmarshal(body, &categories)
	if err != nil {
		return nil
	}
	// TODO: correspond to medium, small
	catMap := map[string]string{}
	for _, cat := range categories.Result.Large {
		catMap[cat.CategoryID] = cat.CategoryName
	}
	return catMap
}

func choiceRecipe(categoryID int) string {

	return ""
}
