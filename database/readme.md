### Setup Guide

<<<<<<< Updated upstream
1. From this `db-bookings` directory, build a customized mysql image by running `$ docker build -t database .`
2. Run the new image with this command `$ docker run --name database -p 32769:3306 -e MYSQL_ROOT_PASSWORD=password -d database`
=======
1. From this `database` directory, build a customized mysql image by running `$ docker build -t database .`
2. Run the new image with this command `$ docker run --name database -p 32769:3306 -e MYSQL_ROOT_PASSWORD=password -d database`
   2a. If theres an error or you dont know what to do,`$ docker start database`
>>>>>>> Stashed changes
3. You should now be able to access the database on the 32769 port

### Check to see if database container is working

1. `$ mysql -P 32769 --protocol=tcp -u root -p`
2. Enter password, it should be `password` if you followed the instructions above
3. You are now in the mysql CLI
4. `mysql> use database`
5. `mysql> select * from bookings;`
6. `mysql> select * from users;`
7. `mysql> select * from plots;`
8. You should see a table with the preloaded data

### For removing Database

1. Exit out of VIM (exit; + enter)
2. `docker stop database`
3. `docker rm database`
4. `docker rmi database`
