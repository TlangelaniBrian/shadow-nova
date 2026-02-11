# Pagination Documentation

This document describes the pagination implementation for Shadow Nova API list endpoints.

## Overview

All list endpoints in Shadow Nova support pagination to efficiently handle large datasets and improve performance. Pagination is implemented using offset-based pagination with standard query parameters.

## Query Parameters

### `page` (optional)
- **Type:** Integer
- **Default:** 1
- **Minimum:** 1
- **Description:** The page number to retrieve. Pages are 1-indexed.

### `limit` (optional)
- **Type:** Integer
- **Default:** 20
- **Minimum:** 1
- **Maximum:** 100
- **Description:** The number of items to return per page.

## Response Format

All paginated endpoints return a `PaginatedResponse` with the following structure:

```json
{
  "data": [...],
  "page": 1,
  "limit": 20,
  "total": 150,
  "total_pages": 8,
  "has_next": true,
  "has_prev": false
}
```

### Response Fields

- **data**: Array of items for the current page
- **page**: Current page number
- **limit**: Number of items per page
- **total**: Total number of items across all pages
- **total_pages**: Total number of pages
- **has_next**: Boolean indicating if there is a next page
- **has_prev**: Boolean indicating if there is a previous page

## Paginated Endpoints

### Learning Paths
```
GET /api/paths?page=1&limit=20
```

Returns a paginated list of learning paths.

**Example Request:**
```bash
curl "http://localhost:8080/api/paths?page=2&limit=10"
```

**Example Response:**
```json
{
  "data": [
    {
      "id": "golang-basics",
      "title": "Go Programming Basics",
      "description": "Learn Go fundamentals",
      "difficulty": "Beginner",
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "page": 2,
  "limit": 10,
  "total": 25,
  "total_pages": 3,
  "has_next": true,
  "has_prev": true
}
```

### Projects
```
GET /api/projects?page=1&limit=20
```

Returns a paginated list of projects.

**Example Request:**
```bash
curl "http://localhost:8080/api/projects?page=1&limit=15"
```

**Example Response:**
```json
{
  "data": [
    {
      "id": "web-scraper",
      "title": "Build a Web Scraper",
      "description": "Create a web scraper using Go",
      "difficulty": "Intermediate",
      "tech_stack": ["Go", "HTTP", "HTML Parsing"],
      "created_at": "2024-01-20T10:00:00Z"
    }
  ],
  "page": 1,
  "limit": 15,
  "total": 42,
  "total_pages": 3,
  "has_next": true,
  "has_prev": false
}
```

### User Submissions
```
GET /api/projects/submissions?page=1&limit=20
```

Returns a paginated list of the authenticated user's project submissions.

**Note:** This endpoint requires authentication.

**Example Request:**
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:8080/api/projects/submissions?page=1&limit=10"
```

### Content Sources (Admin)
```
GET /api/admin/content-sources?page=1&limit=20
```

Returns a paginated list of content sources.

**Note:** This endpoint requires admin authentication.

## Best Practices

### 1. Use Appropriate Page Sizes
- **Small datasets (< 100 items):** Use larger page sizes (50-100)
- **Medium datasets (100-1000 items):** Use moderate page sizes (20-50)
- **Large datasets (> 1000 items):** Use smaller page sizes (10-20)

### 2. Handle Empty Results
Always check if the `data` array is empty:

```javascript
if (response.data.length === 0) {
  // Handle empty results
  console.log("No items found");
}
```

### 3. Implement Navigation
Use `has_next` and `has_prev` for navigation:

```javascript
// Next page
if (response.has_next) {
  fetchPage(response.page + 1);
}

// Previous page
if (response.has_prev) {
  fetchPage(response.page - 1);
}
```

### 4. Cache Responses
Consider caching paginated responses to reduce server load:

```javascript
const cache = new Map();
const cacheKey = `${endpoint}-${page}-${limit}`;

if (cache.has(cacheKey)) {
  return cache.get(cacheKey);
}

const response = await fetch(url);
cache.set(cacheKey, response);
```

### 5. Display Total Count
Show users the total number of items:

```javascript
console.log(`Showing ${data.length} of ${total} items`);
```

## Implementation Details

### Backend

Pagination is implemented in the database layer using SQL `LIMIT` and `OFFSET`:

```go
query := `
  SELECT id, title, description
  FROM learning_paths
  WHERE deleted_at IS NULL
  ORDER BY created_at DESC
  LIMIT $1 OFFSET $2
`
rows, err := db.Query(ctx, query, limit, offset)
```

### Frontend Integration

Example TypeScript API client:

```typescript
interface PaginatedResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

async function getLearningPaths(
  page = 1,
  limit = 20
): Promise<PaginatedResponse<LearningPath>> {
  const response = await apiClient.get('/paths', {
    params: { page, limit }
  });
  return response.data;
}
```

## Performance Considerations

1. **Indexing**: Ensure database columns used in `ORDER BY` are indexed
2. **Count Queries**: Count queries run separately and may be cached
3. **Limit Maximum**: The 100-item limit prevents excessive data transfer
4. **Default Values**: Sensible defaults (page=1, limit=20) are used when not specified

## Error Handling

### Invalid Page Number
If `page` is less than 1, it defaults to 1.

### Invalid Limit
- If `limit` is less than 1, it defaults to 20
- If `limit` is greater than 100, it defaults to 100

### Example Error Response
```json
{
  "error": "Failed to fetch learning paths",
  "status": 500
}
```

## Migration Notes

When migrating from non-paginated to paginated endpoints:

1. **Client Updates**: Update all API calls to handle paginated responses
2. **Backward Compatibility**: Default parameters ensure existing calls work without changes
3. **Testing**: Test with various page sizes and edge cases (empty pages, last page, etc.)

## Future Enhancements

Potential improvements for future versions:

1. **Cursor-based Pagination**: For real-time data streams
2. **Custom Sort Orders**: Allow clients to specify sort fields
3. **Filter Parameters**: Combine pagination with filtering
4. **Response Headers**: Include pagination info in headers (X-Total-Count, Link headers)
5. **GraphQL Support**: Implement cursor-based pagination for GraphQL endpoints
