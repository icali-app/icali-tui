package icshelper

import (
	"time"

	ics "github.com/arran4/golang-ical"
)

func FindEventsForDay(cal *ics.Calendar, day time.Time) []*ics.VEvent {
	res := make([]*ics.VEvent, 0)

	for _, event := range cal.Events() {
		start, err := event.GetStartAt()

		if err == nil && isSameDay(day, start) {
			res = append(res, event)
			continue
		}

		// TODO: we will ignore everything else for now
		// end, err := event.GetStartAt()
		// if err == nil && isSameDay(cell.info.day, end) {
		// 	res = append(res, event)
		// 	continue
		// }
	}

	return res
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}
