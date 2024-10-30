package main

import (
	"fmt"
	"strconv"
	"strings"
)

func toFloat(value string) float64 {
	result, _ := strconv.ParseFloat(value, 64)
	return result
}

func parseResultPerline(result string) []string {
	resultSplitPerLine := strings.Split(result, "\n")
	// remove empty lines from the result
	for i, line := range resultSplitPerLine {
		if line == "" {
			resultSplitPerLine = append(resultSplitPerLine[:i], resultSplitPerLine[i+1:]...)
		}
	}

	return resultSplitPerLine
}

func formatResult(result string, index int) string {
	resultSplit := strings.Split(result, ":")
	return resultSplit[index]
}

// get value data size without unit 10.04TB -> 10.04
func getValueDataSizeWithoutUnit(value string) string {

	const units = "KMGTPEZY"

	for i := 0; i < len(value); i++ {
		if strings.Contains(units, string(value[i])) {

			dataSizeFormat := value[:i]

			// Converto to byte if the unit is not byte
			switch string(value[i]) {
			case "K":
				dataSizeFormat = fmt.Sprintf("%.2f", (toFloat(value[:i]) * 1024))
			case "M":
				dataSizeFormat = fmt.Sprintf("%.2f", (toFloat(value[:i]) * 1024 * 1024))
			case "G":
				dataSizeFormat = fmt.Sprintf("%.2f", (toFloat(value[:i]) * 1024 * 1024 * 1024))
			case "T":
				dataSizeFormat = fmt.Sprintf("%.2f", (toFloat(value[:i]) * 1024 * 1024 * 1024 * 1024))
			}

			return dataSizeFormat
		}
	}

	return value
}
