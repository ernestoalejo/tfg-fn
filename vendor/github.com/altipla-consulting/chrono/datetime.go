package chrono

import (
	"time"

	pb "github.com/golang/protobuf/ptypes/timestamp"
)

func Now() time.Time {
	return time.Now().UTC()
}

func DateTimeToProto(t time.Time) *pb.Timestamp {
	nanos := time.Duration(t.UnixNano()) * time.Nanosecond

	return &pb.Timestamp{
		Seconds: int64(nanos / time.Second),
		Nanos:   int32(nanos % time.Second),
	}
}

func DateTimeFromProto(proto *pb.Timestamp) time.Time {
	if proto == nil {
		return time.Time{}
	}

	return time.Unix(proto.Seconds, int64(proto.Nanos)).UTC()
}
