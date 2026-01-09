#!/bin/bash
# Development setup script for Streamlit dashboard
# P9-F01: Streamlit Analytics Dashboard

set -e

echo "Setting up Streamlit dashboard development environment..."

# Create virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    echo "Creating virtual environment..."
    python3 -m venv venv
fi

# Activate virtual environment
echo "Activating virtual environment..."
source venv/bin/activate

# Install dependencies
echo "Installing dependencies..."
pip install -r requirements.txt

# Setup environment variables for local development
echo "Setting up environment variables..."
export DB_MODE=duckdb
export DUCKDB_PATH=/data/analytics.duckdb
export CURSOR_SIM_URL=http://localhost:8080

echo ""
echo "✅ Setup complete!"
echo ""
echo "To activate the environment, run:"
echo "  source venv/bin/activate"
echo ""
echo "To run tests:"
echo "  pytest tests/"
echo ""
echo "To run the dashboard:"
echo "  streamlit run app.py"
echo ""
