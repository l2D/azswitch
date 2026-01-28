package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/l2D/azswitch/internal/azure"
)

// ViewType represents the current view.
type ViewType int

const (
	ViewSubscriptions ViewType = iota
	ViewDirectories
)

// State represents the application state.
type State int

const (
	StateLoading State = iota
	StateReady
	StateError
	StateSwitching
	StateSuccess
)

// Model represents the TUI model.
type Model struct {
	// Azure client
	client azure.Client

	// Current state
	state State

	// Current view
	view ViewType

	// Data
	account       *azure.Account
	subscriptions []azure.Subscription
	tenants       []azure.Tenant

	// UI state
	cursor       int
	tenantCursor int
	err          error
	message      string

	// Components
	spinner spinner.Model
	help    help.Model
	keys    KeyMap

	// Window size
	width  int
	height int

	// Show help
	showHelp bool

	// Quit flag
	quitting bool
}

// Messages
type (
	// errMsg is sent when an error occurs.
	errMsg struct{ err error }

	// dataLoadedMsg is sent when data is loaded.
	dataLoadedMsg struct {
		account       *azure.Account
		subscriptions []azure.Subscription
		tenants       []azure.Tenant
	}

	// switchedMsg is sent when a switch operation completes.
	switchedMsg struct {
		message string
	}
)

// NewModel creates a new TUI model.
func NewModel(client azure.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	h := help.New()
	h.ShowAll = false

	return Model{
		client:  client,
		state:   StateLoading,
		view:    ViewSubscriptions,
		spinner: s,
		help:    h,
		keys:    DefaultKeyMap(),
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadData(),
	)
}

// loadData loads the Azure data.
func (m Model) loadData() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		account, err := m.client.GetCurrentAccount(ctx)
		if err != nil {
			return errMsg{err}
		}

		subs, err := m.client.ListSubscriptions(ctx)
		if err != nil {
			return errMsg{err}
		}

		tenants, err := m.client.ListTenants(ctx)
		if err != nil {
			return errMsg{err}
		}

		return dataLoadedMsg{
			account:       account,
			subscriptions: subs,
			tenants:       tenants,
		}
	}
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case errMsg:
		m.state = StateError
		m.err = msg.err
		return m, nil

	case dataLoadedMsg:
		m.state = StateReady
		m.account = msg.account
		m.subscriptions = msg.subscriptions
		m.tenants = msg.tenants
		// Set cursor to current subscription
		for i := range m.subscriptions {
			if m.subscriptions[i].IsDefault {
				m.cursor = i
				break
			}
		}
		return m, nil

	case switchedMsg:
		m.state = StateSuccess
		m.message = msg.message
		return m, m.loadData()
	}

	return m, nil
}

// handleKeyMsg handles keyboard input.
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Always allow quit
	if msg.String() == "ctrl+c" {
		m.quitting = true
		return m, tea.Quit
	}

	// Don't handle keys while loading or switching
	if m.state == StateLoading || m.state == StateSwitching {
		return m, nil
	}

	switch {
	case key.Matches(msg, m.keys.Quit):
		m.quitting = true
		return m, tea.Quit

	case key.Matches(msg, m.keys.Help):
		m.showHelp = !m.showHelp
		m.help.ShowAll = m.showHelp
		return m, nil

	case key.Matches(msg, m.keys.Tab):
		if m.view == ViewSubscriptions {
			m.view = ViewDirectories
		} else {
			m.view = ViewSubscriptions
		}
		return m, nil

	case key.Matches(msg, m.keys.Up):
		if m.view == ViewSubscriptions {
			if m.cursor > 0 {
				m.cursor--
			}
		} else {
			if m.tenantCursor > 0 {
				m.tenantCursor--
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.Down):
		if m.view == ViewSubscriptions {
			if m.cursor < len(m.subscriptions)-1 {
				m.cursor++
			}
		} else {
			if m.tenantCursor < len(m.tenants)-1 {
				m.tenantCursor++
			}
		}
		return m, nil

	case key.Matches(msg, m.keys.Select):
		return m.handleSelect()

	case key.Matches(msg, m.keys.Refresh):
		m.state = StateLoading
		return m, tea.Batch(m.spinner.Tick, m.loadData())
	}

	return m, nil
}

// handleSelect handles the selection.
func (m Model) handleSelect() (tea.Model, tea.Cmd) {
	if m.view == ViewSubscriptions && len(m.subscriptions) > 0 {
		sub := m.subscriptions[m.cursor]
		if sub.IsDefault {
			// Already selected
			return m, nil
		}
		m.state = StateSwitching
		return m, tea.Batch(
			m.spinner.Tick,
			m.switchSubscription(sub.ID),
		)
	} else if m.view == ViewDirectories && len(m.tenants) > 0 {
		tenant := m.tenants[m.tenantCursor]
		m.state = StateSwitching
		return m, tea.Batch(
			m.spinner.Tick,
			m.switchTenant(tenant.TenantID),
		)
	}
	return m, nil
}

// switchSubscription switches to the specified subscription.
func (m Model) switchSubscription(id string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		if err := m.client.SetSubscription(ctx, id); err != nil {
			return errMsg{err}
		}
		return switchedMsg{message: "Subscription switched successfully"}
	}
}

// switchTenant switches to the specified tenant using interactive login.
func (m Model) switchTenant(id string) tea.Cmd {
	cmd := exec.Command("az", "login", "--tenant", id)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errMsg{err}
		}
		return switchedMsg{message: "Directory switched successfully"}
	})
}

// View renders the UI.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Header with current account
	s.WriteString(m.renderHeader())
	s.WriteString("\n")

	// Main content
	switch m.state {
	case StateLoading:
		s.WriteString(m.renderLoading())
	case StateError:
		s.WriteString(m.renderError())
	case StateSwitching:
		s.WriteString(m.renderSwitching())
	default:
		s.WriteString(m.renderTabs())
		s.WriteString("\n")
		if m.view == ViewSubscriptions {
			s.WriteString(m.renderSubscriptions())
		} else {
			s.WriteString(m.renderDirectories())
		}
	}

	// Help
	s.WriteString("\n")
	s.WriteString(HelpStyle.Render(m.help.View(m.keys)))

	return s.String()
}

// renderHeader renders the header section.
func (m Model) renderHeader() string {
	if m.account == nil {
		return TitleStyle.Render("Azure Account Switcher")
	}

	var content strings.Builder
	content.WriteString(TitleStyle.Render("Azure Account Switcher"))
	content.WriteString("\n")
	content.WriteString(fmt.Sprintf("  %s %s\n", MutedStyle.Render("User:"), m.account.User.Name))
	content.WriteString(fmt.Sprintf("  %s %s\n", MutedStyle.Render("Tenant:"), m.account.TenantDisplayName))
	content.WriteString(fmt.Sprintf("  %s %s", MutedStyle.Render("Subscription:"), CurrentStyle.Render(m.account.Name)))

	return HeaderBoxStyle.Render(content.String())
}

// renderTabs renders the tab bar.
func (m Model) renderTabs() string {
	subsTab := "Subscriptions"
	dirsTab := "Directories"

	if m.view == ViewSubscriptions {
		subsTab = ActiveTabStyle.Render(subsTab)
		dirsTab = InactiveTabStyle.Render(dirsTab)
	} else {
		subsTab = InactiveTabStyle.Render(subsTab)
		dirsTab = ActiveTabStyle.Render(dirsTab)
	}

	return fmt.Sprintf("  %s  |  %s", subsTab, dirsTab)
}

// renderLoading renders the loading state.
func (m Model) renderLoading() string {
	return fmt.Sprintf("\n  %s Loading...", m.spinner.View())
}

// renderSwitching renders the switching state.
func (m Model) renderSwitching() string {
	return fmt.Sprintf("\n  %s Switching...", m.spinner.View())
}

// renderError renders the error state.
func (m Model) renderError() string {
	return fmt.Sprintf("\n  %s %s", ErrorStyle.Render("Error:"), m.err.Error())
}

// renderSubscriptions renders the subscriptions list.
func (m Model) renderSubscriptions() string {
	if len(m.subscriptions) == 0 {
		return MutedStyle.Render("\n  No subscriptions found")
	}

	var s strings.Builder
	s.WriteString("\n")

	for i := range m.subscriptions {
		sub := &m.subscriptions[i]
		cursor := "  "
		if i == m.cursor {
			cursor = CursorStyle.Render("> ")
		}

		name := sub.Name
		switch {
		case sub.IsDefault:
			name = CurrentStyle.Render(name + " ✓")
		case i == m.cursor:
			name = SelectedStyle.Render(name)
		default:
			name = NormalStyle.Render(name)
		}

		s.WriteString(fmt.Sprintf("%s%s\n", cursor, name))
		s.WriteString(fmt.Sprintf("    %s\n", MutedStyle.Render(sub.ID)))
	}

	return s.String()
}

// renderDirectories renders the directories (tenants) list with their subscriptions.
func (m Model) renderDirectories() string {
	if len(m.tenants) == 0 {
		return MutedStyle.Render("\n  No directories found")
	}

	// Group subscriptions by tenant ID
	subsByTenant := make(map[string][]azure.Subscription)
	for i := range m.subscriptions {
		sub := &m.subscriptions[i]
		subsByTenant[sub.TenantID] = append(subsByTenant[sub.TenantID], *sub)
	}

	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(WarningStyle.Render("  ⚠ Switching directories will open browser for re-authentication"))
	s.WriteString("\n\n")

	for i := range m.tenants {
		tenant := &m.tenants[i]
		cursor := "  "
		if i == m.tenantCursor {
			cursor = CursorStyle.Render("> ")
		}

		name := tenant.Title()
		isCurrent := m.account != nil && tenant.TenantID == m.account.TenantID

		switch {
		case isCurrent:
			name = CurrentStyle.Render(name + " ✓")
		case i == m.tenantCursor:
			name = SelectedStyle.Render(name)
		default:
			name = NormalStyle.Render(name)
		}

		s.WriteString(fmt.Sprintf("%s%s\n", cursor, name))

		// Show subscriptions for this directory
		if subs, ok := subsByTenant[tenant.TenantID]; ok && len(subs) > 0 {
			for j := range subs {
				subName := subs[j].Name
				if subs[j].IsDefault {
					subName = CurrentStyle.Render("• " + subName)
				} else {
					subName = MutedStyle.Render("• " + subName)
				}
				s.WriteString(fmt.Sprintf("    %s\n", subName))
			}
		} else {
			s.WriteString(fmt.Sprintf("    %s\n", MutedStyle.Render("(no subscriptions)")))
		}
	}

	return s.String()
}
