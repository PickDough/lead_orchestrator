package domain

import "time"

type Client struct {
	Id                int64         `db:"id"`
	Name              string        `db:"name"`
	WorkingHoursStart time.Duration `db:"working_hours_start"`
	WorkingHoursEnd   time.Duration `db:"working_hours_end"`
	LeadCapacity      int           `db:"lead_capacity"`
	Priority          int           `db:"priority"`
}
