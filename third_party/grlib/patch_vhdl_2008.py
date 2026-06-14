import sys
import re

def comment_ranges(lines, ranges):
    ranges.sort(reverse=True)
    for start, end in ranges:
        # index is line_number - 1
        for i in range(start-1, end):
            lines[i] = "-- " + lines[i]

def patch_vhd(input_path, output_path, mode):
    with open(input_path, 'r', encoding='latin-1') as f:
        lines = f.read().splitlines()
    
    if mode == "stdio":
        # Comment out conflicting ones
        ranges = [(39, 42), (44, 46), (48, 50), (52, 55), (57, 59), (61, 65), (67, 71), (73, 77)]
        ranges += [(141, 191), (193, 201), (203, 213), (215, 223), (225, 235), (237, 260), (262, 269), (271, 281)]
        comment_ranges(lines, ranges)
        
        # Add aliases in the spec (around line 80)
        alias_idx = 79 # after the commented out spec
        lines.insert(alias_idx, "   alias HRead is IEEE.Std_Logic_1164.HRead [Line, Std_ULogic_Vector, Boolean];")
        lines.insert(alias_idx + 1, "   alias HRead is IEEE.Std_Logic_1164.HRead [Line, Std_ULogic_Vector];")
        lines.insert(alias_idx + 2, "   alias HRead is IEEE.Std_Logic_1164.HRead [Line, Std_Logic_Vector, Boolean];")
        lines.insert(alias_idx + 3, "   alias HRead is IEEE.Std_Logic_1164.HRead [Line, Std_Logic_Vector];")
        lines.insert(alias_idx + 4, "   alias HWrite is IEEE.Std_Logic_1164.HWrite [Line, Std_ULogic_Vector, SIDE, WIDTH];")
        lines.insert(alias_idx + 5, "   alias HWrite is IEEE.Std_Logic_1164.HWrite [Line, Std_Logic_Vector, SIDE, WIDTH];")
        lines.insert(alias_idx + 6, "   alias Write is IEEE.Std_Logic_1164.Write [Line, Std_ULogic, SIDE, WIDTH];")

    elif mode == "testlib":
        ranges = [(58, 58), (110, 114), (236, 247), (391, 420)]
        comment_ranges(lines, ranges)
        
    content = "\n".join(lines)
    
    # Fix Text type and Read/Write calls
    content = re.sub(r'\bfile\s+(readfile|ReadFile)\s*:\s*text\b', r'file f_\1 : Std.TextIO.Text', content, flags=re.I)
    content = re.sub(r'\b(readline|endfile|read|file_open|file_close|hread|hwrite)\s*\(\s*(readfile|ReadFile)\b', r'\1(f_\2', content, flags=re.I)

    with open(output_path, 'w', encoding='latin-1') as f:
        f.write(content)

if __name__ == "__main__":
    mode = sys.argv[1]
    input_vhd = sys.argv[2]
    output_vhd = sys.argv[3]
    patch_vhd(input_vhd, output_vhd, mode)
