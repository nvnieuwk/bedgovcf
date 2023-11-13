package bedgovcf

import (
	"fmt"
	"log"
	"math"
	"slices"
	"strconv"
	"strings"
)

func resolveField(configValues []string, bedValues []string, bedHeader []string) string {

	input := []string{}
	for _, v := range configValues {
		if strings.HasPrefix(v, "$") {
			var headerIndex int
			for j, w := range bedHeader {
				if w == v[1:] {
					headerIndex = j
					break
				}
			}
			input = append(input, bedValues[headerIndex])
			continue
		} else {
			input = append(input, v)
		}
	}

	function := ""
	if strings.HasPrefix(input[0], "~") {
		function = configValues[0][1:]
	} else {
		return strings.Join(input, " ")
	}

	switch function {
	case "round":
		// ~round <value>
		float, err := strconv.ParseFloat(input[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse the value (%v) to a float: %v", input[1], err)
		}
		round := math.Round(float)
		if round == -0 {
			round = 0
		}
		return fmt.Sprintf("%v", round)
	case "sum":
		// ~sum <value1> <value2> ...
		var sum float64
		for _, v := range input[1:] {
			float, err := strconv.ParseFloat(v, 64)
			if err != nil {
				log.Fatalf("Failed to parse the value (%v) to a float: %v", v, err)
			}
			sum += float
		}

		return strconv.FormatFloat(sum, 'g', -1, 64)
	case "min":
		// ~min <startValue> <valueToSubstract1> <valueToSubstract2> ...
		min, err := strconv.ParseFloat(input[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse the value (%v) to a float: %v", input[1], err)
		}
		for _, v := range input[2:] {
			float, err := strconv.ParseFloat(v, 64)
			if err != nil {
				log.Fatalf("Failed to parse the value (%v) to a float: %v", v, err)
			}
			min -= float
		}
		return strconv.FormatFloat(min, 'g', -1, 64)
	case "if":
		// ~if <value1> <operator> <value2> <value_if_true> <value_if_false>
		// supported operators: > < >= <= ==
		v1 := input[1]
		operator := input[2]
		v2 := input[3]
		vTrue := input[4]
		vFalse := input[5:]

		floatV1, err1 := strconv.ParseFloat(v1, 64)
		floatV2, err2 := strconv.ParseFloat(v2, 64)

		floatOperators := []string{"<", ">", "<=", ">="}
		if slices.Contains(floatOperators, operator) && (err1 != nil || err2 != nil) {
			log.Fatalf("Failed to parse the values (%v and %v) to a float: %v and %v", v1, v2, err1, err2)
		}

		vFalseResolved := ""
		if strings.HasPrefix(vFalse[0], "~") {
			vFalseResolved = resolveField(vFalse, bedValues, bedHeader)
		} else {
			vFalseResolved = strings.Join(vFalse, " ")
		}

		switch operator {
		case "<":
			if floatV1 < floatV2 {
				return vTrue
			} else {
				return vFalseResolved
			}
		case ">":
			if floatV1 > floatV2 {
				return vTrue
			} else {
				return vFalseResolved
			}
		case ">=":
			if floatV1 >= floatV2 {
				return vTrue
			} else {
				return vFalseResolved
			}
		case "<=":
			if floatV1 <= floatV2 {
				return vTrue
			} else {
				return vFalseResolved
			}
		case "==":
			if err1 == nil && err2 == nil {
				if floatV1 == floatV2 {
					return vTrue
				} else {
					return vFalseResolved
				}
			} else {
				if v1 == v2 {
					return vTrue
				} else {
					return vFalseResolved
				}
			}
		}
	}

	return ""
}
