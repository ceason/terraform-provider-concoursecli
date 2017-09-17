## Motivation
> todo: write this

## Usage
This provider has a few prerequsites right now:
- `fly` must exist on your $PATH (may change in future)
- Download the plugin from the [Releases](https://github.com/ceason/terraform-provider-concoursecli/releases) page
- [Install](https://terraform.io/docs/plugins/basics.html) it

Example pipeline resource
```hcl
provider concoursecli {
  target = "flycli-target-name"
}

resource concoursecli_pipeline hello_world {
  name   = "hello-demo"
  paused = true
  public = false
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
```


#### Logging
This provider currently uses `glog`. To configure it, you must pass flags via environment variable like this:
```shell
export TF_PROVIDER_CONCOURSECLI_FLAGS="-v=4 -log_dir=/tmp/glog"
```

## Todo
- More tests
    - test fly wrapper via object lifecycle (create, unpause, pause, show, hide, update, delete)
- Input validation
    - Ensure specified fly target exists
- More flexible provider config
    - Perhaps via environment variables, etc
    - Maybe download `fly` if not present. Could even specify version
- Implement Teams
    - How to handle authentication config? may not be too bad..
    - Team switching will be more involved
        - Will require creating targets for each team and switching between
