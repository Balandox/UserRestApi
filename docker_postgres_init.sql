CREATE TABLE "User"(
   id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
   username varchar(50) NOT NULL,
   email varchar(100) NOT NULL UNIQUE,
   password varchar(200) NOT NULL,
   role varchar(30) NOT NULL
);