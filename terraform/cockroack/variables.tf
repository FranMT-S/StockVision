variable "cockroach_apikey" {
  description = "Api key of cockroachdb"
  type        = string
  sensitive   = true
}

variable "cluster_name" {
  description = "Name of the cluster"
  type        = string
  default     = "stockvisioncluster"
}

variable "db_user_name" {
  description = "Name of the sql user"
  type        = string
  default     = "appuser"
}

variable "db_user_password" {
  description = "Password of the sql user"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "Region of the cluster"
  type        = string
  default     = "us-east-1"
}

variable "db_name" {
  description = "Name of the database"
  type        = string
  default     = "defaultdb"
}
