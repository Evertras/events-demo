package stream

type EventID string

const EventIDUserRegistered EventID = "UserRegistered"

const headerKeyEventType = "evtype"
const headerKeySpanContext = "uber-trace-id"
