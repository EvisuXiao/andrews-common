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
}

func NewCron() *Cron {
	return &Cron{cron: cron.New()}
}

func (c *Cron) AddJob(period string, handler func(), preload bool) (cron.EntryID, error) {
	id, err := c.cron.AddFunc(period, handler)
	if utils.HasErr(err) {
		return 0, err
	}
	if preload {
		go handler()
	}
	c.Jobs = append(c.Jobs, &Job{id, period, handler})
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
	//if !c.HasJob() {
	//	return
	//}
	c.cron.Start()
}

func (c *Cron) Stop() {
	c.cron.Stop()
}
