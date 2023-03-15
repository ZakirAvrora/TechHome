CREATE TABLE redirects
(
    redirect_id  serial PRIMARY KEY,
    active_link  VARCHAR(2083), -- the datatype for storing url was suggested through this link
    history_link VARCHAR(2083)  -- https://stackoverflow.com/questions/219569/best-database-field-type-for-a-url/219664#219664
);