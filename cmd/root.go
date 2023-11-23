/*
Copyright © 2023 Sergey Kruglov <srgy.krglv@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergey-kruglov/gotempl/lib"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotempl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: generateFiles,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotempl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Config struct {
	Templates []Template
}

type Template struct {
	Name  string
	Files []TemplateFile
}

type TemplateFile struct {
	Name         string
	TemplatePath string
	Parameters   []TemplateFileParameter
}

type TemplateFileParameter struct {
	Name  string
	Value int8
	Case  string
}

// Generate files function
func generateFiles(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatalln("No arguments provided!")
	}

	if len(args) < 3 {
		log.Fatalln("At least 3 arguments required, ex.: gotempl controller Todo ...params src/controllers/")
	}

	config := getConfig()
	if config.Templates == nil {
		log.Fatalln("Incorrect config file.")
	}

	outPathArg := args[len(args)-1]

	template := getTemplate(args[0], config)
	for _, file := range template.Files {
		pwd, _ := os.Getwd()
		path := filepath.Join(pwd, file.TemplatePath)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalln("Cannot read template path file. Check the config file.")
		}

		contentText := string(content)
		fileName := file.Name
		for _, param := range file.Parameters {
			argValue := args[param.Value+1]
			contentText = strings.Replace(contentText, param.Name, argValue, -1)
			fileName = strings.Replace(fileName, param.Name, argValue, -1)
		}

		outPath := filepath.Join(pwd, outPathArg, fileName)
		os.WriteFile(outPath, []byte(contentText), os.FileMode(0644))
	}
}

func getConfig() Config {
	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, lib.ConfigFileName)
	configJson, _ := os.ReadFile(path)

	var config Config
	if err := json.Unmarshal(configJson, &config); err != nil {
		log.Fatalln(err)
	}

	return config
}

func getTemplate(templateName string, config Config) Template {
	var template Template
	found := false
	for _, t := range config.Templates {
		if t.Name == templateName {
			template = t
			found = true
		}
	}

	if !found {
		log.Fatalln("Unknown template type. Check the config file.")
	}

	return template
}
