// LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type VerilatorJSON struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Modulesp []Module `json:"modulesp"`
	Miscsp   []Misc   `json:"miscsp"`
}

type Misc struct {
	Type   string  `json:"type"`
	Typesp []DType `json:"typesp"`
}

type DType struct {
	Type   string `json:"type"`
	Addr   string `json:"addr"`
	Name   string `json:"name"`
	Range  string `json:"range,omitempty"`
	Signed bool   `json:"signed,omitempty"`
}

type Module struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	OrigName string `json:"origName"`
	Stmtsp   []Stmt `json:"stmtsp"`
}

type Stmt struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	OrigName  string `json:"origName"`
	Direction string `json:"direction,omitempty"`
	DtypeName string `json:"dtypeName,omitempty"`
	Dtypep    string `json:"dtypep,omitempty"`
}

var boundsRegex = regexp.MustCompile(`\[(\d+):(\d+)\]`)

func main() {
	if len(os.Args) < 5 || len(os.Args) > 6 {
		log.Fatalf("Usage: %s <input.json> <out_proxy.vhdl> <out_bindings.hpp> <top_module> [out_arch.vhdl]", os.Args[0])
	}

	jsonPath := os.Args[1]
	vhdlPath := os.Args[2]
	hppPath := os.Args[3]
	topModuleName := os.Args[4]
	
	archPath := ""
	if len(os.Args) == 6 {
		archPath = os.Args[5]
	}

	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		log.Fatalf("Error opening JSON: %v", err)
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Error reading JSON: %v", err)
	}

	var vJSON VerilatorJSON
	if err := json.Unmarshal(jsonData, &vJSON); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	if len(vJSON.Modulesp) == 0 {
		log.Fatal("No modules found in JSON")
	}

	// Build type map
	typeMap := make(map[string]DType)
	for _, misc := range vJSON.Miscsp {
		if misc.Type == "TYPETABLE" {
			for _, t := range misc.Typesp {
				typeMap[t.Addr] = t
			}
		}
	}

	if archPath != "" {
		// Separate entity and architecture
		entityOut, err := os.Create(vhdlPath)
		if err != nil {
			log.Fatalf("Error creating VHDL entity: %v", err)
		}
		defer entityOut.Close()

		archOut, err := os.Create(archPath)
		if err != nil {
			log.Fatalf("Error creating VHDL architecture: %v", err)
		}
		defer archOut.Close()

		generateVHDL(vJSON.Modulesp, typeMap, entityOut, archOut, topModuleName)
	} else {
		// Combined entity and architecture
		vhdlOut, err := os.Create(vhdlPath)
		if err != nil {
			log.Fatalf("Error creating VHDL: %v", err)
		}
		defer vhdlOut.Close()
		generateVHDL(vJSON.Modulesp, typeMap, vhdlOut, vhdlOut, topModuleName)
	}

	hppOut, err := os.Create(hppPath)
	if err != nil {
		log.Fatalf("Error creating HPP: %v", err)
	}
	defer hppOut.Close()
	generateHPP(vJSON.Modulesp, hppOut, topModuleName)
}

func mapVerilogToVHDLType(direction, vType string, dt *DType) string {
	dir := "in"
	if direction == "OUTPUT" {
		dir = "out"
	} else if direction == "INOUT" {
		dir = "inout"
	}

	// Try DType from typetable first
	if dt != nil && dt.Range != "" {
		parts := strings.Split(dt.Range, ":")
		if len(parts) == 2 {
			return fmt.Sprintf("%s std_logic_vector(%s downto %s)", dir, parts[0], parts[1])
		}
	}

	// Parse vector bounds from the type string (fallback), e.g., "logic [7:0]"
	matches := boundsRegex.FindStringSubmatch(vType)
	if len(matches) == 3 {
		return fmt.Sprintf("%s std_logic_vector(%s downto %s)", dir, matches[1], matches[2])
	}

	return fmt.Sprintf("%s std_logic", dir)
}

func generateVHDL(modules []Module, typeMap map[string]DType, entityOut, archOut io.Writer, topModuleName string) {
	const vhdlEntityTemplate = `library ieee;
use ieee.std_logic_1164.all;
use ieee.numeric_std.all;

{{range .Modules}}
{{if eq .Name $.TargetModule}}
entity {{.Name}} is
  generic (
    INSTANCE_ID : integer
  );
  port (
{{- range $i, $port := .Ports}}
    {{$port.OrigName}} : {{$port.VhdlType}}{{if not $port.IsLast}};{{end}}
{{- end}}
  );
end entity;
{{end}}
{{end}}
`

	const vhdlArchTemplate = `
architecture proxy of {{.TargetModule}} is
  procedure step_verilator(id: integer; path: string);
  attribute foreign of step_verilator : procedure is "VHPIDIRECT verilator_step_call_{{.TargetModule}}";
  procedure step_verilator(id: integer; path: string) is
  begin
    report "VHPIDIRECT binding failed! C function not called." severity failure;
  end procedure;
begin
  process({{- range $i, $mod := .Modules}}{{if eq $mod.Name $.TargetModule}}{{- range $j, $port := $mod.InputPorts}}{{$port.OrigName}}{{if not $port.IsLastInput}}, {{end}}{{end}}{{end}}{{end}})
  begin
    step_verilator(INSTANCE_ID, {{.TargetModule}}'path_name);
  end process;
end architecture;
`

	type PortData struct {
		OrigName    string
		VhdlType    string
		Direction   string
		IsLast      bool
		IsLastInput bool
	}

	type ModuleData struct {
		Type       string
		Name       string
		Ports      []PortData
		InputPorts []PortData
		ClkPort    string
	}

	var data struct {
		TargetModule string
		Modules      []ModuleData
	}
	data.TargetModule = topModuleName

	for _, mod := range modules {
		if mod.Type != "MODULE" || mod.Name != topModuleName {
			continue
		}

		var mData ModuleData
		mData.Type = mod.Type
		mData.Name = mod.Name

		var ports []Stmt
		for _, stmt := range mod.Stmtsp {
			if stmt.Type == "VAR" && (stmt.Direction == "INPUT" || stmt.Direction == "OUTPUT" || stmt.Direction == "INOUT") {
				ports = append(ports, stmt)
			}
		}

		for i, port := range ports {
			var dt *DType
			if d, ok := typeMap[port.Dtypep]; ok {
				dt = &d
			}
			vhdlType := mapVerilogToVHDLType(port.Direction, port.DtypeName, dt)

			pData := PortData{
				OrigName:  port.OrigName,
				VhdlType:  vhdlType,
				Direction: port.Direction,
				IsLast:    i == len(ports)-1,
			}
			mData.Ports = append(mData.Ports, pData)

			if port.Direction == "INPUT" || port.Direction == "INOUT" {
				mData.InputPorts = append(mData.InputPorts, pData)
			}
		}

		for i := range mData.InputPorts {
			mData.InputPorts[i].IsLastInput = (i == len(mData.InputPorts)-1)
		}

		data.Modules = append(data.Modules, mData)
	}

	entityTmpl := template.Must(template.New("vhdlEntity").Parse(vhdlEntityTemplate))
	if err := entityTmpl.Execute(entityOut, data); err != nil {
		log.Fatalf("Error executing VHDL entity template: %v", err)
	}

	archTmpl := template.Must(template.New("vhdlArch").Parse(vhdlArchTemplate))
	if err := archTmpl.Execute(archOut, data); err != nil {
		log.Fatalf("Error executing VHDL architecture template: %v", err)
	}
}

func generateHPP(modules []Module, out io.Writer, topModuleName string) {
	const hppTemplate = `#pragma once

#include "vhpi_user.h"
#include <string>
#include <unordered_map>
#include <memory>
#include <iostream>
#include "V{{.TopModule.Name}}.h"

struct InstanceState {
    std::unique_ptr<VerilatedContext> context;
    std::unique_ptr<V{{.TopModule.Name}}> dut;
    std::string path_prefix;
    std::unordered_map<std::string, vhpiHandleT> handles;
};

inline vhpiHandleT get_vhpi_handle(InstanceState& state, const std::string& name) {
    if (state.handles.find(name) != state.handles.end()) return state.handles[name];
    std::string full_name = state.path_prefix + ":" + name;
    vhpiHandleT h = vhpi_handle_by_name(full_name.c_str(), nullptr);
    if (h) {
        state.handles[name] = h;
        return h;
    }
    fprintf(stderr, "[VHPI] vhpi_handle_by_name failed for %s\n", full_name.c_str());
    return nullptr;
}

inline void init_bindings(int id, InstanceState& state, const std::string& path_prefix) {
    state.context = std::make_unique<VerilatedContext>();
    state.dut = std::make_unique<V{{.TopModule.Name}}>(state.context.get());
    state.path_prefix = path_prefix;
}

inline void sync_inputs(InstanceState& state) {
{{- range .TopModule.Stmtsp}}
{{- if eq .Type "VAR"}}
{{- if eq .Direction "INPUT"}}
    {
        vhpiHandleT net_handle = get_vhpi_handle(state, "{{.UpperName}}");
        if (net_handle) {
            int size = vhpi_get(vhpiSizeP, net_handle);
            vhpiValueT val;
            if (size == 1) {
                val.format = vhpiLogicVal;
                vhpi_get_value(net_handle, &val);
                state.dut->{{.OrigName}} = (val.value.enumv == 3) ? 1 : 0;
            } else {
                vhpiEnumT vec[64]; // Support up to 64 bits
                val.format = vhpiLogicVecVal;
                val.bufSize = size * sizeof(vhpiEnumT);
                val.value.enumvs = vec;
                vhpi_get_value(net_handle, &val);
                uint64_t int_val = 0;
                for (int i = 0; i < size; i++) {
                	if (vec[i] == 3) { // 3 is vhpi1
                		int_val |= (1ULL << (size - 1 - i));
                	}
                }
                state.dut->{{.OrigName}} = int_val;
            }
        }
    }
{{- end}}
{{- end}}
{{- end}}
}

inline void sync_outputs(InstanceState& state) {
{{- range .TopModule.Stmtsp}}
{{- if eq .Type "VAR"}}
{{- if eq .Direction "OUTPUT"}}
    {
        vhpiHandleT net_handle = get_vhpi_handle(state, "{{.UpperName}}");
        if (net_handle) {
            int size = vhpi_get(vhpiSizeP, net_handle);
            vhpiValueT val;
            uint64_t int_val = state.dut->{{.OrigName}};
            if (size == 1) {
                val.format = vhpiLogicVal;
                val.value.enumv = int_val ? 3 : 2;
                vhpi_put_value(net_handle, &val, vhpiForcePropagate);
            } else {
                vhpiEnumT vec[64]; // Support up to 64 bits
                for (int i = 0; i < size; i++) {
                    if (int_val & (1ULL << (size - 1 - i))) {
                        vec[i] = 3; // vhpi1
                    } else {
                        vec[i] = 2; // vhpi0
                    }
                }
                val.format = vhpiLogicVecVal;
                val.bufSize = size * sizeof(vhpiEnumT);
                val.value.enumvs = vec;
                vhpi_put_value(net_handle, &val, vhpiForcePropagate);
            }
        } else {
            fprintf(stderr, "[VHPI] Output %s handle NOT FOUND\n", "{{.UpperName}}");
        }
    }
{{- end}}
{{- end}}
{{- end}}
}

inline void eval_model(InstanceState& state) {
    state.dut->eval();
}
`

	// We only expect one top module for now.
	var topModule *Module
	for _, mod := range modules {
		if mod.Type == "MODULE" && mod.Name == topModuleName {
			topModule = &mod
			break
		}
	}

	if topModule == nil {
		return
	}

	// Create a wrapper struct for the template to include augmented data if needed
	type TemplateStmt struct {
		Type      string
		OrigName  string
		UpperName string
		Direction string
	}

	type TemplateModule struct {
		Name   string
		Stmtsp []TemplateStmt
	}

	var tMod TemplateModule
	tMod.Name = topModule.Name
	for _, s := range topModule.Stmtsp {
		tMod.Stmtsp = append(tMod.Stmtsp, TemplateStmt{
			Type:      s.Type,
			OrigName:  s.OrigName,
			UpperName: strings.ToUpper(s.OrigName),
			Direction: s.Direction,
		})
	}

	data := struct {
		TopModule TemplateModule
	}{
		TopModule: tMod,
	}

	tmpl := template.Must(template.New("hpp").Parse(hppTemplate))
	if err := tmpl.Execute(out, data); err != nil {
		log.Fatalf("Error executing HPP template: %v", err)
	}
}
