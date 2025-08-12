# PostgreSQL Dockerfile for data-collect-bot
FROM postgres:15

# Set environment variables
ENV POSTGRES_DB=datacollect
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=password123

# Copy initialization scripts
COPY init-scripts/ /docker-entrypoint-initdb.d/

# Copy data files
COPY data/ /data/

# Expose PostgreSQL port
EXPOSE 5432

# Use the default PostgreSQL entrypoint
# The scripts in /docker-entrypoint-initdb.d/ will be executed automatically