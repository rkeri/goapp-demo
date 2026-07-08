output "namespace" {
    value = kubernetes_namespace.goapp-demo.metadata[0].name
}

output "release_name" {
    value = helm_release.goapp-demo.name
}
