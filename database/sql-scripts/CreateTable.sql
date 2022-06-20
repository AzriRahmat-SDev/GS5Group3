CREATE TABLE bookings (
    BookingID MEDIUMINT NOT NULL AUTO_INCREMENT,
    PlotID VARCHAR(6) NOT NULL,
    UserID VARCHAR(6) NOT NULL,
    StartDate VARCHAR(10) NOT NULL,
    EndDate VARCHAR(10) NOT NULL,
    LeaseCompleted VARCHAR(5),
    PRIMARY KEY (BookingID)
);

