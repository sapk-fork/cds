package sdk

import (
	"encoding/json"
	"fmt"
)

const (
	WebHookModelName           = "WebHook"
	RepositoryWebHookModelName = "RepositoryWebHook"
	SchedulerModelName         = "Scheduler"
	GitPollerModelName         = "Git Repository Poller"
)

var (
	WebHookModel = WorkflowHookModel{
		Author:     "CDS",
		Type:       WorkflowHookModelBuiltin,
		Identifier: "github.com/ovh/cds/hook/builtin/webhook",
		Name:       WebHookModelName,
		Icon:       "Linkify",
		DefaultConfig: WorkflowNodeHookConfig{
			"method": {
				Value:        "POST",
				Configurable: true,
			},
		},
	}

	RepositoryWebHookModel = WorkflowHookModel{
		Author:     "CDS",
		Type:       WorkflowHookModelBuiltin,
		Identifier: "github.com/ovh/cds/hook/builtin/repositorywebhook",
		Name:       RepositoryWebHookModelName,
		Icon:       "Linkify",
		DefaultConfig: WorkflowNodeHookConfig{
			"method": {
				Value:        "POST",
				Configurable: false,
			},
		},
	}

	GitPollerModel = WorkflowHookModel{
		Author:     "CDS",
		Type:       WorkflowHookModelBuiltin,
		Identifier: "github.com/ovh/cds/hook/builtin/poller",
		Name:       GitPollerModelName,
		Icon:       "git square",
	}

	SchedulerModel = WorkflowHookModel{
		Author:     "CDS",
		Type:       WorkflowHookModelBuiltin,
		Identifier: "github.com/ovh/cds/hook/builtin/scheduler",
		Name:       SchedulerModelName,
		Icon:       "fa-clock-o",
		DefaultConfig: WorkflowNodeHookConfig{
			"cron": {
				Value:        "0 * * * *",
				Configurable: true,
			},
			"timezone": {
				Value:        "UTC",
				Configurable: true,
			},
			"payload": {
				Value:        "{}",
				Configurable: true,
			},
		},
	}
)

// Hook used to link a git repository to a given pipeline
type Hook struct {
	ID            int64    `json:"id"`
	UID           string   `json:"uid"`
	Pipeline      Pipeline `json:"pipeline"`
	ApplicationID int64    `json:"application_id"`
	Kind          string   `json:"kind"`
	Host          string   `json:"host"`
	Project       string   `json:"project"`
	Repository    string   `json:"repository"`
	Enabled       bool     `json:"enabled"`
	Link          string   `json:"link"`
}

// AddHook creates a new hook between a pipeline and a repository
func AddHook(a *Application, p *Pipeline, host string, project string, repository string) (*Hook, error) {
	h := Hook{
		Pipeline:      *p,
		ApplicationID: a.ID,
		Kind:          "stash",
		Host:          host,
		Project:       project,
		Repository:    repository,
	}

	data, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("/project/%s/application/%s/pipeline/%s/hook", p.ProjectKey, a.Name, p.Name)
	data, _, err = Request("POST", uri, data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &h); err != nil {
		return nil, err
	}

	return &h, nil
}

// GetHooks lists all hooks related to a pipeline
func GetHooks(project, application, pipeline string) ([]Hook, error) {
	var hooks []Hook

	uri := fmt.Sprintf("/project/%s/application/%s/pipeline/%s/hook", project, application, pipeline)

	data, code, err := Request("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	if code > 300 {
		return nil, fmt.Errorf("HTTP %d", code)
	}

	err = json.Unmarshal(data, &hooks)
	if err != nil {
		return nil, err
	}

	return hooks, nil
}

// DeleteHook remove a hook previously created
func DeleteHook(project, application, pipeline string, id int64) error {
	uri := fmt.Sprintf("/project/%s/application/%s/pipeline/%s/hook/%d", project, application, pipeline, id)

	_, code, err := Request("DELETE", uri, nil)
	if err != nil {
		return err
	}

	if code >= 300 {
		return fmt.Errorf("HTTP %d", code)
	}

	return nil
}

// GetDefaultHookModel return the workflow hook model by its name
func GetDefaultHookModel(modelName string) WorkflowHookModel {
	switch modelName {
	case SchedulerModelName:
		return SchedulerModel
	case RepositoryWebHookModelName:
		return RepositoryWebHookModel
	case WebHookModelName:
		return WebHookModel
	case GitPollerModelName:
		return GitPollerModel
	}

	return WebHookModel
}
