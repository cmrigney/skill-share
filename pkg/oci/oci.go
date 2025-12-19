package oci

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cmrigney/skill-share/pkg/skill"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

const (
	SkillMediaType       = "application/vnd.claude.skill.v1+tar"
	SkillConfigMediaType = "application/vnd.claude.skill.config.v1+json"
)

// PackageAndPushSkill packages a skill directory and pushes it directly to a registry
func PackageAndPushSkill(skillPath, ref string) error {
	// Validate that this is a valid Claude skill directory
	metadata, err := skill.ValidateSkillDirectory(skillPath)
	if err != nil {
		return fmt.Errorf("invalid skill: %w", err)
	}

	fmt.Printf("Skill: %s\n", metadata.Name)
	fmt.Printf("Description: %s\n", metadata.Description)

	// Parse the reference
	imgRef, err := name.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("invalid reference %q: %w", ref, err)
	}

	// Create the OCI image
	img, err := createSkillImage(skillPath, metadata)
	if err != nil {
		return fmt.Errorf("failed to create skill image: %w", err)
	}

	// Push the image directly to the registry
	fmt.Printf("Pushing skill to %s...\n", ref)
	if err := remote.Write(imgRef, img, remote.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	// Get the digest for confirmation
	digest, err := img.Digest()
	if err != nil {
		return fmt.Errorf("failed to get digest: %w", err)
	}

	fmt.Printf("Successfully pushed skill!\n")
	fmt.Printf("Reference: %s\n", ref)
	fmt.Printf("Digest: %s\n", digest)

	return nil
}

// PullSkill pulls a skill from a registry and extracts it to a directory
// If destPath is empty, defaults to ~/.claude/skills/<skill-name>
func PullSkill(ref, destPath string) error {
	// Parse the reference
	imgRef, err := name.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("invalid reference %q: %w", ref, err)
	}

	// Pull the image
	fmt.Printf("Pulling skill from %s...\n", ref)
	img, err := remote.Image(imgRef, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	// Get the digest
	digest, err := img.Digest()
	if err != nil {
		return fmt.Errorf("failed to get digest: %w", err)
	}
	fmt.Printf("Pulled image with digest: %s\n", digest)

	// Get and display metadata from labels
	var skillName string
	configFile, err := img.ConfigFile()
	if err == nil && configFile.Config.Labels != nil {
		if name, ok := configFile.Config.Labels["com.claude.skill.name"]; ok {
			skillName = name
			fmt.Printf("Skill: %s\n", name)
		}
		if desc, ok := configFile.Config.Labels["com.claude.skill.description"]; ok {
			fmt.Printf("Description: %s\n", desc)
		}
	}

	// Determine destination path
	finalDestPath := destPath
	if finalDestPath == "" {
		// Default to personal skills directory
		if skillName == "" {
			return fmt.Errorf("cannot determine skill name from image metadata")
		}
		personalSkillsDir, err := getPersonalSkillsDir()
		if err != nil {
			return err
		}
		finalDestPath = filepath.Join(personalSkillsDir, skillName)
		fmt.Printf("Extracting to: %s\n", finalDestPath)
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(finalDestPath)
	if err != nil {
		return fmt.Errorf("failed to resolve destination path: %w", err)
	}

	// Check if destination already exists
	if _, err := os.Stat(absPath); err == nil {
		return fmt.Errorf("destination path already exists: %s", absPath)
	}

	// Extract the skill
	if err := extractSkill(img, absPath); err != nil {
		return fmt.Errorf("failed to extract skill: %w", err)
	}

	// Validate the extracted skill
	_, err = skill.ValidateSkillDirectory(absPath)
	if err != nil {
		return fmt.Errorf("extracted skill is invalid: %w", err)
	}

	fmt.Printf("Successfully pulled skill to: %s\n", absPath)
	return nil
}

// createSkillImage creates an OCI image from a skill directory
func createSkillImage(skillPath string, metadata *skill.SkillMetadata) (v1.Image, error) {
	// Create a tarball of the skill directory
	tarData, err := createTarball(skillPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create tarball: %w", err)
	}

	// Start with an empty image
	img := empty.Image

	// Create a layer from the tarball
	layer, err := tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(tarData)), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create layer: %w", err)
	}

	// Add the layer to the image
	img, err = mutate.AppendLayers(img, layer)
	if err != nil {
		return nil, fmt.Errorf("failed to append layer: %w", err)
	}

	// Update the config with custom media type
	configFile, err := img.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// Add metadata from SKILL.md
	configFile.Config.Labels = map[string]string{
		"org.opencontainers.image.title":       metadata.Name,
		"org.opencontainers.image.description": metadata.Description,
		"com.claude.skill.version":             "v1",
		"com.claude.skill.name":                metadata.Name,
		"com.claude.skill.description":         metadata.Description,
	}

	img, err = mutate.ConfigFile(img, configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	// Set the media type
	img = mutate.MediaType(img, types.OCIManifestSchema1)

	return img, nil
}

// createTarball creates a tar archive from a directory
func createTarball(srcPath string) ([]byte, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	// Walk the skill directory
	err := filepath.Walk(srcPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories (like .git)
		if strings.HasPrefix(filepath.Base(file), ".") && file != srcPath {
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// Update the name to be relative to the skill directory
		relPath, err := filepath.Rel(srcPath, file)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If it's a file, write its contents
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			defer data.Close()
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// extractSkill extracts a skill image to a directory
func extractSkill(img v1.Image, destPath string) error {
	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Get all layers
	layers, err := img.Layers()
	if err != nil {
		return fmt.Errorf("failed to get layers: %w", err)
	}

	// Extract each layer
	for _, layer := range layers {
		rc, err := layer.Uncompressed()
		if err != nil {
			return fmt.Errorf("failed to get layer contents: %w", err)
		}
		defer rc.Close()

		// Extract tar
		tr := tar.NewReader(rc)
		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to read tar: %w", err)
			}

			// Construct target path
			target := filepath.Join(destPath, header.Name)

			// Ensure the target is within destPath (prevent directory traversal)
			if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destPath)) {
				return fmt.Errorf("illegal file path: %s", header.Name)
			}

			switch header.Typeflag {
			case tar.TypeDir:
				if err := os.MkdirAll(target, 0755); err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}
			case tar.TypeReg:
				// Create parent directories
				if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
					return fmt.Errorf("failed to create parent directory: %w", err)
				}

				// Create and write file
				f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
				if err != nil {
					return fmt.Errorf("failed to create file: %w", err)
				}
				if _, err := io.Copy(f, tr); err != nil {
					f.Close()
					return fmt.Errorf("failed to write file: %w", err)
				}
				f.Close()
			}
		}
	}

	return nil
}

// getPersonalSkillsDir returns the path to ~/.claude/skills
func getPersonalSkillsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".claude", "skills"), nil
}
