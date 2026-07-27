package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// isPortFree checks if a port is free by temporarily binding to it
func isPortFree(port uint) bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	l.Close()
	return true
}

// FindAvailablePorts finds the first contiguous block of `count` free ports
// within [start, end] (inclusive) and returns the start and end ports.
func FindAvailablePorts(start, end, count uint) (uint, uint, error) {

	if start > end {
		err := fmt.Errorf("start port must not exceed end port")
		return 0, 0, err
	}
	if start == 0 {
		err := fmt.Errorf("start port must be greater than zero")
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

	var runStart uint
	var runLength uint = 0

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

	return 0, 0, fmt.Errorf("no contiguous range of %d ports found", count)
}

// Color and column style variables
var (
	rangePortStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	returnPortStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	helpTitleStyle  = lipgloss.NewStyle().Foreground(charmtone.Charple).Bold(true)
)

// Flag variables
type config struct {
	ports struct {
		start uint
		end   uint
		count uint
	}
}

var cfg config

func main() {

	// Create new flag sets
	portsFS := &pflag.FlagSet{}

	// CLI Create a new cobra Command
	cobraCmd := &cobra.Command{
		Use: "findcports",
		Long: fmt.Sprintf(`%s

Searches within a specified port range and returns the
start and end ports of the first contiguous block of
available ports with the requested size.`,
			helpTitleStyle.Render("DESCRIPTION")),
		RunE: runApp,
	}

	// =============================================
	// CLI Flags
	// =============================================

	// ports
	portsFS.UintVarP(&cfg.ports.start, "start", "s", 20000,
		"First port in the search range",
	)

	portsFS.UintVarP(&cfg.ports.end, "end", "e", 20999,
		"Last port in the search range",
	)

	portsFS.UintVarP(&cfg.ports.count, "count", "c", 10,
		"Size of the contiguous port block to find",
	)

	cobraCmd.Flags().AddFlagSet(portsFS)

	// =============================================
	// Execute app
	// =============================================

	// Set new flag display settings
	cobraCmd.Flags().SetInterspersed(false)
	cobraCmd.Flags().SortFlags = false
	cobraCmd.Flags().PrintDefaults()

	// Remove completion and help commands (--help is kept)
	cobraCmd.CompletionOptions.DisableDefaultCmd = true
	cobraCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// CLI Version
	cobraCmd.Version = version.Print("findcports")

	// Execute cobra command
	if err := fang.Execute(context.Background(), cobraCmd); err != nil {
		os.Exit(1)
	}
}

func runApp(cobraCmd *cobra.Command, args []string) error {

	startPort, endPort, err := FindAvailablePorts(
		cfg.ports.start,
		cfg.ports.end,
		cfg.ports.count)
	if err != nil {
		return err
	}

	if startPort == endPort {
		fmt.Printf(
			"Found the first available port between %s and %s:\n%s\n",
			rangePortStyle.Render(strconv.FormatUint(uint64(cfg.ports.start), 10)),
			rangePortStyle.Render(strconv.FormatUint(uint64(cfg.ports.end), 10)),
			returnPortStyle.Render(strconv.FormatUint(uint64(startPort), 10)),
		)
	} else {
		fmt.Printf(
			"Found the first available continuous ports between %s and %s:\n%s-%s (%d ports)\n",
			rangePortStyle.Render(strconv.FormatUint(uint64(cfg.ports.start), 10)),
			rangePortStyle.Render(strconv.FormatUint(uint64(cfg.ports.end), 10)),
			returnPortStyle.Render(strconv.FormatUint(uint64(startPort), 10)),
			returnPortStyle.Render(strconv.FormatUint(uint64(endPort), 10)),
			cfg.ports.count,
		)
	}
	return nil
}
