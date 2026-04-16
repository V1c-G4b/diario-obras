resource "helm_release" "diario_obra" {
  name       = "diario-obra"
  namespace  = "diario-obra"
  chart      = "${path.module}/../chart"
  depends_on = [helm_release.monitoring]

  set {
    name  = "database.user"
    value = var.database_user
  }

  set {
    name  = "database.password"
    value = var.database_password
  }

  set {
    name  = "storage.user"
    value = var.storage_user
  }

  set {
    name  = "storage.password"
    value = var.storage_password
  }
}
