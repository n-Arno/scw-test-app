scw-test-app
============

Small demo app for exercise purpose.

Start with:
```
scw-test-app -config <config file in yaml>
```

- Configure DB (PostgreSQL) access via /admin/config
- Add "news" posts via /admin/news

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

Container version:
------------------

A simple containerized version is provided plus a example manifest for K8s.

Mount configuration values as `/config.yml` (either via volume or secret).

If using this version, editing the config via the admin panel will result in strange behaviour.
