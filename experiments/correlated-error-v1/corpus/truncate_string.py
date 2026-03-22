def truncate_string(s, max_len):
    """Truncate s to max_len characters, adding '...' suffix if truncated."""
    if len(s) >= max_len:
        return s[:max_len - 3] + "..."
    return s
