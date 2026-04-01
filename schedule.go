package gocron

import (
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

type InitializeFunction func(job Schedule) map[string]string
type RuntimeFunction func(job Schedule) bool

type Gocron struct {
	cron               *cron.Cron
	items              map[string]*JobState
	cfg                map[string]string
	initializeFunction InitializeFunction
	runtimeFunction    RuntimeFunction
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

func (g *Gocron) SetInitFn(f InitializeFunction) {
	g.initializeFunction = f
}

func (g *Gocron) SetRuntimeFn(f RuntimeFunction) {
	g.runtimeFunction = f
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
	slog.Info("reload schedule start", "id", id)
	g.Remove(id)
	delete(g.items, id)
	g.Add(state.Job)
	slog.Info("reload schedule done", "id", id)
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
		slog.Warn("[gocron] job exists", "id", job.GetId())
		return
	}

	if g.initializeFunction != nil {
		cfg := g.initializeFunction(job)
		job.SetConfig(cfg)
	}

	state := &JobState{Job: job, State: Ready, Id: job.GetId()}
	entries := make([]JobStateItem, 0)
	for i := 0; i < job.Duplicate(); i++ {
		cronJob := job.ToCronJob()
		if cronJob == nil {
			cronJob = GenCronJobWithCanRun(job, g.runtimeFunction)
		}
		if entryId, err := g.cron.AddJob(job.GetSpec(), cronJob); err == nil {
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
			slog.Info("[gocron] add job success", "duplicate", i, "id", id, "config", job.Config())
		} else {
			state.State = AddFailed
			state.Err = err
			slog.Error("[gocron] add job fail", "duplicate", i, "id", id, "err", err)
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
