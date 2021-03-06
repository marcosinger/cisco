package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/http"
)

type Heroku struct {
	Status struct {
		Production  string
		Development string
	}

	Issues []struct {
		Updates []struct {
			Title      string
			Contents   string
			Created_at string
		}
	}
}

func main() {
	heroku := new(Heroku)
	res := heroku.Call()

	fmt.Println(res)
}

func (h *Heroku) ProductionStatus() string {
	return fmt.Sprintf("%s: %s", "Production", h.colorFor(h.Status.Production))
}

func (h *Heroku) DevelopmentStatus() string {
	return fmt.Sprintf("%s: %s", "Development", h.colorFor(h.Status.Development))
}

func (h *Heroku) colorFor(c string) string {
	cfn := colorFn(c)

	switch c {
	case "green":
		return cfn("[ok]")

	default:
		issue := h.Issues[0].Updates[0]
		info := fmt.Sprintf("[%s] \n%s", issue.Title, issue.Contents)
		return cfn(info)
	}
}

func colorFn(c string) func(...interface{}) string {
	var colorAttr color.Attribute

	switch c {
	case "yellow":
		colorAttr = color.FgYellow
	case "red":
		colorAttr = color.FgRed
	default:
		colorAttr = color.FgGreen
	}

	return color.New(colorAttr).SprintFunc()
}

func (h *Heroku) String() string {
	blue := color.New(color.FgBlue).SprintFunc()
	return fmt.Sprintf("%s\n\n%s\n%s", blue("[Heroku Status]"), h.ProductionStatus(), h.DevelopmentStatus())
}

func (h *Heroku) Url() string {
	return "https://status.heroku.com/api/v3/current-status"
}

func (h *Heroku) Call() (heroku *Heroku) {
	res, err := http.Get(h.Url())

	if err != nil {
		log.Printf("Heroku error when calling the API: %s", err.Error())
		return
	}

	decoder := json.NewDecoder(res.Body)
	ok := decoder.Decode(&heroku)

	if ok != nil {
		log.Printf("Heroku error when parsing JSON: %s", ok.Error())
		return
	}

	return

}
