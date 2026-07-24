// Package apiclient is a compatibility shim over skyvisor-go-shared/apiclient
// so existing Aviation-tracker imports keep working.
package apiclient

import (
	"context"
	"io"

	shared "github.com/FACorreiaa/skyvisor-go-shared/apiclient"
	"github.com/FACorreiaa/skyvisor-go-shared/domain"
)

var (
	ErrNotFound        = shared.ErrNotFound
	ErrUnauthorized    = shared.ErrUnauthorized
	ErrPaymentRequired = shared.ErrPaymentRequired
	ErrConflict        = shared.ErrConflict
)

type (
	Client                    = shared.Client
	Me                        = shared.Me
	CheckoutSession           = shared.CheckoutSession
	APIError                  = shared.APIError
	LiveFlightQuery           = shared.LiveFlightQuery
	AirportBoardQuery         = shared.AirportBoardQuery
	AnalyticsQuery            = shared.AnalyticsQuery
	FlightLive                = domain.FlightLive
	Flight                    = domain.Flight
	SourceProvenance          = domain.SourceProvenance
	DataFreshness             = domain.DataFreshness
	InboundAircraft           = domain.InboundAircraft
	TripSegment               = domain.TripSegment
	ConnectionRisk            = domain.ConnectionRisk
	AutoWatchSkipped          = domain.AutoWatchSkipped
	AutoWatchResult           = domain.AutoWatchResult
	Trip                      = domain.Trip
	CreateTrip                = domain.CreateTrip
	ShareLink                 = domain.ShareLink
	PublicShare               = domain.PublicShare
	Watch                     = domain.Watch
	PAT                       = domain.PAT
	CreatedPAT                = domain.CreatedPAT
	CreateWatch               = domain.CreateWatch
	AssistantResponse         = domain.AssistantResponse
	Entitlements              = domain.Entitlements
	AirportBoard              = domain.AirportBoard
	AnalyticsReport           = domain.AnalyticsReport
	WhatIfRequest             = domain.WhatIfRequest
	WhatIfResult              = domain.WhatIfResult
	LogisticsQuery            = domain.LogisticsQuery
	LogisticsOverview         = domain.LogisticsOverview
	LogisticsDisruption       = domain.LogisticsDisruption
	Team                      = domain.Team
	CreateTeam                = domain.CreateTeam
	JoinTeam                  = domain.JoinTeam
	OperationsDashboard       = domain.OperationsDashboard
	OperationsSummary         = domain.OperationsSummary
	OperationsAttention       = domain.OperationsAttention
	OperationsWatch           = domain.OperationsWatch
	OperationsConnection      = domain.OperationsConnection
	OperationsFreshness       = domain.OperationsFreshness
	UsageSnapshot             = domain.UsageSnapshot
	OperationalCase           = domain.OperationalCase
	OperationalAlternative    = domain.OperationalAlternative
	CreateOperationalCase     = domain.CreateOperationalCase
	UpdateOperationalCase     = domain.UpdateOperationalCase
	DecisionRecord            = domain.DecisionRecord
	CreateDecisionRecord      = domain.CreateDecisionRecord
	RecordDecisionAction      = domain.RecordDecisionAction
	RecordDecisionOutcome     = domain.RecordDecisionOutcome
	OperationalCaseDetail     = domain.OperationalCaseDetail
	DecisionTrustMetrics      = domain.DecisionTrustMetrics
	TrustBreakdown            = domain.TrustBreakdown
	CaseAuditEvent            = domain.CaseAuditEvent
	WebhookIntegration        = domain.WebhookIntegration
	CreateWebhookIntegration  = domain.CreateWebhookIntegration
	WebhookIntegrationCreated = domain.WebhookIntegrationCreated
	WebhookDelivery           = domain.WebhookDelivery
)

func New(baseURL string) (*Client, error) {
	return shared.New(baseURL)
}

// Compile-time helpers so callers that used local methods keep working via Client alias.
var (
	_ = context.Background
	_ = io.EOF
)
