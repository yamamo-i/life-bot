package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func commandList() string {
	return `specify command: [help, recipe]`
}

func help() string {
	return `
help : list commands
recipe [category]: show recommended recipe
					if does not specify category: list categories.
					specify category: show recommended recipe.
	`
}

func messageHandling(msg string, typ string) string {
	log.Printf("[DEBUG] msg: %s, type: %s", msg, typ)
	words := strings.Split(strings.TrimSpace(msg), " ")
	wakeupKeywords := []string{env.BotName, "<@" + env.BotID + ">"}
	for _, keyWord := range wakeupKeywords {
		if words[0] == keyWord {
			return _messageHandling(words[1:])
		}
	}
	return ""
}

func _messageHandling(words []string) string {
	if len(words) == 0 {
		return commandList()
	}
	switch words[0] {
	case "recipe":
		var result string
		// TODO: 2patterの返却値を返す
		catMap, _ := getCategories("large")
		if len(words) == 1 {
			keyList := make([]int, len(catMap))
			for key := range catMap {
				keyList = append(keyList, key)
			}
			for key := range keyList {
				if val, ok := catMap[key]; ok {
					result += val + "(" + strconv.Itoa(key) + ")\n"
				}
			}
		} else if len(words) == 2 {
			var catID int
			var err error
			var ok bool
			var recipeList []recipeSummary
			if catID, err = strconv.Atoi(words[1]); err == nil {
				if _, ok := catMap[catID]; ok {
					recipeList, _ = choiceRecipe(catID)
				} else {
					return "category id is not found."
				}
			} else if catID, ok = contains(catMap, words[1]); ok {
				recipeList, _ = choiceRecipe(catID)
			} else if ok == false {
				return "category is not found."
			}
			for _, recipe := range recipeList {
				result += recipe.Title + ": " + recipe.URL + "\n"
			}
		} else {
			return "illeagal command. recipe [categoryID]."
		}
		return result
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

type resRecipe struct {
	Result []recipeResult `json:"result"`
}

type recipeResult struct {
	RecipeTitle string `json:"recipeTitle"`
	RecipeURL   string `json:"recipeUrl"`
}

func getCategories(categoryType string) (map[int]string, error) {
	resp, err := http.Get(fmt.Sprintf("https://app.rakuten.co.jp/services/api/Recipe/CategoryList/20170426?applicationId=%s&categoryType=%s", env.RakutenID, categoryType))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var categories resCategory
	err = json.Unmarshal(body, &categories)
	if err != nil {
		return nil, err
	}
	// TODO: correspond to medium, small
	catMap := map[int]string{}
	for _, cat := range categories.Result.Large {
		_id, _ := strconv.Atoi(cat.CategoryID)
		catMap[_id] = cat.CategoryName
	}
	return catMap, nil
}

type recipeSummary struct {
	Title string
	URL   string
}

func choiceRecipe(categoryID int) ([]recipeSummary, error) {

	// TODO: 共通化したい
	resp, err := http.Get(fmt.Sprintf("https://app.rakuten.co.jp/services/api/Recipe/CategoryRanking/20170426?applicationId=%s&categoryId=%d", env.RakutenID, categoryID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var recipes resRecipe
	err = json.Unmarshal(body, &recipes)
	if err != nil {
		return nil, err
	}
	recipeList := []recipeSummary{}
	for _, recipe := range recipes.Result {
		result := recipeSummary{Title: recipe.RecipeTitle, URL: recipe.RecipeURL}
		recipeList = append(recipeList, result)
	}
	return recipeList, nil
}

func contains(s map[int]string, e string) (int, bool) {
	for k, v := range s {
		if e == v {
			return k, true
		}
	}
	return 0, false
}
