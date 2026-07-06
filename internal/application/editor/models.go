package editor

import (
	"context"
	"time"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	rt "github.com/Ajayvtl/devserver/internal/runtime"
)

// Document represents an open file in the editor.
type Document struct {
	URI      string `json:"uri"`
	Language string `json:"language"`
	IsDirty  bool   `json:"isDirty"`
}

// Position represents a line and character position.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range represents a text selection or range.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// LayoutConfig represents editor split states.
type LayoutConfig struct {
	Orientation string  `json:"orientation"`
	Sizes       []int   `json:"sizes"`
	Groups      []Group `json:"groups"`
}

// Group represents a split pane in the editor.
type Group struct {
	ActiveTab string   `json:"activeTab"`
	Tabs      []string `json:"tabs"` // URIs
}

// TerminalInfo represents a terminal instance in the editor.
type TerminalInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // e.g., bash, zsh, powershell
}

// Thread represents an AI conversation thread context.
type Thread struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// WorkspaceSession defines the persistent desktop-like state of an editor session.
type WorkspaceSession struct {
	WorkspaceID     common.WorkspaceID `json:"workspaceId"`
	OpenTabs        []Document         `json:"openTabs"`
	ActiveTab       string             `json:"activeTab"`
	PinnedTabs      []string           `json:"pinnedTabs"`
	CursorPosition  Position           `json:"cursorPosition"`
	Selections      []Range            `json:"selections"`
	SplitLayout     LayoutConfig       `json:"splitLayout"`
	ExplorerWidth   int                `json:"explorerWidth"`
	PanelSizes      map[string]int     `json:"panelSizes"`
	ExpandedFolders []string           `json:"expandedFolders"`
	RecentFiles     []string           `json:"recentFiles"`
	SearchHistory   []string           `json:"searchHistory"`
	CommandHistory  []string           `json:"commandHistory"`
	Terminals       []TerminalInfo     `json:"terminals"`
	AIHistory       []Thread           `json:"aiHistory"`
	EditorSettings  map[string]any     `json:"editorSettings"`
}

// SessionManager persists and synchronizes WorkspaceSession.
type SessionManager interface {
	GetSession(ctx context.Context, workspaceID common.WorkspaceID) (*WorkspaceSession, error)
	SaveSession(ctx context.Context, session *WorkspaceSession) error
}

// ProcessManager manages the lifecycle of editor instances (e.g., OpenVSCode Server).
type ProcessManager interface {
	Start(ctx context.Context, workspaceID common.WorkspaceID) error
	Stop(ctx context.Context, workspaceID common.WorkspaceID) error
	Status(ctx context.Context, workspaceID common.WorkspaceID) rt.Status
}

// ExtensionManager maps VS Code extensions to DevServer Providers.
type ExtensionManager interface {
	Install(ctx context.Context, extensionID string) error
	Remove(ctx context.Context, extensionID string) error
	List(ctx context.Context) ([]string, error)
}

// ProxyManager handles HTTP and WebSocket routing to the editor instance.
type ProxyManager interface {
	RegisterRoute(workspaceID common.WorkspaceID, targetURL string) error
	RemoveRoute(workspaceID common.WorkspaceID) error
}

// EventBridge translates editor IPC events onto the DevServer Event Bus.
type EventBridge interface {
	// Start begins listening to the editor IPC channel.
	Start(ctx context.Context, workspaceID common.WorkspaceID) error
	Stop(ctx context.Context, workspaceID common.WorkspaceID) error
	Send(ctx context.Context, workspaceID common.WorkspaceID, eventType string, payload any) error
	Receive(ctx context.Context, workspaceID common.WorkspaceID, payload []byte) error
}
