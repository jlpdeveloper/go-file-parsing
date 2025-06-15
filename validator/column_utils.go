package validator

import (
	"go-file-parsing/utils"
)

// PreprocessColumns trims whitespace from all columns in the input slice.
// It checks if trimming is necessary by checking if the first or last character is whitespace.
// This function is used to optimize string operations by trimming each column only once.
func PreprocessColumns(cols []string) []string {
	trimmedCols := make([]string, len(cols))
	for i, col := range cols {
		trimmedCols[i] = utils.TrimIfNeeded(col)
	}
	return trimmedCols
}
