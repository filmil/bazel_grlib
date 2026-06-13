library ieee;
use ieee.std_logic_1164.all;
use ieee.numeric_std.all;

library adder_bridge;
library counter_bridge;

entity top_tb is
end entity;

architecture sim of top_tb is
  signal clk    : std_logic := '0';
  signal rst    : std_logic := '1';
  signal enable : std_logic := '0';
  signal count  : std_logic_vector(7 downto 0);

  signal a   : std_logic_vector(7 downto 0) := x"05";
  signal b   : std_logic_vector(7 downto 0) := x"0A";
  signal sum : std_logic_vector(8 downto 0);

begin
  clk <= not clk after 5 ns;

  process
  begin
    report "Starting simulation...";
    wait for 15 ns;
    rst <= '0';
    report "Reset deasserted at 15ns";
    wait for 10 ns;
    enable <= '1';
    report "Enable asserted at 25ns";

    wait for 100 ns;
    report "Finished waiting 100ns";

    -- Test adder
    a <= x"FF";
    b <= x"01";
    report "Testing adder inputs FF + 01";
    wait for 15 ns; -- Wait past clock edge

    -- When C++ VPI linking is implemented, uncomment:
    assert sum = "100000000" report "Verification failed: sum != 100" severity error;
    assert count /= x"00" report "Verification failed: counter is zero" severity error;

    report "Reached end of stimulus process";

    assert false report "Simulation finished successfully" severity note;
    std.env.finish;
    wait;
  end process;

  adder_inst : entity adder_bridge.adder
    generic map (
      INSTANCE_ID => 1
    )
    port map (
      a   => a,
      b   => b,
      sum => sum
    );

  counter_inst : entity counter_bridge.counter
    generic map (
      INSTANCE_ID => 2
    )
    port map (
      clk    => clk,
      rst    => rst,
      enable => enable,
      count  => count
    );

end architecture;
