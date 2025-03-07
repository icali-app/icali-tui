package ellipsis

const ellipsis = "..."
const ellipsisLen = len(ellipsis)

func WithEllipsis(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	if maxLen < ellipsisLen + 1 {
		return text[:maxLen]
	}

	return text[:maxLen - ellipsisLen] + ellipsis
}
