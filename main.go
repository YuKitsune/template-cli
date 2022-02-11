package main

import (
	"errors"
	"fmt"
	"github.com/traefik/paerser/parser"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		_, _ = fmt.Fprintf(os.Stderr, "err: %s", err.Error())
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

		// Read in the data
		data, err := ioutil.ReadFile(v)
		if err != nil {
			return fmt.Errorf("error reading \"%s\": %s", v, err.Error())
		}

		// Figure out where we want to write our results to
		var writer io.Writer
		if dryRun {
			writer = os.Stdout
		} else {
			path := filepath.Join(outputDir, v)

			// All good, write to the file
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
			if err != nil {
				return fmt.Errorf("error opening file for \"%s\": %s", path, err.Error())
			}

			// Ensure it's empty first
			err = file.Truncate(0)
			if err != nil {
				return err
			}

			writer = file
		}

		origStr := string(data)

		// Todo: What should the name be?
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

	// Ensure we're allowed to overwrite if the output file is the same as the input file
	absOutPath, err := filepath.Abs(outputDir)
	if err != nil {
		// todo: improve error message
		return fmt.Errorf("error determining output directory \"%s\": %s", absOutPath, err.Error())
	}

	// Ensure files are not overwriten unless the --overwrite flag is specifie
	for _, inputFile := range inputFiles {

		absInFile, err := filepath.Abs(inputFile)
		if err != nil {
			// todo: improve error message
			return fmt.Errorf("error determining input file \"%s\": %s", absInFile, err.Error())
		}

		absInPath := filepath.Dir(absInFile)
		if !overwrite && absInPath == absOutPath {
			return fmt.Errorf("execution would overwrite input files, use the --overwrite flag to allow for the input files to be overridden")
		}
	}

	return nil
}

func readValues(v interface{}) error {

	// First, read from the values file if it exists
	if len(valuesFile) > 0 && fileExists(valuesFile) {
		data, err := ioutil.ReadFile(valuesFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(data, v)
		if err != nil {
			return err
		}
	}

	// Second, read from the flags, flags should overwrite any values defined in the file
	// vals := getValues(values)
	// Todo: Extremely scuffed, maybe write your own parser?
	for _, value := range values {
		m := make(map[string]string)
		parts := strings.Split(value, "=")
		name := parts[0]
		val := parts[1]

		m[name] = val

		firstPart := strings.Split(name, ".")[0]

		err := parser.Decode(m, v, firstPart)
		if err != nil {
			return err
		}
	}

	return nil
}

func getValues(values []string) map[string]string {

	m := make(map[string]string)

	for _, value := range values {
		parts := strings.Split(value, "=")
		name := parts[0]
		val := parts[1]

		m[name] = val
	}

	return m
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
