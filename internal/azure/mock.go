package azure

import (
	"context"
)

// MockClient is a mock implementation of the Azure Client interface for testing.
type MockClient struct {
	// CheckCLIFunc is called when CheckCLI is invoked.
	CheckCLIFunc func(ctx context.Context) error

	// CheckLoginFunc is called when CheckLogin is invoked.
	CheckLoginFunc func(ctx context.Context) error

	// GetCurrentAccountFunc is called when GetCurrentAccount is invoked.
	GetCurrentAccountFunc func(ctx context.Context) (*Account, error)

	// ListSubscriptionsFunc is called when ListSubscriptions is invoked.
	ListSubscriptionsFunc func(ctx context.Context) ([]Subscription, error)

	// ListTenantsFunc is called when ListTenants is invoked.
	ListTenantsFunc func(ctx context.Context) ([]Tenant, error)

	// SetSubscriptionFunc is called when SetSubscription is invoked.
	SetSubscriptionFunc func(ctx context.Context, subscriptionIDOrName string) error

	// LoginToTenantFunc is called when LoginToTenant is invoked.
	LoginToTenantFunc func(ctx context.Context, tenantID string) error

	// Calls tracks function call history.
	Calls struct {
		CheckCLI          int
		CheckLogin        int
		GetCurrentAccount int
		ListSubscriptions int
		ListTenants       int
		SetSubscription   []string
		LoginToTenant     []string
	}
}

// NewMockClient creates a new mock client with default implementations.
func NewMockClient() *MockClient {
	return &MockClient{
		CheckCLIFunc: func(_ context.Context) error {
			return nil
		},
		CheckLoginFunc: func(_ context.Context) error {
			return nil
		},
		GetCurrentAccountFunc: func(_ context.Context) (*Account, error) {
			return &Account{
				Name:              "Test Subscription",
				ID:                "00000000-0000-0000-0000-000000000001",
				TenantID:          "00000000-0000-0000-0000-000000000002",
				TenantDisplayName: "Test Tenant",
				IsDefault:         true,
				State:             "Enabled",
				User: User{
					Name: "test@example.com",
					Type: "user",
				},
			}, nil
		},
		ListSubscriptionsFunc: func(_ context.Context) ([]Subscription, error) {
			return []Subscription{
				{
					Name:              "Test Subscription 1",
					ID:                "00000000-0000-0000-0000-000000000001",
					TenantID:          "00000000-0000-0000-0000-000000000002",
					TenantDisplayName: "Test Tenant",
					IsDefault:         true,
					State:             "Enabled",
				},
				{
					Name:              "Test Subscription 2",
					ID:                "00000000-0000-0000-0000-000000000003",
					TenantID:          "00000000-0000-0000-0000-000000000002",
					TenantDisplayName: "Test Tenant",
					IsDefault:         false,
					State:             "Enabled",
				},
			}, nil
		},
		ListTenantsFunc: func(_ context.Context) ([]Tenant, error) {
			return []Tenant{
				{
					DisplayName:   "Test Tenant",
					TenantID:      "00000000-0000-0000-0000-000000000002",
					DefaultDomain: "test.onmicrosoft.com",
				},
				{
					DisplayName:   "Another Tenant",
					TenantID:      "00000000-0000-0000-0000-000000000004",
					DefaultDomain: "another.onmicrosoft.com",
				},
			}, nil
		},
		SetSubscriptionFunc: func(_ context.Context, _ string) error {
			return nil
		},
		LoginToTenantFunc: func(_ context.Context, _ string) error {
			return nil
		},
	}
}

// CheckCLI implements Client.
func (m *MockClient) CheckCLI(ctx context.Context) error {
	m.Calls.CheckCLI++
	return m.CheckCLIFunc(ctx)
}

// CheckLogin implements Client.
func (m *MockClient) CheckLogin(ctx context.Context) error {
	m.Calls.CheckLogin++
	return m.CheckLoginFunc(ctx)
}

// GetCurrentAccount implements Client.
func (m *MockClient) GetCurrentAccount(ctx context.Context) (*Account, error) {
	m.Calls.GetCurrentAccount++
	return m.GetCurrentAccountFunc(ctx)
}

// ListSubscriptions implements Client.
func (m *MockClient) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	m.Calls.ListSubscriptions++
	return m.ListSubscriptionsFunc(ctx)
}

// ListTenants implements Client.
func (m *MockClient) ListTenants(ctx context.Context) ([]Tenant, error) {
	m.Calls.ListTenants++
	return m.ListTenantsFunc(ctx)
}

// SetSubscription implements Client.
func (m *MockClient) SetSubscription(ctx context.Context, subscriptionIDOrName string) error {
	m.Calls.SetSubscription = append(m.Calls.SetSubscription, subscriptionIDOrName)
	return m.SetSubscriptionFunc(ctx, subscriptionIDOrName)
}

// LoginToTenant implements Client.
func (m *MockClient) LoginToTenant(ctx context.Context, tenantID string) error {
	m.Calls.LoginToTenant = append(m.Calls.LoginToTenant, tenantID)
	return m.LoginToTenantFunc(ctx, tenantID)
}

// Ensure MockClient implements Client.
var _ Client = (*MockClient)(nil)
