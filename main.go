package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// isPortFree checks if a port is free by temporarily binding to it
func isPortFree(port uint16) bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Printf("failed to close Listener: %v", err)
		}
	}()

	return true
}

// FindAvailablePorts finds the first contiguous block of `count` free ports
// within [start, end] (inclusive) and returns the start and end ports.
func FindAvailablePorts(start, end, count uint16) (uint16, uint16, error) {

	if start > end {
		err := fmt.Errorf("start port must not exceed end port")
		return 0, 0, err
	}
	if end == 0 {
		err := fmt.Errorf("end port must be greater than zero")
		return 0, 0, err
	}
	if count == 0 {
		err := fmt.Errorf("count must be greater than zero")
		return 0, 0, err
	}

	var runStart uint16
	var runLength uint16 = 0

	for p := start; p <= end; p++ {
		if isPortFree(p) {
			if runLength == 0 {
				runStart = p
			}
			runLength++

			if runLength == count {
				return runStart, runStart + runLength - 1, nil
			}
		} else {
			runLength = 0
		}
	}

	return 0, 0, fmt.Errorf(
		"no contiguous range of %d ports found between %d and %d",
		count, start, end)
}

// Set version
var version = "dev"

// Color style variables
var (
	helpTitleStyle = lipgloss.NewStyle().Foreground(charmtone.Charple).
		Bold(true)
)

// Flag variables
type config struct {
	ports struct {
		start uint16
		end   uint16
		count uint16
	}
}

// Setting the port limit
var portLimit uint16 = math.MaxUint16

var cfg config

func main() {

	// Create new flag sets
	portsFS := &pflag.FlagSet{}

	// CLI Create a new cobra Command
	cobraCmd := &cobra.Command{
		Use: "findbaseport",
		Long: fmt.Sprintf(`%s

Searches within a specified port range (inclusive) and returns the
start port of the first contiguous block of available
ports with the requested size.`,
			helpTitleStyle.Render("DESCRIPTION")),
		RunE: runApp,
	}

	// =============================================
	// CLI Flags
	// =============================================

	// ports
	portsFS.Uint16VarP(&cfg.ports.start, "start", "s", 0,
		"First port in the search range",
	)

	portsFS.Uint16VarP(&cfg.ports.end, "end", "e", portLimit,
		"Last port in the search range",
	)

	portsFS.Uint16VarP(&cfg.ports.count, "count", "c", 0,
		"Size of the contiguous port block to find",
	)

	cobraCmd.Flags().AddFlagSet(portsFS)

	// =============================================
	// Execute app
	// =============================================

	// Set new flag display settings
	cobraCmd.Flags().SetInterspersed(false)
	cobraCmd.Flags().SortFlags = false

	// Remove completion and help commands (--help is kept)
	cobraCmd.CompletionOptions.DisableDefaultCmd = true
	cobraCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Execute cobra command
	if err := fang.Execute(
		context.Background(),
		cobraCmd,
		fang.WithVersion(version)); err != nil {
		os.Exit(1)
	}
}

func runApp(cmd *cobra.Command, _ []string) error {

	// Check if --count was set with either --start or --end
	if !cmd.Flags().Changed("count") ||
		(!cmd.Flags().Changed("start") && !cmd.Flags().Changed("end")) {
		return fmt.Errorf(
			"flag --count with either --start or --end is required")
	}

	startPort, _, err := FindAvailablePorts(
		cfg.ports.start,
		cfg.ports.end,
		cfg.ports.count)
	if err != nil {
		return err
	} else {
		fmt.Printf("%d\n", startPort)
	}

	return nil
}
