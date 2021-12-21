package cron

import (
	"github.com/robfig/cron/v3"

	"github.com/EvisuXiao/andrews-common/utils"
)

type Cron struct {
	cron *cron.Cron
	Jobs []*Job
}

type Job struct {
	id      cron.EntryID
	period  string
	handler func()
	preload bool
}

var c = &Cron{}

func NewCron() *Cron {
	if utils.IsEmpty(c.cron) {
		c.cron = cron.New()
	}
	return c
}

func (c *Cron) AddJob(period string, handler func(), preload bool) (cron.EntryID, error) {
	id, err := c.cron.AddFunc(period, handler)
	if utils.HasErr(err) {
		return 0, err
	}
	c.Jobs = append(c.Jobs, &Job{id, period, handler, preload})
	return id, nil
}

func (c *Cron) RemoveJob(id cron.EntryID) {
	c.cron.Remove(id)
	var jobs []*Job
	for _, j := range c.Jobs {
		if j.id != id {
			jobs = append(jobs, j)
		}
	}
	c.Jobs = jobs
}

func (c *Cron) HasJob() bool {
	return !utils.IsEmpty(c.Jobs)
}

func (c *Cron) Start() {
	if !c.HasJob() {
		return
	}
	for _, j := range c.Jobs {
		if j.preload {
			go j.handler()
		}
	}
	c.cron.Start()
}

func (c *Cron) Stop() {
	c.cron.Stop()
}
