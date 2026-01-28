package azure

import (
	"context"
	"testing"
)

func TestMockClient_GetCurrentAccount(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	account, err := client.GetCurrentAccount(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if account.Name != "Test Subscription" {
		t.Errorf("expected name 'Test Subscription', got '%s'", account.Name)
	}

	if account.User.Name != "test@example.com" {
		t.Errorf("expected user 'test@example.com', got '%s'", account.User.Name)
	}

	if client.Calls.GetCurrentAccount != 1 {
		t.Errorf("expected 1 call, got %d", client.Calls.GetCurrentAccount)
	}
}

func TestMockClient_ListSubscriptions(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	subs, err := client.ListSubscriptions(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(subs) != 2 {
		t.Errorf("expected 2 subscriptions, got %d", len(subs))
	}

	if subs[0].Name != "Test Subscription 1" {
		t.Errorf("expected 'Test Subscription 1', got '%s'", subs[0].Name)
	}

	if !subs[0].IsDefault {
		t.Error("expected first subscription to be default")
	}
}

func TestMockClient_ListTenants(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	tenants, err := client.ListTenants(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tenants) != 2 {
		t.Errorf("expected 2 tenants, got %d", len(tenants))
	}

	if tenants[0].DisplayName != "Test Tenant" {
		t.Errorf("expected 'Test Tenant', got '%s'", tenants[0].DisplayName)
	}
}

func TestMockClient_SetSubscription(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	err := client.SetSubscription(ctx, "test-subscription-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(client.Calls.SetSubscription) != 1 {
		t.Errorf("expected 1 call, got %d", len(client.Calls.SetSubscription))
	}

	if client.Calls.SetSubscription[0] != "test-subscription-id" {
		t.Errorf("expected 'test-subscription-id', got '%s'", client.Calls.SetSubscription[0])
	}
}

func TestMockClient_LoginToTenant(t *testing.T) {
	client := NewMockClient()
	ctx := context.Background()

	err := client.LoginToTenant(ctx, "test-tenant-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(client.Calls.LoginToTenant) != 1 {
		t.Errorf("expected 1 call, got %d", len(client.Calls.LoginToTenant))
	}

	if client.Calls.LoginToTenant[0] != "test-tenant-id" {
		t.Errorf("expected 'test-tenant-id', got '%s'", client.Calls.LoginToTenant[0])
	}
}

func TestSubscription_Methods(t *testing.T) {
	sub := Subscription{
		Name: "My Subscription",
		ID:   "00000000-0000-0000-0000-000000000001",
	}

	if sub.Title() != "My Subscription" {
		t.Errorf("expected Title() to return 'My Subscription', got '%s'", sub.Title())
	}

	if sub.Description() != "00000000-0000-0000-0000-000000000001" {
		t.Errorf("expected Description() to return subscription ID")
	}

	if sub.FilterValue() != "My Subscription" {
		t.Errorf("expected FilterValue() to return name")
	}
}

func TestTenant_Methods(t *testing.T) {
	t.Run("with display name", func(t *testing.T) {
		tenant := Tenant{
			DisplayName:   "My Tenant",
			DefaultDomain: "mytenant.onmicrosoft.com",
			TenantID:      "00000000-0000-0000-0000-000000000002",
		}

		if tenant.Title() != "My Tenant" {
			t.Errorf("expected Title() to return 'My Tenant', got '%s'", tenant.Title())
		}

		if tenant.FilterValue() != "My Tenant" {
			t.Errorf("expected FilterValue() to return display name")
		}
	})

	t.Run("without display name", func(t *testing.T) {
		tenant := Tenant{
			DefaultDomain: "mytenant.onmicrosoft.com",
			TenantID:      "00000000-0000-0000-0000-000000000002",
		}

		if tenant.Title() != "mytenant.onmicrosoft.com" {
			t.Errorf("expected Title() to return domain, got '%s'", tenant.Title())
		}

		if tenant.FilterValue() != "mytenant.onmicrosoft.com" {
			t.Errorf("expected FilterValue() to return domain")
		}
	})

	t.Run("description returns tenant ID", func(t *testing.T) {
		tenant := Tenant{
			TenantID: "00000000-0000-0000-0000-000000000002",
		}

		if tenant.Description() != "00000000-0000-0000-0000-000000000002" {
			t.Errorf("expected Description() to return tenant ID")
		}
	})
}

func TestMockClient_CustomBehavior(t *testing.T) {
	client := NewMockClient()

	// Override default behavior
	client.CheckCLIFunc = func(_ context.Context) error {
		return ErrAzureCLINotInstalled
	}

	err := client.CheckCLI(context.Background())
	if err != ErrAzureCLINotInstalled {
		t.Errorf("expected ErrAzureCLINotInstalled, got %v", err)
	}
}
