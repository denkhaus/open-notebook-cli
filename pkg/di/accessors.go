package di

import (
	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/samber/do/v2"
)

// Service getter helpers for easy access (only implemented services)

func GetConfig(injector do.Injector) config.Service {
	return do.MustInvoke[config.Service](injector)
}

func GetLogger(injector do.Injector) shared.Logger {
	return do.MustInvoke[shared.Logger](injector)
}

func GetAuth(injector do.Injector) shared.Auth {
	return do.MustInvoke[shared.Auth](injector)
}

func GetHTTPClient(injector do.Injector) shared.HTTPClient {
	return do.MustInvoke[shared.HTTPClient](injector)
}

func GetNotebookService(injector do.Injector) shared.NotebookService {
	return do.MustInvoke[shared.NotebookService](injector)
}

func GetSearchService(injector do.Injector) shared.SearchService {
	return do.MustInvoke[shared.SearchService](injector)
}

// Repository getters (only implemented ones)
func GetNotebookRepository(injector do.Injector) shared.NotebookRepository {
	return do.MustInvoke[shared.NotebookRepository](injector)
}

func GetNoteRepository(injector do.Injector) shared.NoteRepository {
	return do.MustInvoke[shared.NoteRepository](injector)
}

func GetSearchRepository(injector do.Injector) shared.SearchRepository {
	return do.MustInvoke[shared.SearchRepository](injector)
}

func GetSourceRepository(injector do.Injector) shared.SourceRepository {
	return do.MustInvoke[shared.SourceRepository](injector)
}

func GetSourceService(injector do.Injector) shared.SourceService {
	return do.MustInvoke[shared.SourceService](injector)
}

func GetChatRepository(injector do.Injector) shared.ChatRepository {
	return do.MustInvoke[shared.ChatRepository](injector)
}

func GetModelRepository(injector do.Injector) shared.ModelRepository {
	return do.MustInvoke[shared.ModelRepository](injector)
}

func GetModelService(injector do.Injector) shared.ModelService {
	return do.MustInvoke[shared.ModelService](injector)
}

func GetJobRepository(injector do.Injector) shared.JobRepository {
	return do.MustInvoke[shared.JobRepository](injector)
}

func GetJobService(injector do.Injector) shared.JobService {
	return do.MustInvoke[shared.JobService](injector)
}

func GetPodcastRepository(injector do.Injector) shared.PodcastRepository {
	return do.MustInvoke[shared.PodcastRepository](injector)
}

func GetPodcastService(injector do.Injector) shared.PodcastService {
	return do.MustInvoke[shared.PodcastService](injector)
}
