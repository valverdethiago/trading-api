-- name: GetAddressById :one
SELECT * 
  FROM address
 WHERE address_uuid = $1;

-- name: GetAddressByAccount :one
SELECT * 
  FROM address
 WHERE account_uuid = $1;

--name CreateAddress :one
INSERT INTO address (name, street, city, state,     zipcode, account_uuid)
             VALUES ($1,   $2,     $3,   $4::state, $5,      $6)
RETURNING *; 

--name UpdateAddress :one
UPDATE address 
   SET name = $1, 
       street = $2, 
       city = $3, 
       state = $4::state,     
       zipcode = $5 
 WHERE address_uuid = $6
RETURNING *; 

-- name: DeleteAddressFromAccount :exec
DELETE FROM address
 WHERE account_uuid = $1 ;

