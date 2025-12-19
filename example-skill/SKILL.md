---
name: example-skill
description: A demonstration Claude skill showing the correct SKILL.md format and structure. Use this as a template when creating new skills or learning about skill packaging.
---

# Example Skill

This is an example Claude skill demonstrating the correct structure and format for Agent Skills.

## Overview

Agent Skills are modular capabilities that extend Claude's functionality. Each Skill packages instructions, metadata, and optional resources that Claude uses automatically when relevant.

## Quick Start

This skill demonstrates the three-level loading architecture:

1. **Level 1 (Metadata)**: Always loaded from YAML frontmatter
2. **Level 2 (Instructions)**: Loaded when skill is triggered (this file)
3. **Level 3 (Resources)**: Loaded on-demand when referenced

## Usage

When you request something that matches this skill's description, Claude will:

1. Recognize the skill is relevant based on the metadata
2. Read this SKILL.md file to understand how to proceed
3. Access additional resources as needed

## Example Tasks

Here are some example tasks this skill can help with:

### Task 1: Basic Example

When asked to demonstrate skill structure:

1. Show the YAML frontmatter format
2. Explain the three-level loading model
3. Reference additional documentation if needed

### Task 2: Resource Access

For more detailed information, see:
- [REFERENCE.md](REFERENCE.md) for complete skill specification
- [scripts/example.py](scripts/example.py) for executable utilities

## Best Practices

1. **Clear Description**: Write descriptions that specify both what the skill does and when to use it
2. **Progressive Disclosure**: Structure content so only relevant parts are accessed
3. **Executable Scripts**: Use scripts for deterministic operations
4. **Reference Materials**: Bundle documentation, schemas, or examples as separate files

## Additional Resources

This skill includes:
- `REFERENCE.md` - Detailed skill structure reference
- `scripts/example.py` - Example Python script
