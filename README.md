[![Go Report Card](https://goreportcard.com/badge/github.com/mlycore/asgard)](https://goreportcard.com/report/github.com/mlycore/asgard)

# Asgard

Asgard is an implementation of artifacts repository in golang. Decision to implement new tiny repo come to me when i realized that nexus requires 4gb RAM machine at minimum. Check the [Nexus memory requrements](https://help.sonatype.com/display/NXRM3/System+Requirements#SystemRequirements-Memory). That is too much, especially when you need it only for small personal projects. This project gives you minimal, but complete nexus functionality, with ability to save artifacts on filesystem or to s3 and basic auth. For most of usecases that is more than enought.

## How to use it ?

```
docker run -d -v ${HOME}/config.yml:/config.yml -p 8080:8080 mworks92/asgard
```

## Configuration 

Create config.yml in your ${HOME} directory or in the same directory where you run binary.

For s3:
```yml
---
http:
  addr: ":443"
  username: "myuser"
  password: "mypassword"
  https: true
  crt: "/certs/domain.crt"
  key: "/certs/domain.key"

storage:
  type: "s3"
  bucket_name: "asgardtest"
  access_key: "*******************"
  secret_key: "**************************************"
```

And for file system:
```yml
---
http:
  addr: ":8080"
  username: "myuser"
  password: "mypassword"

storage:
  type: "fs"
  base_dir: "/applications"
```

## How to build it ?

Make sure to install golang, set all env variables etc.
Clone project to your go-workspace.
Cd to the project folder and run:

```
make build
```

And it will compile app.

Run:

```
make run
```

And it will run app locally on port 8080 by default.

## License

[GPL](LICENSE)

[license-url]: LICENSE

[license-image]: https://img.shields.io/github/license/mashape/apistatus.svg

[capture]: capture.png
