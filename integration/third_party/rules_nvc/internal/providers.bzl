NVCInfo = provider(
    doc = "Information on how to run NVC for VHDL analysis, elaboration and simulation.",
    fields = {
      "analyzer": "The NVC analyzer executable file.",
      "artifacts_dir": "The directory containing NVC standard libraries and artifacts.",
      "deps": "List[File]: Additional toolchain dependencies (e.g. shared libraries).",
    }
)

VHDLLibraryProvider = provider(
    doc = "Contains the information about the binary files in a compiled VHDL library.",
    fields = {
        "libraries": "List[(string, File)]: A mapping from a library name to its directory location. Contains both this library and its dependencies, ensuring no duplicate keys.",
        "entities": "List[string]: The entities emphasized in this library.",
        "library_name": "string: The name of the library (e.g., `ieee`, `work`).",
        "library_dir": """File: The container directory where the library is located.
        NVC will not create a library in an existing directory, so a container directory is used.
        The actual library directory is `$library_dir/$library_name`.""",
        "includes": "List[string]: List of directories to include for Verilog.",
        "hdrs": "depset[File]: List of include files for Verilog.",
        "vpi_plugins": "depset[File]: List of VPI plugins required for simulation.",
    },
)

ElaborateProvider = provider(
    doc = "Provides information about an elaborated VHDL entity.",
    fields = {
        "entity": "string: The name of the elaborated entity.",
    }
)

