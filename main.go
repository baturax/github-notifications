package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Type0Diabet struct {
	Unread     bool       `json:"unread"`
	Reason     string     `json:"reason"`
	UpdatedAt  string     `json:"updated_at"`
	Subject    Subject    `json:"subject"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	FullName string `json:"full_name"`
}

type Subject struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Type  string `json:"type"`
}

func main() {

	home := os.Getenv("HOME")
	cache := filepath.Join(home, ".cache")
	tokenPath := filepath.Join(cache, "token")
	tokenByte, err := os.ReadFile(tokenPath)
	Err(err)

	token := strings.TrimSpace(string(tokenByte))

	req, err := http.NewRequest("GET", "https://api.github.com/notifications", nil)
	Err(err)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	Err(err)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalln(resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	Err(err)

	var nots []Type0Diabet
	err = json.Unmarshal(data, &nots)
	Err(err)

	if len(nots) == 0 {
		fmt.Println("No notification, see ya")
	}

	for i, v := range nots {
		fmt.Printf(`
Notification [%d]:
 * Is it unread? %t
 * Reason:       %s
 * Updated At:   %s
 * Title:        %s
 * Url:          %s
 * Type:         %s
 * Full Repo:    %s
`, i+1, v.Unread, v.Reason, v.UpdatedAt, v.Subject.Title, v.Subject.Url, v.Subject.Type, v.Repository.FullName)
	}

}

func Err(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}
