#!/bin/bash

# Creamos el directorio de trabajo si no existe
[ -d output ] || mkdir output

DEBIAN_VERSION=(11.4 11.5)
GABY_VERSION=$(git describe --abbrev=0 --tags)

for dv in ${DEBIAN_VERSION[*]}; do 
    [ -f gaby ] && rm gaby
    podman run --rm -it -v $PWD:/usr/local/src/ alortiz/constructor-golang:$dv-1.18 go build
    tar czvf output/gaby-${GABY_VERSION}-$dv.tgz builder/*instalador.sh gaby --transform 's,^builder/,,' --transform 's,^,gaby/,' 
done
