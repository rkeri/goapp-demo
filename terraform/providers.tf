terraform {
    required_version = ">= 1.5"

    required_providers {
        kubernetes = {
            source  = "hashicorp/kubernetes"
            version = "~> 2.30"
        }
        helm = {
            source  = "hashicorp/helm"
            version = "~> 3.0"
        }
    }
}

provider "kubernetes" {
    config_path    = "~/.kube/config"
    config_context = var.kube_context
}

provider "helm" {
    kubernetes = {
        config_path    = "~/.kube/config"
        config_context = var.kube_context
    }
}
