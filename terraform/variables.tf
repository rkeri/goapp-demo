variable "namespace" {
    type    = string
    default = "production"
}

variable "environment" {
    type    = string
    default = "prod"
}

variable "kube_context" {
    type    = string
    default = "kind-kind"
}
