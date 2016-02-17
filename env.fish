if not set -q GOROOT
  echo "GOROOT is not set."
else
  set -x GOPATH $PWD/fake-go-path
  set -x GO15VENDOREXPERIMENT 1

  # setup fake go path for intellij
  mkdir -p $GOPATH vendor
  ln -sf ../vendor $GOPATH/src
  rm -f vendor/vendor
end
