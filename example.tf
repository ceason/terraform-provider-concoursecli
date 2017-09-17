
resource concoursecli_pipeline hello_world {
  name   = "hello-demo"
  paused = true
  public = false
  // language=yaml
  yaml   = <<EOF
jobs:
- name: hello-world
  plan:
  - task: say-hello
    config:
      platform: linux
      image_resource:
        type: docker-image
        source: {repository: ubuntu}
      run:
        path: echo
        args: ["Hello, world!"]
EOF
}