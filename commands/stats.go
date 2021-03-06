package commands

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/gaudit/analyze"
	"github.com/hashicorp/gaudit/appends"
	"github.com/hashicorp/gaudit/config"
	"github.com/hashicorp/gaudit/state"
)

func Stats(options config.Options) {

	// get last audit
	audit, err := state.Load(options.Storage)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	// get rules
	rules, err := analyze.Load(options)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
	var rulesList []string
	for _, rule := range rules {
		rulesList = append(rulesList, rule.Name)
	}
	sort.Strings(rulesList)

	// get appends
	appendList, err := appends.Load(options)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}

	// init stat structure
	stats := make(map[string]int)
	stats["TOTAL"] = len(audit.Repos)
	stats["blank_description"] = 0
	stats["blank_language"] = 0
	stats["blank_topics"] = 0
	stats["private"] = 0
	stats["archived"] = 0
	stats["disabled"] = 0
	stats["blank_license"] = 0
	stats["total_stars"] = 0
	stats["total_forks"] = 0
	stats["total_watchers"] = 0
	stats["updated_0-30_days"] = 0
	stats["updated_31-60_days"] = 0
	stats["updated_61-90_days"] = 0
	stats["updated_91-365_days"] = 0
	stats["updated_>_365_days"] = 0

	// output
	for _, k := range audit.Index {
		repo := audit.Repos[k]

		if strings.TrimSpace(repo.Description) == "" {
			stats["blank_description"]++
		}

		if strings.TrimSpace(repo.Language) == "" {
			stats["blank_language"]++
		}

		if len(repo.Topics) == 0 {
			stats["blank_topics"]++
		}

		if repo.Private {
			stats["private"]++
		}

		if repo.Archived {
			stats["archived"]++
		}

		if repo.Disabled {
			stats["disabled"]++
		}

		if strings.TrimSpace(repo.License) == "" {
			stats["blank_license"]++
		}

		stats["total_stars"] += repo.Stargazers

		stats["total_forks"] += repo.Forks

		stats["total_watchers"] += repo.Watchers

		hrs := int(time.Since(repo.Updated).Hours())

		if hrs <= (30 * 24) {
			stats["updated_0-30_days"]++
		}

		if hrs <= (60*24) && hrs > (30*24) {
			stats["updated_31-60_days"]++
		}

		if hrs <= (90*24) && hrs > (60*24) {
			stats["updated_61-90_days"]++
		}

		if hrs <= (365*24) && hrs > (90*24) {
			stats["updated_91-365_days"]++
		}

		if hrs > (365 * 24) {
			stats["updated_>_365_days"]++
		}

		for _, rule := range audit.Results[k].Rules {
			if rule.Status == "success" {
				stats[rule.Name] = stats[rule.Name] + 1
			}
		}

		for _, a := range appendList {
			if a.Name == repo.FullName {
				// do something
			}
		}

	}

	// print results
	printStat("TOTAL", stats)
	printStat("total_stars", stats)
	printStat("total_forks", stats)
	printStat("total_watchers", stats)
	printStat("blank_description", stats)
	printStat("blank_language", stats)
	printStat("blank_topics", stats)
	printStat("private", stats)
	printStat("archived", stats)
	printStat("disabled", stats)
	printStat("blank_license", stats)
	printStat("updated_0-30_days", stats)
	printStat("updated_31-60_days", stats)
	printStat("updated_61-90_days", stats)
	printStat("updated_91-365_days", stats)
	printStat("updated_>_365_days", stats)
	for _, r := range rulesList {
		printStat(r, stats)
	}

}

func printStat(stat string, stats map[string]int) {
	fmt.Print(strings.Replace(stat, "_", " ", -1))
	fmt.Print(":")
	fmt.Print(strings.Repeat(" ", 22-len(stat)))
	fmt.Printf("%d", stats[stat])
	fmt.Print(strings.Repeat(" ", 8-len(strconv.Itoa(stats[stat]))))
	if stats[stat] < stats["TOTAL"] {
		fmt.Printf("  (%5.2f %%)", (float64(stats[stat]) / float64(stats["TOTAL"]) * 100))
	}
	fmt.Print("\n")
}
