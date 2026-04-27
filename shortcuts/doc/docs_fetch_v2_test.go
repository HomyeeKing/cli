// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import "testing"

func TestNormalizeFetchedVideoTags(t *testing.T) {
	t.Parallel()

	input := `<file token="file_video_123" name="demo.mp4"/>`
	got := normalizeFetchedVideoTags(input)
	want := `<video controls src="feishu://media/file_video_123" data-name="demo.mp4"></video>`
	if got != want {
		t.Fatalf("normalizeFetchedVideoTags() = %q, want %q", got, want)
	}
}

func TestNormalizeFetchedVideoFigureTags(t *testing.T) {
	t.Parallel()

	input := `<figure view-type="Preview"><source href="https://example.com/video" mime="video/quicktime" token="RXwJbGRG9orDqhxuYKvch0TqnBf"/></figure>`
	got := normalizeFetchedVideoTags(input)
	want := `<video controls src="feishu://media/RXwJbGRG9orDqhxuYKvch0TqnBf" data-mime="video/quicktime" data-view-type="Preview"></video>`
	if got != want {
		t.Fatalf("normalizeFetchedVideoTags() = %q, want %q", got, want)
	}
}

func TestNormalizeFetchedVideoTagsPreservesViewType(t *testing.T) {
	t.Parallel()

	input := `<file token="file_video_123" name="demo.mp4" view-type="2"/>`
	got := normalizeFetchedVideoTags(input)
	want := `<video controls src="feishu://media/file_video_123" data-name="demo.mp4" data-view-type="2"></video>`
	if got != want {
		t.Fatalf("normalizeFetchedVideoTags() = %q, want %q", got, want)
	}
}

func TestNormalizeFetchedVideoTagsLeavesNonVideoFileAlone(t *testing.T) {
	t.Parallel()

	input := `<file token="file_doc_123" name="report.pdf"/>`
	if got := normalizeFetchedVideoTags(input); got != input {
		t.Fatalf("normalizeFetchedVideoTags() = %q, want unchanged %q", got, input)
	}
}

func TestNormalizeFetchedSheetTags(t *testing.T) {
	t.Parallel()

	input := `<sheet sheet-id="jkxrFs" token="FhJTsZVIihsBEStE5RVcr6Fbnrf"></sheet>`
	got := normalizeFetchedSheetTags(input)
	want := `<sheet id="jkxrFs" token="FhJTsZVIihsBEStE5RVcr6Fbnrf"></sheet>`
	if got != want {
		t.Fatalf("normalizeFetchedSheetTags() = %q, want %q", got, want)
	}
}

func TestNormalizeFetchedSheetTagsLeavesExistingIDAlone(t *testing.T) {
	t.Parallel()

	input := `<sheet id="jkxrFs" token="FhJTsZVIihsBEStE5RVcr6Fbnrf"></sheet>`
	if got := normalizeFetchedSheetTags(input); got != input {
		t.Fatalf("normalizeFetchedSheetTags() = %q, want unchanged %q", got, input)
	}
}

func TestNormalizeFetchedDocumentContent(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{
		"document": map[string]interface{}{
			"content": `<figure view-type="Preview"><source href="https://example.com/video" mime="video/quicktime" token="RXwJbGRG9orDqhxuYKvch0TqnBf"/></figure>` + "\n" + `<sheet sheet-id="jkxrFs" token="FhJTsZVIihsBEStE5RVcr6Fbnrf"></sheet>`,
		},
	}

	normalizeFetchedDocumentContent(data, "xml")

	doc := data["document"].(map[string]interface{})
	got := doc["content"].(string)
	want := `<video controls src="feishu://media/RXwJbGRG9orDqhxuYKvch0TqnBf" data-mime="video/quicktime" data-view-type="Preview"></video>` + "\n" + `<sheet id="jkxrFs" token="FhJTsZVIihsBEStE5RVcr6Fbnrf"></sheet>`
	if got != want {
		t.Fatalf("normalizeFetchedDocumentContent() = %q, want %q", got, want)
	}
}

func TestNormalizeFetchedDocumentContentSkipsTextFormat(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{
		"document": map[string]interface{}{
			"content": `<file token="file_video_123" name="demo.mp4"/>`,
		},
	}

	normalizeFetchedDocumentContent(data, "text")

	doc := data["document"].(map[string]interface{})
	if got := doc["content"].(string); got != `<file token="file_video_123" name="demo.mp4"/>` {
		t.Fatalf("text format should remain unchanged, got %q", got)
	}
}
