package checkins

import (
	"fmt"
	"github.com/statping/statping/utils"
	"time"
)

func (c *Checkin) Expected() time.Duration {
	last := c.LastHit()
	now := utils.Now()
	lastDir := now.Sub(last.CreatedAt)
	sub := time.Duration(c.Period() - lastDir)
	return sub
}

func (c *Checkin) Period() time.Duration {
	duration, _ := time.ParseDuration(fmt.Sprintf("%vs", c.Interval))
	return duration
}

// Grace will return the duration of the Checkin Grace Period (after service hasn't responded, wait a bit for a response)
func (c *Checkin) Grace() time.Duration {
	duration, _ := time.ParseDuration(fmt.Sprintf("%vs", c.GracePeriod))
	return duration
}

// Start will create a channel for the checkin checking go routine
func (c *Checkin) Start() {
	c.Running = make(chan bool)
}

// Close will stop the checkin routine
func (c *Checkin) Close() {
	if c.IsRunning() {
		close(c.Running)
	}
}

// IsRunning returns true if the checkin go routine is running
func (c *Checkin) IsRunning() bool {
	if c.Running == nil {
		return false
	}
	select {
	case <-c.Running:
		return false
	default:
		return true
	}
}

// String will return a Checkin API string
func (c *Checkin) String() string {
	return c.ApiKey
}

func (c *Checkin) Link() string {
	return fmt.Sprintf("%v/checkin/%v", "DOMAINHERE", c.ApiKey)
}
