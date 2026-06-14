package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func main() {
	workspaceDir := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if workspaceDir != "" {
		if err := os.Chdir(workspaceDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error changing to workspace directory: %v\n", err)
			os.Exit(1)
		}
	}

	grlibPath := "."
	if len(os.Args) > 1 {
		grlibPath = os.Args[1]
	}

	outPath := "third_party/grlib/grlib.BUILD"
	if len(os.Args) > 2 {
		outPath = os.Args[2]
	}
	
	nvcPath := "tool.nvc/BUILD.bazel"

	libs, _ := readLines(filepath.Join(grlibPath, "lib/libs.txt"))
	var activeLibs []string
	libStds := make(map[string]string)
	for _, l := range libs {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		parts := strings.Fields(l)
		libName := parts[0]
		activeLibs = append(activeLibs, libName)
		for _, p := range parts[1:] {
			if strings.HasPrefix(p, "vhdlstd=") {
				libStds[libName] = strings.TrimPrefix(p, "vhdlstd=")
			}
		}
	}

	var discoveredLibs []string
	entries, _ := os.ReadDir(filepath.Join(grlibPath, "lib"))
	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "tech" {
				techEntries, _ := os.ReadDir(filepath.Join(grlibPath, "lib/tech"))
				for _, te := range techEntries {
					if te.IsDir() {
						if _, err := os.Stat(filepath.Join(grlibPath, "lib/tech", te.Name(), "dirs.txt")); err == nil {
							discoveredLibs = append(discoveredLibs, "tech/"+te.Name())
						}
					}
				}
			} else {
				if _, err := os.Stat(filepath.Join(grlibPath, "lib", entry.Name(), "dirs.txt")); err == nil {
					discoveredLibs = append(discoveredLibs, entry.Name())
				}
			}
		}
	}

	var allLibs []string
	added := make(map[string]bool)
	addToAll := func(l string) {
		if !added[l] {
			allLibs = append(allLibs, l)
			added[l] = true
		}
	}
	addToAll("grlib")
	addToAll("techmap")
	for _, l := range activeLibs { addToAll(l) }
	for _, l := range discoveredLibs { addToAll(l) }

	libPathMap := make(map[string]string)
	for _, l := range allLibs {
		libPathMap[filepath.Base(l)] = l
	}

	libDeps := make(map[string][]string)
	libFiles := make(map[string][]string)
	libNoelvCfg := make(map[string]bool)

	for _, lib := range allLibs {
		libSourcePath := filepath.Join(grlibPath, "lib", lib)
		dirs, err := readLines(filepath.Join(libSourcePath, "dirs.txt"))
		if err != nil { continue }

		if lib == "gaisler" {
			dirs = append(dirs, "l5nv/shared", "noelv/pkg", "noelv/core", "noelv/subsys", "noelv")
		}

		var files []string
		depsSet := make(map[string]bool)
		for _, dir := range dirs {
			dir = strings.TrimSpace(dir)
			if dir == "" || strings.HasPrefix(dir, "#") { continue }

			for _, dp := range strings.Fields(dir) {
				if strings.HasPrefix(dp, "#") { break }
				subdirPath := filepath.Join(libSourcePath, dp)
				linesSyn, _ := readLines(filepath.Join(subdirPath, "vhdlsyn.txt"))
				linesSim, _ := readLines(filepath.Join(subdirPath, "vhdlsim.txt"))

				if len(linesSyn) == 0 && len(linesSim) == 0 {
					entries, _ := os.ReadDir(subdirPath)
					var pkgFiles []string
					var otherFiles []string
					for _, entry := range entries {
						if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".vhd") {
							name := entry.Name()
							if strings.Contains(name, "_cfg") || strings.Contains(name, "_types") || strings.Contains(name, "_pkg") || name == "noelv.vhd" {
								pkgFiles = append(pkgFiles, name)
							} else {
								otherFiles = append(otherFiles, name)
							}
						}
					}
					sort.Strings(pkgFiles)
					sort.Strings(otherFiles)
					linesSyn = append(pkgFiles, otherFiles...)
				}

				for _, fLine := range append(linesSyn, linesSim...) {
					fLine = strings.TrimSpace(fLine)
					if fLine == "" || strings.HasPrefix(fLine, "#") { continue }
					fParts := strings.Fields(fLine)
					fileName := fParts[0]
					
					relVhdPath := filepath.Join("lib", lib, dp, fileName)
					fullVhdPath := filepath.Join(grlibPath, relVhdPath)

					if _, err := os.Stat(fullVhdPath); err == nil {
						if fileName == "grtachom.vhd" { continue }
						if relVhdPath == "lib/grlib/sparc/cpu_disas.vhd" { continue }
						
						if fileName == "noelv_cfg_32.vhd" || fileName == "noelv_cfg_64.vhd" {
							libNoelvCfg[lib] = true
							continue
						}

						finalRef := relVhdPath
						if relVhdPath == "lib/grlib/stdlib/config.vhd" {
							finalRef = "@grlib//third_party/grlib:config.vhd"
						}

						found := false
						for _, existing := range files {
							if existing == finalRef { found = true; break }
						}
						if !found {
							files = append(files, finalRef)
							fileDeps := parseLibDeps(fullVhdPath)
							for _, d := range fileDeps {
								if d != filepath.Base(lib) { depsSet[d] = true }
							}
						}
					}
				}
			}
		}

		if len(files) > 0 || libNoelvCfg[lib] {
			libFiles[lib] = files
			var deps []string
			for d := range depsSet {
				if p, ok := libPathMap[d]; ok { deps = append(deps, p) }
			}
			libDeps[lib] = deps
		}
	}

	// 1. Generate grlib.BUILD (Tool agnostic filegroups)
	os.MkdirAll(filepath.Dir(outPath), 0755)
	gb, _ := os.Create(outPath)
	fmt.Fprintln(gb, "exports_files(glob([\"**/*.vhd\"]))")
	fmt.Fprintln(gb, "")
	fmt.Fprintln(gb, "filegroup(")
	fmt.Fprintln(gb, "    name = \"grlib_srcs_all\",")
	fmt.Fprintln(gb, "    srcs = glob([\"**\"]),")
	fmt.Fprintln(gb, "    visibility = [\"//visibility:public\"],")
	fmt.Fprintln(gb, ")")
	fmt.Fprintln(gb, "")
	
	for _, lib := range allLibs {
		files := libFiles[lib]
		if len(files) == 0 && !libNoelvCfg[lib] { continue }
		libBase := filepath.Base(lib)
		fmt.Fprintln(gb, "# do not sort")
		fmt.Fprintln(gb, "filegroup(")
		fmt.Fprintf(gb, "    name = \"%s_files\",\n", libBase)
		fmt.Fprintln(gb, "    # do not sort")
		if libBase == "grlib" {
			fmt.Fprint(gb, "    srcs = []")
			for _, f := range libFiles[lib] {
				base := filepath.Base(f)
				if base == "stdio.vhd" {
					fmt.Fprintln(gb, " + select({")
					fmt.Fprintln(gb, "        \"@grlib//:std_2008\": [\"@grlib//third_party/grlib:lib/grlib/stdlib/stdio_2008.vhd\"],")
					fmt.Fprintln(gb, "        \"@grlib//:std_2019\": [\"@grlib//third_party/grlib:lib/grlib/stdlib/stdio_2008.vhd\"],")
					fmt.Fprintln(gb, "        \"//conditions:default\": [\"lib/grlib/stdlib/stdio.vhd\"],")
					fmt.Fprint(gb, "    })")
				} else if base == "testlib.vhd" {
					fmt.Fprintln(gb, " + select({")
					fmt.Fprintln(gb, "        \"@grlib//:std_2008\": [\"@grlib//third_party/grlib:lib/grlib/stdlib/testlib_2008.vhd\"],")
					fmt.Fprintln(gb, "        \"@grlib//:std_2019\": [\"@grlib//third_party/grlib:lib/grlib/stdlib/testlib_2008.vhd\"],")
					fmt.Fprintln(gb, "        \"//conditions:default\": [\"lib/grlib/stdlib/testlib.vhd\"],")
					fmt.Fprint(gb, "    })")
				} else {
					fmt.Fprintf(gb, " + [\"%s\"]", f)
				}
			}
			fmt.Fprintln(gb, ",")
		} else if libNoelvCfg[lib] {
			fmt.Fprintln(gb, "    srcs = select({")
			fmt.Fprintln(gb, "        \"@grlib//:NOELV_RV64\": [\"lib/gaisler/noelv/pkg/noelv_cfg_64.vhd\"],")
			fmt.Fprintln(gb, "        \"//conditions:default\": [\"lib/gaisler/noelv/pkg/noelv_cfg_32.vhd\"],")
			fmt.Fprintln(gb, "    }) + [")
			for _, f := range files { fmt.Fprintf(gb, "        \"%s\",\n", f) }
			fmt.Fprintln(gb, "    ],")
		} else {
			fmt.Fprintln(gb, "    srcs = [")
			for _, f := range files { fmt.Fprintf(gb, "        \"%s\",\n", f) }
			fmt.Fprintln(gb, "    ],")
		}
		fmt.Fprintln(gb, "    visibility = [\"//visibility:public\"],")
		fmt.Fprintln(gb, ")")
		fmt.Fprintln(gb, "")
	}
	gb.Close()

	// 2. Generate tool.nvc/BUILD.bazel (NVC-specific rules)
	os.MkdirAll(filepath.Dir(nvcPath), 0755)
	nv, _ := os.Create(nvcPath)
	fmt.Fprintln(nv, "load(\"@rules_nvc//nvc:rules.bzl\", \"vhdl_elaborate\", \"vhdl_library\", \"vhdl_test\")")
	fmt.Fprintln(nv, "")
	
	// Helper to get transitive deps
	getTransitiveDeps := func(lib string) []string {
		res := make(map[string]bool)
		var visit func(string)
		visit = func(curr string) {
			for _, d := range libDeps[curr] {
				depBase := filepath.Base(d)
				if (len(libFiles[d]) > 0 || libNoelvCfg[d]) && !res[depBase] {
					res[depBase] = true
					visit(d)
				}
			}
		}
		visit(lib)
		var sorted []string
		for d := range res { sorted = append(sorted, d) }
		sort.Strings(sorted)
		return sorted
	}

	for _, lib := range allLibs {
		files := libFiles[lib]
		if len(files) == 0 && !libNoelvCfg[lib] { continue }
		libBase := filepath.Base(lib)
		
		transDeps := getTransitiveDeps(lib)

		fmt.Fprintln(nv, "vhdl_library(")
		fmt.Fprintf(nv, "    name = \"%s\",\n", libBase)
		fmt.Fprintf(nv, "    srcs = [\"@grlib_srcs//:%s_files\"],\n", libBase)
		fmt.Fprintln(nv, "    standard = select({")
		fmt.Fprintln(nv, "        \"@grlib//:std_1987\": \"1987\",")
		fmt.Fprintln(nv, "        \"@grlib//:std_1993\": \"1993\",")
		fmt.Fprintln(nv, "        \"@grlib//:std_2002\": \"2002\",")
		fmt.Fprintln(nv, "        \"@grlib//:std_2008\": \"2008\",")
		fmt.Fprintln(nv, "        \"@grlib//:std_2019\": \"2019\",")
		std := libStds[lib]
		if std == "" || std == "93" { std = "1993" } else if std == "08" { std = "2008" }
		fmt.Fprintf(nv, "        \"//conditions:default\": \"%s\",\n", std)
		fmt.Fprintln(nv, "    }),")
		fmt.Fprintln(nv, "    deps = [")
		for _, d := range transDeps {
			if d != libBase {
				fmt.Fprintf(nv, "        \":%s\",\n", d)
			}
		}
		fmt.Fprintln(nv, "    ],")
		
		// Also mark manual if any direct dep was missing (handled by transDeps being filtered)
		hasMissingDirectDep := false
		for _, d := range libDeps[lib] {
			if len(libFiles[d]) == 0 && !libNoelvCfg[d] {
				hasMissingDirectDep = true; break
			}
		}
		if hasMissingDirectDep || libBase == "gsi" {
			fmt.Fprintln(nv, "    tags = [\"manual\"],")
		}

		fmt.Fprintln(nv, "    visibility = [\"//visibility:public\"],")
		fmt.Fprintln(nv, ")")
		fmt.Fprintln(nv, "")
	}

	fmt.Fprintln(nv, "# Integration tests moved to tool.nvc")
	fmt.Fprintln(nv, "vhdl_test(")
	fmt.Fprintln(nv, "    name = \"noelv_simulation\",")
	fmt.Fprintln(nv, "    srcs = [\"//tool.nvc/integration:noelv_tb.vhd\"],")
	fmt.Fprintln(nv, "    deps = [")
	fmt.Fprintln(nv, "        \":grlib\",")
	fmt.Fprintln(nv, "        \":gaisler\",")
	fmt.Fprintln(nv, "        \":techmap\",")
	fmt.Fprintln(nv, "    ],")
	fmt.Fprintln(nv, "    entity = \"noelv_tb\",")
	fmt.Fprintln(nv, "    standard = select({")
	fmt.Fprintln(nv, "        \"@grlib//:std_1987\": \"1987\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_1993\": \"1993\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_2002\": \"2002\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_2008\": \"2008\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_2019\": \"2019\",")
	fmt.Fprintln(nv, "        \"//conditions:default\": \"1993\",")
	fmt.Fprintln(nv, "    }),")
	fmt.Fprintln(nv, ")")
	fmt.Fprintln(nv, "")

	fmt.Fprintln(nv, "vhdl_library(")
	fmt.Fprintln(nv, "    name = \"noelvsys_repro_lib\",")
	fmt.Fprintln(nv, "    srcs = [\"//tool.nvc/integration:noelvsys_repro_tb.vhd\"],")
	fmt.Fprintln(nv, "    deps = [")
	fmt.Fprintln(nv, "        \":grlib\",")
	fmt.Fprintln(nv, "        \":gaisler\",")
	fmt.Fprintln(nv, "        \":techmap\",")
	fmt.Fprintln(nv, "    ],")
	fmt.Fprintln(nv, "    standard = select({")
	fmt.Fprintln(nv, "        \"@grlib//:std_1987\": \"1987\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_1993\": \"1993\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_2002\": \"2002\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_2008\": \"2008\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_2019\": \"2019\",")
	fmt.Fprintln(nv, "        \"//conditions:default\": \"1993\",")
	fmt.Fprintln(nv, "    }),")
	fmt.Fprintln(nv, "    tags = [\"manual\"],")
	fmt.Fprintln(nv, ")")
	fmt.Fprintln(nv, "")

	fmt.Fprintln(nv, "vhdl_elaborate(")
	fmt.Fprintln(nv, "    name = \"noelvsys_repro_tb\",")
	fmt.Fprintln(nv, "    library = \":noelvsys_repro_lib\",")
	fmt.Fprintln(nv, "    standard = select({")
	fmt.Fprintln(nv, "        \"@grlib//:std_1987\": \"1987\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_1993\": \"1993\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_2002\": \"2002\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_2008\": \"2008\",")
	fmt.Fprintln(nv, "        \"@grlib//:std_2019\": \"2019\",")
	fmt.Fprintln(nv, "        \"//conditions:default\": \"1993\",")
	fmt.Fprintln(nv, "    }),")
	fmt.Fprintln(nv, "    tags = [\"manual\"],")
	fmt.Fprintln(nv, ")")

	nv.Close()
	fmt.Println("Generated grlib.BUILD and tool.nvc/BUILD.bazel")
}

var libRegex = regexp.MustCompile(`(?i)^\s*library\s+([^;]+);`)

func parseLibDeps(path string) []string {
	f, err := os.Open(path)
	if err != nil { return nil }
	defer f.Close()
	var deps []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		match := libRegex.FindStringSubmatch(scanner.Text())
		if len(match) > 1 {
			libs := strings.Split(match[1], ",")
			for _, l := range libs {
				dep := strings.ToLower(strings.TrimSpace(l))
				if dep != "ieee" && dep != "std" && dep != "work" && dep != "" {
					found := false
					for _, existing := range deps {
						if existing == dep { found = true; break }
					}
					if !found {
						deps = append(deps, dep)
					}
				}
			}
		}
	}
	return deps
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil { return nil, err }
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { lines = append(lines, scanner.Text()) }
	return lines, scanner.Err()
}
