package cmd

import (
	"os"
	"regexp"

	"github.com/lovromazgon/chew"
	"github.com/lovromazgon/chew/funcmap"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(functionsCmd)

	functionsCmd.Flags().StringVarP(&filter, "filter", "f", "", "Get documentation for functions that match the regex")
}

var functionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "Print the documentation for custom added functions",
	Long:  `Chew has many custom functions that can be used inside templates. The documentation for each function can be viewed with this command.`,
	RunE:  functionsRun,
}

// ----------------------------------------------------------------

var (
	filter string
)

func functionsRun(cmd *cobra.Command, args []string) error {
	template := chew.New("main")

	_, err := template.Parse(funcmap.FuncDocTemplates)
	if err != nil {
		return err
	}

	regex, err := regexp.Compile(filter)
	if err != nil {
		return err
	}

	chewable := functionsToChewable(template.Functions, regex)
	return template.ExecuteChewable(&chew.WriterWrapper{os.Stdout}, *chewable)
}

func functionsToChewable(functions funcmap.Functions, filter *regexp.Regexp) *chew.Chewable {
	chewable := &chew.Chewable{
		Data: make([]chew.ChewableData, len(functions)),
	}

	for i, fun := range functions {
		if filter != nil && !filter.Match([]byte(fun.Doc.Name)) {
			continue
		}

		localData, err := chew.ToMap(fun.Doc)
		if err != nil {
			panic(err)
		}

		data := chew.ChewableData{
			Templates: map[string]string{
				fun.Doc.Template(): fun.Doc.Name + ".out",
			},
			Local: localData,
		}
		chewable.Data[i] = data
	}

	return chewable
}
