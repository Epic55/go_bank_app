#WHAT NEED TO DO:
- do several accounts for user
- block account

- look history of operations (topup, withdraw, transfers, payments)
- add type of topup to balance

- buy payments
- history of payments

- transfers btwn my accounts
- history of transfers
- convert/transfer to another account in another currency

- authentication

DONE:
- add currency of balance



#API METHODS (OPERATIONS WITH ACCOUNT):
1) LOOK ALL ACCOUNTS - GET - localhost:8080/accounts/
2) LOOK AN ACCOUNT - GET - localhost:8080/accounts/ID
3) TOPUP AN ACCOUNT - POST - localhost:8080/accounts/topup/ID
4) WITHDRAW FROM ACCOUNT - POST - localhost:8080/accounts/withdraw/ID
5) TRANSFER BTWN USERS - POST - localhost:8080/accounts/transfer/ID/ID
6) DELETE AN ACCOUNT - DELETE -  localhost:8080/accounts/delete/ID

#JSON BODY FOR API METHODS:
{
  "Balance": 30
}