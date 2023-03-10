package between

import "time"

type Between struct {
	start time.Time
	end   time.Time
}

func NewBetween(start, end time.Time) *Between {
	return &Between{
		start: start,
		end:   end,
	}
}

func (b *Between) ListDays() []time.Time {

	start := b.start
	end := b.end

	var dates []time.Time

	for {
		if start.After(end) {
			return dates
		}

		date := start
		start = start.AddDate(0, 0, 1)

		dates = append(dates, date)
	}
}
