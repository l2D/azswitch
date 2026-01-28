package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/l2D/azswitch/internal/azure"
)

func TestNewModel(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)

	if model.state != StateLoading {
		t.Errorf("expected state to be StateLoading, got %v", model.state)
	}

	if model.view != ViewSubscriptions {
		t.Errorf("expected view to be ViewSubscriptions, got %v", model.view)
	}
}

func TestModel_Update_DataLoaded(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)

	msg := dataLoadedMsg{
		account: &azure.Account{
			Name:              "Test Sub",
			TenantDisplayName: "Test Tenant",
			User:              azure.User{Name: "test@test.com"},
		},
		subscriptions: []azure.Subscription{
			{Name: "Sub 1", ID: "id-1", IsDefault: false},
			{Name: "Sub 2", ID: "id-2", IsDefault: true},
		},
		tenants: []azure.Tenant{
			{DisplayName: "Tenant 1", TenantID: "tid-1"},
		},
	}

	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	if m.state != StateReady {
		t.Errorf("expected state to be StateReady, got %v", m.state)
	}

	if m.account == nil {
		t.Fatal("expected account to be set")
	}

	if len(m.subscriptions) != 2 {
		t.Errorf("expected 2 subscriptions, got %d", len(m.subscriptions))
	}

	// Cursor should be on the default subscription (index 1)
	if m.cursor != 1 {
		t.Errorf("expected cursor to be 1 (default sub), got %d", m.cursor)
	}
}

func TestModel_Update_KeyNavigation(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)

	// Load data first
	model.state = StateReady
	model.subscriptions = []azure.Subscription{
		{Name: "Sub 1", ID: "id-1"},
		{Name: "Sub 2", ID: "id-2"},
		{Name: "Sub 3", ID: "id-3"},
	}
	model.cursor = 0

	// Test down navigation
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	if m.cursor != 1 {
		t.Errorf("expected cursor to be 1 after down, got %d", m.cursor)
	}

	// Test up navigation
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.Update(msg)
	m = newModel.(Model)

	if m.cursor != 0 {
		t.Errorf("expected cursor to be 0 after up, got %d", m.cursor)
	}
}

func TestModel_Update_TabSwitch(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)
	model.state = StateReady

	if model.view != ViewSubscriptions {
		t.Errorf("expected initial view to be ViewSubscriptions")
	}

	// Press tab
	msg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	if m.view != ViewTenants {
		t.Errorf("expected view to be ViewTenants after tab, got %v", m.view)
	}

	// Press tab again
	newModel, _ = m.Update(msg)
	m = newModel.(Model)

	if m.view != ViewSubscriptions {
		t.Errorf("expected view to be ViewSubscriptions after second tab, got %v", m.view)
	}
}

func TestModel_Update_ErrorMsg(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)

	msg := errMsg{err: azure.ErrNotLoggedIn}
	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	if m.state != StateError {
		t.Errorf("expected state to be StateError, got %v", m.state)
	}

	if m.err != azure.ErrNotLoggedIn {
		t.Errorf("expected error to be ErrNotLoggedIn")
	}
}

func TestModel_Update_WindowSize(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	if m.width != 100 {
		t.Errorf("expected width to be 100, got %d", m.width)
	}

	if m.height != 50 {
		t.Errorf("expected height to be 50, got %d", m.height)
	}
}

func TestModel_View_Loading(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)
	model.state = StateLoading

	view := model.View()

	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestModel_View_Ready(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)
	model.state = StateReady
	model.account = &azure.Account{
		Name:              "Test Sub",
		TenantDisplayName: "Test Tenant",
		User:              azure.User{Name: "test@test.com"},
	}
	model.subscriptions = []azure.Subscription{
		{Name: "Sub 1", ID: "id-1", IsDefault: true},
	}

	view := model.View()

	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestModel_View_Quitting(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)
	model.quitting = true

	view := model.View()

	if view != "" {
		t.Error("expected empty view when quitting")
	}
}

func TestModel_CursorBounds(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)
	model.state = StateReady
	model.subscriptions = []azure.Subscription{
		{Name: "Sub 1", ID: "id-1"},
	}
	model.cursor = 0

	// Try to go up when at top
	msg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	if m.cursor != 0 {
		t.Errorf("cursor should stay at 0, got %d", m.cursor)
	}

	// Try to go down when at bottom
	msg = tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = m.Update(msg)
	m = newModel.(Model)

	if m.cursor != 0 {
		t.Errorf("cursor should stay at 0 (only 1 item), got %d", m.cursor)
	}
}

func TestModel_HelpToggle(t *testing.T) {
	client := azure.NewMockClient()
	model := NewModel(client)
	model.state = StateReady

	if model.showHelp {
		t.Error("expected showHelp to be false initially")
	}

	// Press ?
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, _ := model.Update(msg)
	m := newModel.(Model)

	if !m.showHelp {
		t.Error("expected showHelp to be true after pressing ?")
	}

	// Press ? again
	newModel, _ = m.Update(msg)
	m = newModel.(Model)

	if m.showHelp {
		t.Error("expected showHelp to be false after pressing ? again")
	}
}
