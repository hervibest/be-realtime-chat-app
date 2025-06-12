package message

const (
	Success = "Successfuly"

	InternalUserAuthNotFound = "user authentication data not found in context"

	ClientInvalidEmailOrPassword = "Make sure you have provide valid email or password"
	ClientUserAlreadyExist       = "Username or email already been used, please use another"
	ClientUnauthenticated        = "Unauthenticated, please try login again"
	ClientPermissionDenied       = "Permission denied for accessing this resource"

	RoomNotFound        = "Room not found for the given id/uuid"
	PeriodAlreadyClosed = "Attendance period already closed"

	InternalGracefulError = "Something wrong happened. Please try again"

	InternalNoRowsAffected = "no rows affected"
)
