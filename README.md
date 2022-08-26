# gaby
Scrapper de los datos de sensores EATON

## Construcción
La construcción del binario debe hacerse en un sistema independiente al servidor en producción. Un servidor en producción no debería tener herramientas para compilación de paquetes.
Pues contrario a la publicidad, Go no es tan independiente como dice ser. Si no se puede usar un sistema igual (O compatible) al destino, pues podría usarse podman (o docker) para que funcione
```bash
podman run  -it  --rm -v "$PWD":/go/src/myapp -w /go/src/myapp golang:1.18-bullseye go build .
```

## Instalación
El binario se envía al servidor destino. SCP bastaría
scp gaby root@servidor:/usr/local/sbin

## Configuracion
```bash
# Crear el fichero para guardar la hora
echo 1600000000000 > /var/lib/gaby

# Configurar las siguientes variables del sistema
URL_ORIGEN_DATOS=http://user:pass@10.0.0.9
GABY_ORGANIZACION="sanidad"
GABY_BUCKET="ambientales"
export GABY_ENDPOINT="http://stats.dominio.com"
export GABY_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export URL_ORIGEN_DATOS="http://user:password@10.0.0.97"
```
