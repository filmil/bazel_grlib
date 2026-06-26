-- NOEL-V / GRLIB integration testbench.
--
-- This testbench demonstrates consuming the `grlib` Bazel module from an
-- external workspace:
--
--   * it links against the generated `grlib`, `gaisler` and `techmap`
--     libraries provided by the module;
--   * it reads the promoted NOEL-V configuration constants out of the
--     `grlib.config` package (CFG_PROC_NUM, CFG_NOELV_XLEN, ...), which are
--     generated from Kconfig when ACTIVE_DESIGN_PREFIX=LIB_GAISLER_NOELV_NOELV
--     (see .bazelrc); and
--   * it elaborates and simulates a real GRLIB core under NVC.
--
-- It instantiates a single `apbuart` on a minimal APB bus. The full NOEL-V
-- subsystem (gaisler.noelvsys) cannot currently be *elaborated* by any released
-- NVC version -- see docs/compilation_challenges.md -- so this testbench
-- exercises the same toolchain and libraries with a core NVC can handle.

library ieee;
use ieee.std_logic_1164.all;

library grlib;
use grlib.config_types.all;
use grlib.config.all;
use grlib.amba.all;
use grlib.stdlib.all;

library gaisler;
use gaisler.uart.all;

entity noelv_tb is
end entity;

architecture behav of noelv_tb is
  signal clk   : std_ulogic := '0';
  signal rstn  : std_ulogic := '0';
  signal done  : boolean    := false;

  -- Minimal APB fabric: the UART is the only slave.
  signal apbi  : apb_slv_in_type  := apb_slv_in_none;
  signal apbo  : apb_slv_out_type;

  signal uarti : uart_in_type := (rxd => '1', ctsn => '0', extclk => '0');
  signal uarto : uart_out_type;
begin

  -- Free-running clock that stops once the stimulus is finished so the
  -- simulation quiesces and NVC exits cleanly with status 0.
  clkgen : process
  begin
    while not done loop
      clk <= '0';
      wait for 5 ns;
      clk <= '1';
      wait for 5 ns;
    end loop;
    wait;
  end process;

  -- Device under test: a GRLIB APB UART.
  uart0 : apbuart
    generic map (
      pindex   => 0,
      paddr    => 0,
      pirq     => 0,
      console  => 0,
      fifosize => 1)
    port map (
      rst   => rstn,
      clk   => clk,
      apbi  => apbi,
      apbo  => apbo,
      uarti => uarti,
      uarto => uarto);

  stim : process
  begin
    report "NOEL-V integration testbench started";
    report "  CFG_AHBDW      = " & integer'image(CFG_AHBDW);
    report "  CFG_PROC_NUM   = " & integer'image(CFG_PROC_NUM);
    report "  CFG_NOELV_XLEN = " & integer'image(CFG_NOELV_XLEN);
    report "  CFG_DOMAINS_NUM= " & integer'image(CFG_DOMAINS_NUM);
    report "  CFG_EIID_NUM   = " & integer'image(CFG_EIID_NUM);

    -- Hold reset, then release and let the UART run for a while.
    rstn <= '0';
    wait for 100 ns;
    rstn <= '1';
    wait for 1000 ns;

    report "Testbench finished successfully";
    done <= true;
    wait;
  end process;

end architecture;
