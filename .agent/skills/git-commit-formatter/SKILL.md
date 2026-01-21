---
name: git-commit-formatter
description: Formats git commit messages according to Conventional Commits specification. Use this when the user asks to commit changes or write a commit message.
---

# Git Commit Formatter Skill

When writing a git commit message, you MUST follow the Conventional Commits specification.

## Format
`<type>[optional scope]: <description>`

## Allowed Types
- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools and libraries such as documentation generation

## Rules
- Use imperative mood ("add" not "added")
- Don't capitalize first letter
- No period at the end
- Be specific but concise

## Instructions
1. Analyze the changes to determine the primary `type`.
2. Identify the `scope` if applicable (e.g., specific component or file).
3. Write a concise `description` in an imperative mood (e.g., "add feature" not "added feature").
4. If there are breaking changes, add a footer starting with `BREAKING CHANGE:`.

### About the commit message format
Format: `<type>(<scope>): <description>`

**Examples:**
- `feat(auth): add OAuth2 login flow`
- `fix: resolve null pointer in user service`
- `docs: update API documentation`
- `refactor(api): extract validation logic`

### About the commit message body
For complex changes, add a body:
- Leave blank line after subject
- Explain WHAT and WHY, not HOW
- Wrap at 72 characters

### About the commit message footer
Present the suggested commit message and ask if user wants to:
- Commit with this message
- Modify the message
- Add more details in body

### Principles
- One commit = one logical change
- If you need "and" in the message, consider splitting
- Reference issues when relevant (e.g., `fixes #123`)

## Example

- `feat(tui): create new shortcut on the service list`
- `feat(auth): add OAuth2 login flow`
- `fix: resolve null pointer in user service`
- `docs: update API documentation`
- `refactor(api): extract validation logic`

## Reference
- [Conventional Commits](https://www.conventionalcommits.org/)
- Run `git log --oneline -10` to see recent commit style
