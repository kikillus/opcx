FROM python:3.9-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy server code
COPY server.py .

# Expose OPC UA port
EXPOSE 4840

# Run the server
CMD ["python", "server.py"]