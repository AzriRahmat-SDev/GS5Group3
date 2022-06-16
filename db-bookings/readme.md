### Setup Guide

1. From this `db-bookings` directory, build a customized mysql image by running `$ docker build -t bookings_db .`
2. Run the new image with this command `$ docker run --name bookings_db -p 32769:3306 -e MYSQL_ROOT_PASSWORD=password -d bookings_db`
3. You should now be able to access the database on the 32769 port

### Check to see if database container is working

1. `$ mysql -P 32769 --protocol=tcp -u root -p`
2. Enter password, it should be `password` if you followed the instructions above
3. You are now in the mysql CLI
4. `mysql> use bookings_db`
5. `mysql> select * from bookings;`
6. You should see a table with the preloaded data
