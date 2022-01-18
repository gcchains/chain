

package syncer

type DoneEvent struct{}
type StartEvent struct{}
type FailedEvent struct{ Err error }
