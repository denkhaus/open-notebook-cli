package services

import (
	"context"
	"io"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// Core service interfaces that will be implemented privately and exposed publicly

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Sync() error
	With(fields ...interface{}) Logger
	WithContext(ctx context.Context) Logger
}

// Auth interface for authentication and token management
type Auth interface {
	Authenticate(ctx context.Context) error
	GetToken(ctx context.Context) (string, error)
	InvalidateToken(ctx context.Context) error
	IsAuthenticated(ctx context.Context) bool
	RefreshToken(ctx context.Context) error
	SetPassword(password string)
}

// HTTPClient interface for API communication
type HTTPClient interface {
	Get(ctx context.Context, endpoint string) (*models.Response, error)
	Post(ctx context.Context, endpoint string, body interface{}) (*models.Response, error)
	Put(ctx context.Context, endpoint string, body interface{}) (*models.Response, error)
	Delete(ctx context.Context, endpoint string) (*models.Response, error)
	PostMultipart(ctx context.Context, endpoint string, fields map[string]string, files map[string]io.Reader) (*models.Response, error)
	Stream(ctx context.Context, endpoint string, body interface{}) (<-chan []byte, error)
	SetAuth(token string)
	WithTimeout(timeout time.Duration) HTTPClient
}

// Repository interfaces for domain operations

// NotebookRepository interface for notebook data operations
type NotebookRepository interface {
	List(ctx context.Context) ([]*models.Notebook, error)
	Create(ctx context.Context, notebook *models.NotebookCreate) (*models.Notebook, error)
	Get(ctx context.Context, id string) (*models.Notebook, error)
	Update(ctx context.Context, id string, notebook *models.NotebookUpdate) (*models.Notebook, error)
	Delete(ctx context.Context, id string) error
	AddSource(ctx context.Context, notebookID, sourceID string) error
	RemoveSource(ctx context.Context, notebookID, sourceID string) error
}

// NoteRepository interface for note data operations
type NoteRepository interface {
	List(ctx context.Context, notebookID string, limit, offset int) ([]*models.Note, error)
	Create(ctx context.Context, note *models.NoteCreate) (*models.Note, error)
	Get(ctx context.Context, id string) (*models.Note, error)
	Update(ctx context.Context, id string, note *models.NoteUpdate) (*models.Note, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, notebookID, query string) ([]*models.Note, error)
}

// SearchRepository interface for search operations
type SearchRepository interface {
	Search(ctx context.Context, req *models.SearchRequest) (*models.SearchResponse, error)
	Ask(ctx context.Context, req *models.AskRequest) (<-chan *models.StreamChunk, error)
	AskSimple(ctx context.Context, req *models.AskRequest) (*models.AskResponse, error)
}

// SourceRepository interface for source management
type SourceRepository interface {
	List(ctx context.Context, limit, offset int) ([]*models.SourceListResponse, error)
	Create(ctx context.Context, source *models.SourceCreate) (*models.Source, error)
	CreateFromJSON(ctx context.Context, source *models.SourceCreate) (*models.Source, error)
	Get(ctx context.Context, id string) (*models.Source, error)
	Update(ctx context.Context, id string, source *models.SourceUpdate) (*models.Source, error)
	Delete(ctx context.Context, id string) error
	GetStatus(ctx context.Context, id string) (*models.SourceStatusResponse, error)
	Download(ctx context.Context, id string) (io.ReadCloser, error)
}

// ModelRepository interface for AI model management
type ModelRepository interface {
	List(ctx context.Context) ([]*models.Model, error)
	Create(ctx context.Context, model *models.ModelCreate) (*models.Model, error)
	Delete(ctx context.Context, id string) error
	GetDefaults(ctx context.Context) (*models.DefaultModelsResponse, error)
	SetDefaults(ctx context.Context, defaults *models.DefaultModelsResponse) error
	GetProviders(ctx context.Context) (*models.ProviderAvailabilityResponse, error)
}

// TransformationRepository interface for transformation management
type TransformationRepository interface {
	List(ctx context.Context) ([]*models.Transformation, error)
	Create(ctx context.Context, transformation *models.TransformationCreate) (*models.Transformation, error)
	Get(ctx context.Context, id string) (*models.Transformation, error)
	Update(ctx context.Context, id string, transformation *models.TransformationUpdate) (*models.Transformation, error)
	Delete(ctx context.Context, id string) error
	Execute(ctx context.Context, req *models.TransformationExecuteRequest) (*models.TransformationExecuteResponse, error)
}

// EmbeddingRepository interface for embedding operations
type EmbeddingRepository interface {
	Embed(ctx context.Context, req *models.EmbedRequest) (*models.EmbedResponse, error)
	Rebuild(ctx context.Context, req *models.RebuildRequest) (*models.RebuildResponse, error)
	GetRebuildStatus(ctx context.Context, commandID string) (*models.RebuildStatusResponse, error)
}

// SettingsRepository interface for application settings
type SettingsRepository interface {
	Get(ctx context.Context) (*models.SettingsResponse, error)
	Update(ctx context.Context, settings *models.SettingsUpdate) (*models.SettingsResponse, error)
}

// ContextRepository interface for context operations
type ContextRepository interface {
	Get(ctx context.Context, req *models.ContextRequest) (*models.ContextResponse, error)
}

// InsightsRepository interface for insights operations
type InsightsRepository interface {
	GetSourceInsights(ctx context.Context, sourceID string) ([]*models.SourceInsightResponse, error)
	CreateSourceInsight(ctx context.Context, sourceID string, req *models.CreateSourceInsightRequest) (*models.SourceInsightResponse, error)
	SaveInsightAsNote(ctx context.Context, insightID string, req *models.SaveAsNoteRequest) (*models.Note, error)
	GetDefaultPrompt(ctx context.Context) (*models.DefaultPromptResponse, error)
	UpdateDefaultPrompt(ctx context.Context, req *models.DefaultPromptUpdate) (*models.DefaultPromptResponse, error)
}

// Service interfaces for business logic

// NotebookService interface for notebook business logic
type NotebookService interface {
	Repository() NotebookRepository
	ListNotebooks(ctx context.Context) ([]*models.Notebook, error)
	CreateNotebook(ctx context.Context, name, description string) (*models.Notebook, error)
	GetNotebook(ctx context.Context, id string) (*models.Notebook, error)
	UpdateNotebook(ctx context.Context, id string, name, description *string, archived *bool) (*models.Notebook, error)
	DeleteNotebook(ctx context.Context, id string) error
	AddSourceToNotebook(ctx context.Context, notebookID, sourceID string) error
	RemoveSourceFromNotebook(ctx context.Context, notebookID, sourceID string) error
}

// NoteService interface for note business logic
type NoteService interface {
	Repository() NoteRepository
	ListNotes(ctx context.Context, notebookID string, limit, offset int) ([]*models.Note, error)
	CreateNote(ctx context.Context, notebookID, content, title string, noteType string) (*models.Note, error)
	GetNote(ctx context.Context, id string) (*models.Note, error)
	UpdateNote(ctx context.Context, id string, content, title *string) (*models.Note, error)
	DeleteNote(ctx context.Context, id string) error
	SearchNotes(ctx context.Context, notebookID, query string) ([]*models.Note, error)
}

// SearchService interface for search business logic
type SearchService interface {
	Repository() SearchRepository
	Search(ctx context.Context, query string, options *models.SearchOptions) (*models.SearchResponse, error)
	Ask(ctx context.Context, question string, options *models.AskOptions) (<-chan *models.StreamChunk, error)
	AskSimple(ctx context.Context, question string, options *models.AskOptions) (*models.AskResponse, error)
}

// SourceService interface for source business logic
type SourceService interface {
	Repository() SourceRepository
	ListSources(ctx context.Context, limit, offset int) ([]*models.SourceListResponse, error)
	AddSourceFromLink(ctx context.Context, link string, options *models.SourceOptions) (*models.Source, error)
	AddSourceFromUpload(ctx context.Context, filename string, file io.Reader, options *models.SourceOptions) (*models.Source, error)
	AddSourceFromText(ctx context.Context, text, title string, options *models.SourceOptions) (*models.Source, error)
	GetSource(ctx context.Context, id string) (*models.Source, error)
	UpdateSource(ctx context.Context, id string, title string, topics []string) (*models.Source, error)
	DeleteSource(ctx context.Context, id string) error
	GetSourceStatus(ctx context.Context, id string) (*models.SourceStatusResponse, error)
	DownloadSource(ctx context.Context, id string) (io.ReadCloser, error)
}

// ModelService interface for model business logic
type ModelService interface {
	Repository() ModelRepository
	ListModels(ctx context.Context) ([]*models.Model, error)
	AddModel(ctx context.Context, model *models.ModelCreate) (*models.Model, error)
	DeleteModel(ctx context.Context, id string) error
	GetDefaults(ctx context.Context) (*models.DefaultModelsResponse, error)
	SetDefaults(ctx context.Context, defaults *models.DefaultModelsResponse) error
	GetProviders(ctx context.Context) (*models.ProviderAvailabilityResponse, error)
}

// TransformationService interface for transformation business logic
type TransformationService interface {
	Repository() TransformationRepository
	ListTransformations(ctx context.Context) ([]*models.Transformation, error)
	CreateTransformation(ctx context.Context, transformation *models.TransformationCreate) (*models.Transformation, error)
	GetTransformation(ctx context.Context, id string) (*models.Transformation, error)
	UpdateTransformation(ctx context.Context, id string, transformation *models.TransformationUpdate) (*models.Transformation, error)
	DeleteTransformation(ctx context.Context, id string) error
	ExecuteTransformation(ctx context.Context, req *models.TransformationExecuteRequest) (*models.TransformationExecuteResponse, error)
}

// EmbeddingService interface for embedding business logic
type EmbeddingService interface {
	Repository() EmbeddingRepository
	EmbedItem(ctx context.Context, itemID, itemType string, async bool) (*models.EmbedResponse, error)
	RebuildEmbeddings(ctx context.Context, mode string, includeSources, includeNotes, includeInsights bool) (*models.RebuildResponse, error)
	GetRebuildStatus(ctx context.Context, commandID string) (*models.RebuildStatusResponse, error)
}

// SettingsService interface for settings business logic
type SettingsService interface {
	Repository() SettingsRepository
	GetSettings(ctx context.Context) (*models.SettingsResponse, error)
	UpdateSettings(ctx context.Context, settings *models.SettingsUpdate) (*models.SettingsResponse, error)
}

// ContextService interface for context business logic
type ContextService interface {
	Repository() ContextRepository
	GetContext(ctx context.Context, notebookID string, config *models.ContextConfig) (*models.ContextResponse, error)
}

// InsightsService interface for insights business logic
type InsightsService interface {
	Repository() InsightsRepository
	GetSourceInsights(ctx context.Context, sourceID string) ([]*models.SourceInsightResponse, error)
	CreateSourceInsight(ctx context.Context, sourceID string, transformationID string, modelID *string) (*models.SourceInsightResponse, error)
	SaveInsightAsNote(ctx context.Context, insightID string, notebookID string) (*models.Note, error)
	GetDefaultPrompt(ctx context.Context) (*models.DefaultPromptResponse, error)
	UpdateDefaultPrompt(ctx context.Context, instructions string) (*models.DefaultPromptResponse, error)
}