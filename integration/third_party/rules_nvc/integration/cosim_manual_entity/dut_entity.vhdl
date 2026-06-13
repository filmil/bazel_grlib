-- LICENSE sha256: c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4

library ieee;
use ieee.std_logic_1164.all;
use ieee.numeric_std.all;

entity dut is
  generic (
    INSTANCE_ID : integer
  );
  port (
    clk : in std_logic;
    d : in std_logic;
    q : out std_logic
  );
end entity;
