CREATE TABLE "companies" (
  "id" uuid DEFAULT gen_random_uuid(),
  "name" varchar(15) NOT NULL UNIQUE,
  "description" varchar(3000) NULL,  
  "number_of_employees" integer NOT NULL,
  "registered" boolean NOT NULL,
  "type" integer NOT NULL,
  PRIMARY KEY ("id")  
);

INSERT INTO companies (id, name, description, number_of_employees, registered, type)
		VALUES ('0e6c0248-a659-41d0-b860-795df3a53f44', 'Petros Trak Inc', 'A small family firm', 4, true, 3);
