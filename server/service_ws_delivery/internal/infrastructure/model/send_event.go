package model

/**
 * Action send event model
 */
type ActionSendEvent struct {
	Action int         `json:"action"`
	Data   interface{} `json:"data"`
}
