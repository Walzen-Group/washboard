package werrors

import "fmt"

type SyncError struct {
	Err     error
	Context string
}

func (e *SyncError) Error() string {
	return e.Context + ": " + e.Err.Error()
}

func NewSyncError(err error, context string) *SyncError {
	return &SyncError{
		Err:     err,
		Context: context,
	}
}

type CannotInsertError struct {
	Context string
	Err     error
}

func (w *CannotInsertError) Error() string {
	return fmt.Sprintf("%s: %v", w.Context, w.Err)
}

func NewCannotInsertError(err error, info string) *CannotInsertError {
	return &CannotInsertError{
		Context: info,
		Err:     err,
	}
}

type DoesNotExistError struct {
	Context string
	Err     error
}

func (w *DoesNotExistError) Error() string {
	return fmt.Sprintf("%s: %v", w.Context, w.Err)
}

func NewDoesNotExistError(err error, info string) *DoesNotExistError {
	return &DoesNotExistError{
		Context: info,
		Err:     err,
	}
}

type AlreadyInProgressError struct {
	Context string
	Err     error
}

func (w *AlreadyInProgressError) Error() string {
	return fmt.Sprintf("%s: %v", w.Context, w.Err)
}

func NewAlreadyInProgressError(err error, info string) *AlreadyInProgressError {
	return &AlreadyInProgressError{
		Context: info,
		Err:     err,
	}
}
