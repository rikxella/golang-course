package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Repository struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	CreatedAt       string `json:"created_at"`
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Printf("Слишком много аргументов. Введите только ссылку на репозиторий.\n")
		return
	}

	pattern := "https://github.com/*/*"
	matched, err := filepath.Match(pattern, args[0])

	if err != nil {
		fmt.Printf("Произошла ошибка: %v\n", err)
		return
	} else if !matched {
		fmt.Printf("Некорректная ссылка: %v\n", args[0])
		return
	}

	segments := strings.Split(strings.Trim(args[0], "/"), "/")
	apiUrl := "https://api.github.com/repos/" + segments[3] + "/" + segments[4]

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Printf("Произошла ошибка: %v\n", err)
		return
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		fmt.Printf("Произошла ошибка: %v\n", err)
		return
	}
	if res.StatusCode > 299 {
		fmt.Printf("Произошла ошибка, статус ответа: %d\n", res.StatusCode)
		return
	}

	var repo Repository
	err = json.Unmarshal(body, &repo)
	if err != nil {
		fmt.Printf("Ошибка декодирования JSON: %v", err)
		return
	}

	fmt.Printf("Имя репозитория: %s\n", repo.Name)
	fmt.Printf("Описание: %s\n", repo.Description)
	fmt.Printf("Количество звезд: %d\n", repo.StargazersCount)
	fmt.Printf("Количество форков: %d\n", repo.ForksCount)

	t, err := time.Parse(time.RFC3339, repo.CreatedAt)
	if err != nil {
		fmt.Printf("Ошибка парсинга даты: %v\n", err)
		return
	}
	fmt.Printf("Дата создания: %s\n", t.Format("02.01.2006 15:04:05"))
}
