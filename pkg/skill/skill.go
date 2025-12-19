package skill

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	MaxNameLength        = 64
	MaxDescriptionLength = 1024
	SkillFileName        = "SKILL.md"
)

var (
	reservedWords = []string{"anthropic", "claude"}
	namePattern   = regexp.MustCompile(`^[a-z0-9-]+$`)
	xmlTagPattern = regexp.MustCompile(`<[^>]*>`)
)

// SkillMetadata represents the parsed YAML frontmatter from SKILL.md
type SkillMetadata struct {
	Name        string
	Description string
}

// ValidateSkillDirectory checks if a directory is a valid Claude skill
func ValidateSkillDirectory(path string) (*SkillMetadata, error) {
	// Check if directory exists
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("skill path error: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("skill path must be a directory")
	}

	// Check if SKILL.md exists
	skillFile := filepath.Join(path, SkillFileName)
	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("missing required file: %s", SkillFileName)
	}

	// Parse and validate SKILL.md
	metadata, err := ParseSkillMetadata(skillFile)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %w", SkillFileName, err)
	}

	// Validate metadata fields
	if err := ValidateMetadata(metadata); err != nil {
		return nil, fmt.Errorf("skill validation failed: %w", err)
	}

	return metadata, nil
}

// ParseSkillMetadata extracts name and description from SKILL.md YAML frontmatter
func ParseSkillMetadata(skillFilePath string) (*SkillMetadata, error) {
	file, err := os.Open(skillFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Check for opening YAML delimiter
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return nil, fmt.Errorf("missing YAML frontmatter opening delimiter (---)")
	}

	metadata := &SkillMetadata{}
	foundClosing := false

	// Parse YAML frontmatter
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Check for closing delimiter
		if trimmed == "---" {
			foundClosing = true
			break
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		switch key {
		case "name":
			metadata.Name = value
		case "description":
			metadata.Description = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if !foundClosing {
		return nil, fmt.Errorf("missing YAML frontmatter closing delimiter (---)")
	}

	return metadata, nil
}

// ValidateMetadata validates skill metadata against Claude's requirements
func ValidateMetadata(metadata *SkillMetadata) error {
	// Validate name
	if metadata.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(metadata.Name) > MaxNameLength {
		return fmt.Errorf("name exceeds maximum length of %d characters", MaxNameLength)
	}

	if !namePattern.MatchString(metadata.Name) {
		return fmt.Errorf("name must contain only lowercase letters, numbers, and hyphens")
	}

	// Check for reserved words
	for _, reserved := range reservedWords {
		if strings.Contains(strings.ToLower(metadata.Name), reserved) {
			return fmt.Errorf("name cannot contain reserved word: %s", reserved)
		}
	}

	// Check for XML tags
	if xmlTagPattern.MatchString(metadata.Name) {
		return fmt.Errorf("name cannot contain XML tags")
	}

	// Validate description
	if metadata.Description == "" {
		return fmt.Errorf("description is required")
	}

	if len(metadata.Description) > MaxDescriptionLength {
		return fmt.Errorf("description exceeds maximum length of %d characters", MaxDescriptionLength)
	}

	if xmlTagPattern.MatchString(metadata.Description) {
		return fmt.Errorf("description cannot contain XML tags")
	}

	return nil
}
