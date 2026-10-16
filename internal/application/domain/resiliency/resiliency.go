package resiliency

import "google.golang.org/grpc/codes"

var StatusCodeMap = map[uint32]codes.Code{
	0:  codes.OK,
	1:  codes.Canceled,
	2:  codes.Unknown,
	3:  codes.InvalidArgument,
	4:  codes.DeadlineExceeded,
	5:  codes.NotFound,
	6:  codes.AlreadyExists,
	7:  codes.PermissionDenied,
	8:  codes.ResourceExhausted,
	9:  codes.FailedPrecondition,
	10: codes.Aborted,
	11: codes.OutOfRange,
	12: codes.Unimplemented,
	13: codes.Internal,
	14: codes.Unavailable,
	15: codes.DataLoss,
	16: codes.Unauthenticated,
}
