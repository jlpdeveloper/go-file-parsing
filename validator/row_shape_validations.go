package validator

import "fmt"

func isValidSize(ctx *RowValidatorContext, cols []string) (map[string]string, error) {
	if len(cols) != ctx.Config.ExpectedColumns {
		return nil, fmt.Errorf("expected %d columns, got %d", ctx.Config.ExpectedColumns, len(cols))
	}
	return nil, nil
}
