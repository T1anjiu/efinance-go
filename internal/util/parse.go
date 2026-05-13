package util

import (
	"strconv"
	"strings"
)

func ParseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")
	if s == "" || s == "-" || s == "--" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func ParseInt(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" || s == "--" {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0
		}
		return int64(f)
	}
	return v
}

func ParseKline(line string) (date, open, close, high, low string, volume, amount string, amplitude, changeRate, changeAmt, turnover string) {
	parts := strings.Split(line, ",")
	if len(parts) < 11 {
		return
	}
	date = parts[0]
	open = parts[1]
	close_ := parts[2]
	high = parts[3]
	low = parts[4]
	volume = parts[5]
	amount = parts[6]
	amplitude = parts[7]
	changeRate = parts[8]
	changeAmt = parts[9]
	turnover = parts[10]
	return date, open, close_, high, low, volume, amount, amplitude, changeRate, changeAmt, turnover
}

func ParseNDaysKline(line string) (date, open, close, high, low string, volume, amount string) {
	parts := strings.Split(line, ",")
	if len(parts) < 7 {
		return
	}
	return parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6]
}

func ParseBill(line string) (date string, mainNetInflow, smallNetInflow, midNetInflow, largeNetInflow, hugeNetInflow string, extra []string) {
	parts := strings.Split(line, ",")
	if len(parts) < 6 {
		return
	}
	date = parts[0]
	mainNetInflow = parts[1]
	smallNetInflow = parts[2]
	midNetInflow = parts[3]
	largeNetInflow = parts[4]
	hugeNetInflow = parts[5]
	if len(parts) > 6 {
		extra = parts[6:]
	}
	return
}
