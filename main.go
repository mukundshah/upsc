// package main

// import (
// 	"github.com/mukundshah/upsc/cmd/root"
// )

// func main() {
// 	root.Execute()
// }

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type Framework struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Color       string `json:"color"`
	TemplateId  string `json:"templateId"`
}

type Toolings struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Color       string `json:"color"`
}

type Language struct {
	Name        string      `json:"name"`
	DisplayName string      `json:"displayName"`
	Color       string      `json:"color"`
	Frameworks  []Framework `json:"frameworks"`
	Toolings    []Toolings  `json:"toolings"`
	TemplateId  string      `json:"templateId"`
}

type Project struct {
	Name       string
	Dir        string
	Lang       string
	Framework  string
	Toolings   []string
	TemplateId string
}

var languages []Language
var langNames []string
var frameworks []Framework
var frameworkNames []string
var toolings []Toolings
var toolingNames []string

func getLanguages() {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &languages)

}

func getFrameworks(language Language) {
	frameworks = language.Frameworks
}

func getToolings(language Language) {
	toolings = language.Toolings
}

func getLanguageNames() {

	for _, lang := range languages {
		langNames = append(langNames, lang.DisplayName)
	}

}

func getLanguageByName(name string) Language {
	for _, lang := range languages {
		if lang.DisplayName == name {
			return lang
		}
	}
	return Language{}
}

func getFrameworkNames() {
	for _, framework := range frameworks {
		frameworkNames = append(frameworkNames, framework.DisplayName)
	}
}

func getFrameworkByName(name string) Framework {
	for _, framework := range frameworks {
		if framework.DisplayName == name {
			return framework
		}
	}
	return Framework{}
}

func getToolingsNames() {
	for _, tooling := range toolings {
		toolingNames = append(toolingNames, tooling.DisplayName)
	}
}

func generateProject(project Project) {
	fmt.Println(project)
}

func main() {

	var project Project

	getLanguages()

	var cmdInit = &cobra.Command{
		Use:                   "init",
		Short:                 "Initialize project from template",
		Long:                  "Initialize project from template",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			survey.AskOne(&survey.Input{
				Message: "Project name:",
				Default: "my-project",
			}, &project.Name, survey.WithValidator(survey.Required))

			survey.AskOne(&survey.Input{
				Message: "Project directory:",
				Default: "./",
			}, &project.Dir, survey.WithValidator(survey.Required))

			getLanguageNames()
			survey.AskOne(&survey.Select{
				Message: "Select language:",
				Options: langNames,
			}, &project.Lang, survey.WithValidator(survey.Required))

			language := getLanguageByName(project.Lang)
			getFrameworks(language)
			getToolings(language)

			if len(frameworks) > 0 {
				getFrameworkNames()
				survey.AskOne(&survey.Select{
					Message: "Select framework:",
					Options: frameworkNames,
				}, &project.Framework, survey.WithValidator(survey.Required))
				framework := getFrameworkByName(project.Framework)
				project.TemplateId = framework.TemplateId

			} else {
				project.TemplateId = language.TemplateId
			}

			if len(toolings) > 0 {
				getToolingsNames()
				survey.AskOne(&survey.MultiSelect{
					Message: "Select toolings:",
					Options: toolingNames,
				}, &project.Toolings, survey.WithValidator(survey.Required))
			}

			generateProject(project)
		},
	}

	var rootCmd = &cobra.Command{Use: "upsc",
		Short: "Universal Project Starter CLI (UPSC)",
		Long:  "Universal Project Starter CLI (UPSC)\nA Fast and flexible CLI for initializing projects from templates",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	rootCmd.AddCommand(cmdInit)
	rootCmd.Execute()
}
