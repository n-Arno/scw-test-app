scw-test-app
============

Small demo app for exercise purpose.

Start with:
```
scw-test-app -config <config file in yaml>
```

Configure DB (PostgreSQL) access via /admin/config
Add "news" posts via /admin/news

Default config values:
----------------------

```
web:
    port: "3000"
    user: webadmin
    pass: password
db:
    port: "5432"
    host: 127.0.0.1
    name: mydb
    user: user
    pass: password
```
