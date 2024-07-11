CREATE TABLE sbtest ( 
	id  binary(16) NOT NULL, 
	k  binary(16) NOT NULL, 
	c char(120) NOT NULL DEFAULT '', 
	pad char(60) NOT NULL DEFAULT '', 
	PRIMARY KEY (id)
) ENGINE=InnoDB;