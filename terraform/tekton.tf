resource "null_resource" "tekton_pipelines" {
  provisioner "local-exec" {
    command = "kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml"
  }
}
