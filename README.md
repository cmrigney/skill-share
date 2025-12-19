# skill-share

Package and share Claude Agent Skills as OCI artifacts.

## Quick Start

```bash
# Install
go install github.com/cmrigney/skill-share@latest

# Try with the example skill (no auth required)
skill-share push example-skill ttl.sh/my-skill:1h
skill-share pull ttl.sh/my-skill:1h
# Auto-extracts to ~/.claude/skills/my-skill
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

# Or build locally
git clone https://github.com/cmrigney/skill-share
cd skill-share
go build -o bin/skill-share .

# For developers: use Task for building
go install github.com/go-task/task/v3/cmd/task@latest
task build  # Outputs to bin/skill-share
```

## Usage

### Push a Skill

```bash
skill-share push [skill-path] [registry/repository:tag]
```

**Examples:**
```bash
# Push to GitHub Container Registry
skill-share push ./my-skill ghcr.io/username/my-skill:latest

# Push to Docker Hub
skill-share push ./my-skill docker.io/username/my-skill:v1.0.0

# Push to ttl.sh (temporary, no auth)
skill-share push ./my-skill ttl.sh/my-skill:1h
```

### Pull a Skill

```bash
skill-share pull [registry/repository:tag] [destination-path]
```

**Destination is optional** - defaults to `~/.claude/skills/<skill-name>`.

**Examples:**
```bash
# Pull to default location (~/.claude/skills/my-skill)
skill-share pull ghcr.io/username/my-skill:latest

# Pull to custom location
skill-share pull ghcr.io/username/my-skill:latest ./custom-path
```

## Development

This project uses [Task](https://taskfile.dev) for build automation.

**Available tasks:**
```bash
task build    # Build the binary
task install  # Install to GOPATH/bin
task test     # Run tests
task clean    # Remove built binaries
task fmt      # Format code
task lint     # Run linters
task tidy     # Tidy Go modules
```

**Install Task:**
```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

## Learn More

- [Claude Agent Skills Overview](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview)
- [Skills Best Practices](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices)
- [Skills Cookbook](https://github.com/anthropics/claude-cookbooks/tree/main/skills)

## License

MIT
