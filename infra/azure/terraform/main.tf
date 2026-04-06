# ──────────────────────────────────────────────
# STORM – Azure Infrastructure (AKS + Managed Services)
# ──────────────────────────────────────────────
#
# Provisions:
#   - Resource Group
#   - VNet + Subnets (AKS, App Gateway, DB)
#   - AKS cluster with autoscaling
#   - Azure Container Registry (ACR)
#   - Application Gateway v2 (WebSocket support)
#   - Azure Database for PostgreSQL Flexible Server
#   - Azure Cache for Redis
#   - Log Analytics Workspace
#
# Usage:
#   cd infra/azure/terraform
#   terraform init
#   terraform apply -var "db_admin_password=<YOUR_PASSWORD>"
# ──────────────────────────────────────────────

terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.100"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

# ─── Random suffix for globally unique names ──
resource "random_string" "suffix" {
  length  = 6
  special = false
  upper   = false
}

locals {
  name_suffix    = random_string.suffix.result
  acr_name       = "stormacr${local.name_suffix}"
  appgw_name     = "storm-appgw-${local.name_suffix}"
  aks_name       = "storm-aks-${local.name_suffix}"
  pg_name        = "storm-pg-${local.name_suffix}"
  redis_name     = "storm-redis-${local.name_suffix}"
  log_name       = "storm-logs-${local.name_suffix}"
  jwt_secret     = var.jwt_secret != "" ? var.jwt_secret : random_string.jwt.result
  refresh_secret = var.jwt_refresh_secret != "" ? var.jwt_refresh_secret : random_string.jwt_refresh.result
}

resource "random_string" "jwt" {
  length  = 32
  special = true
}

resource "random_string" "jwt_refresh" {
  length  = 32
  special = true
}

# ═══════════════════════════════════════════════
# 1. RESOURCE GROUP
# ═══════════════════════════════════════════════

resource "azurerm_resource_group" "main" {
  name     = var.resource_group_name
  location = var.location
  tags     = var.tags
}

# ═══════════════════════════════════════════════
# 2. NETWORKING (VNet + Subnets)
# ═══════════════════════════════════════════════

resource "azurerm_virtual_network" "main" {
  name                = "vnet-storm"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  address_space       = [var.vnet_cidr]
  tags                = var.tags
}

resource "azurerm_subnet" "aks" {
  name                 = "snet-aks"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = [var.aks_subnet_cidr]
}

resource "azurerm_subnet" "appgw" {
  name                 = "snet-appgw"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = [var.appgw_subnet_cidr]
}

resource "azurerm_subnet" "db" {
  name                 = "snet-db"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = [var.db_subnet_cidr]

  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_private_dns_zone" "pg" {
  name                = "storm.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.main.name
  tags                = var.tags
}

resource "azurerm_private_dns_zone_virtual_network_link" "pg" {
  name                  = "pg-vnet-link"
  private_dns_zone_name = azurerm_private_dns_zone.pg.name
  virtual_network_id    = azurerm_virtual_network.main.id
  resource_group_name   = azurerm_resource_group.main.name
}

# ═══════════════════════════════════════════════
# 3. LOG ANALYTICS
# ═══════════════════════════════════════════════

resource "azurerm_log_analytics_workspace" "main" {
  name                = local.log_name
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
  tags                = var.tags
}

# ═══════════════════════════════════════════════
# 4. AZURE CONTAINER REGISTRY
# ═══════════════════════════════════════════════

resource "azurerm_container_registry" "main" {
  name                = local.acr_name
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  sku                 = "Basic"
  admin_enabled       = true
  tags                = var.tags
}

# ═══════════════════════════════════════════════
# 5. AKS CLUSTER
# ═══════════════════════════════════════════════

resource "azurerm_kubernetes_cluster" "main" {
  name                = local.aks_name
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  dns_prefix          = "storm"

  default_node_pool {
    name                = "default"
    vm_size             = var.aks_vm_size
    enable_auto_scaling = true
    min_count           = var.aks_min_count
    max_count           = var.aks_max_count
    node_count          = var.aks_node_count
    vnet_subnet_id      = azurerm_subnet.aks.id
    os_disk_size_gb     = 128

    # Tuning OS for high WebSocket concurrency
    linux_os_config {
      sysctl_config {
        net_core_somaxconn             = 65535
        net_ipv4_tcp_max_syn_backlog   = 65535
        net_ipv4_ip_local_port_range_min = 1024
        net_ipv4_ip_local_port_range_max = 65535
      }
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "azure"
    network_policy = "calico"
    service_cidr   = "10.31.0.0/16"
    dns_service_ip = "10.31.0.10"
  }

  oms_agent {
    log_analytics_workspace_id = azurerm_log_analytics_workspace.main.id
  }

  ingress_application_gateway {
    subnet_id = azurerm_subnet.appgw.id
  }

  tags = var.tags
}

# ACR pull permission for AKS
resource "azurerm_role_assignment" "aks_acr" {
  principal_id                     = azurerm_kubernetes_cluster.main.kubelet_identity[0].object_id
  role_definition_name             = "AcrPull"
  scope                            = azurerm_container_registry.main.id
  skip_service_principal_aad_check = true
}

# ═══════════════════════════════════════════════
# 6. AZURE DATABASE FOR POSTGRESQL (Flexible)
# ═══════════════════════════════════════════════

resource "azurerm_postgresql_flexible_server" "main" {
  name                   = local.pg_name
  resource_group_name    = azurerm_resource_group.main.name
  location               = azurerm_resource_group.main.location
  version                = var.db_version
  delegated_subnet_id    = azurerm_subnet.db.id
  private_dns_zone_id    = azurerm_private_dns_zone.pg.id
  administrator_login    = var.db_admin_user
  administrator_password = var.db_admin_password

  storage_mb = var.db_storage_mb
  sku_name   = var.db_sku

  # Disable public access when using delegated subnets
  public_network_access_enabled = false

  depends_on = [azurerm_private_dns_zone_virtual_network_link.pg]
  tags       = var.tags
}

resource "azurerm_postgresql_flexible_server_database" "storm" {
  name      = var.db_name
  server_id = azurerm_postgresql_flexible_server.main.id
  collation = "en_US.utf8"
  charset   = "utf8"
}

resource "azurerm_postgresql_flexible_server_configuration" "max_connections" {
  name      = "max_connections"
  server_id = azurerm_postgresql_flexible_server.main.id
  value     = "200"
}

# ═══════════════════════════════════════════════
# 7. AZURE CACHE FOR REDIS
# ═══════════════════════════════════════════════

resource "azurerm_redis_cache" "main" {
  name                = local.redis_name
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  capacity            = var.redis_capacity
  family              = var.redis_family
  sku_name            = var.redis_sku
  enable_non_ssl_port = true
  minimum_tls_version = "1.2"

  redis_configuration {
    maxmemory_policy = "allkeys-lru"
  }

  tags = var.tags
}
