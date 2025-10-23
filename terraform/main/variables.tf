# AWS variables
variable "aws_region" { 
  type = string 
  default = "us-east-1" 

  validation {
    condition     = var.aws_region != ""
    error_message = "not must be empty"
  }
}

variable "vpc_cidr"   { 
  type = string
  default = "10.0.0.0/16" 

  validation {
    condition     = var.vpc_cidr != ""
    error_message = "not must be empty"
  }
}

# ECR variables
variable "ecr_api_image" {
  type    = string

  validation {
    condition     = var.ecr_api_image != ""
    error_message = "not must be empty"
  }
}

variable "ecr_app_image" {
  type    = string

  validation {
    condition     = var.ecr_app_image != ""
    error_message = "not must be empty"
  }
}

# Database variables
variable "db_host" {
  type    = string

  validation {
    condition     = var.db_host != ""
    error_message = "not must be empty"
  }
}

variable "db_name" {
  type    = string

  validation {
    condition     = var.db_name != ""
    error_message = "not must be empty"
  }
}

variable "db_user" {
  type    = string

  validation {
    condition     = var.db_user != ""
    error_message = "not must be empty"
  }
}

variable "db_password" {
  type    = string
  sensitive = true

  validation {
    condition     = length(var.db_password) > 5
    error_message = "must be have length greater than 5"
  }
}

variable "db_port" {
  type    = string

  validation {
    condition     = var.db_port != ""
    error_message = "not must be empty"
  }
}

variable "db_ssl" {
  type    = string
  default = "true"

  validation {
    condition     = var.db_ssl != ""
    error_message = "not must be empty"
  }
}

variable "stock_api_url" {
  type    = string

  validation {
    condition     = var.stock_api_url != ""
    error_message = "not must be empty"
  }
}

variable "stock_api_token" {
  type    = string

  validation {
    condition     = var.stock_api_token != ""
    error_message = "not must be empty"
  }
}

variable "financial_base_url" {
  type    = string

  validation {
    condition     = var.financial_base_url != ""
    error_message = "not must be empty"
  }
}

variable "financial_token" {
  type    = string

  validation {
    condition     = var.financial_token != ""
    error_message = "not must be empty"
  }
}

variable "finhub_base_url" {
  type    = string

  validation {
    condition     = var.finhub_base_url != ""
    error_message = "not must be empty"
  }
}

variable "finhub_token" {
  type    = string

  validation {
    condition     = var.finhub_token != ""
    error_message = "not must be empty"
  }
}

variable "gemini_api_key" {
  type    = string

  validation {
    condition     = var.gemini_api_key != ""
    error_message = "not must be empty"
  }
}



