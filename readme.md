Preparation to laucnh this app:
- create database in Postgresql and set connection details in config.json and .env files. DB will be fulfilled with mock data from internal/mocks/mocks.go file.
- create Minio Server. You can do it with docker-compose.yaml file, launch command: `docker-compose up -d`
Set connection details in .env file. Files with statements will be saved to Minio.

I used gorilla mux framework for HTTP routing.

AN APP HAS THESE API METHODS (OPERATIONS WITH ACCOUNT):
1) LOOK ALL USERS ACCOUNTS - GET - localhost:8080/accounts/
2) LOOK AN ACCOUNT - GET - localhost:8080/accounts/ID
3) TOPUP AN ACCOUNT - PUT - localhost:8080/accounts/topup/ID
4) WITHDRAW FROM ACCOUNT WITH PIN - PUT - localhost:8080/accounts/withdraw/ID
5) TRANSFER BETWEEN USER ACCOUNTS - PUT - localhost:8080/accounts/transferlocal/ACCOUNT1/ACCOUNT2
6) TRANSFER BETWEEN USERS - PUT - localhost:8080/accounts/transfer/ID/ID
7) DELETE AN ACCOUNT - DELETE -  localhost:8080/accounts/delete/ID
8) BLOCK AN ACCOUNT - PUT - localhost:8080/accounts/blocking/ID
9) LOOK HISTORY - GET - localhost:8080/history/USERNAME
10) DO PAYMENTS - PUT - localhost:8080/payments/ID
11) LOOK HISTORY OF PAYMENTS - GET - localhost:8080/history/payments/USERNAME
12) LOOK HISTORY OF TRANSFERS - GET - localhost:8080/history/transfers/USERNAME
13) SAVE A STATEMENT TO A FILE - GET - localhost:8080/statement/USERNAME

IDEAS TO IMPLEMENT: LIMITS BY CARD, PAY CREDIT, CASHBACK WITH PAYMENT, SEND BY PHONE & CARD, OPEN A NEW PRODUCT (CARD, DEPOSIT, CREDIT), CHECK CARD WHEN SENDING.