package kizcool

import "encoding/json"

import "fmt"

// Event is an interface for any event
type Event interface {
	event()
}

// GenericEvent is the minimal set of fields shared by all events
type GenericEvent struct {
	Timestamp int    `json:"timestamp,omitempty"`
	Name      string `json:"name,omitempty"`
}

// event does nothing, just implements the interface
func (e *GenericEvent) event() { return }

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
		return err
	}
	for _, r := range raw {
		// unamrshal into a map to check the "Name" field
		var obj map[string]interface{}
		err := json.Unmarshal(r, &obj)
		if err != nil {
			return err
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
		case "GatewaySynchronizationStartedEvent":
			actual = &GatewaySynchronizationStartedEvent{}
		case "GatewaySynchronizationEndedEvent":
			actual = &GatewaySynchronizationEndedEvent{}
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
			return err
		}
		*events = append(*events, actual)
	}
	return nil
}
