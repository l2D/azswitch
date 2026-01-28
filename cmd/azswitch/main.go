// Package main is the entry point for azswitch.
package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/l2D/azswitch/internal/azure"
	"github.com/l2D/azswitch/internal/tui"
	"github.com/l2D/azswitch/internal/version"
)

var (
	// Flags
	flagList         bool
	flagCurrent      bool
	flagSubscription string
	flagTenant       string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "azswitch",
	Short: "Switch Azure tenants, directories, and subscriptions",
	Long: `azswitch is a TUI application for switching Azure tenants, 
directories, and subscriptions.

Run without flags to enter interactive mode.`,
	Version: version.Short(),
	RunE:    run,
}

func init() {
	rootCmd.Flags().BoolVarP(&flagList, "list", "l", false, "List all subscriptions")
	rootCmd.Flags().BoolVarP(&flagCurrent, "current", "c", false, "Show current account")
	rootCmd.Flags().StringVarP(&flagSubscription, "subscription", "s", "", "Switch to subscription by ID or name")
	rootCmd.Flags().StringVarP(&flagTenant, "tenant", "t", "", "Switch to tenant by ID")

	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func run(_ *cobra.Command, _ []string) error {
	client := azure.NewCLIClient()
	ctx := context.Background()

	// Check if Azure CLI is installed
	if err := client.CheckCLI(ctx); err != nil {
		//nolint:staticcheck // ST1005: Azure CLI is a proper noun
		return fmt.Errorf("Azure CLI is not installed. Install from: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli")
	}

	// Check if logged in
	if err := client.CheckLogin(ctx); err != nil {
		return fmt.Errorf("not logged in to Azure CLI. Run: az login")
	}

	// Handle non-interactive flags
	if flagCurrent {
		return showCurrent(ctx, client)
	}

	if flagList {
		return listSubscriptions(ctx, client)
	}

	if flagSubscription != "" {
		return switchSubscription(ctx, client, flagSubscription)
	}

	if flagTenant != "" {
		return switchTenant(ctx, client, flagTenant)
	}

	// Interactive mode
	return runInteractive(client)
}

func showCurrent(ctx context.Context, client azure.Client) error {
	account, err := client.GetCurrentAccount(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Current Azure Account:")
	fmt.Printf("  User:         %s\n", account.User.Name)
	fmt.Printf("  Tenant:       %s (%s)\n", account.TenantDisplayName, account.TenantID)
	fmt.Printf("  Subscription: %s\n", account.Name)
	fmt.Printf("  ID:           %s\n", account.ID)
	fmt.Printf("  State:        %s\n", account.State)

	return nil
}

func listSubscriptions(ctx context.Context, client azure.Client) error {
	subs, err := client.ListSubscriptions(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Available Subscriptions:")
	for i := range subs {
		sub := &subs[i]
		indicator := "  "
		if sub.IsDefault {
			indicator = "* "
		}
		fmt.Printf("%s%s\n", indicator, sub.Name)
		fmt.Printf("    ID:    %s\n", sub.ID)
		fmt.Printf("    State: %s\n", sub.State)
	}

	return nil
}

func switchSubscription(ctx context.Context, client azure.Client, subscription string) error {
	fmt.Printf("Switching to subscription: %s\n", subscription)

	if err := client.SetSubscription(ctx, subscription); err != nil {
		return fmt.Errorf("failed to switch subscription: %w", err)
	}

	fmt.Println("Successfully switched subscription")
	return showCurrent(ctx, client)
}

func switchTenant(ctx context.Context, client azure.Client, tenant string) error {
	fmt.Printf("Switching to tenant: %s\n", tenant)
	fmt.Println("This will open a browser for authentication...")

	if err := client.LoginToTenant(ctx, tenant); err != nil {
		return fmt.Errorf("failed to switch tenant: %w", err)
	}

	fmt.Println("Successfully switched tenant")
	return showCurrent(ctx, client)
}

func runInteractive(client azure.Client) error {
	model := tui.NewModel(client)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}
