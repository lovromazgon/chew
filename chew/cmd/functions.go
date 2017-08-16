package cmd

import (
	"bitbucket.org/lovromazgon/chew"
	"bitbucket.org/lovromazgon/chew/funcmap"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	RootCmd.AddCommand(functionsCmd)

	functionsCmd.Flags().StringVarP(&function, "func", "f", "", "Get documentation for only one function")
}

var functionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "Print the documentation for custom added functions",
	Long:  `Chew has many custom functions that can be used inside templates. The documentation for each function can be viewed with this command.`,
	RunE:  functionsRun,
}

// ----------------------------------------------------------------

var (
	function string
)

func functionsRun(cmd *cobra.Command, args []string) error {
	template := chew.New("main")

	_, err := template.Parse(funcmap.FuncDocTemplates)
	if err != nil {
		return err
	}

	chewable := functionsToChewable(template.Functions, function)
	return template.ExecuteChewable(&chew.WriterWrapper{os.Stdout}, *chewable)
}

func functionsToChewable(functions funcmap.Functions, filter string) *chew.Chewable {
	chewable := &chew.Chewable{
		Data: make([]chew.ChewableData, len(functions)),
	}

	for i,fun := range functions {
		if filter != "" && filter != fun.Doc.Name {
			continue
		}

		localData, err := chew.ToMap(fun.Doc)
		if err != nil {
			panic(err)
		}

		data := chew.ChewableData{
			Templates: map[string]string {
				fun.Doc.Template():fun.Doc.Name + ".out",
			},
			Local: localData,
		}
		chewable.Data[i] = data
	}

	return chewable
}