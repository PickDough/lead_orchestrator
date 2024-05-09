package errors

type DomainError struct {
}

func (de *DomainError) Error() string {
	return "domain error"
}

type WorkingHoursEndBeforeStartError struct {
	DomainError
}

func (de *WorkingHoursEndBeforeStartError) Error() string {
	return "working hours end must be greater than working hours start"
}
