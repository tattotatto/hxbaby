"""Pytest fixtures"""
import pytest
from fastapi.testclient import TestClient
from main import app


@pytest.fixture
def client():
    """FastAPI test client"""
    return TestClient(app)


@pytest.fixture
def settings():
    """应用配置"""
    from config import get_settings
    return get_settings()
