library ieee;
use ieee.std_logic_1164.all;
use ieee.numeric_std.all;

library cosim_bridge;

entity top_tb is
end entity top_tb;

architecture sim of top_tb is
  -- The bridged component will have generics fixed by the verilator build,
  -- but we can still declare them if needed for our own testbench logic.
  -- Alternatively, the generated proxy will just have fixed width ports
  -- based on the elaborated verilog parameters. Let's look at the generated
  -- proxy port widths. If we parameterized it with WIDTH=16, then a, b are 16 bits.
  signal clk : std_logic := '0';
  signal a, b : std_logic_vector(15 downto 0) := (others => '0');
  signal sum : std_logic_vector(16 downto 0);

  component adder is
    generic (
      INSTANCE_ID : integer
    );
    port (
      clk : in std_logic;
      a : in std_logic_vector(15 downto 0);
      b : in std_logic_vector(15 downto 0);
      sum : out std_logic_vector(16 downto 0)
    );
  end component;

begin

  dut_inst: component adder
    generic map (
      INSTANCE_ID => 0
    )
    port map (
      clk => clk,
      a => a,
      b => b,
      sum => sum
    );

  clk_proc: process
  begin
    wait for 5 ns;
    clk <= not clk;
  end process;

  stim_proc: process
  begin
    wait until rising_edge(clk);
    a <= x"0010";
    b <= x"0020";
    wait until rising_edge(clk);
    wait until rising_edge(clk);
    wait for 2 ps; -- allow VPI to evaluate

    report "Value of sum: " & to_string(sum);
    assert sum = '0' & x"0030" report "Addition failed! Expected: " & to_string(std_logic_vector'('0' & x"0030")) & ", Actual: " & to_string(sum) severity failure;

    a <= x"1234";
    b <= x"4321";
    wait until rising_edge(clk);
    wait until rising_edge(clk);
    wait for 2 ps; -- allow VPI to evaluate

    report "Value of sum: " & to_string(sum);
    assert sum = '0' & x"5555" report "Addition failed! Expected: " & to_string(std_logic_vector'('0' & x"5555")) & ", Actual: " & to_string(sum) severity failure;
    report "Test completed successfully.";
    std.env.stop;
  end process;

end architecture sim;
