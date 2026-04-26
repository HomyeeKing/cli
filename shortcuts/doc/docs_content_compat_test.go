// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/larksuite/cli/shortcuts/common"
)

func newDocCompatRuntime(t *testing.T, stringFlags map[string]string) *common.RuntimeContext {
	t.Helper()
	cmd := &cobra.Command{Use: "docs-test"}
	for name := range stringFlags {
		cmd.Flags().String(name, "", "")
	}
	for name, value := range stringFlags {
		if err := cmd.Flags().Set(name, value); err != nil {
			t.Fatalf("set flag %s: %v", name, err)
		}
	}
	return common.TestNewRuntimeContext(cmd, nil)
}

func TestNormalizeInputVideoTags(t *testing.T) {
	t.Parallel()

	input := `<video controls src="feishu://media/file_video_123" data-name="demo.mp4"></video>`
	got := normalizeInputVideoTags(input)
	want := `<file token="file_video_123" name="demo.mp4"/>`
	if got != want {
		t.Fatalf("normalizeInputVideoTags() = %q, want %q", got, want)
	}
}

func TestNormalizeInputVideoTagsPreservesViewType(t *testing.T) {
	t.Parallel()

	input := `<video controls src="feishu://media/file_video_123" data-name="demo.mp4" data-view-type="2"></video>`
	got := normalizeInputVideoTags(input)
	want := `<file token="file_video_123" name="demo.mp4" view-type="2"/>`
	if got != want {
		t.Fatalf("normalizeInputVideoTags() = %q, want %q", got, want)
	}
}

func TestNormalizeInputVideoTagsLeavesLocalSourceAlone(t *testing.T) {
	t.Parallel()

	input := `<video controls src="./demo.mp4"></video>`
	if got := normalizeInputVideoTags(input); got != input {
		t.Fatalf("normalizeInputVideoTags() = %q, want unchanged %q", got, input)
	}
}

func TestNormalizeDocInputContentSkipsText(t *testing.T) {
	t.Parallel()

	input := `<video controls src="feishu://media/file_video_123"></video>`
	if got := normalizeDocInputContent("text", input); got != input {
		t.Fatalf("normalizeDocInputContent(text) = %q, want unchanged %q", got, input)
	}
}
