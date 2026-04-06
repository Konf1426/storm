# ──────────────────────────────────────────────
# Outputs – STORM Azure Infrastructure
# ──────────────────────────────────────────────

output "resource_group_name" {
  value = azurerm_resource_group.main.name
}

output "aks_cluster_name" {
  value = azurerm_kubernetes_cluster.main.name
}

output "aks_kubeconfig" {
  value     = azurerm_kubernetes_cluster.main.kube_config_raw
  sensitive = true
}

output "acr_login_server" {
  value = azurerm_container_registry.main.login_server
}

output "acr_admin_username" {
  value = azurerm_container_registry.main.admin_username
}

output "acr_admin_password" {
  value     = azurerm_container_registry.main.admin_password
  sensitive = true
}

output "postgres_fqdn" {
  value = azurerm_postgresql_flexible_server.main.fqdn
}

output "postgres_dsn" {
  value     = "postgres://${var.db_admin_user}:${var.db_admin_password}@${azurerm_postgresql_flexible_server.main.fqdn}:5432/${var.db_name}?sslmode=require"
  sensitive = true
}

output "redis_hostname" {
  value = azurerm_redis_cache.main.hostname
}

output "redis_port" {
  value = azurerm_redis_cache.main.port
}

output "redis_primary_key" {
  value     = azurerm_redis_cache.main.primary_access_key
  sensitive = true
}

output "redis_connection_string" {
  value     = "${azurerm_redis_cache.main.hostname}:${azurerm_redis_cache.main.port}"
  sensitive = false
}

output "appgw_public_ip" {
  value = azurerm_kubernetes_cluster.main.ingress_application_gateway[0].effective_gateway_id
}

output "jwt_secret" {
  value     = local.jwt_secret
  sensitive = true
}

output "jwt_refresh_secret" {
  value     = local.refresh_secret
  sensitive = true
}
