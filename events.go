package kizcool

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Event is an interface for any event
type Event interface {
	Kind() string
}

// GenericEvent is the minimal set of fields shared by all events
type GenericEvent struct {
	Timestamp int    `json:"timestamp,omitempty"`
	Name      string `json:"name,omitempty"`
}

// Kind returns a partial text description of the event
func (e *GenericEvent) Kind() string {
	return fmt.Sprintf("%v", e)
}

// ExecutionEvent is the minimal set of fields shared by all Execution events
type ExecutionEvent struct {
	ExecID   ExecID `json:"execID,omitempty"`
	SetupOID string `json:"setupOID,omitempty"`
	SubType  int    `json:"subType,omitempty"`
	Type     int    `json:"type,omitempty"`
}

// ExecutionRegisteredEvent indicates an execution has been registered
type ExecutionRegisteredEvent struct {
	GenericEvent
	ExecutionEvent
	Label     string   `json:"label,omitempty"`
	Metadata  string   `json:"metadata,omitempty"`
	TriggerID string   `json:"triggerId,omitempty"`
	Actions   []Action `json:"actions,omitempty"`
}

// ExecutionStateChangedEvent indicates a change in the state of an execution
type ExecutionStateChangedEvent struct {
	GenericEvent
	ExecutionEvent
	NewState        string `json:"newState,omitempty"`
	OldState        string `json:"oldState,omitempty"`
	OwnerKey        string `json:"ownerKey,omitempty"`
	TimeToNextState int    `json:"timeToNextState,omitempty"`
}

// CommandExecutionStateChangedEvent indicates a change in the state of the execution of a command
type CommandExecutionStateChangedEvent struct {
	GenericEvent
	DeviceURL       DeviceURL `json:"deviceURL,omitempty"`
	ExecID          ExecID    `json:"execID,omitempty"`
	SetupOID        string    `json:"setupOID,omitempty"`
	NewState        string    `json:"newState,omitempty"`
	FailureType     string    `json:"failureType,omitempty"`
	FailureTypeCode int       `json:"failureTypeCode,omitempty"`
	Rank            int       `json:"rank,omitempty"`
}

// GatewayEvent indicates an event related to a gateway
type GatewayEvent struct {
	GatewayID string `json:"gatewayId,omitempty"`
}

// GatewaySynchronizationStartedEvent indicates the start of synchronization of a gateway
type GatewaySynchronizationStartedEvent struct {
	GenericEvent
	GatewayEvent
}

// GatewaySynchronizationEndedEvent indicates the end of synchronization of a gateway
type GatewaySynchronizationEndedEvent struct {
	GenericEvent
	GatewayEvent
}

// GatewayDownEvent indicates the gateway has become unreachable
type GatewayDownEvent struct {
	GenericEvent
	GatewayEvent
}

// GatewayAliveEvent indicates the gateway is accessible again
type GatewayAliveEvent struct {
	GenericEvent
	GatewayEvent
}

// RefreshAllDevicesStatesCompletedEvent indicates the end of a request to get the state of all devices
type RefreshAllDevicesStatesCompletedEvent struct {
	GenericEvent
	GatewayEvent
	ProtocolType int `json:"protocolType,omitempty"`
}

// DeviceStateChangedEvent indicates a change in the state of a device
type DeviceStateChangedEvent struct {
	GenericEvent
	SetupOID     string        `json:"setupOID,omitempty"`
	DeviceURL    DeviceURL     `json:"deviceURL,omitempty"`
	DeviceStates []DeviceState `json:"deviceStates,omitempty"`
}

// EndUserLoginEvent happens when a user authenticates
type EndUserLoginEvent struct {
	GenericEvent
	SetupOID      string `json:"setupOID,omitempty"`
	UserID        string `json:"userId,omitempty"`
	UserAgentType string `json:"userAgentType,omitempty"`
}

// Events is a slide of Event, used for unmarshalling several events of unknown type from json
type Events []Event

// UnmarshalJSON unmarshals an event from json, detecting the right event type
func (events *Events) UnmarshalJSON(data []byte) error {
	// this just splits up the JSON array into the raw JSON for each object
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		log.WithFields(log.Fields{
			"data": data,
			"err":  err,
		}).Debug("Error splitting into raw items")
		return fmt.Errorf("Error splitting json into raw items. %w", err)
	}
	for _, r := range raw {
		// unamrshal into a map to check the "Name" field
		var obj map[string]interface{}
		err := json.Unmarshal(r, &obj)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"data": r,
			}).Debug("Error retrieving Name field")
			return fmt.Errorf("Error retrieving Name field. %w", err)
		}

		eventType := ""
		if t, ok := obj["name"].(string); ok {
			eventType = t
		}

		// unmarshal again into the correct type
		var actual Event
		switch eventType {
		case "ExecutionRegisteredEvent":
			actual = &ExecutionRegisteredEvent{}
		case "ExecutionStateChangedEvent":
			actual = &ExecutionStateChangedEvent{}
		case "CommandExecutionStateChangedEvent":
			actual = &CommandExecutionStateChangedEvent{}
		case "GatewaySynchronizationStartedEvent":
			actual = &GatewaySynchronizationStartedEvent{}
		case "GatewaySynchronizationEndedEvent":
			actual = &GatewaySynchronizationEndedEvent{}
		case "GatewayDownEvent":
			actual = &GatewayDownEvent{}
		case "GatewayAliveEvent":
			actual = &GatewayAliveEvent{}
		case "RefreshAllDevicesStatesCompletedEvent":
			actual = &RefreshAllDevicesStatesCompletedEvent{}
		case "DeviceStateChangedEvent":
			actual = &DeviceStateChangedEvent{}
		case "EndUserLoginEvent":
			actual = &EndUserLoginEvent{}
		default:
			return fmt.Errorf("Cannot unmarshal unknown event of type %s: %v", eventType, obj)
		}

		err = json.Unmarshal(r, actual)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"data": r,
			}).Debug("Error unmarshalling into struct")
			return fmt.Errorf("Error unmarshalling into struct. %w", err)
		}
		*events = append(*events, actual)
	}
	return nil
}
