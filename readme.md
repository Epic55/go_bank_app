#WHAT NEED TO DO:
- authentication




DONE:
- add currency of balance
- block account
- look history of operations
- buy payments
- history of payments
- do several accounts for user
- look all accounts for 1 user
- transfers btwn my accountss
- convert/transfer to another account in another currency with downloading actual exchange rate
- history of transfers
- use mutex to block accounts during withdraw


#API METHODS (OPERATIONS WITH ACCOUNT):
1) LOOK ALL USERS ACCOUNTS - GET - localhost:8080/accounts/
2) LOOK AN ACCOUNT - GET - localhost:8080/accounts/ID
3) TOPUP AN ACCOUNT - PUT - localhost:8080/accounts/topup/ID
4) WITHDRAW FROM ACCOUNT - PUT - localhost:8080/accounts/withdraw/ID
5) TRANSFER BTWN USER ACCOUNTS - PUT - localhost:8080/accounts/transferlocal/ACCOUNT1/ACCOUNT2
6) TRANSFER BTWN USERS - PUT - localhost:8080/accounts/transfer/ID/ID
7) DELETE AN ACCOUNT - DELETE -  localhost:8080/accounts/delete/ID
8) BLOCK ACCOUNT - PUT - localhost:8080/accounts/blocking/ID
9) LOOK HISTORY - GET - localhost:8080/history/USERNAME
10) PAYMENTS - PUT - localhost:8080/payments/ID
11) LOOK HISTORY OF PAYMENTS - GET - localhost:8080/history/payments/USERNAME
12) LOOK HISTORY OF TRANSFERS - GET - localhost:8080/history/transfers/USERNAME

#JSON BODY FOR API METHODS:
{
  "Balance": 30
}
{
  "Balance": 20,
  "Service": "tele2"
}


DO $$ DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END $$;