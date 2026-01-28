package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Common errors.
var (
	ErrAzureCLINotInstalled = errors.New("azure CLI is not installed")
	ErrNotLoggedIn          = errors.New("not logged in to Azure CLI")
	ErrCommandFailed        = errors.New("azure CLI command failed")
)

// Client defines the interface for Azure CLI operations.
type Client interface {
	// CheckCLI verifies that Azure CLI is installed.
	CheckCLI(ctx context.Context) error

	// CheckLogin verifies that the user is logged in.
	CheckLogin(ctx context.Context) error

	// GetCurrentAccount returns the current Azure account information.
	GetCurrentAccount(ctx context.Context) (*Account, error)

	// ListSubscriptions returns all available subscriptions.
	ListSubscriptions(ctx context.Context) ([]Subscription, error)

	// ListTenants returns all available tenants.
	ListTenants(ctx context.Context) ([]Tenant, error)

	// SetSubscription switches to the specified subscription.
	SetSubscription(ctx context.Context, subscriptionIDOrName string) error

	// LoginToTenant logs in to a specific tenant.
	LoginToTenant(ctx context.Context, tenantID string) error
}

// CLIClient implements Client using the Azure CLI.
type CLIClient struct {
	// azPath is the path to the az CLI binary.
	azPath string
}

// NewCLIClient creates a new Azure CLI client.
func NewCLIClient() *CLIClient {
	return &CLIClient{
		azPath: "az",
	}
}

// CheckCLI verifies that Azure CLI is installed.
func (c *CLIClient) CheckCLI(ctx context.Context) error {
	_, err := exec.LookPath(c.azPath)
	if err != nil {
		return ErrAzureCLINotInstalled
	}
	return nil
}

// CheckLogin verifies that the user is logged in.
func (c *CLIClient) CheckLogin(ctx context.Context) error {
	_, err := c.GetCurrentAccount(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Please run 'az login'") ||
			strings.Contains(err.Error(), "not logged in") {
			return ErrNotLoggedIn
		}
		return err
	}
	return nil
}

// GetCurrentAccount returns the current Azure account information.
func (c *CLIClient) GetCurrentAccount(ctx context.Context) (*Account, error) {
	output, err := c.runCommand(ctx, "account", "show", "--output", "json")
	if err != nil {
		return nil, err
	}

	var account Account
	if err := json.Unmarshal(output, &account); err != nil {
		return nil, fmt.Errorf("failed to parse account: %w", err)
	}

	return &account, nil
}

// ListSubscriptions returns all available subscriptions.
func (c *CLIClient) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	output, err := c.runCommand(ctx, "account", "list", "--output", "json")
	if err != nil {
		return nil, err
	}

	var subscriptions []Subscription
	if err := json.Unmarshal(output, &subscriptions); err != nil {
		return nil, fmt.Errorf("failed to parse subscriptions: %w", err)
	}

	return subscriptions, nil
}

// ListTenants returns all available tenants.
func (c *CLIClient) ListTenants(ctx context.Context) ([]Tenant, error) {
	output, err := c.runCommand(ctx, "account", "tenant", "list", "--output", "json")
	if err != nil {
		return nil, err
	}

	var tenants []Tenant
	if err := json.Unmarshal(output, &tenants); err != nil {
		return nil, fmt.Errorf("failed to parse tenants: %w", err)
	}

	return tenants, nil
}

// SetSubscription switches to the specified subscription.
func (c *CLIClient) SetSubscription(ctx context.Context, subscriptionIDOrName string) error {
	_, err := c.runCommand(ctx, "account", "set", "--subscription", subscriptionIDOrName)
	return err
}

// LoginToTenant logs in to a specific tenant.
func (c *CLIClient) LoginToTenant(ctx context.Context, tenantID string) error {
	_, err := c.runCommand(ctx, "login", "--tenant", tenantID, "--output", "none")
	return err
}

// runCommand executes an Azure CLI command and returns the output.
func (c *CLIClient) runCommand(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, c.azPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		return nil, fmt.Errorf("%w: %s", ErrCommandFailed, strings.TrimSpace(errMsg))
	}

	return stdout.Bytes(), nil
}
