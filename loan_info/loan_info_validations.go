package loan_info

import (
	"errors"
	"go-file-parsing/validator"
	"strconv"
)

func hasValidLoanAmount(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	loanAmount, err := strconv.Atoi(cols[2])
	if err != nil {
		return nil, errors.New("loan amount is not a number")
	}
	if loanAmount <= 0 {
		return nil, errors.New("loan amount is not a positive number")
	}
	fundingAmount, err := strconv.Atoi(cols[3])
	if err != nil {
		return nil, errors.New("funding amount is not a number")
	}
	if fundingAmount <= 0 {
		return nil, errors.New("funding amount is not a positive number")
	}
	fundingInvAmt, err := strconv.Atoi(cols[4])
	if err != nil {
		return nil, errors.New("funding inv amt is not a number")
	}
	if fundingInvAmt <= 0 {
		return nil, errors.New("funding inv amt is not a positive number")
	}
	if fundingInvAmt != fundingAmount {
		return nil, errors.New("funding inv amt is not equal to funding amount")
	}
	return map[string]string{
		"loanAmount":    cols[2],
		"fundingAmount": cols[3],
		"fundingInvAmt": cols[4],
	}, nil
}

func hasValidInterestRate(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	rate, err := strconv.ParseFloat(cols[6], 64)
	if err != nil {
		return nil, errors.New("interest rate is not a number")
	}
	if rate < 5 || rate > 35 {
		return nil, errors.New("interest rate is not between 5% and 35%")
	}
	return map[string]string{
		"interestRate": cols[6],
	}, nil
}
