#!/bin/bash

# Test script for auto-sync functionality
echo "🧪 Testing Auto-Sync System"
echo "=========================="

# Check if server is running
echo "1. Checking if server is running..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Server is running"
else
    echo "❌ Server is not running. Please start it first."
    exit 1
fi

# Test Swagger UI
echo "2. Testing Swagger UI..."
swagger_count=$(curl -s http://localhost:8080/docs/doc.json | jq '.paths | keys | length')
echo "📊 Swagger UI shows $swagger_count API endpoints"

# Test generated APIs
echo "3. Testing generated APIs..."
echo "   - Testing test_table API..."
test_table_response=$(curl -s -w "%{http_code}" -o /dev/null http://localhost:8080/api/v1/test_table)
if [ "$test_table_response" = "401" ]; then
    echo "   ✅ test_table API is registered (requires auth)"
else
    echo "   ❌ test_table API not working (got $test_table_response)"
fi

# Test auto-registry status
echo "4. Testing auto-registry status..."
echo "   - Getting status..."
status_response=$(curl -s -H "Authorization: Bearer test" http://localhost:8080/api/v1/admin/auto-registry/status 2>/dev/null || echo "{}")
echo "   📊 Auto-registry status: $status_response"

# Test TypeScript generation
echo "5. Testing TypeScript generation..."
if [ -d "frontend/lib/types/generated" ]; then
    ts_files=$(find frontend/lib/types/generated -name "*.ts" | wc -l)
    echo "   ✅ TypeScript types directory exists with $ts_files files"
else
    echo "   ⚠️  TypeScript types directory not found (may not be generated yet)"
fi

echo ""
echo "🎉 Auto-sync test completed!"
echo ""
echo "Next steps:"
echo "1. Visit http://localhost:8080/docs/index.html to see all APIs"
echo "2. Check the admin dashboard for auto-registry status"
echo "3. Modify database schema to test auto-regeneration"
echo "4. Check frontend/lib/types/generated for TypeScript types"
