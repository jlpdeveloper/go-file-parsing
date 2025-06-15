package loan_info

import "go-file-parsing/validator"

func passExtraData(_ *validator.RowValidatorContext, cols []string) (map[string]string, error) {
	result := make(map[string]string)

	// Check if we have enough columns
	if len(cols) < 93 { // avg_cur_bal is at index 92, which is the highest index we need
		return result, nil
	}

	// Add avg_cur_bal (index 92)
	if cols[92] != "" {
		result["avg_cur_bal"] = cols[92]
	}

	// Add application_type (index 69)
	if cols[69] != "" {
		result["application_type"] = cols[69]

		// Add annual_inc_joint (index 70) only if application_type is "Joint App"
		if cols[69] == "Joint App" && cols[70] != "" {
			result["annual_inc_joint"] = cols[70]
		}
	}

	// Add tot_coll_amt (index 74)
	if cols[74] != "" {
		result["tot_coll_amt"] = cols[74]
	}

	// Add acc_now_delinq (index 73)
	if cols[73] != "" {
		result["acc_now_delinq"] = cols[73]
	}

	return result, nil
}
