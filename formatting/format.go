package format

import (
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"fmt"
	"strconv"
	"strings"
)

func Name(fullName string, mp *dtos.MpWithExpenseCategories) {
	var firstName string
	var lastName string

	switch fullName {

	case "Vacant":
		firstName = constants.VACANT
		lastName = constants.VACANT

	default:
		var formattedStr string

		if strings.Contains(fullName, "Right Hon. ") {
			formattedStr = strings.ReplaceAll(fullName, "Right Hon. ", "")
		} else if strings.Contains(fullName, "Hon. ") {
			formattedStr = strings.ReplaceAll(fullName, "Hon. ", "")
		} else {
			formattedStr = fullName
		}
		names := strings.Split(formattedStr, ", ")

		firstName = names[1]
		lastName = names[0]
	}
	mp.MpName.FirstName = firstName
	mp.MpName.LastName = lastName
}

func ExpenseToFloat(expenseTotal string) (float64, error) {
	trimmedExpense := strings.TrimPrefix(expenseTotal, "$")
	trimmedExpense = strings.ReplaceAll(trimmedExpense, ",", "")

	expenseFloat, err := strconv.ParseFloat(trimmedExpense, 64)
	if err != nil {
		return -1, fmt.Errorf("Failed to parse string: %v", err)
	}

	return float64(expenseFloat), nil
}
