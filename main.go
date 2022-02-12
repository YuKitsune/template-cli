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
	values []string

	inputFiles []string
	outputDir  string

	overwrite bool
	dryRun    bool
)

func init() {
	rootCmd.Flags().StringArrayVarP(&values, "value", "v", make([]string, 0), "a value to be passed into the template, format should be \"name=value\"")

	rootCmd.Flags().StringArrayVarP(&inputFiles, "input", "i", make([]string, 0), "the path to a template file")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "the path to a directory where the files will be placed after the templates have been applied")

	rootCmd.Flags().BoolVar(&overwrite, "overwrite", false, "allows the input files should be overwritten if necessary after applying the templates")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "simulates a template execution, printing the results to stdout")
}

var rootCmd = &cobra.Command{
	Use:   "populate-template",
	Short: "Populate Template is a tool used to populate a set of files with values using Go templates",
	RunE:  run,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(_ *cobra.Command, _ []string) error {

	// Validate flags
	err := validateFlags()
	if err != nil {
		return err
	}

	// Load in the values
	var values map[string]interface{}
	err = parseValues(&values)
	if err != nil {
		return err
	}

	// Iterate through input files
	for _, inputFile := range inputFiles {

		// Read in the data
		data, err := ioutil.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read \"%s\": %s", inputFile, err.Error())
		}

		origStr := string(data)

		// Figure out where we want to write our results to
		writer, err := getResultWriter(inputFile)
		if err != nil {
			return err
		}

		// Todo: What should the name be?
		tmpl, err := template.New("test").Parse(origStr)
		if err != nil {
			return fmt.Errorf("failed to create template for \"%s\": %s", inputFile, err.Error())
		}

		err = tmpl.Execute(writer, values)
		if err != nil {
			return fmt.Errorf("failed to execute template for \"%s\": %s", inputFile, err.Error())
		}
	}

	return nil
}

func validateFlags() error {

	// Validate values
	if len(values) == 0 {
		return fmt.Errorf("at least one value must be specified via the --value (-v) flag")
	}

	// Ensure we're allowed to overwrite if the output file is the same as the input file
	absOutPath, err := filepath.Abs(outputDir)
	if err != nil {
		// todo: improve error message
		return fmt.Errorf("failed to find output directory \"%s\": %s", absOutPath, err.Error())
	}

	// Ensure files are not overwritten unless the --overwrite flag is specified
	for _, inputFile := range inputFiles {

		absInFile, err := filepath.Abs(inputFile)
		if err != nil {
			// todo: improve error message
			return fmt.Errorf("failed to find input file \"%s\": %s", absInFile, err.Error())
		}

		absInPath := filepath.Dir(absInFile)
		if !overwrite && absInPath == absOutPath {
			return fmt.Errorf("execution would overwrite input files, use the --overwrite flag to allow for the input files to be overwritten")
		}
	}

	return nil
}

func parseValues(v interface{}) error {

	// Todo: Use a custom parser
	rootNode := "root"
	valueMap := getValues(values, rootNode)
	err := parser.Decode(valueMap, v, rootNode)
	if err != nil {
		return err
	}

	return nil
}

func getValues(values []string, rootNode string) map[string]string {

	m := make(map[string]string)

	for _, value := range values {
		parts := strings.Split(value, "=")

		// Need to prefix everything with the root node name
		// This doesn't actually make it into the values interface, it's something traefik specific that we don't need
		key := fmt.Sprintf("%s.%s", rootNode, parts[0])
		val := parts[1]

		m[key] = val
	}

	return m
}

func getResultWriter(inputFile string) (io.Writer, error) {

	if dryRun {
		// Todo: Custom formatting
		return os.Stdout, nil
	}

	inputFileName := path.Base(inputFile)
	outFileName := filepath.Join(outputDir, inputFileName)

	// All good, write to the file
	file, err := os.OpenFile(outFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file \"%s\": %s", outFileName, err.Error())
	}

	return file, nil
}
