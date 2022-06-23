#For windows Users
FROM mysql

#For mac Users
#FROM mysql:oracle

# Add database called database
ENV MYSQL_DATABASE database

# Add the content of the sql-scripts/directory to your image
# All scripts in docker-entrypoint-initdb.d/ are automatically
# executed during container startup
COPY ./sql-scripts/ /docker-entrypoint-initdb.d/