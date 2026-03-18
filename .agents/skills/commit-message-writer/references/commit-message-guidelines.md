## Commit Message Writing Guidelines

Use conversation history and staged code changes together whenever possible. Conversation explains intent, while `git diff --staged` shows the final implementation that will actually be committed.

### Recommended Workflow

#### Preferred Approach: Conversation + Staged Diff

Use this approach when the current conversation reflects the work that is about to be committed.

**Workflow:**
1. Finish the task and stage the intended changes.
2. Review the conversation context for the goal, constraints, and notable implementation decisions.
3. Review `git diff --staged` to confirm the final committed behavior.
4. Draft the commit message from both sources.

**Example Prompt:**
> "We're done. Please generate a commit message for me."

#### Alternate Approach: Staged Diff Only

Use this approach when there is no reliable conversation context, such as work completed offline.

**Workflow:**
1. Stage the intended changes.
2. Run `git diff --staged`.
3. Draft the commit message from the staged diff alone.

**Example Prompt:**
> "Please generate a commit message for the following `git diff`:"
> ```diff
> diff --git a/src/utils/math.js b/src/utils/math.js
> index 6e9b2f7..8b4e6ad 100644
> --- a/src/utils/math.js
> +++ b/src/utils/math.js
> @@ -1,5 +1,9 @@
>  function add(a, b) {
>    return a + b;
>  }
> +
> +function subtract(a, b) {
> +  return a - b;
> +}
>
> -module.exports = { add };
> +module.exports = { add, subtract };
> ```

### Required Format

Follow the Conventional Commits 1.0.0 specification:
https://www.conventionalcommits.org/en/v1.0.0/

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Writing Rules

- Write the message in English.
- Keep the subject line under 72 characters.
- Use imperative mood, such as `add feature` instead of `added feature`.
- Make the subject describe the outcome of the change, not a file-by-file edit list.
- Add a body only when it improves clarity.
- Reference issues in footers when relevant, for example `Fixes #123` or `Relates to #456`.
- Include `Co-Authored-By` only when it is actually appropriate for the commit.
- Never list AI, agents, or automated tools in `Co-Authored-By`.
- `Co-Authored-By` entries must refer to human collaborators only.
- Review grammar, spelling, and punctuation before finalizing the message.

### Git Safety Rules

- NEVER commit secrets. Use environment-specific configuration files.
- NEVER update git config via automation.
- NEVER run destructive git commands such as `push --force` or `hard reset` without explicit user request.
- NEVER skip git hooks such as `--no-verify` or `--no-gpg-sign` unless explicitly requested.

### Example

```text
feat(user): add password reset endpoint

Implements user story US-042 for self-service password reset.
Adds new POST /api/user/reset-password endpoint with email validation.

BREAKING CHANGE: /api/user/password endpoint removed; use /api/user/reset-password

Fixes #123
```
