go-dockerclient-commit
======================

This is a tiny test repo that recreates an issue I've been having with committing containers using
[foo](http://github.com/fsouza-dockerclient).  The when I create an image by commiting a container
via the library, the resulting image does not have a CMD.

### Usage

    go test github.com/pmorie/go-dockerclient-commit