load("@rules_nvc//nvc:rules.bzl", "vhdl_library")

exports_files(glob(["**/*.vhd"]))

filegroup(
    name = "grlib_srcs_all",
    srcs = glob(["**"]),
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "grlib_files",
    # do not sort
    srcs = [
        "lib/grlib/stdlib/version.vhd",
        "lib/grlib/stdlib/config_types.vhd",
        "@grlib//third_party/grlib:config.vhd",
        "lib/grlib/stdlib/stdlib.vhd",
        "lib/grlib/stdlib/stdio.vhd",
        "lib/grlib/stdlib/testlib.vhd",
        "lib/grlib/util/util.vhd",
        "lib/grlib/sparc/sparc.vhd",
        "lib/grlib/sparc/sparc_disas.vhd",
        "lib/grlib/riscv/riscv.vhd",
        "lib/grlib/riscv/riscv_disas.vhd",
        "lib/grlib/riscv/cpu_disas.vhd",
        "lib/grlib/modgen/multlib.vhd",
        "lib/grlib/modgen/leaves.vhd",
        "lib/grlib/amba/amba.vhd",
        "lib/grlib/amba/devices.vhd",
        "lib/grlib/amba/defmst.vhd",
        "lib/grlib/amba/apbctrl.vhd",
        "lib/grlib/amba/apbctrlx.vhd",
        "lib/grlib/amba/apbctrlsp.vhd",
        "lib/grlib/amba/apbctrldp.vhd",
        "lib/grlib/amba/apbctrl3p.vhd",
        "lib/grlib/amba/apbctrl4p.vhd",
        "lib/grlib/amba/ahbctrl.vhd",
        "lib/grlib/amba/dma2ahb_pkg.vhd",
        "lib/grlib/amba/dma2ahb.vhd",
        "lib/grlib/amba/ahbmst.vhd",
        "lib/grlib/amba/ahblitm2ahbm.vhd",
        "lib/grlib/amba/dma2ahb_tp.vhd",
        "lib/grlib/amba/amba_tp.vhd",
        "lib/grlib/dftlib/dftlib.vhd",
        "lib/grlib/dftlib/trstmux.vhd",
        "lib/grlib/dftlib/synciotest.vhd",
        "lib/grlib/generic_bm/generic_bm_pkg.vhd",
        "lib/grlib/generic_bm/ahb_be.vhd",
        "lib/grlib/generic_bm/axi4_be.vhd",
        "lib/grlib/generic_bm/bmahbmst.vhd",
        "lib/grlib/generic_bm/bm_fre.vhd",
        "lib/grlib/generic_bm/bm_me_rc.vhd",
        "lib/grlib/generic_bm/bm_me_wc.vhd",
        "lib/grlib/generic_bm/fifo_control_rc.vhd",
        "lib/grlib/generic_bm/fifo_control_wc.vhd",
        "lib/grlib/generic_bm/generic_bm_ahb.vhd",
        "lib/grlib/generic_bm/generic_bm_axi.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "grlib",
    # do not sort
    srcs = [] + ["lib/grlib/stdlib/version.vhd"] + ["lib/grlib/stdlib/config_types.vhd"] + ["@grlib//third_party/grlib:config.vhd"] + ["lib/grlib/stdlib/stdlib.vhd"] + select({
        "@grlib//:std_2008": ["@grlib//third_party/grlib:lib/grlib/stdlib/stdio_2008.vhd"],
        "@grlib//:std_2019": ["@grlib//third_party/grlib:lib/grlib/stdlib/stdio_2008.vhd"],
        "//conditions:default": ["lib/grlib/stdlib/stdio.vhd"],
    }) + select({
        "@grlib//:std_2008": ["@grlib//third_party/grlib:lib/grlib/stdlib/testlib_2008.vhd"],
        "@grlib//:std_2019": ["@grlib//third_party/grlib:lib/grlib/stdlib/testlib_2008.vhd"],
        "//conditions:default": ["lib/grlib/stdlib/testlib.vhd"],
    }) + ["lib/grlib/util/util.vhd"] + ["lib/grlib/sparc/sparc.vhd"] + ["lib/grlib/sparc/sparc_disas.vhd"] + ["lib/grlib/riscv/riscv.vhd"] + ["lib/grlib/riscv/riscv_disas.vhd"] + ["lib/grlib/riscv/cpu_disas.vhd"] + ["lib/grlib/modgen/multlib.vhd"] + ["lib/grlib/modgen/leaves.vhd"] + ["lib/grlib/amba/amba.vhd"] + ["lib/grlib/amba/devices.vhd"] + ["lib/grlib/amba/defmst.vhd"] + ["lib/grlib/amba/apbctrl.vhd"] + ["lib/grlib/amba/apbctrlx.vhd"] + ["lib/grlib/amba/apbctrlsp.vhd"] + ["lib/grlib/amba/apbctrldp.vhd"] + ["lib/grlib/amba/apbctrl3p.vhd"] + ["lib/grlib/amba/apbctrl4p.vhd"] + ["lib/grlib/amba/ahbctrl.vhd"] + ["lib/grlib/amba/dma2ahb_pkg.vhd"] + ["lib/grlib/amba/dma2ahb.vhd"] + ["lib/grlib/amba/ahbmst.vhd"] + ["lib/grlib/amba/ahblitm2ahbm.vhd"] + ["lib/grlib/amba/dma2ahb_tp.vhd"] + ["lib/grlib/amba/amba_tp.vhd"] + ["lib/grlib/dftlib/dftlib.vhd"] + ["lib/grlib/dftlib/trstmux.vhd"] + ["lib/grlib/dftlib/synciotest.vhd"] + ["lib/grlib/generic_bm/generic_bm_pkg.vhd"] + ["lib/grlib/generic_bm/ahb_be.vhd"] + ["lib/grlib/generic_bm/axi4_be.vhd"] + ["lib/grlib/generic_bm/bmahbmst.vhd"] + ["lib/grlib/generic_bm/bm_fre.vhd"] + ["lib/grlib/generic_bm/bm_me_rc.vhd"] + ["lib/grlib/generic_bm/bm_me_wc.vhd"] + ["lib/grlib/generic_bm/fifo_control_rc.vhd"] + ["lib/grlib/generic_bm/fifo_control_wc.vhd"] + ["lib/grlib/generic_bm/generic_bm_ahb.vhd"] + ["lib/grlib/generic_bm/generic_bm_axi.vhd"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "techmap_files",
    # do not sort
    srcs = [
        "lib/techmap/gencomp/gencomp.vhd",
        "lib/techmap/gencomp/netcomp.vhd",
        "lib/techmap/alltech/allclkgen.vhd",
        "lib/techmap/alltech/allddr.vhd",
        "lib/techmap/alltech/allmem.vhd",
        "lib/techmap/alltech/allmul.vhd",
        "lib/techmap/alltech/allpads.vhd",
        "lib/techmap/alltech/alltap.vhd",
        "lib/techmap/inferred/memory_inferred.vhd",
        "lib/techmap/inferred/ddr_inferred.vhd",
        "lib/techmap/inferred/mul_inferred.vhd",
        "lib/techmap/inferred/ddr_phy_inferred.vhd",
        "lib/techmap/inferred/ddrphy_datapath.vhd",
        "lib/techmap/inferred/fifo_inferred.vhd",
        "lib/techmap/inferred/sim_pll.vhd",
        "lib/techmap/inferred/lpddr2_phy_inferred.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "techmap",
    # do not sort
    srcs = [":techmap_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "spw_files",
    # do not sort
    srcs = [
        "lib/spw/comp/spwcomp.vhd",
        "lib/spw/wrapper/grspw_gen.vhd",
        "lib/spw/wrapper/grspw2_gen.vhd",
        "lib/spw/wrapper/grspw_codec_gen.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "spw",
    # do not sort
    srcs = [":spw_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
        ":techmap",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "eth_files",
    # do not sort
    srcs = [
        "lib/eth/comp/ethcomp.vhd",
        "lib/eth/core/greth_pkg.vhd",
        "lib/eth/core/eth_rstgen.vhd",
        "lib/eth/core/eth_edcl_ahb_mst.vhd",
        "lib/eth/core/eth_ahb_mst.vhd",
        "lib/eth/core/greth_tx.vhd",
        "lib/eth/core/greth_rx.vhd",
        "lib/eth/core/grethc.vhd",
        "lib/eth/wrapper/greth_gen.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "eth",
    # do not sort
    srcs = [":eth_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
        ":techmap",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "opencores_files",
    # do not sort
    srcs = [
        "lib/opencores/can/cancomp.vhd",
        "lib/opencores/can/can_top.vhd",
        "lib/opencores/i2c/i2c_master_bit_ctrl.vhd",
        "lib/opencores/i2c/i2c_master_byte_ctrl.vhd",
        "lib/opencores/i2c/i2coc.vhd",
        "lib/opencores/ge_1000baseX/ge_1000baseX_comp.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "opencores",
    # do not sort
    srcs = [":opencores_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "gaisler_files",
    # do not sort
    srcs = select({
        "@@//:NOELV_RV64": ["lib/gaisler/noelv/pkg/noelv_cfg_64.vhd"],
        "//conditions:default": ["lib/gaisler/noelv/pkg/noelv_cfg_32.vhd"],
    }) + [
        "lib/gaisler/arith/arith.vhd",
        "lib/gaisler/arith/mul32.vhd",
        "lib/gaisler/arith/div32.vhd",
        "lib/gaisler/memctrl/memctrl.vhd",
        "lib/gaisler/memctrl/sdctrl.vhd",
        "lib/gaisler/memctrl/sdctrl64.vhd",
        "lib/gaisler/memctrl/sdmctrl.vhd",
        "lib/gaisler/memctrl/srctrl.vhd",
        "lib/gaisler/srmmu/mmuconfig.vhd",
        "lib/gaisler/srmmu/mmuiface.vhd",
        "lib/gaisler/srmmu/libmmu.vhd",
        "lib/gaisler/srmmu/mmutlbcam.vhd",
        "lib/gaisler/srmmu/mmulrue.vhd",
        "lib/gaisler/srmmu/mmulru.vhd",
        "lib/gaisler/srmmu/mmutlb.vhd",
        "lib/gaisler/srmmu/mmutw.vhd",
        "lib/gaisler/srmmu/mmu.vhd",
        "lib/gaisler/leon3/leon3.vhd",
        "lib/gaisler/leon3/grfpushwx.vhd",
        "lib/gaisler/leon3v3/tbufmem.vhd",
        "lib/gaisler/leon3v3/tbufmem_2p.vhd",
        "lib/gaisler/leon3v3/dsu3x.vhd",
        "lib/gaisler/leon3v3/dsu3.vhd",
        "lib/gaisler/leon3v3/dsu3_mb.vhd",
        "lib/gaisler/leon3v3/libfpu.vhd",
        "lib/gaisler/leon3v3/libiu.vhd",
        "lib/gaisler/leon3v3/libcache.vhd",
        "lib/gaisler/leon3v3/libleon3.vhd",
        "lib/gaisler/leon3v3/regfile_3p_l3.vhd",
        "lib/gaisler/leon3v3/mmu_acache.vhd",
        "lib/gaisler/leon3v3/mmu_icache.vhd",
        "lib/gaisler/leon3v3/mmu_dcache.vhd",
        "lib/gaisler/leon3v3/cachemem.vhd",
        "lib/gaisler/leon3v3/mmu_cache.vhd",
        "lib/gaisler/leon3v3/grfpwx.vhd",
        "lib/gaisler/leon3v3/grlfpwx.vhd",
        "lib/gaisler/leon3v3/iu3.vhd",
        "lib/gaisler/leon3v3/proc3.vhd",
        "lib/gaisler/leon3v3/grfpwxsh.vhd",
        "lib/gaisler/leon3v3/leon3x.vhd",
        "lib/gaisler/leon3v3/leon3cg.vhd",
        "lib/gaisler/leon3v3/leon3s.vhd",
        "lib/gaisler/leon3v3/leon3sh.vhd",
        "lib/gaisler/leon3v3/l3stat.vhd",
        "lib/gaisler/leon3v3/cmvalidbits.vhd",
        "lib/gaisler/leon4/leon4.vhd",
        "lib/gaisler/irqmp/irqmp.vhd",
        "lib/gaisler/irqmp/irqamp.vhd",
        "lib/gaisler/irqmp/irqmp_bmode.vhd",
        "lib/gaisler/l2cache/pkg/l2cache.vhd",
        "lib/gaisler/can/can.vhd",
        "lib/gaisler/can/can_mod.vhd",
        "lib/gaisler/can/can_oc.vhd",
        "lib/gaisler/can/can_mc.vhd",
        "lib/gaisler/can/canmux.vhd",
        "lib/gaisler/can/can_rd.vhd",
        "lib/gaisler/canfd/canfd.vhd",
        "lib/gaisler/axi/axi.vhd",
        "lib/gaisler/axi/ahbm2axi.vhd",
        "lib/gaisler/axi/ahbm2axi3.vhd",
        "lib/gaisler/axi/ahbm2axi4.vhd",
        "lib/gaisler/axi/axinullslv.vhd",
        "lib/gaisler/axi/ahb2axib.vhd",
        "lib/gaisler/axi/ahb2axi3b.vhd",
        "lib/gaisler/axi/ahb2axi4b.vhd",
        "lib/gaisler/axi/ahb2axi_l.vhd",
        "lib/gaisler/axi/axis_buffer.vhd",
        "lib/gaisler/axi/axis_gearbox.vhd",
        "lib/gaisler/axi/axi4_resize.vhd",
        "lib/gaisler/axi/axi2ahb.vhd",
        "lib/gaisler/misc/misc.vhd",
        "lib/gaisler/misc/rstgen.vhd",
        "lib/gaisler/misc/gptimer.vhd",
        "lib/gaisler/misc/ahbram.vhd",
        "lib/gaisler/misc/ahbdpram.vhd",
        "lib/gaisler/misc/ahbtrace_mmb.vhd",
        "lib/gaisler/misc/ahbtrace_mb.vhd",
        "lib/gaisler/misc/ahbtrace.vhd",
        "lib/gaisler/misc/grgpio.vhd",
        "lib/gaisler/misc/ahbstat.vhd",
        "lib/gaisler/misc/logan.vhd",
        "lib/gaisler/misc/apbps2.vhd",
        "lib/gaisler/misc/charrom_package.vhd",
        "lib/gaisler/misc/charrom.vhd",
        "lib/gaisler/misc/apbvga.vhd",
        "lib/gaisler/misc/svgactrl.vhd",
        "lib/gaisler/misc/grsysmon.vhd",
        "lib/gaisler/misc/gracectrl.vhd",
        "lib/gaisler/misc/grgpreg.vhd",
        "lib/gaisler/misc/ahb_mst_iface.vhd",
        "lib/gaisler/misc/grgprbank.vhd",
        "lib/gaisler/misc/grversion.vhd",
        "lib/gaisler/misc/apb3cdc.vhd",
        "lib/gaisler/misc/ahbsmux.vhd",
        "lib/gaisler/misc/ahbmmux.vhd",
        "lib/gaisler/ambatest/ahbtbp.vhd",
        "lib/gaisler/ambatest/ahbtbm.vhd",
        "lib/gaisler/net/net.vhd",
        "lib/gaisler/pci/pci.vhd",
        "lib/gaisler/pci/pcipads.vhd",
        "lib/gaisler/pci/grpci2/pcilib2.vhd",
        "lib/gaisler/pci/grpci2/grpci2_ahb_mst.vhd",
        "lib/gaisler/pci/grpci2/grpci2_phy.vhd",
        "lib/gaisler/pci/grpci2/grpci2_phy_wrapper.vhd",
        "lib/gaisler/pci/grpci2/grpci2_cdc_gate.vhd",
        "lib/gaisler/pci/grpci2/grpci2.vhd",
        "lib/gaisler/pci/grpci2/wrapper/grpci2_gen.vhd",
        "lib/gaisler/pci/ptf/pt_pkg.vhd",
        "lib/gaisler/pci/ptf/pt_pci_master.vhd",
        "lib/gaisler/pci/ptf/pt_pci_target.vhd",
        "lib/gaisler/pci/ptf/pt_pci_arb.vhd",
        "lib/gaisler/uart/uart.vhd",
        "lib/gaisler/uart/libdcom.vhd",
        "lib/gaisler/uart/apbuart.vhd",
        "lib/gaisler/uart/apbuart_16550.vhd",
        "lib/gaisler/uart/dcom.vhd",
        "lib/gaisler/uart/dcom_uart.vhd",
        "lib/gaisler/uart/ahbuart.vhd",
        "lib/gaisler/sim/sim.vhd",
        "lib/gaisler/sim/sram.vhd",
        "lib/gaisler/sim/sram16.vhd",
        "lib/gaisler/sim/phy.vhd",
        "lib/gaisler/sim/ser_phy.vhd",
        "lib/gaisler/sim/ahbrep.vhd",
        "lib/gaisler/sim/delay_wire.vhd",
        "lib/gaisler/sim/pwm_check.vhd",
        "lib/gaisler/sim/ramback.vhd",
        "lib/gaisler/sim/slavecheck_slv.vhd",
        "lib/gaisler/sim/ddrram.vhd",
        "lib/gaisler/sim/ddr2ram.vhd",
        "lib/gaisler/sim/ddr3ram.vhd",
        "lib/gaisler/sim/sdrtestmod.vhd",
        "lib/gaisler/sim/ahbram_sim.vhd",
        "lib/gaisler/sim/aximem.vhd",
        "lib/gaisler/sim/axirep.vhd",
        "lib/gaisler/sim/axixmem.vhd",
        "lib/gaisler/sim/sramtestmod.vhd",
        "lib/gaisler/sim/delay_wire2.vhd",
        "lib/gaisler/sim/uartprint.vhd",
        "lib/gaisler/sim/dfi_phy_sim.vhd",
        "lib/gaisler/sim/dfi_phy_sim_fr.vhd",
        "lib/gaisler/sim/htif.vhd",
        "lib/gaisler/jtag/jtag.vhd",
        "lib/gaisler/jtag/libjtagcom.vhd",
        "lib/gaisler/jtag/jtagcom.vhd",
        "lib/gaisler/jtag/bscanregs.vhd",
        "lib/gaisler/jtag/bscanregsbd.vhd",
        "lib/gaisler/jtag/jtagcom2.vhd",
        "lib/gaisler/jtag/ahbjtag.vhd",
        "lib/gaisler/jtag/ahbjtag_exttap.vhd",
        "lib/gaisler/jtag/ahbjtag_bsd.vhd",
        "lib/gaisler/jtag/jtagcomrv.vhd",
        "lib/gaisler/jtag/ahbjtagrv.vhd",
        "lib/gaisler/jtag/jtagtst.vhd",
        "lib/gaisler/jtag/jtag_rv.vhd",
        "lib/gaisler/greth/ethernet_mac.vhd",
        "lib/gaisler/greth/greth.vhd",
        "lib/gaisler/greth/greth_mb.vhd",
        "lib/gaisler/greth/greth_gbit.vhd",
        "lib/gaisler/greth/greths.vhd",
        "lib/gaisler/greth/greth_gbit_mb.vhd",
        "lib/gaisler/greth/greths_mb.vhd",
        "lib/gaisler/greth/grethm.vhd",
        "lib/gaisler/greth/grethm_mb.vhd",
        "lib/gaisler/greth/adapters/rgmii.vhd",
        "lib/gaisler/greth/adapters/rgmii_kc705.vhd",
        "lib/gaisler/greth/adapters/rgmii_series7.vhd",
        "lib/gaisler/greth/adapters/rgmii_series6.vhd",
        "lib/gaisler/greth/adapters/comma_detect.vhd",
        "lib/gaisler/greth/adapters/sgmii.vhd",
        "lib/gaisler/greth/adapters/elastic_buffer.vhd",
        "lib/gaisler/greth/adapters/gmii_to_mii.vhd",
        "lib/gaisler/greth/adapters/word_aligner.vhd",
        "lib/gaisler/spacewire/spacewire.vhd",
        "lib/gaisler/hssl/hssl.vhd",
        "lib/gaisler/usb/grusb.vhd",
        "lib/gaisler/ddr/ddrpkg.vhd",
        "lib/gaisler/ddr/ddrintpkg.vhd",
        "lib/gaisler/ddr/ddrphy_wrap.vhd",
        "lib/gaisler/ddr/ddr2spax_ahb.vhd",
        "lib/gaisler/ddr/ddr2spax_ddr.vhd",
        "lib/gaisler/ddr/ddr2buf.vhd",
        "lib/gaisler/ddr/ddr2spax.vhd",
        "lib/gaisler/ddr/ddr2spa.vhd",
        "lib/gaisler/ddr/ddr1spax.vhd",
        "lib/gaisler/ddr/ddr1spax_ddr.vhd",
        "lib/gaisler/ddr/ddrspa.vhd",
        "lib/gaisler/ddr/ahb2mig_7series_pkg.vhd",
        "lib/gaisler/ddr/ahb2mig_7series.vhd",
        "lib/gaisler/ddr/ahb2mig_7series_ddr2_dq16_ad13_ba3.vhd",
        "lib/gaisler/ddr/ahb2mig_7series_ddr3_dq16_ad15_ba3.vhd",
        "lib/gaisler/ddr/ahb2mig_7series_cpci_xc7k.vhd",
        "lib/gaisler/ddr/ahb2axi_mig_7series.vhd",
        "lib/gaisler/ddr/axi_mig_7series.vhd",
        "lib/gaisler/ddr/ahb2avl_async.vhd",
        "lib/gaisler/ddr/ahb2avl_async_be.vhd",
        "lib/gaisler/gr1553b/gr1553b_pkg.vhd",
        "lib/gaisler/gr1553b/gr1553b_pads.vhd",
        "lib/gaisler/gr1553b/gr1553b_nlw.vhd",
        "lib/gaisler/gr1553b/gr1553b_stdlogic.vhd",
        "lib/gaisler/gr1553b/simtrans1553.vhd",
        "lib/gaisler/iommu/iommu.vhd",
        "lib/gaisler/i2c/i2c.vhd",
        "lib/gaisler/i2c/i2cmst.vhd",
        "lib/gaisler/i2c/i2cmst_gen.vhd",
        "lib/gaisler/i2c/i2cslv.vhd",
        "lib/gaisler/i2c/i2c2ahbx.vhd",
        "lib/gaisler/i2c/i2c2ahb.vhd",
        "lib/gaisler/i2c/i2c2ahb_apb.vhd",
        "lib/gaisler/i2c/i2c2ahb_gen.vhd",
        "lib/gaisler/i2c/i2c2ahb_apb_gen.vhd",
        "lib/gaisler/spi/spi.vhd",
        "lib/gaisler/spi/spimctrl.vhd",
        "lib/gaisler/spi/spictrlx.vhd",
        "lib/gaisler/spi/spictrl.vhd",
        "lib/gaisler/spi/spi2ahbx.vhd",
        "lib/gaisler/spi/spi2ahb.vhd",
        "lib/gaisler/spi/spi2ahb_apb.vhd",
        "lib/gaisler/spi/spi_flash.vhd",
        "lib/gaisler/nand/nandpkg.vhd",
        "lib/gaisler/grdmac/grdmac_pkg.vhd",
        "lib/gaisler/grdmac/apbmem.vhd",
        "lib/gaisler/grdmac/grdmac_ahbmst.vhd",
        "lib/gaisler/grdmac/grdmac_alignram.vhd",
        "lib/gaisler/grdmac/grdmac.vhd",
        "lib/gaisler/grdmac/grdmac_1p.vhd",
        "lib/gaisler/grdmac2/grdmac2_pkg.vhd",
        "lib/gaisler/grdmac2/grdmac2_apb.vhd",
        "lib/gaisler/grdmac2/mem2buf.vhd",
        "lib/gaisler/grdmac2/buf2mem.vhd",
        "lib/gaisler/grdmac2/grdmac2_ctrl.vhd",
        "lib/gaisler/grdmac2/grdmac2.vhd",
        "lib/gaisler/grdmac2/grdmac2_ahb.vhd",
        "lib/gaisler/grdmac2/grdmac2_acc.vhd",
        "lib/gaisler/grdmac2/grdmac2_axi.vhd",
        "lib/gaisler/subsys/subsys.vhd",
        "lib/gaisler/subsys/leon_dsu_stat_base.vhd",
        "lib/gaisler/plic/plic.vhd",
        "lib/gaisler/plic/grplic.vhd",
        "lib/gaisler/plic/plic_encoder.vhd",
        "lib/gaisler/plic/plic_gateway.vhd",
        "lib/gaisler/plic/plic_target.vhd",
        "lib/gaisler/plic/grplic_ahb.vhd",
        "lib/gaisler/aplic/aplic.vhd",
        "lib/gaisler/aplic/graplic_ahb.vhd",
        "lib/gaisler/aplic/aplic_encoder.vhd",
        "lib/gaisler/leon5/leon5.vhd",
        "lib/gaisler/leon5v0/leon5int.vhd",
        "lib/gaisler/leon5v0/leon5sys.vhd",
        "lib/gaisler/leon5v0/irqmp5.vhd",
        "lib/gaisler/leon5v0/l5stat.vhd",
        "lib/gaisler/leon5v0/nanofpu.vhd",
        "lib/gaisler/leon5v0/cpucore5int.vhd",
        "lib/gaisler/leon5v0/tbufmem5.vhd",
        "lib/gaisler/leon5v0/itbufmem5.vhd",
        "lib/gaisler/leon5v0/bht_pap.vhd",
        "lib/gaisler/leon5v0/btb.vhd",
        "lib/gaisler/leon5v0/inst_text.vhd",
        "lib/gaisler/leon5v0/iu5.vhd",
        "lib/gaisler/leon5v0/cctrl5.vhd",
        "lib/gaisler/leon5v0/tcmwrap5.vhd",
        "lib/gaisler/leon5v0/cachemem5.vhd",
        "lib/gaisler/leon5v0/regfile5_ram.vhd",
        "lib/gaisler/leon5v0/regfile5_dff.vhd",
        "lib/gaisler/leon5v0/cpucore5.vhd",
        "lib/gaisler/leon5v0/dbgmod5.vhd",
        "lib/gaisler/l2c_lite/l2c_lite.vhd",
        "lib/gaisler/l2c_lite/l2c_lite_core.vhd",
        "lib/gaisler/l2c_lite/l2c_lite_ahb.vhd",
        "lib/gaisler/l2c_lite/l2c_lite_axi3.vhd",
        "lib/gaisler/l2c_lite/l2c_lite_axi4.vhd",
        "lib/gaisler/nandfctrl2/nandfctrl2_pkg.vhd",
        "lib/gaisler/l5nv/shared/busif5_types.vhd",
        "lib/gaisler/l5nv/shared/l5nv_shared.vhd",
        "lib/gaisler/l5nv/shared/dmnv_ic_dmaport.vhd",
        "lib/gaisler/l5nv/shared/dmnv_ic_busport.vhd",
        "lib/gaisler/l5nv/shared/dmnv_ic_ebp.vhd",
        "lib/gaisler/l5nv/shared/dmnv_ic.vhd",
        "lib/gaisler/l5nv/shared/dmnv_ahbs.vhd",
        "lib/gaisler/l5nv/shared/tbufmemnv_mbus.vhd",
        "lib/gaisler/l5nv/shared/dmnv_trace.vhd",
        "lib/gaisler/l5nv/shared/dmnv_trace_ahb.vhd",
        "lib/gaisler/l5nv/shared/dmnv_reg_step.vhd",
        "lib/gaisler/l5nv/shared/l5tsc.vhd",
        "lib/gaisler/l5nv/shared/busif5x.vhd",
        "lib/gaisler/l5nv/shared/busif5rdb.vhd",
        "lib/gaisler/l5nv/shared/busif5.vhd",
        "lib/gaisler/l5nv/shared/tcmwrap5.vhd",
        "lib/gaisler/l5nv/shared/cachemem5.vhd",
        "lib/gaisler/l5nv/shared/snoopmem5.vhd",
        "lib/gaisler/noelv/pkg/noelv.vhd",
        "lib/gaisler/noelv/pkg/noelv_cpu_cfg.vhd",
        "lib/gaisler/noelv/pkg/nvnlconfig.vhd",
        "lib/gaisler/noelv/core/utilnv.vhd",
        "lib/gaisler/noelv/core/noelvtypes.vhd",
        "lib/gaisler/noelv/core/noelvint.vhd",
        "lib/gaisler/noelv/core/nvsupport.vhd",
        "lib/gaisler/noelv/core/mmuconfig.vhd",
        "lib/gaisler/noelv/core/bhtnv.vhd",
        "lib/gaisler/noelv/core/btbnv.vhd",
        "lib/gaisler/noelv/core/btbdmnv.vhd",
        "lib/gaisler/noelv/core/rasnv.vhd",
        "lib/gaisler/noelv/core/tbufmemnv.vhd",
        "lib/gaisler/noelv/core/fputilnv.vhd",
        "lib/gaisler/noelv/core/mul64.vhd",
        "lib/gaisler/noelv/core/div64.vhd",
        "lib/gaisler/noelv/core/regfile64sramnv.vhd",
        "lib/gaisler/noelv/core/regfile64dffnv.vhd",
        "lib/gaisler/noelv/core/alunv.vhd",
        "lib/gaisler/noelv/core/rvvi.vhd",
        "lib/gaisler/noelv/core/iunv.vhd",
        "lib/gaisler/noelv/core/itracenv.vhd",
        "lib/gaisler/noelv/core/cctrl5nv.vhd",
        "lib/gaisler/noelv/core/mulfp.vhd",
        "lib/gaisler/noelv/core/nanofpunv.vhd",
        "lib/gaisler/noelv/core/cpucorenvbc.vhd",
        "lib/gaisler/noelv/core/cpucorenvb.vhd",
        "lib/gaisler/noelv/core/cpucorenv.vhd",
        "lib/gaisler/noelv/core/interrupt_file.vhd",
        "lib/gaisler/noelv/core/imsic_int_files.vhd",
        "lib/gaisler/noelv/core/fpsimutilnv.vhd",
        "lib/gaisler/noelv/core/kanatalog.vhd",
        "lib/gaisler/noelv/subsys/noelvcpu.vhd",
        "lib/gaisler/noelv/subsys/dummy_pnp.vhd",
        "lib/gaisler/noelv/subsys/noelvcfgmap.vhd",
        "lib/gaisler/noelv/subsys/noelvsys.vhd",
        "lib/gaisler/noelv/aclint/clint.vhd",
        "lib/gaisler/noelv/aclint/clint_ahb.vhd",
        "lib/gaisler/noelv/aclint/aclint_ahb.vhd",
        "lib/gaisler/noelv/imsic/imsic_ahb.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "gaisler",
    # do not sort
    srcs = [":gaisler_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":eth",
        ":grlib",
        ":opencores",
        ":techmap",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "esa_files",
    # do not sort
    srcs = [
        "lib/esa/memoryctrl/memoryctrl.vhd",
        "lib/esa/memoryctrl/mctrl.vhd",
        "lib/esa/pci/pcicomp.vhd",
        "lib/esa/pci/pci_arb_pkg.vhd",
        "lib/esa/pci/pci_arb.vhd",
        "lib/esa/pci/pciarb.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "esa",
    # do not sort
    srcs = [":esa_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":gaisler",
        ":grlib",
        ":techmap",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "fmf_files",
    # do not sort
    srcs = [
        "lib/fmf/utilities/conversions.vhd",
        "lib/fmf/utilities/gen_utils.vhd",
        "lib/fmf/flash/flash.vhd",
        "lib/fmf/flash/s25fl064a.vhd",
        "lib/fmf/flash/m25p80.vhd",
        "lib/fmf/fifo/idt7202.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "fmf",
    # do not sort
    srcs = [":fmf_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "gsi_files",
    # do not sort
    srcs = [
        "lib/gsi/ssram/functions.vhd",
        "lib/gsi/ssram/core_burst.vhd",
        "lib/gsi/ssram/g880e18bt.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "gsi",
    # do not sort
    srcs = [":gsi_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "cypress_files",
    # do not sort
    srcs = [
        "lib/cypress/ssram/components.vhd",
        "lib/cypress/ssram/package_utility.vhd",
        "lib/cypress/ssram/cy7c1354b.vhd",
        "lib/cypress/ssram/cy7c1380d.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "cypress",
    # do not sort
    srcs = [":cypress_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "micron_files",
    # do not sort
    srcs = [
        "lib/micron/sdram/components.vhd",
        "lib/micron/sdram/mt48lc16m16a2.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "micron",
    # do not sort
    srcs = [":micron_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "atc18_files",
    # do not sort
    srcs = [
        "lib/tech/atc18/components/atmel_components.vhd",
        "lib/tech/atc18/components/atmel_simprims.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "atc18",
    # do not sort
    srcs = [":atc18_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "dware_files",
    # do not sort
    srcs = [
        "lib/tech/dware/simprims/DWpackages.vhd",
        "lib/tech/dware/simprims/DW_Foundation_arith.vhd",
        "lib/tech/dware/simprims/DW_Foundation_comp.vhd",
        "lib/tech/dware/simprims/DW_Foundation_comp_arith.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "dware",
    # do not sort
    srcs = [":dware_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "ec_files",
    # do not sort
    srcs = [
        "lib/tech/ec/orca/orcacomp.vhd",
        "lib/tech/ec/orca/global.vhd",
        "lib/tech/ec/orca/orca.vhd",
        "lib/tech/ec/orca/orca_ecmem.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "ec",
    # do not sort
    srcs = [":ec_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "eclipsee_files",
    # do not sort
    srcs = [
        "lib/tech/eclipsee/simprims/eclipse.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "eclipsee",
    # do not sort
    srcs = [":eclipsee_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "simprim_files",
    # do not sort
    srcs = [
        "lib/tech/simprim/vcomponents/vcomponents.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "simprim",
    # do not sort
    srcs = [":simprim_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "umc18_files",
    # do not sort
    srcs = [
        "lib/tech/umc18/components/umc_components.vhd",
        "lib/tech/umc18/components/umc_simprims.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "umc18",
    # do not sort
    srcs = [":umc18_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "virage_files",
    # do not sort
    srcs = [
        "lib/tech/virage/vcomponents/virage_vcomponents.vhd",
        "lib/tech/virage/simprims/virage_simprims.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "virage",
    # do not sort
    srcs = [":virage_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "testgrouppolito_files",
    # do not sort
    srcs = [
        "lib/testgrouppolito/pr/dprc_pkg.vhd",
        "lib/testgrouppolito/pr/dprc.vhd",
        "lib/testgrouppolito/pr/sync_dprc.vhd",
        "lib/testgrouppolito/pr/async_dprc.vhd",
        "lib/testgrouppolito/pr/d2prc.vhd",
        "lib/testgrouppolito/pr/d2prc_edac.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "testgrouppolito",
    # do not sort
    srcs = [":testgrouppolito_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":grlib",
        ":unisim",
        ":techmap",
    ],
    visibility = ["//visibility:public"],
)

# do not sort
filegroup(
    name = "work_files",
    # do not sort
    srcs = [
        "lib/work/debug/debug.vhd",
        "lib/work/debug/grtestmod.vhd",
        "lib/work/debug/cpu_disas.vhd",
    ],
    visibility = ["//visibility:public"],
)

vhdl_library(
    name = "work",
    # do not sort
    srcs = [":work_files"],
    standard = select({
        "@grlib//:std_1987": "1987",
        "@grlib//:std_1993": "1993",
        "@grlib//:std_2002": "2002",
        "@grlib//:std_2008": "2008",
        "@grlib//:std_2019": "2019",
        "//conditions:default": "1993",
    }),
    deps = [
        ":gaisler",
        ":grlib",
    ],
    visibility = ["//visibility:public"],
)

