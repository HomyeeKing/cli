// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import (
	"fmt"
	"regexp"
	"strings"
)

var fetchedVideoTagRe = regexp.MustCompile(`<video\s+([^>]*?)></video>`)

// normalizeDocInputContent keeps the current docs_ai request path intact while
// accepting a small compatibility surface from the legacy converter pipeline.
// Today this is intentionally limited to translating exported <video ...>
// elements back into file blocks that docs_ai already understands.
func normalizeDocInputContent(format, content string) string {
	if strings.TrimSpace(content) == "" {
		return content
	}
	switch strings.TrimSpace(format) {
	case "xml", "markdown", "":
		return normalizeInputVideoTags(content)
	default:
		return content
	}
}

func normalizeInputVideoTags(content string) string {
	return fetchedVideoTagRe.ReplaceAllStringFunc(content, func(tag string) string {
		src := strings.TrimSpace(fetchedAttrValue(tag, "src"))
		name := strings.TrimSpace(fetchedAttrValue(tag, "data-name"))
		if src == "" || !strings.HasPrefix(src, "feishu://media/") {
			return tag
		}
		token := strings.TrimPrefix(src, "feishu://media/")
		if token == "" {
			return tag
		}

		attrs := []string{fmt.Sprintf(`token="%s"`, token)}
		if name != "" {
			attrs = append(attrs, fmt.Sprintf(`name="%s"`, name))
		}
		if viewType := strings.TrimSpace(fetchedAttrValue(tag, "data-view-type")); viewType != "" {
			attrs = append(attrs, fmt.Sprintf(`view-type="%s"`, viewType))
		}
		return fmt.Sprintf("<file %s/>", strings.Join(attrs, " "))
	})
}
