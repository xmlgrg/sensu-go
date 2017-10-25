// Package pipelined provides the traditional Sensu event pipeline.
package pipelined

import (
	"github.com/sensu/sensu-go/types"

	"github.com/Knetic/govaluate"
)

func evaluateEventFilterStatement(event *types.Event, statement string) bool {
	expr, err := govaluate.NewEvaluableExpression(statement)
	if err != nil {
		logger.Warn("failed to create evaluable expression")
		return false
	}

	result, err := expr.Evaluate(map[string]interface{}{"event": event})
	if err != nil {
		logger.Warn("failed to evaluate filter")
		return false
	}

	match, ok := result.(bool)
	if !ok {
		logger.Warn("filters must evaluate to boolean values")
	}

	return match
}

func evaluateEventFilter(event *types.Event, filter types.EventFilter) bool {
	for _, statement := range filter.Statements {
		match := evaluateEventFilterStatement(event, statement)

		// Allow - One of the statements did not match, filter the event
		if filter.Action == types.EventFilterActionAllow && !match {
			return true
		}

		// Deny - One of the statements did not match, do not filter the event
		if filter.Action == types.EventFilterActionDeny && !match {
			return false
		}
	}

	// Allow - All of the statements matched, do not filter the event
	if filter.Action == types.EventFilterActionAllow {
		return false
	}

	// Deny - All of the statements matched, filter the event
	if filter.Action == types.EventFilterActionDeny {
		return true
	}

	// Something weird happened, let's not filter the event and log a warning message
	logger.Warn("pipelined not filtering event due to unhandled case")
	return false
}

// filterEvent filters a Sensu event, determining if it will continue
// through the Sensu pipeline.
func (p *Pipelined) filterEvent(handler *types.Handler, event *types.Event) bool {
	incident := p.isIncident(event)
	metrics := p.hasMetrics(event)

	// Do not filter the event if the event has metrics
	if metrics {
		return false
	}

	// Filter the event if it is not an incident
	if !incident {
		return true
	}

	// Do not filter the event if the handler has no event filters
	if len(handler.Filters) == 0 {
		return false
	}

	// Loop through all of the handler's event filters, if any filter evaluates to false then
	// do not filter the event
	for _, filter := range handler.Filters {
		filtered := evaluateEventFilter(event, filter)
		if !filtered {
			return false
		}
	}

	return true
}

// isIncident determines if an event indicates an incident.
func (p *Pipelined) isIncident(event *types.Event) bool {
	if event.Check.Status != 0 {
		return true
	}

	return false
}

// hasMetrics determines if an event has metric data.
func (p *Pipelined) hasMetrics(event *types.Event) bool {
	if event.Metrics != nil {
		return true
	}

	return false
}
