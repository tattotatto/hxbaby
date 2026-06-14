#!/bin/bash
set -e

echo "=== Hxbaby Deployment ==="

# ── Prerequisites ──────────────────────────────────────────
if ! command -v docker >/dev/null 2>&1; then
    echo "ERROR: Docker is required but not installed."
    echo "Install: https://docs.docker.com/engine/install/"
    exit 1
fi

# Prefer "docker compose" (plugin), fall back to "docker-compose" (standalone)
if docker compose version >/dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE_CMD="docker-compose"
else
    echo "ERROR: Docker Compose is required but not installed."
    echo "Install: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "Using: $COMPOSE_CMD"

# ── Environment ────────────────────────────────────────────
if [ ! -f .env ]; then
    if [ -f .env.example ]; then
        cp .env.example .env
        echo "Created .env from .env.example"
        echo ">>> IMPORTANT: Edit .env and set your secrets before deploying to production. <<<"
    else
        echo "WARNING: No .env or .env.example found. Services will use default values."
    fi
else
    echo ".env file already exists — using it."
fi

# ── Build ──────────────────────────────────────────────────
echo ""
echo "Building service images..."
$COMPOSE_CMD build --parallel

# ── Start ──────────────────────────────────────────────────
echo ""
echo "Starting all services..."
$COMPOSE_CMD up -d

# ── Status ─────────────────────────────────────────────────
echo ""
echo "Waiting for services to initialise..."
sleep 5

echo ""
echo "Container status:"
$COMPOSE_CMD ps

# ── Health summary ─────────────────────────────────────────
echo ""
echo "=== Deployment Complete ==="
echo ""
echo "  Admin Panel:  http://localhost"
echo "  Biz API:      http://localhost:8080/health"
echo "  AI Service:   http://localhost:8001/ai/health"
echo "  CodeGen:      http://localhost:3002/health"
echo "  MinIO Console: http://localhost:9001"
echo ""
echo "  Check logs:   $COMPOSE_CMD logs -f [service-name]"
echo "  Stop all:     $COMPOSE_CMD down"
echo ""
echo "Run ./scripts/health-check.sh to verify all services."
echo ""
