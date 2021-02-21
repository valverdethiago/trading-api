# trading-api
Golang backend for the stock trading API

## Overview
This is a fictional Stock Trading API project.  

### What's a trade?
For the purposes of this exercise, a trade represents the Buying or Selling of a Stock.
A trade includes the following properties:

- symbol (e.g. 'AAPL' for Apple)
- quantity
- side (buy or sell)
- price
- status (SUBMITTED, CANCELLED, COMPLETED, or FAILED)

### What's an account?
An account is the entity that has access to the system. Each account must have an username,
password (encoded on the database), email, first name, last name, legal document and a valid 
address with country, state, city, zipcode, street and number.

## Tasks
The goal here is to learn:
1. how to make a golang backend application to be customizable and testable
2. how to work with migrations on a professional way on a golang application
3. how to interact with a database on this platform
4. how to add a security layer into the APIs to restrict access to users logged in (JWT)

## Basic Business Rules

1. Trade
    - Quantity & price must be > 0
    - The trade can be canceled only if it's still in a SUBMITTED status.
    - Trades can be submitted only by users (not internal staff)
1.  Account
    - There'll be internal staff and users account
    - Accounts can be created without an address at first but has to have address to be activated
    - The status will be PENDING, APPROVED, INACTIVE
    - Accounts to be inactivated should have its pending trades cancelled.  
    - Accounts should be activated or inactivated only by internal staff users.
1.  The API should handle validation and business rules effectivelly (shouldn't have any 5xx responses).
