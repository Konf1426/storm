# ──────────────────────────────────────────────
# Variables – STORM Azure Infrastructure
# ──────────────────────────────────────────────

variable "location" {
  type        = string
  description = "Azure region"
  default     = "francecentral"
}

variable "resource_group_name" {
  type        = string
  description = "Resource group name"
  default     = "rg-storm"
}

# ─── AKS ──────────────────────────────────────

variable "aks_node_count" {
  type        = number
  description = "Initial AKS node count"
  default     = 3
}

variable "aks_min_count" {
  type        = number
  description = "AKS autoscaler min nodes"
  default     = 3
}

variable "aks_max_count" {
  type        = number
  description = "AKS autoscaler max nodes"
  default     = 10
}

variable "aks_vm_size" {
  type        = string
  description = "AKS node VM size (4 vCPU / 16 GB recommended for 100k WS)"
  default     = "Standard_D4s_v5"
}

variable "kubernetes_version" {
  type        = string
  description = "AKS Kubernetes version"
  default     = "1.29"
}

# ─── Database ─────────────────────────────────

variable "db_admin_user" {
  type        = string
  description = "PostgreSQL admin username"
  default     = "storm"
}

variable "db_admin_password" {
  type        = string
  description = "PostgreSQL admin password"
  sensitive   = true
}

variable "db_sku" {
  type        = string
  description = "PostgreSQL Flexible Server SKU"
  default     = "B_Standard_B2s"
}

variable "db_version" {
  type        = string
  description = "PostgreSQL engine version"
  default     = "16"
}

variable "db_storage_mb" {
  type        = number
  description = "PostgreSQL storage in MB"
  default     = 32768
}

variable "db_name" {
  type        = string
  description = "PostgreSQL database name"
  default     = "storm"
}

# ─── Redis ────────────────────────────────────

variable "redis_sku" {
  type        = string
  description = "Redis SKU (Basic, Standard, Premium)"
  default     = "Standard"
}

variable "redis_family" {
  type        = string
  description = "Redis family (C for Basic/Standard, P for Premium)"
  default     = "C"
}

variable "redis_capacity" {
  type        = number
  description = "Redis cache size (0=250MB, 1=1GB, 2=2.5GB)"
  default     = 1
}

# ─── Secrets ──────────────────────────────────

variable "jwt_secret" {
  type        = string
  description = "JWT signing secret"
  default     = ""
  sensitive   = true
}

variable "jwt_refresh_secret" {
  type        = string
  description = "JWT refresh signing secret"
  default     = ""
  sensitive   = true
}

# ─── Networking ───────────────────────────────

variable "vnet_cidr" {
  type        = string
  description = "Virtual network CIDR"
  default     = "10.30.0.0/16"
}

variable "aks_subnet_cidr" {
  type        = string
  description = "AKS subnet CIDR"
  default     = "10.30.1.0/24"
}

variable "appgw_subnet_cidr" {
  type        = string
  description = "Application Gateway subnet CIDR"
  default     = "10.30.2.0/24"
}

variable "db_subnet_cidr" {
  type        = string
  description = "Database delegated subnet CIDR"
  default     = "10.30.3.0/24"
}

# ─── Tags ─────────────────────────────────────

variable "tags" {
  type        = map(string)
  description = "Tags for all resources"
  default = {
    project     = "storm"
    environment = "load-test"
    managed_by  = "terraform"
  }
}
