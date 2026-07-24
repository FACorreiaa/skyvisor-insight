package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
	mcpdocs "github.com/FACorreiaa/Aviation-tracker/app/view/mcpdocs"
)

func (h *Handler) MCPPlaygroundPage(w http.ResponseWriter, r *http.Request) error {
	accessToken, err := h.apiAccessToken(r)
	if err != nil || accessToken == "" || h.service.API() == nil {
		http.Redirect(w, r, "/login?return_to=/mcp", http.StatusSeeOther)
		return nil
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	usage, err := h.service.API().GetUsage(ctx, accessToken)
	available := err == nil
	message := ""
	if err != nil {
		message = "Current agent usage is temporarily unavailable. Connection instructions remain valid."
		usage = apiclient.UsageSnapshot{}
	}
	return h.CreateLayout(w, r, "MCP and agents", mcpdocs.PlaygroundPage(usage, available, message)).Render(r.Context(), w)
}
