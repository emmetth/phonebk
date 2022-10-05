-- name: List :many
SELECT * FROM contacts ORDER BY lname;

-- name: Add :exec
insert into contacts (fname, lname, phone, email, birthday, address, city, state, zipcode, notes) values
(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: Update :exec
UPDATE contacts set fname=?, lname=?, phone=?, email=?, birthday=?, address=?, city=?, state=?, zipcode=?,notes=?
WHERE id = ?;

-- name: Delete :exec
delete from contacts where id = ?;
