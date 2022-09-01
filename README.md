# gaby
Scrapper de los datos de sensores EATON Powerware

## Construcción
La versión de Go a usar es la 1.18.

Por conveniencia, la construcción del binario debe hacerse en un sistema independiente al servidor en producción. Para construirlo, basta con ejecutar el siguiente comando desde el directorio

```bash
go build .
```

Si lo prefiere, puede usarse podman (o docker) para tal tarea
```bash
podman run  -it  --rm -v "$PWD":/go/src/myapp -w /go/src/myapp golang:1.18-bullseye go build .
```

## Instalación
El binario se envía al servidor destino. Por ejemplo, con ```scp```
```bash
scp gaby root@servidor:/usr/local/sbin
```

## Configuracion
```bash
### Creamos el directorio para guardar el último envio
mkdir /var/lib/gaby

### Tenemos un fichero de configuración de tipo YAML para guardar las credenciales de acceso al backend InfluxDB2:
cat <<MAFI >/etc/gaby.yaml
destino: http://stats.sanidad.gob.sv:8086
destino-token: 2-qbR-mdKDF6f9qO-QW-UftFSeuGnXUoc_R2W_UKEw6mC1mbndISAbnKyw40dCdgaQtfQH2dYFHlRtV0gWpHgA==
destino-organizacion: sanidad
destino-bucket: ambientales
MAFI
```

## Uso
### Desde consola
Considerando que el fichero anterior se haya configurado correctamente, y que ambos servidores tienen las mismas credenciales, el único cambio sería la dirección IP, con lo cual el comando a usar quedaría de la siguiente forma:
```bash
# Obtiene los datos para el último muestreo del sensor 
gaby --origen 10.0.0.98

# Obtiene los datos para el ultimo muestreo y los envía al servidor influxDB
gaby --origen 10.0.0.98 --envio 

```
El comando cuenta con una pequeña ayuda para revisar todos los parámetros que es posible cambiar

```bash
gaby --help 
```

### Como una tarea mediante systemd
```bash
# El siguiente es un template cuya variable es la IP del sensor EATON destino
cat <<MAFI > /lib/systemd/system/gaby@.service 
[Unit]
Description=Gaby: Scrapper para dispositivo EATON

[Service]
ExecStart=/usr/local/sbin/gaby --origen %i --envio

[Install]
WantedBy=multi-user.target
Also=gaby.timer
MAFI

# Un target nos permite llamar a varios servicios
cat <<MAFI > /lib/systemd/system/gaby.target 
[Unit]
Description=Target para los servicios de Gaby
BindsTo=gaby@10.0.0.97.service gaby@10.0.0.98.service
After=gaby@10.0.0.97.service gaby@10.0.0.98.service

[Install]
WantedBy=timers.target
MAFI

cat <<MAFI > /lib/systemd/system/gaby.timer 
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
