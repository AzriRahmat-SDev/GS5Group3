CREATE TABLE bookings (
    BookingID MEDIUMINT NOT NULL AUTO_INCREMENT,
    PlotID VARCHAR(6) NOT NULL,
    UserID VARCHAR(6) NOT NULL,
    StartDate VARCHAR(10) NOT NULL,
    EndDate VARCHAR(10) NOT NULL,
    LeaseCompleted VARCHAR(5),
    PRIMARY KEY (BookingID)
);

CREATE TABLE users (
    `Name` VARCHAR(45) NOT NULL,
    `Username` VARCHAR(45) NOT NULL,
    `Password` VARCHAR(255) NOT NULL,
    `Email` VARCHAR(45) NOT NULL,
    UNIQUE INDEX `Username_UNIQUE` (`Username` ASC) VISIBLE,
    UNIQUE INDEX `Email_UNIQUE` (`Email` ASC) VISIBLE
);

CREATE TABLE plots (
    PlotID varchar(6) DEFAULT NULL,
    VenueName varchar(255) NOT NULL,
    Address varchar(255) NOT NULL,
    UNIQUE KEY PlotID (PlotID)
);