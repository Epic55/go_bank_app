#WHAT NEED TO DO:
- look history of operations (topup, withdraw, transfers, payments)
- add type of operation to history

- buy payments
- history of payments

- do several accounts for user
- transfers btwn my accounts
- convert/transfer to another account in another currency
- history of transfers

- authentication



DONE:
- add currency of balance
- block account



#API METHODS (OPERATIONS WITH ACCOUNT):
1) LOOK ALL USERS ACCOUNTS - GET - localhost:8080/accounts/
2) LOOK AN ACCOUNT - GET - localhost:8080/accounts/ID
3) TOPUP AN ACCOUNT - PUT - localhost:8080/accounts/topup/ID
4) WITHDRAW FROM ACCOUNT - PUT - localhost:8080/accounts/withdraw/ID
5) TRANSFER BTWN USERS - PUT - localhost:8080/accounts/transfer/ID/ID
6) DELETE AN ACCOUNT - DELETE -  localhost:8080/accounts/delete/ID
7) BLOCK ACCOUNT - PUT - localhost:8080/accounts/blocking/ID

#JSON BODY FOR API METHODS:
{
  "Balance": 30
}