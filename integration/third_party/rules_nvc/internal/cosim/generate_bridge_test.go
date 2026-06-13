package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateVHDL(t *testing.T) {
	tests := []struct {
		name          string
		modules       []Module
		expectedParts []string
	}{
		{
			name: "Single module with clk, d, q",
			modules: []Module{
				{
					Type: "MODULE",
					Name: "dut",
					Stmtsp: []Stmt{
						{Type: "VAR", OrigName: "clk", Direction: "INPUT", DtypeName: "logic"},
						{Type: "VAR", OrigName: "d", Direction: "INPUT", DtypeName: "logic"},
						{Type: "VAR", OrigName: "q", Direction: "OUTPUT", DtypeName: "logic"},
					},
				},
			},
			expectedParts: []string{
				"entity dut is",
				"clk : in std_logic;",
				"d : in std_logic;",
				"q : out std_logic",
				"architecture proxy of dut is",
				"process(clk, d)",
				"step_verilator(INSTANCE_ID, dut'path_name);",
			},
		},
		{
			name: "Adder without clock using dynamic vector bounds",
			modules: []Module{
				{
					Type: "MODULE",
					Name: "adder",
					Stmtsp: []Stmt{
						{Type: "VAR", OrigName: "a", Direction: "INPUT", DtypeName: "logic [7:0]"},
						{Type: "VAR", OrigName: "b", Direction: "INPUT", DtypeName: "logic [7:0]"},
						{Type: "VAR", OrigName: "sum", Direction: "OUTPUT", DtypeName: "logic [8:0]"},
					},
				},
			},
			expectedParts: []string{
				"entity adder is",
				"a : in std_logic_vector(7 downto 0);",
				"b : in std_logic_vector(7 downto 0);",
				"sum : out std_logic_vector(8 downto 0)",
				"process(a, b)",
				"step_verilator(INSTANCE_ID, adder'path_name);",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			generateVHDL(tt.modules, nil, &buf, &buf, tt.modules[0].Name)
			output := buf.String()

			for _, part := range tt.expectedParts {
				if !strings.Contains(output, part) {
					t.Errorf("Expected output to contain %q, but it didn't.\nFull Output:\n%s", part, output)
				}
			}
		})
	}
}

func TestGenerateHPP(t *testing.T) {
	modules := []Module{
		{
			Type: "MODULE",
			Name: "dut",
		},
	}

	var buf bytes.Buffer
	generateHPP(modules, &buf, "dut")
	output := buf.String()

	expectedParts := []string{
		"#pragma once",
		"#include \"vhpi_user.h\"",
		"struct InstanceState {",
		"std::unique_ptr<Vdut> dut;",
		"inline void init_bindings(int id, InstanceState& state, const std::string& path_prefix) {",
		"inline void sync_inputs(InstanceState& state) {",
		"inline void sync_outputs(InstanceState& state) {",
		"inline void eval_model(InstanceState& state) {",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("Expected HPP to contain %q", part)
		}
	}
}
