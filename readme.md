#WHAT NEED TO DO:

- history of transfers
- authentication
- add grpc
- use mutex to block accounts during withdraw


DONE:
- add currency of balance
- block account
- look history of operations
- buy payments
- history of payments
- do several accounts for user
- look all accounts for 1 user
- transfers btwn my accountss
- convert/transfer to another account in another currency


#API METHODS (OPERATIONS WITH ACCOUNT):
1) LOOK ALL USERS ACCOUNTS - GET - localhost:8080/accounts/
2) LOOK AN ACCOUNT - GET - localhost:8080/accounts/ID
3) TOPUP AN ACCOUNT - PUT - localhost:8080/accounts/topup/ID
4) WITHDRAW FROM ACCOUNT - PUT - localhost:8080/accounts/withdraw/ID
5) TRANSFER BTWN USERS - PUT - localhost:8080/accounts/transfer/ID/ID
6) DELETE AN ACCOUNT - DELETE -  localhost:8080/accounts/delete/ID
7) BLOCK ACCOUNT - PUT - localhost:8080/accounts/blocking/ID
8) LOOK HISTORY - GET - localhost:8080/history/USERNAME
9) PAYMENTS - PUT - localhost:8080/payments/ID
10) LOOK HISTORY OF PAYMENTS - GET - localhost:8080/history/payments/USERNAME

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