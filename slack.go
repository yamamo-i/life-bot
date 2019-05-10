package main

import (
	"log"
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

// Message_typeのアクションを取得した時に動作するmethod
func handleMsg(msg string) string {
	log.Printf("[DEBUG] msg: %s", msg)
	words := strings.Split(strings.TrimSpace(msg), " ")
	wakeupKeywords := []string{env.BotName, "<@" + env.BotID + ">"}
	for _, keyWord := range wakeupKeywords {
		if words[0] == keyWord {
			return _handleMsg(words[1:])
		}
	}
	if len(words) == 0 {
		return commandList()
	}
	return ""
}

func _handleMsg(words []string) string {
	switch words[0] {
	case "recipe":
		result, err := getRecipeResponse(words[1:])
		if err != nil {
			log.Printf("error has occured: %#v\n", err)
			return err.Error()
		}
		return result
	default:
		return help()
	}
}
