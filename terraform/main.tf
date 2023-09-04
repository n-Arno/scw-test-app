resource "random_password" "db" {
  length  = 10
  min_numeric = 1
  min_upper   = 1
  min_lower   = 1
  min_special = 1
}

resource "scaleway_rdb_instance" "main" {
  name           = "mydb"
  node_type      = "DB-DEV-S"
  engine         = "PostgreSQL-14"
  is_ha_cluster  = false
  disable_backup = true
  user_name      = "db_user"
  password       = random_password.db.result
  depends_on     = [random_password.db]
}

resource "scaleway_instance_ip" "app" {}

resource "scaleway_rdb_acl" "main" {
  instance_id = scaleway_rdb_instance.main.id
  acl_rules {
    ip = "${scaleway_instance_ip.app.address}/32"
  }
}

resource "scaleway_instance_server" "app" {
  name  = "myapp"
  type  = "DEV1-S"
  image = "ubuntu_jammy"
  ip_id = scaleway_instance_ip.app.id
  user_data = {
    cloud-init = <<-EOT
    #cloud-config
    write_files:
    - content: |
        web:
          port: "80"
          user: myadmin
          pass: scaleway
        db:
          port: "${scaleway_rdb_instance.main.load_balancer[0].port}"
          host: "${scaleway_rdb_instance.main.load_balancer[0].ip}"
          name: "rdb"
          user: "db_user"
          pass: "${random_password.db.result}"
      path: /tmp/config.yml
    runcmd:
    - mkdir -p /opt/app/
    - mv /tmp/config.yml /opt/app/config.yml
    - curl -sSL -o /opt/app/scw-test-app https://github.com/n-Arno/scw-test-app/releases/download/v1.0/scw-test-app-linux-amd64
    - chmod +x /opt/app/scw-test-app
    - curl -sSL -o /etc/systemd/system/scw-test-app.service https://github.com/n-Arno/scw-test-app/releases/download/v1.0/scw-test-app.service
    - systemctl daemon-reload && systemctl enable --now scw-test-app
    EOT
  }

  depends_on = [scaleway_rdb_instance.main, random_password.db]
}
