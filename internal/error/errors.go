package internal_errors

type ErrInternal struct {
	Message string
}

func (e *ErrInternal) Error() string {
	return e.Message
}

type ErrNotFound struct {
	ErrInternal
}

func NewErrNotFound(message string) *ErrNotFound {
	return &ErrNotFound{
		ErrInternal: ErrInternal{
			Message: message,
		},
	}
}

type ErrBadRequest struct {
	ErrInternal
}

func NewErrBadRequest(message string) *ErrBadRequest {
	return &ErrBadRequest{
		ErrInternal: ErrInternal{
			Message: message,
		},
	}
}

type ErrForbidden struct {
	ErrInternal
}

func NewErrForbidden(message string) *ErrForbidden {
	return &ErrForbidden{
		ErrInternal: ErrInternal{
			Message: message,
		},
	}
}

type ErrUnauthorized struct {
	ErrInternal
}

func NewErrUnauthorized(message string) *ErrUnauthorized {
	return &ErrUnauthorized{
		ErrInternal: ErrInternal{
			Message: message,
		},
	}
}

type ErrInternalServer struct {
	ErrInternal
}

func NewErrInternalServer(message string) *ErrInternalServer {
	return &ErrInternalServer{
		ErrInternal: ErrInternal{
			Message: message,
		},
	}
}

type ErrBadGateway struct {
	ErrInternal
	Message string
}

func NewErrBadGateway(message string) *ErrBadGateway {
	return &ErrBadGateway{
		ErrInternal: ErrInternal{
			Message: message,
		},
	}
}
