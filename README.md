# Docker Nuke
Docker Nuke is a Docker CLI plugin that:
Removes ALL of the containers and images with NO EXCEPTION!

## Requirements
You need go installed and configured to build it.

## Build
```sh
$ make
```

## Install

For installing, you have 2 options. 

One is to make a link to the built binary with:

```sh
$ make link
```

Another one is to copy the built binary to the plugins' folder with:
```sh
$ make install
```

### Check installation
To check the installation, just run:
``` 
$ docker --help
```
And check if the plugin is listed as a managed command

## Run
Once you have it installed, you can just use that as any other docker cli plugin, like

```sh
$ docker nuke
```

# License
[MIT](LICENSE)