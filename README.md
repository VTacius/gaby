# gaby
Scrapper de los datos de sensores EATON

## Construcción
La versión de Go a usar es la 1.18.
Por conveniencia, la construcción del binario debe hacerse en un sistema independiente al servidor en producción. 

```bash
go build .
```

Si lo prefiere, puede usarse podman (o docker) para tal tarea
```bash
podman run  -it  --rm -v "$PWD":/go/src/myapp -w /go/src/myapp golang:1.18-bullseye go build .
```

## Instalación
El binario se envía al servidor destino. SCP bastaría
```bash
scp gaby root@servidor:/usr/local/sbin
```

## Configuracion
```bash
### Creamos el directorio para guardar el último envio
mkdir /var/lib/gaby

### Crear los ficheros para guardar la hora. Necesitamos uno por cada sensor que tengamos
echo 1600000000000 > /var/lib/gaby/10.10.20.21

### Configurar las variables del sistema en el archivo correspondiente:
cat <<MAFI >/etc/default/gaby
GABY_ENDPOINT="http://stats.dominio.com"
GABY_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
GABY_BUCKET="ambientales"
GABY_ORGANIZACION="sanidad"
MAFI
```

## Uso
### Desde consola
```bash
# Obtiene los datos para el último muestreo del sensor 
gaby --origen 10.10.20.21 --password-origen admin --usuario-origen EATON 

# Obtiene los datos para el ultimo muestreo y los envía al servidor influxDB
gaby --origen 10.10.20.21 --password-origen admin --usuario-origen EATON --envio

```
### Como una tarea mediante systemd
```bash
# El siguiente es un template cuya variable es la IP del sensor EATON destino
cat <<MAFI> /lib/systemd/system/gaby@.service 
[Unit]
Description=Gaby: Scrapper para dispositivo EATON

[Service]
EnvironmentFile=/etc/default/gaby
ExecStart=/usr/local/sbin/gaby --origen %i --password-origen admin --usuario-origen EATON --envio

[Install]
WantedBy=multi-user.target
Also=gaby.timer
MAFI

# Un target nos permite llamar a varios servicios
cat <<MAFI> /lib/systemd/system/gaby.target 
[Unit]
Description=Target para los servicios de Gaby
BindsTo=gaby@10.0.0.97.service gaby@10.0.0.98.service
After=gaby@10.0.0.97.service gaby@10.0.0.98.service

[Install]
WantedBy=timers.target
MAFI

cat <<MAFI> /lib/systemd/system/gaby.timer 
[Unit]
Description=Corre el Scrapper Eaton cada minuto

[Timer]
OnCalendar=*-*-* *:*:00
Unit=gaby.target

[Install]
WantedBy=timers.target
MAFI

# Activamos el timer
systemctl enable --now gaby.timer
```
