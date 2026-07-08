resource "kubernetes_namespace" "goapp-demo" {
    metadata {
        name = var.namespace
    }
}

resource "helm_release" "goapp-demo" {
    name       = "goapp-demo"
    chart      = "../helm"
    namespace  = kubernetes_namespace.goapp-demo.metadata[0].name

    values = [
        file("${path.module}/values.yaml")
    ]
}
