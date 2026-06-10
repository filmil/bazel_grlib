package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	checkMode := false
	workspaceDir := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	
	grlibSrcs := os.Getenv("GRLIB_SRCS")
	if grlibSrcs == "" {
		fmt.Fprintf(os.Stderr, "GRLIB_SRCS environment variable must be set\n")
		os.Exit(1)
	}

	genBuildFiles := ""
	in2kconfig := ""
	genMasterKconfig := ""
	vhdPreprocess := ""
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--gen_build_files" && i+1 < len(os.Args) {
			genBuildFiles = os.Args[i+1]
			if !filepath.IsAbs(genBuildFiles) {
				if abs, err := filepath.Abs(genBuildFiles); err == nil {
					genBuildFiles = abs
				}
			}
			i++
		} else if os.Args[i] == "--in2kconfig" && i+1 < len(os.Args) {
			in2kconfig = os.Args[i+1]
			if !filepath.IsAbs(in2kconfig) {
				if abs, err := filepath.Abs(in2kconfig); err == nil {
					in2kconfig = abs
				}
			}
			i++
		} else if os.Args[i] == "--gen_master_kconfig" && i+1 < len(os.Args) {
			genMasterKconfig = os.Args[i+1]
			if !filepath.IsAbs(genMasterKconfig) {
				if abs, err := filepath.Abs(genMasterKconfig); err == nil {
					genMasterKconfig = abs
				}
			}
			i++
		} else if os.Args[i] == "--vhd_preprocess" && i+1 < len(os.Args) {
			vhdPreprocess = os.Args[i+1]
			if !filepath.IsAbs(vhdPreprocess) {
				if abs, err := filepath.Abs(vhdPreprocess); err == nil {
					vhdPreprocess = abs
				}
			}
			i++
		} else if os.Args[i] == "--check" {
			checkMode = true
		}
	}

	if genBuildFiles == "" || in2kconfig == "" || genMasterKconfig == "" || vhdPreprocess == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s --gen_build_files <path> --in2kconfig <path> --gen_master_kconfig <path> --vhd_preprocess <path> [--check]\n", os.Args[0])
		os.Exit(1)
	}

	if !checkMode && workspaceDir != "" {
		if err := os.Chdir(workspaceDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error changing to workspace directory: %v\n", err)
			os.Exit(1)
		}
	}

	if checkMode {
		tmpDir, err := os.MkdirTemp("", "maintain-check")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating temp dir: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir)

		fmt.Println("Check mode: generating files to temporary directory...")
		
		// 1. Run gen_build_files
		tmpBuildFile := filepath.Join(tmpDir, "grlib.BUILD")
		runScript(genBuildFiles, grlibSrcs, tmpBuildFile, vhdPreprocess)

		// 2. Run in2kconfig
		tmpKconfigDir := filepath.Join(tmpDir, "kconfig")
		runScript(in2kconfig, grlibSrcs, tmpKconfigDir)

		// Compare
		outOfSync := false
		
		if !filesEqual(tmpBuildFile, "third_party/grlib/grlib.BUILD") {
			fmt.Fprintf(os.Stderr, "third_party/grlib/grlib.BUILD is out of sync!\n")
			printDiff(tmpBuildFile, "third_party/grlib/grlib.BUILD")
			outOfSync = true
		}

		err = filepath.Walk(tmpKconfigDir, func(path string, info os.FileInfo, err error) error {
			if err != nil { return err }
			if info.IsDir() { return nil }
			rel, _ := filepath.Rel(tmpKconfigDir, path)
			origPath := filepath.Join("third_party/grlib/kconfig", rel)
			if !filesEqual(path, origPath) {
				fmt.Fprintf(os.Stderr, "%s is out of sync!\n", origPath)
				printDiff(path, origPath)
				outOfSync = true
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking temp kconfig dir: %v\n", err)
			outOfSync = true
		}

		if outOfSync {
			fmt.Fprintf(os.Stderr, "\nFiles are out of sync! Run 'bazel run //:maintain' to update.\n")
			os.Exit(1)
		}
		fmt.Println("All files are up to date.")
	} else {
		fmt.Println("Synchronizing files...")
		runScript(genBuildFiles, grlibSrcs, "third_party/grlib/grlib.BUILD", vhdPreprocess)
		
		kconfigDir := "third_party/grlib/kconfig"
		os.RemoveAll(kconfigDir)
		os.MkdirAll(kconfigDir, 0755)
		runScript(in2kconfig, grlibSrcs, kconfigDir)
		
		runScript(genMasterKconfig)
		fmt.Println("Files synchronized.")
	}
}

func runScript(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s: %v\n", name, err)
		os.Exit(1)
	}
}

func printDiff(path1, path2 string) {
	cmd := exec.Command("diff", "-u", path1, path2)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func filesEqual(path1, path2 string) bool {
	f1, err := os.Open(path1)
	if err != nil { return false }
	defer f1.Close()

	f2, err := os.Open(path2)
	if err != nil { return false }
	defer f2.Close()

	b1 := make([]byte, 4096)
	b2 := make([]byte, 4096)

	for {
		n1, err1 := f1.Read(b1)
		n2, err2 := f2.Read(b2)

		if n1 != n2 || !bytes.Equal(b1[:n1], b2[:n2]) {
			return false
		}

		if err1 == io.EOF && err2 == io.EOF {
			return true
		}

		if err1 != nil || err2 != nil {
			return false
		}
	}
}
