
resource "helm_release" "monitoring" {
  name             = "monitoring"
  namespace        = "monitoring"
  chart            = "kube-prometheus-stack"
  create_namespace = true

  repository = "https://prometheus-community.github.io/helm-charts"

  set {
    name  = "grafana.adminPassword"
    value = var.grafana_admin_password
  }
}
