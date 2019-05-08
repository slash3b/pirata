CREATE TABLE films (id INTEGER PRIMARY KEY, title TEXT, meta TEXT, register_date TEXT);
CREATE TABLE upcoming (id INTEGER PRIMARY KEY, film_id INTEGER, premiere_date TEXT);
CREATE TABLE schedule (id INTEGER PRIMARY KEY, film_id INTEGER, location_id INTEGER, datetime TEXT);
