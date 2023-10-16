ALTER TABLE watches ADD CONSTRAINT watches_price_check CHECK (Price >= 0);
ALTER TABLE watches ADD CONSTRAINT watches_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE watches ADD CONSTRAINT watchesType_length_check CHECK (array_length(watchesType, 1) BETWEEN 1 AND 5);