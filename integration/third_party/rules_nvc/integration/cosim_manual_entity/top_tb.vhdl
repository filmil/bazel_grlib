-- LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4

library ieee;
use ieee.std_logic_1164.all;

library cosim_bridge;

entity top_tb is
end entity;

architecture sim of top_tb is
  signal clk : std_logic := '0';
  signal d   : std_logic := '0';
  signal q   : std_logic;

begin
  clk <= not clk after 5 ns;

  process
  begin
    wait for 10 ns;
    d <= '1';
    wait for 15 ns; -- wait for clock edge + delta
    wait for 2 ps;

    assert q = '1' report "Verification failed: q != '1'. Actual: " & to_string(q) severity error;

    d <= '0';
    wait for 15 ns;
    wait for 2 ps;

    assert q = '0' report "Verification failed: q != '0'. Actual: " & to_string(q) severity error;

    assert false report "Simulation finished successfully" severity note;
    std.env.finish;
    wait;
  end process;

  dut_inst : entity cosim_bridge.dut
    generic map (
      INSTANCE_ID => 1
    )
    port map (
      clk => clk,
      d   => d,
      q   => q
    );

end architecture;
