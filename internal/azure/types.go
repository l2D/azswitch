// Package azure provides Azure CLI wrapper functionality.
package azure

// Account represents the current Azure account information.
type Account struct {
	EnvironmentName   string `json:"environmentName"`
	HomeTenantID      string `json:"homeTenantId"`
	ID                string `json:"id"`
	IsDefault         bool   `json:"isDefault"`
	ManagedByTenants  []any  `json:"managedByTenants"`
	Name              string `json:"name"`
	State             string `json:"state"`
	TenantDisplayName string `json:"tenantDisplayName"`
	TenantID          string `json:"tenantId"`
	User              User   `json:"user"`
}

// User represents the user information within an account.
type User struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Subscription represents an Azure subscription.
type Subscription struct {
	CloudName         string `json:"cloudName"`
	HomeTenantID      string `json:"homeTenantId"`
	ID                string `json:"id"`
	IsDefault         bool   `json:"isDefault"`
	ManagedByTenants  []any  `json:"managedByTenants"`
	Name              string `json:"name"`
	State             string `json:"state"`
	TenantDisplayName string `json:"tenantDisplayName"`
	TenantID          string `json:"tenantId"`
	User              User   `json:"user"`
}

// Tenant represents an Azure AD tenant/directory.
type Tenant struct {
	DefaultDomain   string   `json:"defaultDomain"`
	DisplayName     string   `json:"displayName"`
	ID              string   `json:"id"`
	TenantID        string   `json:"tenantId"`
	TenantCategory  string   `json:"tenantCategory"`
	TenantType      string   `json:"tenantType"`
	Domains         []string `json:"domains,omitempty"`
	CountryCode     string   `json:"countryCode,omitempty"`
	TenantBrandName string   `json:"tenantBrandingLogoUrl,omitempty"`
}

// Title returns a display title for the subscription.
func (s Subscription) Title() string {
	return s.Name
}

// Description returns a description for the subscription.
func (s Subscription) Description() string {
	return s.ID
}

// FilterValue returns the value used for filtering.
func (s Subscription) FilterValue() string {
	return s.Name
}

// Title returns a display title for the tenant.
func (t Tenant) Title() string {
	if t.DisplayName != "" {
		return t.DisplayName
	}
	return t.DefaultDomain
}

// Description returns a description for the tenant.
func (t Tenant) Description() string {
	return t.TenantID
}

// FilterValue returns the value used for filtering.
func (t Tenant) FilterValue() string {
	if t.DisplayName != "" {
		return t.DisplayName
	}
	return t.DefaultDomain
}
