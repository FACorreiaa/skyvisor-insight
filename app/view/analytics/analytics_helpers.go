package analytics

import (
	"fmt"
	"strconv"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func windowValue(requested, reported int) string {
	if requested > 0 {
		return strconv.Itoa(requested)
	}
	if reported > 0 {
		return strconv.Itoa(reported)
	}
	return "7"
}

func exportURL(q apiclient.AnalyticsQuery) string {
	u := "/analytics/export?airline=" + q.Airline + "&airport=" + q.Airport
	if q.WindowDays > 0 {
		u += "&window_days=" + strconv.Itoa(q.WindowDays)
	}
	return u
}

func pct(v float64) string {
	return fmt.Sprintf("%.1f%%", v)
}

func num(v float64) string {
	return fmt.Sprintf("%.1f", v)
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
