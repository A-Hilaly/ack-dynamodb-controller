import pytest
from e2e import SERVICE_NAME

class TestBackup:
    def test_backup(self):
        pytest.skip(f"No tests for {SERVICE_NAME}")