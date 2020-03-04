package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Recipe関係のmsgをbotが取得した時に稼働するメソッド
func getRecipeResponse(words []string) (string, error) {
	var result string
	// カテゴリのID: 名前は共通で必要なので取得する
	catMap, _ := getCategories("large")

	// botname recipe: がコールされた時の挙動
	if len(words) == 0 {
		keyList := make([]int, len(catMap))
		for key := range catMap {
			keyList = append(keyList, key)
		}
		for key := range keyList {
			if val, ok := catMap[key]; ok {
				result += val + "(" + strconv.Itoa(key) + ")\n"
			}
		}
		// botname recipe [categoryID|categoryName]: がコールされた時の挙動
	} else if len(words) == 1 {
		var catID int
		var err error
		var ok bool
		var recipeList []recipeSummary
		if catID, err = strconv.Atoi(words[0]); err == nil {
			_, ok = catMap[catID]
		} else {
			catID, ok = getKeyFromValue(catMap, words[0])
		}
		if ok {
			recipeList, _ = getRecipe(catID)
		} else {
			return "", fmt.Errorf("Error: %s", "category is not found.\n")
		}
		for _, recipe := range recipeList {
			result += recipe.Title + ": " + recipe.URL + "\n"
		}
	} else {
		return "", fmt.Errorf("Error: %s", "illeagal command. recipe [categoryID].\n")
	}
	return result, nil
}

// Category取得APIの返却レスポンスの構造体(第1階層)
// https://webservice.rakuten.co.jp/api/recipecategorylist/
type resCategory struct {
	Result categoryResult `json:"result"`
}

// category取得APIの返却レスポンスの構造体(第2階層)
type categoryResult struct {
	Large  []categories `json:"large"`
	Medium []categories `json:"medium"`
	Small  []categories `json:"small"`
}

// category取得APIの返却レスポンスの構造体(第3階層)
type categories struct {
	CategoryID   string `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	CatrgoryURL  string `json:"categoryUrl"`
}

// Recipeのcategoryの一覧を取得するメソッド
// categoryTypeごとに変化させたいがmedium,smallは数が多すぎるので現在未対応
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

// Recipe取得APIの返却レスポンスの構造体(第1階層)
// https://webservice.rakuten.co.jp/api/recipecategoryranking/
type resRecipe struct {
	Result []recipeResult `json:"result"`
}

// Recipe取得APIの返却レスポンスの構造体(第2階層)
type recipeResult struct {
	RecipeTitle string `json:"recipeTitle"`
	RecipeURL   string `json:"recipeUrl"`
}

// 取得したレシピの返却構造体
type recipeSummary struct {
	Title string
	URL   string
}

// CategoryIDに対応したおすすめRecipeを取得する
func getRecipe(categoryID int) ([]recipeSummary, error) {

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

// Mapのvalue値からkey値を逆引きするメソッド.
func getKeyFromValue(m map[int]string, value string) (int, bool) {
	for k, v := range m {
		if value == v {
			return k, true
		}
	}
	return 0, false
}
