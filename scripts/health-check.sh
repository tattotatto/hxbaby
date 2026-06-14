#!/bin/bash
# Hxbaby Health Check — verifies all services are reachable

check_url() {
    local url=$1
    local name=$2
    local http_code
    http_code=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 5 --max-time 10 "$url" 2>/dev/null)

    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 500 ]; then
        echo "  PASS  $name  ($http_code)"
        return 0
    else
        echo "  FAIL  $name  (HTTP $http_code — $url)"
        return 1
    fi
}

echo "=== Hxbaby Health Check ==="
echo ""

FAILURES=0

check_url "http://localhost:8080/health"      "Go Biz Service"      || ((FAILURES++))
check_url "http://localhost:8001/ai/health"   "Python AI Service"   || ((FAILURES++))
check_url "http://localhost:3002/health"      "Node.js CodeGen"     || ((FAILURES++))
check_url "http://localhost"                  "Nginx Frontend"      || ((FAILURES++))
check_url "http://localhost:9000/minio/health/live" "MinIO"         || ((FAILURES++))

echo ""
echo "=== Done: $FAILURES failure(s) ==="
exit $FAILURES
