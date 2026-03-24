resource "helm_release" "diario_obra" {
  name = "diario-obra"
  namespace = "diario-obra"
chart = "/home/victor/projetos/diario-obras/chart"
}
