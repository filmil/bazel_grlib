# Project Standards and Best Practices

This project follows the standards outlined in the [AI Coding SOP](https://github.com/filmil/ai-coding-sop).

## General Development Principles

- **Lifecycle:** Follow a **Research -> Strategy -> Execution** lifecycle for all tasks.
- **Bug Fixes:** ALWAYS empirically reproduce a bug with a test case before attempting a fix.
- **Verification:** ALL changes must be verified through automated tests and workspace standards.
- **Communication:** Use `update_topic` regularly to keep users informed of progress and strategic shifts during multi-turn tasks.
- **Conventions:** Strictly adhere to local workspace conventions, architectural patterns, and style.
- **Efficiency:** Prioritize surgical edits and utilize specialized sub-agents for repetitive or high-volume tasks to maintain context efficiency.

## Git Commit Standards

- **Conventional Commits:** Use [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) for all git commit messages.
  - Example: `feat: add VHDL 2019 support`
  - Example: `fix: resolve name collision in stdio.vhd`
  - Example: `docs: update README with version flag instructions`
