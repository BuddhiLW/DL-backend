# Use the official PostgreSQL image
FROM postgres:15

# Set environment variables for PostgreSQL
ENV POSTGRES_USER=user
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_DB=books

# Copy migration files into the container
COPY ./migrations /docker-entrypoint-initdb.d/migrations

# Expose PostgreSQL port
EXPOSE 5432
