# Project Standards and Best Practices

This project follows the standards outlined in the [AI Coding SOP](https://github.com/filmil/ai-coding-sop).

## General Development Principles

- **Lifecycle:** Follow a **Research -> Strategy -> Execution** lifecycle for all tasks.
- **Bug Fixes:** ALWAYS empirically reproduce a bug with a test case before attempting a fix.
- **Verification:** ALL changes must be verified through automated tests and workspace standards.
- **Communication:** Use `update_topic` regularly to keep users informed of progress and strategic shifts during multi-turn tasks.
- **Conventions:** Strictly adhere to local workspace conventions, architectural patterns, and style.
- **Efficiency:** Prioritize surgical edits and utilize specialized sub-agents for repetitive or high-volume tasks to maintain context efficiency.
- **Protected Files:** Never auto-modify dotfiles, `*.lock` files, or `*.nix` files.

## Git Standards

- **Conventional Commits:** Use [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) for all git commit messages.
  - Example: `feat: add VHDL 2019 support`
- **Maintenance Scope:** Unless otherwise instructed, only apply maintenance tasks to files in the git index or uncommitted files to avoid redoing work on already committed files.
- **Commit Metadata:** Every commit created by an automated assistant should contain this note as the last line:
  ```
  This commit has been created by an automated coding assistant,
  with human supervision.
  ```
  Also append the prompt used to generate the commit in full.
- **Workflow:** Prefer **rebase** over merge. Use `git rebase --pull origin main` to sync with the main branch.

## Bazel Standards

- **Restricted Files:** Never auto-modify `//tools/bazel`.
- **Formatting:** Do **not** run `buildifier`, as it may disrupt the specific ordering required for VHDL files.
- **Documentation Build:** When updating documentation, verify correctness by running `bazel build //:docs`.
- **Target Maintenance:** When adding Doxygen documentation, ensure the source filegroup targets (or individual source files) are added to the `srcs` attribute of the `doxygen` target named `//:docs`.

## Engineering & Documentation

- **Documentation Standards:** 
  - Use **Doxygen** rules for all documentation.
  - Include all VHDL files, C headers, and any other program source files containing documentation.
  - Maintain up-to-date documentation for the public API of all source files.
- **License Maintenance:**
  - Every source file and `BUILD` file must have a license reference at the beginning.
  - If a license reference is missing, add the appropriate **SPDX label**.
  - Do not modify `*.gtkw` files or files under `//third_party` during license updates.
  - Every subdirectory under `//third_party` must contain a `LICENSE` file copied from its source distribution.
