package gocron

import (
	"github.com/kordar/gologger"
	"github.com/robfig/cron/v3"
	"time"
)

type Gocron struct {
	cron  *cron.Cron
	items map[string]*JobState
	cfg   map[string]string
}

func NewGocron(cfg map[string]string) *Gocron {
	if cfg == nil {
		cfg = map[string]string{}
	}
	return &Gocron{
		cron:  cron.New(),
		items: make(map[string]*JobState),
		cfg:   cfg,
	}
}

func (g *Gocron) Start() {
	g.cron.Start()
}

func (g *Gocron) Cron() *cron.Cron {
	return g.cron
}

func (g *Gocron) Reload(id string) {
	state := g.items[id]
	if state == nil {
		return
	}
	logger.Infof("reload the schedule named %s", id)
	g.Remove(id)
	delete(g.items, id)
	g.Add(state.Job)
	logger.Infof("reload the schedule finished named %s", id)
}

func (g *Gocron) Remove(id string) {
	if g.items[id] == nil {
		return
	}
	state := g.items[id]
	for _, entry := range state.Items {
		g.cron.Remove(entry.EntryId)
	}
	state.State = Shutdown
}

func (g *Gocron) Add(job Schedule) {

	id := job.GetId()
	if g.items[id] != nil {
		logger.Warnf("[gocron] job %s exists", job.GetId())
		return
	}

	state := &JobState{Job: job, State: Ready, Id: job.GetId()}
	entries := make([]JobStateItem, 0)
	for i := 0; i < job.Duplicate(); i++ {
		if entryId, err := g.cron.AddJob(job.GetSpec(), job.ToCronJob()); err == nil {
			entry := JobStateItem{time.Now(), entryId, map[string]interface{}{
				"spec":        job.GetSpec(),
				"description": job.Description(),
				"tag":         job.Tag(),
				"duplicate":   i, // 副本序号
				"node_id":     g.cfg["nodeId"],
				"node_addr":   g.cfg["nodeAddr"],
			}}
			state.State = Running
			entries = append(entries, entry)
			logger.Infof("[gocron] add job-%d %s success, config: %v", i, id, job.Config())
		} else {
			state.State = AddFailed
			state.Err = err
			logger.Errorf("[gocron] add job-%d %s fail, err: %v", i, id, err)
			break
		}

	}
	state.Items = entries
	g.items[id] = state
}

func (g *Gocron) GetEntryItemsById(id string) []JobStateItem {
	state := g.items[id]
	return state.Items
}

func (g *Gocron) Prints() []JobStateItem {
	entries := make([]JobStateItem, 0)
	for _, state := range g.items {
		entries = append(entries, state.Items...)
	}
	return entries
}

func (g *Gocron) State() []JobStateVO {
	data := make([]JobStateVO, 0)
	for id, state := range g.items {
		vo := JobStateVO{JobId: id, Data: state.Items, State: state.State.String()}
		if state.Err != nil {
			vo.Err = state.Err.Error()
		}
		data = append(data, vo)
	}
	return data
}
