-- NVC elaboration-bug reproducer (NOT part of the normal test suite).
--
-- This is the original "full NOEL-V subsystem" testbench. It instantiates
-- gaisler.noelvsys, which no released NVC version can currently elaborate:
--
--   * NVC 1.18.x : "** Fatal: (init): missing body for ... CCTRL5NV.PA_MSB()I"
--   * NVC 1.19.0 - 1.21.0 : "** Fatal: (init): invalid container kind T_ELAB
--                            for NAHBIRQ" (cf. nickg/nvc#1537)
--
-- It is wired to a `manual`-tagged Bazel target so `bazel test //...` stays
-- green; build it explicitly to reproduce the crash:
--
--   bazel build //:noelvsys_repro_noelvsys_repro_tb \
--       --@grlib//:ACTIVE_DESIGN_PREFIX=LIB_GAISLER_NOELV_NOELV
--
-- See docs/compilation_challenges.md for the debugging synopsis and the plan
-- for reducing this to a minimal reproducer.

library ieee;
use ieee.std_logic_1164.all;

library grlib;
use grlib.config_types.all;
use grlib.config.all;
use grlib.amba.all;
use grlib.stdlib.all;
use grlib.devices.all;

library gaisler;
use gaisler.noelv.all;
use gaisler.uart.all;
use gaisler.misc.all;
use gaisler.jtag.all;

entity noelvsys_repro_tb is
end entity;

architecture behav of noelvsys_repro_tb is
  signal clk      : std_ulogic := '0';
  signal rstn     : std_ulogic := '0';

  -- Use promoted names from grlib.config
  constant ncpu   : integer := CFG_PROC_NUM;
  signal gclk     : std_logic_vector(ncpu-1 downto 0) := (others => '1');

  signal ahbmi    : ahb_mst_in_type;
  -- Fixed vector bounds for ahbmo. Formal is (ncpu + nextmst - 1 downto ncpu)
  -- If nextmst = 0, length is 0.
  signal ahbmo    : ahb_mst_out_vector_type(ncpu to ncpu-1);

  signal ahbsi    : ahb_slv_in_type;
  -- nextslv = 1, formal is (nextslv - 1 downto 0) -> (0 downto 0)
  signal ahbso    : ahb_slv_out_vector_type(0 downto 0) := (others => ahbs_none);

  -- ndbgmst = 1, formal is (ndbgmst - 1 downto 0) -> (0 downto 0)
  signal dbgmi    : ahb_mst_in_vector_type(0 downto 0);
  signal dbgmo    : ahb_mst_out_vector_type(0 downto 0) := (others => ahbm_none);

  signal apbi     : apb_slv_in_type;
  signal apbo     : apb_slv_out_vector; -- Constrained in amba.vhd

  signal uarti    : uart_in_type := (rxd => '1', others => '0');
  signal uarto    : uart_out_type;

begin

  -- Clock generation
  clk <= not clk after 10 ns;

  -- Reset generation
  process
  begin
    rstn <= '0';
    wait for 100 ns;
    rstn <= '1';
    wait;
  end process;

  -- Instantiate NOEL-V subsystem
  sys0: entity gaisler.noelvsys
    generic map (
      fabtech  => 0,
      memtech  => 0,
      ncpu     => ncpu,
      nextmst  => 0,
      nextslv  => 1,
      nextapb  => 0,
      ndbgmst  => 1,
      nintdom  => CFG_DOMAINS_NUM,
      neiid    => CFG_EIID_NUM,
      cached   => 0,
      wbmask   => 0,
      busw     => CFG_AHBDW,
      cmemconf => 0,
      rfconf   => 0,
      fpuconf  => 0,
      tcmconf  => 0,
      mulconf  => 0,
      intcconf => 0,
      disas    => 0,
      ahbtrace => 0,
      cfg      => 0,
      devid    => 0,
      nodbus   => 0,
      trace    => 0,
      scantest => 0
    )
    port map (
      clk      => clk,
      gclk     => gclk,
      rstn     => rstn,
      pwrd     => open,
      ahbmi    => ahbmi,
      ahbmo    => ahbmo,
      ahbsi    => ahbsi,
      ahbso    => ahbso,
      dbgmi    => dbgmi,
      dbgmo    => dbgmo,
      apbi     => apbi,
      apbo     => apbo,
      dsuen    => '1',
      dsubreak => '0',
      cpu0errn => open,
      uarti    => uarti,
      uarto    => uarto,
      cnt      => open,
      etso     => open,
      etsi     => (others => nv_etrace_sink_in_none),
      testen   => '0',
      testrst  => '1',
      scanen   => '0',
      testoen  => '1',
      testsig  => (others => '0')
    );

  process
  begin
    report "NOEL-V Integration Testbench started";
    report "XLEN is: " & integer'image(CFG_NOELV_XLEN);

    wait until rstn = '1';
    wait for 1000 ns;

    report "Testbench finished successfully";
    assert false report "Simulation finished" severity failure;
    wait;
  end process;

end architecture;
