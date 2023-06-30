package main

import (
	"strings"

	"1pt.one/pop"
	"1pt.one/shortener"

	"golang.design/x/clipboard"
)

const HELP_TEXT = "Type or paste a URL and press enter to shorten it. It will get copied to your clipboard."
const RESULT_TEXT = "The link has been copied to your clipboard. You can close the launcher."

type Plugin struct {
	shortened_url string
	long_url      string
	short_id      string
}

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	plugin := Plugin{shortened_url: ""}
	requests := make(chan pop.Request)

	go pop.HandleRequests(requests)

	for request := range requests {
		switch request.Type {
		case "Exit":
			return
		case "Search":
			plugin.search(request.Query)
		case "Activate":
			plugin.activate(request.ID)
		}
	}
}

func (plugin *Plugin) search(query string) {
	var split []string = strings.Split(query, " ")

	if len(split) > 0 {
		plugin.long_url = split[0]
	}
	if len(split) > 1 {
		plugin.short_id = split[1]
	}

	if plugin.shortened_url != "" {
		pop.ShowSingleResult(pop.PluginSearchResult{
			Name:        plugin.shortened_url,
			Description: RESULT_TEXT,
		})
	} else {
		pop.ShowSingleResult(pop.PluginSearchResult{
			Name: HELP_TEXT,
		})
	}
}

func (plugin *Plugin) activate(index int) {
	plugin.shortened_url = ""

	pop.ClearInput()
	pop.ShowSingleResult(pop.PluginSearchResult{
		Name: "Shortening...",
	})

	shortened_url, err := shortener.ShortenURL(plugin.long_url, plugin.short_id)
	plugin.shortened_url = shortened_url

	if err != nil {
		pop.ShowSingleResult(pop.PluginSearchResult{
			Name: "Error: " + err.Error(),
		})
	} else {
		clipboard.Write(clipboard.FmtText, []byte(plugin.shortened_url))
		pop.ShowSingleResult(pop.PluginSearchResult{
			Name:        plugin.shortened_url,
			Description: RESULT_TEXT,
		})
	}
}
