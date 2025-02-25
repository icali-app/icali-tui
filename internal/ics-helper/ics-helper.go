package icshelper

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

func NewId() string {
	return fmt.Sprintf("%s@tui-cal", uuid.New()) // TODO: domain?
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}


func FindEventsForDay(cal *ics.Calendar, day time.Time) []*ics.VEvent {
	return findThingForDay(cal, day, func(c *ics.Calendar) []*ics.VEvent {
		return c.Events()
	})
}


func FindTodosForDay(cal *ics.Calendar, day time.Time) []*ics.VTodo {
	return findThingForDay(cal, day, func(c *ics.Calendar) []*ics.VTodo {
		return c.Todos()
	})
}

// separate interface needed since ics.ComponentBase is a struct
type HasGetStartAt interface {
	GetStartAt() (time.Time, error)	
}

func findThingForDay[T HasGetStartAt](cal *ics.Calendar, day time.Time, selector func(cal *ics.Calendar) []T) []T {
	res := make([]T, 0)

	for _, thing := range selector(cal) {
		start, err := thing.GetStartAt()

		if err == nil && isSameDay(day, start) {
			res = append(res, thing)
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
