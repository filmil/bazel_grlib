package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	workspaceDir := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if workspaceDir != "" {
		os.Chdir(workspaceDir)
	}

	rootKconfig, err := os.Create("Kconfig")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Kconfig: %v\n", err)
		os.Exit(1)
	}
	defer rootKconfig.Close()

	fmt.Fprintln(rootKconfig, "# GRLIB Master Kconfig - Flat inclusion of all components")
	fmt.Fprintln(rootKconfig, "mainmenu \"GRLIB Master Configuration\"")
	fmt.Fprintln(rootKconfig, "")
    
    fmt.Fprintln(rootKconfig, "config AHBDW")
    fmt.Fprintln(rootKconfig, "    int \"AHB Data Width\"")
    fmt.Fprintln(rootKconfig, "    default 32")
    fmt.Fprintln(rootKconfig, "")
    fmt.Fprintln(rootKconfig, "config AHB_ACDM")
    fmt.Fprintln(rootKconfig, "    int \"Enable AMBA Compliant Data Muxing\"")
    fmt.Fprintln(rootKconfig, "    default 0")
    fmt.Fprintln(rootKconfig, "")
	fmt.Fprintln(rootKconfig, "config ACTIVE_DESIGN_PREFIX")
	fmt.Fprintln(rootKconfig, "    string \"Active Design Prefix (e.g. DESIGNS_LEON3_MINIMAL)\"")
	fmt.Fprintln(rootKconfig, "    default \"\"")
	fmt.Fprintln(rootKconfig, "")

	fmt.Fprintln(rootKconfig, "config NOELV_RV64")
	fmt.Fprintln(rootKconfig, "    bool \"Build NOEL-V as 64-bit\"")
	fmt.Fprintln(rootKconfig, "    default y")
	fmt.Fprintln(rootKconfig, "")

	err = filepath.Walk("third_party/grlib/kconfig", func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if !info.IsDir() && strings.HasSuffix(path, ".Kconfig") {
			fmt.Fprintf(rootKconfig, "rsource \"%s\"\n", path)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking kconfig dir: %v\n", err)
	}
    
    fmt.Println("Generated flat master Kconfig")
}
