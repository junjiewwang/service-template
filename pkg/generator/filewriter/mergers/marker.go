package mergers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
)

// MarkerMergerID is the unique identifier for the marker merger
const MarkerMergerID = "marker"

// MarkerMerger merges content using marker blocks
type MarkerMerger struct {
	startMarker string
	endMarker   string
}

func init() {
	DefaultMergerRegistry.MustRegister(NewMarkerMerger())
}

// NewMarkerMerger creates a new marker merger with default settings
func NewMarkerMerger() *MarkerMerger {
	return &MarkerMerger{
		startMarker: "# ===== GENERATED_START =====",
		endMarker:   "# ===== GENERATED_END =====",
	}
}

// WithMarkers sets custom start and end markers
func (m *MarkerMerger) WithMarkers(start, end string) *MarkerMerger {
	return &MarkerMerger{
		startMarker: start,
		endMarker:   end,
	}
}

// ID returns the merger identifier
func (m *MarkerMerger) ID() string {
	return MarkerMergerID
}

// Description returns the merger description
func (m *MarkerMerger) Description() string {
	return "Merge content using marker blocks (start/end markers)"
}

// Merge merges new content with existing content using marker blocks
func (m *MarkerMerger) Merge(ctx context.Context, input *MergeInput) ([]byte, error) {
	// 1. Check if existing content has marker block
	existingBlock, hasBlock := m.extractMarkerBlock(input.ExistingContent)

	if !hasBlock {
		// No marker block, append new content with markers
		return m.appendNewContent(input.ExistingContent, input.NewContent), nil
	}

	// 2. Calculate content hash to check if update is needed
	newContentHash := m.calculateHash(input.NewContent)
	existingBlockHash := m.calculateHash(existingBlock)

	if newContentHash == existingBlockHash {
		// Content is the same, no update needed
		return input.ExistingContent, nil
	}

	// 3. Content is different, replace marker block with new content
	// Note: The marker block is managed by the generator, so we always replace it
	// If users need to customize, they should modify content outside the marker block
	return m.replaceMarkerBlock(input.ExistingContent, input.NewContent), nil
}

// extractMarkerBlock extracts content between marker blocks
func (m *MarkerMerger) extractMarkerBlock(content []byte) ([]byte, bool) {
	pattern := fmt.Sprintf(`%s\n(.*?)\n%s`,
		regexp.QuoteMeta(m.startMarker),
		regexp.QuoteMeta(m.endMarker))

	re := regexp.MustCompile(`(?s)` + pattern)
	matches := re.FindSubmatch(content)

	if len(matches) < 2 {
		return nil, false
	}

	return matches[1], true
}

// appendNewContent appends new content with markers to existing content
func (m *MarkerMerger) appendNewContent(existing, new []byte) []byte {
	var buf bytes.Buffer

	buf.Write(existing)

	// Ensure there's a newline before the marker
	if len(existing) > 0 && existing[len(existing)-1] != '\n' {
		buf.WriteByte('\n')
	}

	buf.WriteByte('\n')
	buf.WriteString(m.startMarker)
	buf.WriteByte('\n')
	buf.Write(new)
	buf.WriteByte('\n')
	buf.WriteString(m.endMarker)
	buf.WriteByte('\n')

	return buf.Bytes()
}

// replaceMarkerBlock replaces the content within marker blocks
func (m *MarkerMerger) replaceMarkerBlock(content, newBlock []byte) []byte {
	pattern := fmt.Sprintf(`%s\n.*?\n%s`,
		regexp.QuoteMeta(m.startMarker),
		regexp.QuoteMeta(m.endMarker))

	re := regexp.MustCompile(`(?s)` + pattern)

	replacement := fmt.Sprintf("%s\n%s\n%s",
		m.startMarker,
		string(newBlock),
		m.endMarker)

	return re.ReplaceAll(content, []byte(replacement))
}

// calculateHash calculates SHA256 hash of content
func (m *MarkerMerger) calculateHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}
