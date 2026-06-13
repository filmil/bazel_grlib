entity hello_world is
end entity hello_world;

library STD;      
library IEEE;    
use IEEE.std_logic_1164.all; 
use STD.textio.all;         
use IEEE.std_logic_textio.all;

-- This is the new package we're using.  The symbol name is `value`.
library pkg; use pkg.package_name.value;

architecture test of hello_world is
  subtype word_32 is std_logic_vector(31 downto 0);
  signal four_32 : word_32 := x"00000004";
  signal counter : integer := 1;    
begin  -- where parallel code is placed
  my_print : process is              
               variable my_line : line;
             begin
               write(my_line, value);
               writeline(output, my_line);               -- write to "output"
               wait;
             end process my_print;
end architecture test;
