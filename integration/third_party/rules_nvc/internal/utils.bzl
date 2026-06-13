
# XXX: this should not be necessary, there's a ctx.files.foo attribute for target foo
def get_single_file_from(target):
    """
    Retrieves the single file associated with a target.

    Args:
        target: The target from which to extract the file.

    Returns:
        The single `File` object from the target.
    """
    file_list = target.files.to_list()
    # assert 1 file
    return file_list[0]

def get_nvc_deps(nvc_info):
    """Returns the list of dependencies from nvc_info."""
    if hasattr(nvc_info, "deps"):
        return nvc_info.deps
    return []

def get_nvc_ld_library_path(nvc_info, base_dir, default_env):
    """Constructs the LD_LIBRARY_PATH environment variable."""
    paths = { base_dir + "/lib/x86_64-linux-gnu": True }
    if hasattr(nvc_info, "deps"):
        for dep in nvc_info.deps:
            if ".so" in dep.basename:
                paths[dep.dirname] = True
    ld_path = ":".join(paths.keys())
    if default_env.get("LD_LIBRARY_PATH", ""):
        ld_path += ":" + default_env.get("LD_LIBRARY_PATH", "")
    return ld_path

