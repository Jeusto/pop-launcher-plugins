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
	defer close(requests)

	// Read from stdin and parse possible the received event
	for scanner.Scan() {
		line := scanner.Text()

		switch line {
		case "\"Exit\"":
			requests <- Request{
				Type: "Exit",
			}
		case "\"Interrupt\"":
			requests <- Request{
				Type: "Interrupt",
			}
		default:
			var request map[string]interface{}

			if err := json.Unmarshal([]byte(line), &request); err != nil {
				ShowErrorMessage(err.Error())
				continue
			}

			if _, ok := request["Search"]; ok {
				requests <- Request{
					Type:  "Search",
					Query: cleanQuery(request["Search"].(string)),
				}
			}
			if _, ok := request["Activate"]; ok {
				requests <- Request{
					Type: "Activate",
					ID:   int(request["Activate"].(float64)),
				}
			}
		}
	}
}

func cleanQuery(query string) string {
	query = strings.Join(strings.Split(query, " ")[1:], " ")
	query = strings.TrimSpace(query)
	return query
}
