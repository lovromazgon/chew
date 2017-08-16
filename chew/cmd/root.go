// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"os"

	"encoding/json"
	"io/ioutil"

	"bitbucket.org/lovromazgon/chew"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "chew",
	Short: "A CLI for Go Templates",
	Long: `Chew is a CLI for Go Templates which generates output based on input data.
It parses all templates that it can find in the defined folder and then
generates the output based on the data in the JSON input.`,
	PreRunE: preChew,
	RunE:    chewRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.Flags().StringVarP(&dataPath, "data", "d", "", "Path to input JSON file with data")
	RootCmd.Flags().StringVarP(&templatesPath, "templates", "t", "", "Path to folder with templates (will be read recursively)")
	RootCmd.Flags().StringVarP(&outPath, "out", "o", "", "Path to output folder")

	RootCmd.MarkFlagFilename("data", ".json")
	RootCmd.MarkFlagRequired("data")
	RootCmd.MarkFlagRequired("templates")
	RootCmd.MarkFlagRequired("out")
}

// ----------------------------------------------------------------

var (
	templatesPath string
	dataPath      string
	outPath       string
)

func preChew(cmd *cobra.Command, args []string) error {
	if templatesPath == "" {
		return errors.New("Templates flag is required!")
	} else if dataPath == "" {
		return errors.New("Data flag is required!")
	} else if outPath == "" {
		return errors.New("Out flag is required!")
	} else if _, err := os.Stat(templatesPath); err != nil {
		return err
	} else if _, err := os.Stat(dataPath); err != nil {
		return err
	}

	return nil
}

func chewRun(cmd *cobra.Command, args []string) error {
	dataRaw, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return err
	}

	chewable := &chew.Chewable{}

	err = json.Unmarshal(dataRaw, chewable)
	if err != nil {
		return err
	}

	template := chew.New("main")
	_, err = template.ParseFolder(templatesPath)
	if err != nil {
		return err
	}

	return template.ExecuteChewable(&chew.MultiFileWriter{Out: outPath}, *chewable)
}
