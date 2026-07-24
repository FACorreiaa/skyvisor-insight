package operations

import (
	"fmt"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

func severityClass(value string) string {
	switch value {
	case "critical":
		return "bg-red-500/10 text-red-700 dark:text-red-400"
	case "high":
		return "bg-amber-500/10 text-amber-700 dark:text-amber-400"
	case "medium":
		return "bg-sky-500/10 text-sky-700 dark:text-sky-400"
	default:
		return "bg-muted text-muted-foreground"
	}
}

func stateClass(value string) string {
	switch value {
	case "evaluated", "resolved", "closed":
		return "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400"
	case "action_required", "executed":
		return "bg-amber-500/10 text-amber-700 dark:text-amber-400"
	case "rejected":
		return "bg-red-500/10 text-red-700 dark:text-red-400"
	default:
		return "bg-primary/10 text-primary"
	}
}

func caseCount(items []apiclient.OperationalCase, statuses ...string) int {
	wanted := make(map[string]struct{}, len(statuses))
	for _, status := range statuses {
		wanted[status] = struct{}{}
	}
	count := 0
	for _, item := range items {
		if _, ok := wanted[item.Status]; ok {
			count++
		}
	}
	return count
}

func percent(value float64) string { return fmt.Sprintf("%.1f%%", value*100) }

func money(minor int64, currency string) string {
	if minor == 0 {
		return "—"
	}
	if strings.TrimSpace(currency) == "" {
		currency = "value"
	}
	return fmt.Sprintf("%s %.2f", strings.ToUpper(currency), float64(minor)/100)
}

func flightList(values []string) string {
	if len(values) == 0 {
		return "No flights linked"
	}
	return strings.Join(values, " · ")
}

func relativeDeadline(value *time.Time) string {
	if value == nil || value.IsZero() {
		return "Not set"
	}
	return value.Local().Format("02 Jan 15:04")
}

func probability(value *float64) string {
	if value == nil {
		return "Unknown"
	}
	return percent(*value)
}

func canApprove(item apiclient.DecisionRecord) bool {
	return item.State == "proposed" && item.ApprovalRequired
}

func canExecute(item apiclient.DecisionRecord) bool {
	return item.State == "approved" || (item.State == "proposed" && !item.ApprovalRequired)
}

func canEvaluate(item apiclient.DecisionRecord) bool {
	return item.State == "executed" || (item.State == "proposed" && !item.ApprovalRequired)
}

func actorLabel(value string) string {
	if value == "agent" {
		return "Agent"
	}
	if value == "human" {
		return "Human"
	}
	return "System"
}
