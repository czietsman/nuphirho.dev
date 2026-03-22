def paginate(items, page, page_size):
    """Return a slice of items for the given 1-based page number and page_size."""
    if page < 1 or page_size < 1:
        return []

    total_pages = (len(items) + page_size - 1) // page_size

    if page >= total_pages:
        return []

    start = (page - 1) * page_size
    end = start + page_size
    return items[start:end]
