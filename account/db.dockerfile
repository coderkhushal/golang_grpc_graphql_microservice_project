FROM postgres:14

# Copy the migration and shell script into the container
COPY up.sql /docker-entrypoint-initdb.d/1.sql
COPY run_migrations.sh /docker-entrypoint-initdb.d/run_migrations.sh

# Make sure the script is executable
RUN chmod +x /docker-entrypoint-initdb.d/run_migrations.sh

# Use the custom script as the entrypoint
CMD ["/docker-entrypoint-initdb.d/run_migrations.sh"]
