package pop

import (
	"encoding/json"
	"fmt"
)

type PluginSearchResult struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Keywords    interface{} `json:"keywords"`
	Icon        interface{} `json:"icon"`
	Exec        interface{} `json:"exec"`
	Window      interface{} `json:"window"`
}

type Append struct {
	Result PluginSearchResult `json:"Append"`
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

func ShowErrorMessage(message string) {
	ClearSearchResults()
	AppendResult(PluginSearchResult{
		Name:        "Error",
		Description: message,
	})
	Finish()
}
