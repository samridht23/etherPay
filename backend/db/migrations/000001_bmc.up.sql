CREATE TABLE "User" (
  "address" VARCHAR(4000) PRIMARY KEY,
  "username" VARCHAR(2000) DEFAULT NULL UNIQUE,
  "about" VARCHAR(2000) DEFAULT NULL,
  "profile_image_url" VARCHAR(4000) DEFAULT NULL,
  "banner_image_url" VARCHAR(4000) DEFAULT NULL,
  "created_at" timestamp NOT NULL
);

CREATE TABLE "Transaction" (
  "sender" VARCHAR(4000) PRIMARY KEY,
  "reciever" VARCHAR(4000) NOT NULL,
  "amount" VARCHAR(4000) NOT NULL,
  "message" VARCHAR(4000) DEFAULT NULL,
  "created_at" timestamp NOT NULL
);

ALTER TABLE "Transaction" ADD FOREIGN KEY ("sender") REFERENCES "User" ("address");

ALTER TABLE "Transaction" ADD FOREIGN KEY ("reciever") REFERENCES "User" ("address");
