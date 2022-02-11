package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	values     []string
	valuesFile string

	inputFiles []string
	outputDir  string

	overwrite bool
	dryRun    bool
)

func init() {

	rootCmd.Flags().StringArrayVarP(&values, "value", "v", make([]string, 0), "a value to be substituted, format should be name=value")
	rootCmd.Flags().StringVarP(&valuesFile, "values-file", "f", "", "the path to a yaml formatted file where the values can be sourced from")

	rootCmd.Flags().StringArrayVarP(&inputFiles, "input", "i", make([]string, 0), "the path to a file where the templates live")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "the path to a directory where the files will be placed after the templates have been applied")

	rootCmd.Flags().BoolVar(&overwrite, "overwrite", false, "whether or not the input files should be overwritten after applying the templates")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "whether or not the this is a dry run")
}

var rootCmd = &cobra.Command{
	Use:   "populate-template",
	Short: "Populate Template is a tool used to populate a set of files with values using Go templates",
	RunE:  run,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {

	// Validate flags
	err := validateFlags()
	if err != nil {
		return err
	}

	// Load values
	var values map[string]interface{}
	err = readValues(&values)
	if err != nil {
		return err
	}

	// Iterate through input files
	for _, v := range inputFiles {

		// Figure out where we want to write our results to
		var writer io.Writer
		if dryRun {
			writer = os.Stdout
		} else {
			path := filepath.Join(outputDir, v)

			// Todo: Move to validate func
			// Ensure we're allowed to overwrite if the output file is the same as the input file
			absOutPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("error opening file for \"%s\": %s", absOutPath, err.Error())
			}

			absInPath, err := filepath.Abs(v)
			if err != nil {
				return fmt.Errorf("error opening file for \"%s\": %s", absInPath, err.Error())
			}

			if !overwrite && absInPath == absOutPath {
				return fmt.Errorf("execution would overwrite input files, use the --overwrite flag to allow for the input files to be overridden")
			}

			// All good, write to the file
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return fmt.Errorf("error opening file for \"%s\": %s", v, err.Error())
			}

			// Ensure it's empty first
			file.Truncate(0)
			writer = file
		}

		data, err := ioutil.ReadFile(v)
		if err != nil {
			return fmt.Errorf("error reading \"%s\": %s", v, err.Error())
		}

		origStr := string(data)

		tmpl, err := template.New("test").Parse(origStr)
		if err != nil {
			return fmt.Errorf("error creating template for \"%s\": %s", v, err.Error())
		}

		err = tmpl.Execute(writer, values)
		if err != nil {
			return fmt.Errorf("error executing template for \"%s\": %s", v, err.Error())
		}
	}

	return nil
}

func validateFlags() error {

	// Validate values
	if len(values) == 0 && len(valuesFile) == 0 {
		return fmt.Errorf("values must be specified either via the --values-file or --value (-v) flags")
	}

	// Validate input files
	if len(inputFiles) == 0 {
		return fmt.Errorf("at least one input file must be specified via the --input (-i) flag")
	}

	// Todo: Ensure files are not overwriten unless the --overwrite flag is specifie

	return nil
}

func readValues(v interface{}) error {

	// First, read from the values file if it exists

	// Second, read from the flags, flags should overwrite any values defined in the file

	return nil
}
