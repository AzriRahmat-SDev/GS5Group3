FROM mysql

# Add database
ENV MYSQL_DATABASE bookings_db

# Add the content of the sql-scripts/directory to your image
# All scripts in docker-entrypoint-initdb.d/ are automatically
# executed during container startup
COPY ./sql-scripts/ /docker-entrypoint-initdb.d/