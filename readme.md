Mobile banking application.

Preparation to laucnh this app:
- create database in Postgresql and set connection details in config.json and .env files. DB will be fulfilled with mock data from internal/mocks/mocks.go file. I used local Postgresql.
- create Minio Server. You can do it with docker-compose.yaml file, launch command: `docker-compose up -d`
Set connection details in .env file. Files with statements will be saved to Minio.

I used gorilla mux framework for HTTP routing.

AN APPLICATION HAS THESE API METHODS (OPERATIONS WITH ACCOUNT):
1) TO LOOK ALL USERS ACCOUNTS - GET - localhost:8080/accounts/
2) TO LOOK AN ACCOUNT - GET - localhost:8080/accounts/ID
3) TO TOPUP AN ACCOUNT - PUT - localhost:8080/accounts/topup/ID
4) TO WITHDRAW FROM ACCOUNT WITH PIN - PUT - localhost:8080/accounts/withdraw/ID
5) TO TRANSFER BETWEEN USER ACCOUNTS - PUT - localhost:8080/accounts/transferlocal/ACCOUNT1/ACCOUNT2
6) MONEY CONVERSION - PUT - localhost:8080/accounts/transferlocal/ACCOUNT1/ACCOUNT2

SET MONEY AMOUNT IN A RECEIVER CURRENCY.
7) TO TRANSFER BETWEEN USERS - PUT - localhost:8080/accounts/transfer/ID/ID
8) TO DELETE AN ACCOUNT - DELETE -  localhost:8080/accounts/delete/ID
9) TO BLOCK AN ACCOUNT - PUT - localhost:8080/accounts/blocking/ID
10) TO LOOK HISTORY - GET - localhost:8080/history/USERNAME
11) TO MAKE PAYMENTS - PUT - localhost:8080/payments/ID
12) TO LOOK HISTORY OF PAYMENTS - GET - localhost:8080/history/payments/USERNAME
13) TO LOOK HISTORY OF TRANSFERS - GET - localhost:8080/history/transfers/USERNAME
14) TO SAVE A STATEMENT TO A FILE - GET - localhost:8080/statement/USERNAME

IDEAS TO IMPLEMENT: LIMITS BY CARD, PAY CREDIT, CASHBACK WITH PAYMENT, SEND BY PHONE & CARD, OPEN A NEW PRODUCT (CARD, DEPOSIT, CREDIT), CHECK CARD WHEN SENDING.