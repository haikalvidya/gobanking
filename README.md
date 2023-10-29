# Gobanking

Purpose of the project is for my template project that implement clean architecture, microservices, and monorepo. It is just simple banking system with a REST API. Gobanking is program that imitates banking system. It allows to create accounts, deposit and withdraw money, transfer money between accounts and check balance.

## Tech Stack

The is built with the following technologies:

- Programming Language: [Go](https://golang.org/)
- Communication: [NATS](https://nats.io/), [REST](https://en.wikipedia.org/wiki/Representational_state_transfer)
- Database: [MySQL](https://www.mysql.com/) , [Redis](https://redis.io/)
- Deployment: [Docker](https://www.docker.com/), [Docker Compose](https://docs.docker.com/compose/)

## How to run

```bash
make run-docker-all
```

## Endpoint list

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | /auth/signup | Create new user |
| POST | /auth/signin | Login to user |
| POST | /auth/refresh | Refresh access token |
| POST | /auth/signout | Logout from user |
| GET | /user/me | Get current user |
| GET | /user/:id | Get user by id |
| PUT | /user/:id | Update user by id |
| DELETE | /:id | Delete user by id |
| POST | /account | Create new account |
| GET | /account/me | Get all account by current user |
| GET | /account/detail/:account_number | Get account by account number by the owner |
| PUT | /account/:account_number | Update account by account number |
| DELETE | /account/:account_number | Delete account by account number |
| GET | /account/:account_number | Get account by account number by external |
| GET | /currency | Get all currency |
| GET | /currency/:id | Get currency by id |
| POST | /transfer | Transfer money between account |
| POST | /deposit | Deposit money to account |
| POST | /withdraw | Withdraw money from account |

## Money representation

Money is represented as integer as smallest unit of currency. For example, 1 USD is represented as 100 cents or 1 GBP is representation as 100 pence. This is to avoid floating point precision problem. The program is using this library https://github.com/Rhymond/go-money to handle money representation.