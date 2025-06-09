package validator

import "fmt"

func isValidSize(ctx *RowValidatorContext, cols []string) error {
	if len(cols) != ctx.Config.ExpectedColumns {
		return fmt.Errorf("expected %d columns, got %d", ctx.Config.ExpectedColumns, len(cols))
	}
	return nil
}
