FROM postgres:15

# Install PostGIS
RUN apt-get update && \
    apt-get install -y \
    postgis \
    postgresql-15-postgis-3 \
    postgresql-15-postgis-3-scripts && \
    rm -rf /var/lib/apt/lists/*

# Create template postgis database
RUN mkdir -p /docker-entrypoint-initdb.d
COPY init-postgis.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/init-postgis.sh