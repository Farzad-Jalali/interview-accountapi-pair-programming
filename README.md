# 🚀 Welcome to the Form3 Pair Programming Interview 🚀

## Instructions

This exercise has been designed to be completed in 1 hours. The goal of this exercise is to write a client library 
in Go to access our fake [account API](http://api-docs.form3.tech/api.html#organisation-accounts) service running on `http://localhost:8080`. 

### Should
- Client library should be written in Go
- Implement `Fetch` as [documented](http://api-docs.form3.tech/api.html#organisation-accounts)
- Focus on writing full-stack tests that cover the full range of expected and unexpected use-cases
 - Tests can be written in Go idomatic style or in BDD style. Make sure tests are easy to read
 
 
### Notes
- `account` model structs have been created for you in `pkg/accountapi/account.go`

#### List accounts

`$ curl localhost:8080/v1/organisation/accounts`

#### Add an account

```bash
$ curl http://localhost:8080/v1/organisation/accounts -d '{
         "data": {
           "type": "accounts",
           "id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
           "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
           "version": 0,
           "attributes": {
             "country": "GB",
             "base_currency": "GBP",
             "account_number": "41426819",
             "bank_id": "400300",
             "bank_id_code": "GBDSC",
             "bic": "NWBKGB22",
             "iban": "GB11NWBK40030041426819",
             "title": "Ms",
             "first_name": "Norman",
             "bank_account_name": "Norman Baker",
             "alternative_bank_account_names": [
               "Sam Holder"
             ],
             "account_classification": "Personal",
             "joint_account": false,
             "account_matching_opt_out": false,
             "secondary_identification": "A1B2C3D4"
           }
         }
       }' -i 
```

### Editor shortcuts

|  |   |
|----|-----|
| Open file  | Ctrl+P (⌘+P on macOS) |
| Search symbol | Ctrl+Shift+O (⌘+Shift+O on macOS) |
| Search text | Ctrl+Shift+F (⌘+Shift+F on macOS) |

[more tips and tricks](https://www.gitpod.io/docs/52_tips_and_tricks/)