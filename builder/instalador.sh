#!/bin/bash

# Movemos el binario a su ubicación
cp gaby /usr/local/sbin/

# Creamos el directorio para guardar datos 
mkdir /var/lib/gaby

# Creamos un ejemplo de fichero de configuración
cat <<MAFI >/etc/gaby.yaml
destino: http://stats.sanidad.gob.sv:8086
destino-token: 2-qbR-mdKDF6f9qO-QW-UftFSeuGnXUoc_R2W_UKEw6mC1mbndISAbnKyw40dCdgaQtfQH2dYFHlRtV0gWpHgA==
destino-organizacion: sanidad
destino-bucket: ambientales
MAFI
