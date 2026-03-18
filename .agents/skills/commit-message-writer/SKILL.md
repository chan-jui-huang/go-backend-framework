---
name: commit-message-writer
description: Generate Git commit messages for this repository by combining conversation intent with staged code changes. Use when the user asks for a commit message, wants help summarizing staged changes, or needs an English Conventional Commit based on git diff.
compatibility: Designed for repository-local agents that can inspect git status and staged diffs.
metadata:
  author: workspace
  version: "0.1.0"
---

# Commit Message Writer

Use this skill when the user wants a commit message for finished work in this repository.

## What this skill does

This skill produces a concise English commit message that follows Conventional Commits and the repository's commit-writing rules.

## Inputs

- Current conversation about the change
- Staged git diff when available
- Staged file list and status when useful for quick verification

## Workflow

1. Read [the commit message guidelines](references/commit-message-guidelines.md) before drafting the message.
2. Prefer the full conversation plus staged changes as the source of truth.
3. If conversation context is insufficient, inspect staged changes and summarize the behavioral change from the diff alone.
4. Write the commit message in English using the Conventional Commits format required by this repository.
5. Keep the subject specific and action-oriented. Add a body only when it materially improves clarity.
6. Verify grammar, spelling, and scope before returning the final message.

## Drafting rules

- Prefer `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `build`, `ci`, or `perf` based on the actual change.
- Include a scope when it clarifies the touched area.
- Describe the user-visible or developer-visible outcome, not a low-level file list.
- Do not mention AI tools, prompts, or that the message was generated unless the user explicitly asks for that.
- If the change is breaking, use the Conventional Commits breaking-change form.

## When staged changes are missing

If the user asks for a commit message but there are no staged changes and the conversation alone is not enough to identify the final diff, ask for staged changes or inspect the relevant diff before writing the message.

## Output

Return only the proposed commit message unless the user asks for alternatives or rationale.

## Example prompts

```text
We're done. Please generate a commit message for me.
```

```text
Please generate a commit message for the staged changes.
```
