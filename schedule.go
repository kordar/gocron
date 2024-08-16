package gocron

import (
	"github.com/kordar/gologger"
	"github.com/robfig/cron/v3"
	"time"
)

type Gocron struct {
	cron        *cron.Cron
	items       map[string][]*EntryItem
	Initializer func(job Schedule) map[string]string
	CanRun      func(job Schedule) bool
}

func NewGocron(f1 func(job Schedule) map[string]string, f2 func(job Schedule) bool) *Gocron {
	return &Gocron{
		cron:        cron.New(),
		items:       make(map[string][]*EntryItem),
		Initializer: f1,
		CanRun:      f2,
	}
}

func (g *Gocron) Start() {
	g.cron.Start()
}

func (g *Gocron) GetItemById(id string) []*EntryItem {
	return g.items[id]
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

func (g *Gocron) Add(job Schedule) {
	cfg := g.Initializer(job)
	job.SetConfig(cfg)

	id := job.GetId()

	exec := func() {

		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("[gocron] recover %v", err)
			}
		}()

		if !g.CanRun(job) {
			return
		}

		job.Execute()
	}

	if g.items[id] != nil {
		logger.Warnf("[gocron] job %s exists", job.GetId())
		return
	}

	entries := make([]*EntryItem, 0)
	for i := 0; i < job.Duplicate(); i++ {
		if entryId, err := g.cron.AddFunc(job.GetSpec(), exec); err == nil {
			entry := &EntryItem{job.GetId(), time.Now(), entryId}
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

type EntryItem struct {
	Id      string       `json:"id"`
	RegTime time.Time    `json:"reg_time"`
	EntryId cron.EntryID `json:"entry_id"`
}

type Schedule interface {
	GetId() string
	GetSpec() string
	Execute()
	Config() map[string]string
	SetConfig(cfg map[string]string)
	Duplicate() int
}
