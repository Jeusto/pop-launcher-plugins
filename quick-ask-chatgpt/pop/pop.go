package pop

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type PluginSearchResult struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Keywords    []string    `json:"keywords"`
	Exec        string      `json:"exec"`
	Icon        interface{} `json:"icon"`
	Window      interface{} `json:"window"`
}

type Append struct {
	Result PluginSearchResult `json:"Append"`
}

type Request struct {
	Type  string `json:"type"`
	ID    int    `json:"id,omitempty"`
	Query string `json:"data,omitempty"`
}

func CloseLauncher() {
	fmt.Println("\"Close\"")
}

func ClearSearchResults() {
	fmt.Println("\"Clear\"")
}

func Finish() {
	fmt.Println("\"Finished\"")
}

func AppendResultList(results []PluginSearchResult) {
	for _, result := range results {
		AppendResult(result)
	}
}

func AppendResult(result PluginSearchResult) {
	appendBytes, _ := json.Marshal(Append{Result: result})
	fmt.Printf("%s\n", appendBytes)
}

func ShowSingleResult(result PluginSearchResult) {
	ClearSearchResults()
	AppendResult(result)
	Finish()
}

func ShowErrorMessage(message string) {
	ClearSearchResults()
	AppendResult(PluginSearchResult{
		Name:        "Error",
		Description: message,
	})
	Finish()
}

func ClearInput() {
	output := map[string]interface{}{
		"Fill": "ask ",
	}
	outputBytes, _ := json.Marshal(output)

	fmt.Printf("%s\n", outputBytes)
	Finish()
}

func HandleRequests(requests chan<- Request) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		var request map[string]interface{}
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			ShowErrorMessage(err.Error())
			continue
		}

		if _, ok := request["Search"]; ok {
			query := request["Search"].(string)
			query = strings.Replace(query, "ask", "", 1)
			query = strings.TrimSpace(query)

			requests <- Request{
				Type:  "search",
				Query: query,
			}
		}
		if _, ok := request["Activate"]; ok {
			requests <- Request{
				Type: "activate",
				ID:   int(request["Activate"].(float64)),
			}
		}
	}

	close(requests)
}
