# Package Merge Git Driver
Resolve semantic version conflicts in package.json and package-lock.json. Resolved only conflicts of project version not depencencies.

### How it works
For example if u've got git flow repository:

|                       | master | develop | release |
| ---                   | ---    | ---     | ---     |
| before release finish | 1.0.5  | 1.2.0   |  1.1.0  |
| after                 | 1.1.0  | 1.2.0   |         |
| before release finish | 1.3.0  | 1.2.0   |  1.1.0  |
| after                 | 1.3.0  | 1.3.0   |         |

In simple words higher version always win.

### Build (optional)
* git clone https://github.com/Ksawierek/package-merge-driver.git
* cd package-merge-driver
* docker run --rm -e XDG_CACHE_HOME=/tmp/.cache -u $UID -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.14 go build -o bin/package-merge-driver
* copy package-merge-driver to any path where driver will by used

### How to use it
* change ~/.gitconfig adding:
```
[merge "packagemerge"]
    name = A custom merge driver for Npm's package.json
    driver = {PATH_TO_DRIVER}/package-merge-driver %O %A %B
```

* ad in project .gitattributes:
```
package.json merge=packagemerge
package-lock.json merge=packagemerge
```

### Multistage docker build
If you're using Docker image to build npm project, you can create image that have driver inside docker image.
For example (including nodejs, gitflow) [Dockerfile](Dockerfile).

```
docker build .
```

And now you can use image to build javascript project with npm.

Everything was tested only on Linux System!