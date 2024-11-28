package gocron

import (
	"github.com/kordar/gologger"
	"github.com/robfig/cron/v3"
	"time"
)

type Gocron struct {
	cron        *cron.Cron
	items       map[string][]*EntryItem
	cached      map[string]*CachedJob
	Initializer func(job Schedule) map[string]string
	CanRun      func(job Schedule) bool
}

func NewGocron(f1 func(job Schedule) map[string]string, f2 func(job Schedule) bool) *Gocron {
	return &Gocron{
		cron:        cron.New(),
		items:       make(map[string][]*EntryItem),
		cached:      make(map[string]*CachedJob),
		Initializer: f1,
		CanRun:      f2,
	}
}

func (g *Gocron) Start() {
	g.cron.Start()
}

func (g *Gocron) Stop() {
	g.cron.Stop()
}

func (g *Gocron) Cron() *cron.Cron {
	return g.cron
}

func (g *Gocron) Entry(id cron.EntryID) cron.Entry {
	return g.cron.Entry(id)
}

func (g *Gocron) Entries() []cron.Entry {
	return g.cron.Entries()
}

func (g *Gocron) GetItemById(id string) []*EntryItem {
	return g.items[id]
}

func (g *Gocron) Reload(id string) {
	if g.cached[id] != nil {
		g.Remove(id)
		cachedJob := g.cached[id]
		cachedJob.Job.SetConfig(map[string]string{
			"spec": "@every 3s",
		})
		g.AddWithJob(cachedJob.Job, cachedJob.CronJob)
	}
}

func (g *Gocron) Remove(id string) {
	if g.items[id] != nil {
		entries := g.items[id]
		for _, entry := range entries {
			g.cron.Remove(entry.EntryId)
		}
		delete(g.items, id)
	}
}

func (g *Gocron) DefaultJob(job Schedule) cron.Job {
	return cron.FuncJob(func() {

		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("[gocron] recover %v", err)
			}
		}()

		if !g.CanRun(job) {
			return
		}

		job.Execute()
	})
}

func (g *Gocron) AddWithChain(job Schedule, f func(funcJob cron.Job) cron.Job) {
	defaultJob := g.DefaultJob(job)
	chain := f(defaultJob)
	g.AddWithJob(job, chain)
}

func (g *Gocron) Add(job Schedule) {
	defaultJob := g.DefaultJob(job)
	g.AddWithJob(job, defaultJob)
}

func (g *Gocron) AddWithJob(job Schedule, funcJob cron.Job) {
	cfg := g.Initializer(job)
	job.SetConfig(cfg)

	id := job.GetId()

	if g.items[id] != nil {
		logger.Warnf("[gocron] job %s exists", job.GetId())
		return
	}

	g.cached[job.GetId()] = &CachedJob{job, funcJob}
	entries := make([]*EntryItem, 0)
	for i := 0; i < job.Duplicate(); i++ {
		if entryId, err := g.cron.AddJob(job.GetSpec(), funcJob); err == nil {
			entry := &EntryItem{job.GetId(), time.Now(), entryId, map[string]interface{}{
				"spec":        job.GetSpec(),
				"description": job.Description(),
				"tag":         job.Tag(),
				"duplicate":   i, // 副本序号
				"node_id":     cfg["nodeId"],
				"node_addr":   cfg["nodeAddr"],
			}}
			entries = append(entries, entry)
			logger.Infof("[gocron] add job-%d %s success, config: %v", i, id, job.Config())
		} else {
			logger.Errorf("[gocron] add job-%d %s fail, err: %v", i, id, err)
		}
	}

	g.items[id] = entries
}

func (g *Gocron) Prints() []*EntryItem {
	entries := make([]*EntryItem, 0)
	for _, items := range g.items {
		entries = append(entries, items...)
	}
	return entries
}

func (g *Gocron) State() []StateEntryItem {
	state := make([]StateEntryItem, 0)
	for id := range g.cached {
		s := StateEntryItem{
			Id:   id,
			Data: make([]*EntryItem, 0),
		}
		if g.items[id] == nil {
			s.State = "shutdown"
		} else {
			s.State = "running"
			s.Data = g.items[id]
		}
		state = append(state, s)
	}
	return state
}
