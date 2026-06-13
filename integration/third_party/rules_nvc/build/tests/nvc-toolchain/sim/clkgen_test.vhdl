library ieee; use ieee.std_logic_1164.all;
library sim; use sim.all;

entity tb is
end entity;

architecture tb of tb is

    -- 1GHz clock.
    constant clock_period: time := 1 ns;

    -- Run the simulation for this long.
    constant sim_run_length: time := 30 * clock_period;

    -- Top level signals.
    signal clk, rst: std_logic;

begin

    reset: entity resetgen(sim) port map(rst);

    uut: entity clkgen(sim) generic map (clock_period) port map (rst, clk);

    -- Terminates the simulation after `sim_run_length` time has elapsed.
    sim_end: process
    begin
       wait for sim_run_length;
       std.env.finish; -- VHDL-2008.
    end process;

end architecture;
