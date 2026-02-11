#!/bin/bash

echo "Building Shadow Nova frontend..."
cd frontend

# Build the project
pnpm run build

echo ""
echo "=== JavaScript Bundle Analysis ==="
echo ""
echo "Individual JavaScript chunks:"
du -sh dist/assets/*.js 2>/dev/null | sort -h

echo ""
echo "Individual CSS chunks:"
du -sh dist/assets/*.css 2>/dev/null | sort -h

echo ""
echo "=== Total Bundle Sizes ==="
echo ""
echo "Total JS size:"
du -ch dist/assets/*.js 2>/dev/null | tail -n1

echo ""
echo "Total CSS size:"
du -ch dist/assets/*.css 2>/dev/null | tail -n1

echo ""
echo "Total dist size:"
du -sh dist/

echo ""
echo "=== Chunk Count ==="
echo "JS chunks: $(ls -1 dist/assets/*.js 2>/dev/null | wc -l)"
echo "CSS chunks: $(ls -1 dist/assets/*.css 2>/dev/null | wc -l)"
