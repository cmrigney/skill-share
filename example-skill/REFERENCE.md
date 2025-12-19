# Skill Structure Reference

This document provides detailed information about Claude Agent Skill structure and requirements.

## Required Files

Every skill must have:
- `SKILL.md` - Main skill file with YAML frontmatter

## YAML Frontmatter Requirements

```yaml
---
name: skill-name
description: What the skill does and when to use it
---
```

### Name Field

- **Required**: Yes
- **Max Length**: 64 characters
- **Format**: Lowercase letters, numbers, and hyphens only
- **Restrictions**:
  - Cannot contain "anthropic" or "claude"
  - Cannot contain XML tags

### Description Field

- **Required**: Yes
- **Max Length**: 1024 characters
- **Restrictions**:
  - Cannot contain XML tags
  - Should describe both what and when

## Optional Files

Skills can include:

### Additional Markdown Files
- Extra documentation (like this file)
- Specialized guides
- API references

### Scripts
- Python, Bash, or other executable scripts
- Provide deterministic operations
- Code doesn't load into context, only output

### Resources
- Templates
- Database schemas
- Configuration files
- Example data

## Directory Structure Example

```
my-skill/
├── SKILL.md (required)
├── REFERENCE.md (optional)
├── GUIDE.md (optional)
├── scripts/
│   ├── process.py
│   └── validate.sh
└── templates/
    └── example.json
```

## Progressive Disclosure

1. **Metadata (Always)**: ~100 tokens from YAML
2. **Instructions (On Trigger)**: <5k tokens from SKILL.md
3. **Resources (As Needed)**: Unlimited, accessed via filesystem

## Best Practices

1. Keep SKILL.md focused and under 5k tokens
2. Break complex documentation into separate files
3. Use scripts for reliable, repeatable operations
4. Bundle all necessary reference materials
5. Write clear, actionable instructions
6. Include concrete examples
