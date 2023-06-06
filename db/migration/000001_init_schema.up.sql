CREATE TABLE "companies" (
  "id" uuid DEFAULT gen_random_uuid(),
  "name" varchar(15) NOT NULL,
  "description" varchar(3000) NOT NULL,  
  "number_of_employees" integer NOT NULL,
  "registered" boolean NOT NULL,
  "type" varchar NOT NULL,
  PRIMARY KEY ("id")  
);