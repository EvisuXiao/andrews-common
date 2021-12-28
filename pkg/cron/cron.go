package cron

import (
	"errors"
	"sync"

	"github.com/robfig/cron/v3"

	"github.com/EvisuXiao/andrews-common/utils"
)

type job struct {
	id      cron.EntryID
	period  string
	handler func()
}

var (
	c    = cron.New()
	mu   sync.Mutex
	jobs = make(map[string]*job)
)

func Start() {
	c.Start()
}

func Stop() {
	c.Stop()
}

func AddJob(name, period string, handler func(), preload bool) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := jobs[name]; ok {
		return errors.New("job name is already exist")
	}
	id, err := c.AddFunc(period, handler)
	if utils.HasErr(err) {
		return err
	}
	jobs[name] = &job{id, period, handler}
	if preload {
		go handler()
	}
	return nil
}

func RemoveJob(name string) error {
	mu.Lock()
	defer mu.Unlock()
	j, ok := jobs[name]
	if !ok {
		return errors.New("job name is not found")
	}
	c.Remove(j.id)
	return nil
}

func GetJobNames() []string {
	var names []string
	for name, _ := range jobs {
		utils.SliceAddStringItem(&names, name)
	}
	return names
}
