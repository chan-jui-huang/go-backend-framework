---
name: agent-customization
description: Create an agent skill customization file (SKILL.md) that packages a workflow.
user-invocable: true
metadata:
  version: '0.1.0'
  description: 'Create a reusable SKILL.md that packages a workflow extracted from a conversation.'
  scope: workspace
  applyTo:
    - '.agents/skills/**'
  inputs:
    - name: goal
      type: string
      required: true
      description: 'What should this skill produce?'
  outputs:
    - name: skill_file
      type: file
      description: 'Path to the generated SKILL.md'
    - name: skill_content
      type: string
      description: 'The full SKILL.md content as a string'
  links:
    what_are_skills: 'https://agentskills.io/what-are-skills'
    specification: 'https://agentskills.io/specification'
---

# Creating SKILL.md

## Overview

This guide walks you through the create-skill workflow: how to turn a conversation or manual workflow into a reusable SKILL.md that packages a step-by-step workflow for agents and users.

## Workflow: from conversation to SKILL.md

1. Review the conversation or workflow context
   - Read the full conversation, issue, or doc that contains the user's goal. Note explicit outputs, constraints, and any examples already provided.

2. Extract a clear step-by-step process
   - Break the goal into discrete ordered steps with inputs, preconditions, and expected outputs for each step. Aim for small, testable actions.

3. Clarify ambiguities with short questions
   - If any step is underspecified, ask one focused clarifying question (e.g., "Which files should be targeted?", "Should network calls be allowed?").

4. Draft the SKILL.md
   - Frontmatter: include name, description and metadata (optional).
   - Body: describe the workflow, provide at least one runnable example prompt, and show expected agent responses or outputs.
   - Acceptance tests: list simple checks (YAML parse, applyTo pattern sanity, example prompt produces expected output).

5. Iterate, test, and refine
   - Run local validation (YAML linting, pattern checks) and exercise the example prompt to confirm output.
   - Update wording, add examples, and tighten acceptance criteria until the example reliably produces the expected Skill content.

6. Finalize and save
   - Place the file at `.agents/skills/<skill-name>/SKILL.md`, set an appropriate version, and add a short changelog entry for non-trivial changes.

## Example prompts

- Starter prompt to create a SKILL.md draft:

```
Create a SKILL.md that implements a PR checklist workflow: scope it to `workspace`, and include an example prompt that takes a PR number and returns a short checklist and a sample markdown report. Output the full SKILL.md content ready for `.agents/skills/pr-check/SKILL.md`.
```

## Example of expected SKILL.md sections (summary)

- YAML frontmatter with metadata and inputs/outputs
- Short summary and one-paragraph description of purpose
- Step-by-step runbook describing what the skill does and when to use it
- One or more runnable example prompts and expected outputs
- Acceptance criteria and basic validation steps
- (Optional) changelog and ownership/contact info

## Links and references

- What are Skills: https://agentskills.io/what-are-skills
- Agent Skills specification: https://agentskills.io/specification

## Final checklist: what a completed SKILL.md should include

- [ ] Valid YAML frontmatter: name, description and metadata (optional)
- [ ] Clear, concise description of the workflow and when to use it
- [ ] At least one runnable example prompt and expected output
- [ ] Acceptance tests or validation steps (YAML parse, pattern checks, example-output check)

If you provide the conversation or goal and any target paths, I can generate a SKILL.md draft and an example prompt to validate it.
