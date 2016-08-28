package chrono

import (
	"math"
	"time"

	pb "github.com/altipla-consulting/protobuf-defs/ptypes"
	"github.com/juju/errors"
)

type Date time.Time

func NewDate(year int, month time.Month, day int) Date {
	return Date(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(d).Format(`"2006-01-02"`)), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	t, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.UTC)
	if err != nil {
		return errors.Trace(err)
	}

	*d = Date(t)

	return nil
}

func (d Date) MarshalRQL() (interface{}, error) {
	t := time.Time(d)
	timeVal := float64(t.UnixNano()) / float64(time.Second)

	// use seconds-since-epoch precision if time.Time `t`
	// is before the oldest nanosecond time
	if t.Before(time.Unix(0, math.MinInt64)) {
		timeVal = float64(t.Unix())
	}

	return map[string]interface{}{
		"$reql_type$": "TIME",
		"epoch_time":  timeVal,
		"timezone":    t.Format("-07:00"),
	}, nil
}

func (d *Date) UnmarshalRQL(data interface{}) error {
	*d = Date(data.(time.Time))
	return nil
}

func (d Date) Before(other Date) bool {
	return time.Time(d).Before(time.Time(other))
}

func (d Date) Equal(other Date) bool {
	return time.Time(d).Equal(time.Time(other))
}

func (d Date) ToProto() *pb.Date {
	t := time.Time(d)

	return &pb.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}
}

func (d Date) Format(format string) string {
	return time.Time(d).Format(format)
}

func (d Date) AddDays(days int64) Date {
	return Date(time.Time(d).Add(time.Duration(days)*24*time.Hour + 12*time.Hour))
}

func Today(loc *time.Location) Date {
	year, month, day := time.Now().In(loc).Date()
	return NewDate(year, month, day)
}

func ZeroDate() Date {
	return Date(time.Time{})
}

func DateFromProto(proto *pb.Date) Date {
	if proto == nil {
		return ZeroDate()
	}

	return NewDate(int(proto.Year), time.Month(proto.Month), int(proto.Day))
}

func ParseDate(layout, value string) (Date, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return ZeroDate(), errors.Trace(err)
	}

	return Date(t), nil
}
