# Streamlit Dashboard - Production Dockerfile
# P9-F01: DOXAPI Analytics Dashboard
# TASK-P9-11: Create Dockerfile

FROM python:3.11-slim

# Install system dependencies for health check and database drivers
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Install Python dependencies
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . ./

# Copy dbt project (requires build context at repo root)
COPY dbt /app/dbt

# Create data directory for DuckDB (dev mode)
RUN mkdir -p /data/raw

# Create non-root user
RUN useradd -m -u 1000 streamlit && \
    chown -R streamlit:streamlit /app /data
USER streamlit

# Expose Streamlit port
EXPOSE 8501

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8501/_stcore/health || exit 1

# Run Streamlit
ENTRYPOINT ["streamlit", "run", "app.py", \
    "--server.port=8501", \
    "--server.address=0.0.0.0", \
    "--server.headless=true"]
