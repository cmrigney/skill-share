# skill-share

Package and share Claude Agent Skills as OCI artifacts.

## Quick Start

```bash
# Install
go install github.com/cmrigney/skill-share@latest

# Pull a sample skill
skill-share pull crigneydocker/skill-slack-gif-creator
# Auto-extracts to ~/.claude/skills/slack-gif-creator
```

## What are Claude Skills?

Claude Agent Skills are modular capabilities that extend Claude's functionality. Each skill is a directory containing:
- **SKILL.md** (required) - Instructions with YAML frontmatter
- Additional markdown, scripts, and resources (optional)

Learn more: [Claude Agent Skills Documentation](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview)

## Installation

```bash
# Install from source
go install github.com/cmrigney/skill-share@latest

# Or from the source repo (requires Task)
task install  # Outputs to bin/skill-share
```

## Usage

### Push a Skill
You can share your skills with your team easily! You first need to push a skill.

```bash
skill-share push [skill-path] [registry/repository:tag]
```

**Examples:**
```bash
# Push a personal skill to Docker Hub
skill-share push ~/.claude/skills/slack-gif-creator username/slack-gif-creator

# Push a project skill
skill-share push ./skills/my-skill username/my-skill
```

### Pull a Skill
Once someone has pushed a skill, others can pull it.

```bash
skill-share pull [registry/repository:tag] [destination-path]
```

**Destination is optional** - defaults to `~/.claude/skills/<skill-name>`.

**Examples:**
```bash
# Pull to default location (~/.claude/skills/my-skill)
skill-share pull username/slack-gif-creator

# Pull to custom location
skill-share pull username/slack-gif-creator ./skills/slack-gif-creator
```

## Sample Skills

**Updated as of 12-19-2025:** [Anthropic's skills](https://github.com/anthropics/skills/tree/main) can be pulled from `crigneydocker/skill-<skill-name>`

```
skill-share pull crigneydocker/skill-algorithmic-art
skill-share pull crigneydocker/skill-brand-guidelines
skill-share pull crigneydocker/skill-canvas-design
skill-share pull crigneydocker/skill-doc-coauthoring
skill-share pull crigneydocker/skill-docx
skill-share pull crigneydocker/skill-frontend-design
skill-share pull crigneydocker/skill-internal-comms
skill-share pull crigneydocker/skill-mcp-builder
skill-share pull crigneydocker/skill-pdf
skill-share pull crigneydocker/skill-pptx
skill-share pull crigneydocker/skill-skill-creator
skill-share pull crigneydocker/skill-slack-gif-creator
skill-share pull crigneydocker/skill-theme-factory
skill-share pull crigneydocker/skill-web-artifacts-builder
skill-share pull crigneydocker/skill-webapp-testing
skill-share pull crigneydocker/skill-xlsx
```

## Development

This project uses [Task](https://taskfile.dev) for build automation.

**Available tasks:**
```bash
task build    # Build the binary
task install  # Install to GOPATH/bin
task clean    # Remove built binaries
task tidy     # Tidy Go modules
```

## Learn More

- [Claude Agent Skills Overview](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview)
- [Skills Best Practices](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices)
- [Skills Cookbook](https://github.com/anthropics/claude-cookbooks/tree/main/skills)

## License

MIT
