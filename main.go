package main

// TODO: Don't use any slices, just use maps instead, since all calculations use maps anyway.

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	if err := Execute(); err != nil {
		fmt.Printf("error: \n%v", err)
	}
}

// Execute runs the root command, which in turn runs child commands based on further input.
// It returns an error if the command could not be executed by cobra
func Execute() error {
	return rootCommand.Execute()
}

const (
	inputFlagName   = "input"
	outputFlagName  = "output"
	formulaFlagName = "formula"

	outputFlagDefaultValue = "output.csv"
)

func init() {
	rootCommand.Flags().String(inputFlagName, "", "The filenames of the input files.")
	rootCommand.Flags().String(outputFlagName, outputFlagDefaultValue, "The Filename of the output file.")
	rootCommand.Flags().String(formulaFlagName, "", "The formula.")
}

var (
	rootCommand = &cobra.Command{
		Use:   "csvset",
		Short: "test1", // TODO: set short
		Long:  `test2`, // TODO: set long
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Makefile
			flags := cmd.NonInheritedFlags()

			// Read input flag
			fmt.Printf("Reading %s flag...\n", inputFlagName)
			inputFlagStringValue, err := flags.GetString(inputFlagName)
			if err != nil {
				fmt.Printf("Couldn't read %s flag.\n", inputFlagName)
				return
			}
			if inputFlagStringValue == "" {
				fmt.Printf("\"%s\" flag not included.\n", inputFlagName)
				return
			}
			inputFilenames := strings.Split(inputFlagStringValue, ",")
			fmt.Printf("Successfully read %s flag: %+v\n", inputFlagName, inputFilenames)

			// Read output flag
			fmt.Printf("Reading %s flag...\n", outputFlagName)
			outputFilename, err := flags.GetString(outputFlagName)
			if err != nil {
				fmt.Printf("Couldn't read %s flag.\n", outputFlagName)
				return
			}
			if outputFilename == outputFlagDefaultValue {
				fmt.Printf("\"%s\" flag not specified, set to %s by default.\n", outputFlagName, outputFlagDefaultValue)
			} else {
				fmt.Printf("Successfully read %s flag: %s\n", outputFlagName, outputFilename)
			}

			// Read formula flag
			fmt.Printf("Reading %s flag...\n", formulaFlagName)
			formula, err := flags.GetString(formulaFlagName)
			if err != nil {
				fmt.Printf("Couldn't read %s flag.\n", formulaFlagName)
				return
			}
			if formula == "" {
				fmt.Printf("\"%s\" flag not included.\n", formulaFlagName)
				return
			} else {
				fmt.Printf("Successfully read %s flag: %s\n", formulaFlagName, formula)
			}

			// Read input files
			fmt.Printf("Reading input files %+v...\n", inputFilenames)
			inputs, err := ReadCSVs(inputFilenames)
			if err != nil {
				fmt.Println(err)
				return
			} else {
				fmt.Printf("Successfully read input files.\n")
			}

			// TODO: Make basic checks for number of parentheses?

			// Perform operations
			fmt.Printf("Performing operations...\n")
			operation := NewOperation(&inputs, &Operands, nil, formula)
			err = operation.Execute()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Operations successfully completed.\n")

			// Create and write to output file
			fmt.Printf("Writing result to file \"%s\"...\n", outputFilename)
			DumpCSV(outputFilename, operation.Result)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Output successfully written to file.\n")
			fmt.Printf("Exiting...\n")
		},
	}
)

func DumpCSV(filename string, data []string) error {
	csvFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed to create file: %v", err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	for i := range data {
		err = csvwriter.Write([]string{data[i]})
		if err != nil {
			return fmt.Errorf("Failed to write file %s: %v", filename, err)
		}
	}
	csvwriter.Flush()

	return nil
}

func ReadCSVs(filenames []string) ([][]string, error) {
	result := make([][]string, len(filenames))
	for nameIndex := range filenames {
		tempResult, err := ReadCSV(filenames[nameIndex])
		if err != nil {
			return nil, fmt.Errorf("Failed to read files: %v", err)
		}
		result[nameIndex] = tempResult
	}
	return result, nil
}

func ReadCSV(filename string) ([]string, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file %s: %v", filename, err)
	}
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %s: %v", filename, err)
	}
	result := make([]string, len(csvLines))
	for lineIndex := range csvLines {
		result[lineIndex] = csvLines[lineIndex][0]
	}
	return result, nil
}

type Operation struct {
	// ValueLists is a reference to a list of all value lists that are available.
	ValueLists *[][]string
	// Operands is a reference to a collection of all Operands that are available.
	Operands *OperandCollection
	// Operand is the Operand to be perform the current Operation.
	Operand *Operand
	// StringOperation is the string representation of the current Operation.
	StringOperation string
	// SubOperations are the Operations that have to be performed to perform the current Operation.
	SubOperations []*Operation
	// Result is the result of the current Operation.
	Result []string
}

func NewOperation(valueLists *[][]string, operands *OperandCollection, operand *Operand, stringOperation string) *Operation {
	return &Operation{
		ValueLists:      valueLists,
		Operands:        operands,
		Operand:         operand,
		StringOperation: stringOperation,
	}
}

func (op *Operation) Execute() error {
	// If StringOperation only contains a number, set that ValueList as result and return.
	op.Printf("resolve-check")
	if digitCheck.MatchString(op.StringOperation) {
		index, _ := strconv.Atoi(op.StringOperation)
		op.Result = (*op.ValueLists)[index]
		return nil
	}

	// Attempt to split StringOperation into terms and operands.
	op.Printf("split")
	terms, operand, err := Operands.SplitStringByOperands(op.StringOperation)
	if err != nil {
		return op.Errorf("Failed to split string:\n%v", err)
	}
	// No terms were found.
	op.Printf("term-check")
	if len(terms) == 0 {
		return op.Errorf("No terms found in string:\n%s", op.StringOperation)
	}
	// No operand was found, it may have been a simple term surrounded by parentheses.
	op.Printf("operand-check")
	if operand == nil {
		op.Printf("length-check")
		if len(terms) == 1 {
			op.Printf("term-zero")
			term := terms[0]
			op.Printf("resolve-check")
			if digitCheck.MatchString(term) {
				index, _ := strconv.Atoi(term)
				// Index term was out of bounds.
				if index >= len(*op.ValueLists) {
					return op.Errorf("Index term out of bounds:\nIndex: %d\nBounds: %d", index, len(*op.ValueLists)-1)
				}
				op.Result = (*op.ValueLists)[index]
				return nil
			}
		}
		// Multiple terms were found, but no Operand (that means the coder of this program messed up).
		return op.Errorf("Multiple terms were found, but no operand.\nTerms:\n%+v", terms)
	}

	// There is an operand and >1 terms, another Operation is required.
	op.SubOperations = make([]*Operation, len(terms))
	op.Operand = operand
	for termIndex := range terms {
		op.SubOperations[termIndex] = NewOperation(op.ValueLists, op.Operands, nil, terms[termIndex])
	}
	for subOperationIndex := range op.SubOperations {
		err := op.SubOperations[subOperationIndex].Execute()
		if err != nil {
			return err
		}
	}
	op.Result, err = op.Operand.Execute(op.GetResults())
	if err != nil {
		return op.Errorf("Failed to execute operand %s:\n%v", op.Operand.Name, err)
	}

	return nil
}

func (op *Operation) Errorf(str string, args ...interface{}) error {
	if str == "" {
		return fmt.Errorf("Error in operation: \"%s\"", op.StringOperation)
	}
	baseString := fmt.Sprintf("Error in operation: \"%s\"\n", op.StringOperation)
	return fmt.Errorf(baseString+"Details: "+str+"\n", args...)
}

func (op *Operation) Printf(str string, args ...interface{}) {
	baseString := fmt.Sprintf("Operation \"%s\": ", op.StringOperation)
	fmt.Printf(baseString+str+"\n", args...)
}

func (op *Operation) GetResults() [][]string {
	values := make([][]string, 0, len(op.SubOperations))
	for i := range op.SubOperations {
		values = append(values, op.SubOperations[i].Result)
	}
	return values
}

func GetFirstNumberFromStringAsString(inStr string) (outStr string, i int) {
	for i < len(inStr) && digitCheck.MatchString(string(inStr[i])) {
		outStr += string(inStr[i])
		i++
	}
	return
}

// digitCheck is used to check if a string is a number
var digitCheck = regexp.MustCompile(`^[0-9]+$`)

type Operand struct {
	Name     string      // the name of the operation this operand performs.
	Terms    int         // The number of terms this operand can accept; if <= 0 it can accept any number.
	Function setFunction // The function used to calculate the result for this operand.
}

func (od *Operand) Execute(values [][]string) ([]string, error) {
	if od.Terms > 0 && len(values) != od.Terms {
		return nil, fmt.Errorf("Incorrect number of terms for this operand, expected: %d, recieved: %d", od.Terms, len(values))
	}

	if len(values) == 0 {
		return []string{}, nil
	} else if len(values) == 1 {
		return values[0], nil
	}

	return od.Function(values), nil
}

type OperandCollection map[string]Operand

func (opc OperandCollection) SplitStringByOperands(str string) (terms []string, opd *Operand, err error) {
	for len(str) > 0 {
		fmt.Printf("Splitting \"%s\"...\n", str)
		// Attempt to extract ta number starting at the beginning of the string
		// append the number to terms and continue with the next iteration.
		// If a number is found add it to terms and continue the loop.
		if intStr, index := GetFirstNumberFromStringAsString(str); intStr != "" {
			terms = append(terms, intStr)
			if len(str)-1 < index {
				return
			}
			str = str[index:]
			continue
		}

		// Get current character and check if it is an operand.
		char := string(str[0])
		charOp, charOpOk := opc[char]
		if charOpOk {
			if opd != nil {
				ok := charOp.Name == opd.Name
				// Check if it is a different operand from a previous one on this recursion level.
				if !ok {
					return nil, nil, fmt.Errorf("Multiple different operands found at equal level of recursion.\nCurrent Operand:\n%+v\nPrevious operand(s):\n%+v", opc[char], opd)
				}
			}

			// If it's the first operand or the same as the current one, set the op var and continue the loop.
			if opd == nil {
				opd = new(Operand)
			}
			*opd = opc[char]
			if len(str)-1 < 1 {
				return
			}
			str = str[1:]
			continue
		}

		// If current char is an opening parenthesis, figure out where the corresponding end parenthesis is
		// and save the parenthesed string as a new term (without the parentheses).
		if char == "(" {
			parts := 1
			index := 1
			for parts != 0 {
				indexChar := string(str[index])
				if indexChar == "(" {
					parts++
				} else if indexChar == ")" {
					parts--
				}
				index++
			}
			// Add to terms, excluding parentheses, and continue the loop.
			terms = append(terms, str[1:index-1])
			if len(str)-1 < index+1 {
				return
			}
			str = str[index:]
			continue
		}
	}
	return
}

var Operands = OperandCollection{
	"+": {
		Name:     "union",
		Function: union,
		Terms:    0,
	},
	"*": {
		Name:     "intersection",
		Function: intersection,
		Terms:    0,
	},
	"-": {
		Name:     "set difference",
		Function: setDifference,
		Terms:    0,
	},
	"/": {
		Name:     "symmetric difference",
		Function: symmetricDifference,
		Terms:    0,
	},
}

type setFunction func([][]string) []string

func union(values [][]string) []string {
	return uniqueStrings(values...)
}

func intersection(values [][]string) []string {
	return stringsInAllSlices(values...)
}

func setDifference(values [][]string) []string {
	return stringsInOnlyFirstSlice(values...)
}

func symmetricDifference(values [][]string) []string {
	return stringsInOnlyOneSlice(values...)
}

func stringsInOnlyFirstSlice(stringSlices ...[]string) []string {
	stringMaps := make([]map[string]bool, 0, len(stringSlices))

	for _, stringSlice := range stringSlices {
		stringMaps = append(stringMaps, stringSliceToMap(stringSlice))
	}

	resultMap := make(map[string]bool)

	for str := range stringMaps[0] {
		for _, stringMap := range stringMaps[1:] {
			_, sOK := stringMap[str]
			_, rOK := resultMap[str]
			if sOK || (rOK && resultMap[str] == false) {
				resultMap[str] = false
			} else {
				resultMap[str] = true
			}
		}
	}

	return mapToStringSlice(resultMap)
}

func stringsInOnlyOneSlice(stringSlices ...[]string) []string {
	stringMaps := make([]map[string]bool, 0, len(stringSlices))

	for _, stringSlice := range stringSlices {
		stringMaps = append(stringMaps, stringSliceToMap(stringSlice))
	}

	resultMap := make(map[string]bool)

	for _, stringMap := range stringMaps {
		for str, _ := range stringMap {
			_, ok := resultMap[str]
			if ok {
				resultMap[str] = false
			} else {
				resultMap[str] = true
			}
		}
	}

	return mapToStringSlice(resultMap)
}

func stringsInAllSlices(stringSlices ...[]string) []string {
	if len(stringSlices) == 0 {
		return []string{}
	} else if len(stringSlices) == 1 {
		return stringSlices[0]
	}

	stringMaps := make([]map[string]bool, 0, len(stringSlices))

	for _, stringSlice := range stringSlices {
		stringMaps = append(stringMaps, stringSliceToMap(stringSlice))
	}

	resultMap := make(map[string]bool, 0)

	for str := range stringMaps[0] {
		for _, stringMap := range stringMaps[1:] {
			if _, ok := stringMap[str]; !ok {
				resultMap[str] = false
			} else {
				resultMap[str] = true
			}
		}
	}

	return mapToStringSlice(resultMap)
}

func uniqueStrings(stringSlices ...[]string) []string {
	if len(stringSlices) == 0 {
		return []string{}
	} else if len(stringSlices) == 1 {
		return stringSlices[0]
	}

	uniqueMap := map[string]bool{}

	for _, stringSlice := range stringSlices {
		for _, string := range stringSlice {
			uniqueMap[string] = true
		}
	}

	return mapToStringSlice(uniqueMap)
}

func stringSliceToMap(slice []string) map[string]bool {
	result := make(map[string]bool, 0)

	for _, string := range slice {
		result[string] = true
	}

	return result
}

func mapToStringSlice(stringMap map[string]bool) []string {
	stringSlice := make([]string, 0, len(stringMap))
	for string, value := range stringMap {
		if value {
			stringSlice = append(stringSlice, string)
		}
	}
	return stringSlice
}
