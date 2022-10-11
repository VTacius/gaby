# gaby
Scrapper de los datos de sensores EATON Powerware

## Escogiendo la versión
Los paquetes se encuentran en [la página de lanzamientos](https://github.com/VTacius/gaby/releases/latest). Por ahora, solo hay paquetes para el último lanzamiento de Debian, y de ellos, debe escoger la versión para su sistema según el esquema de versionado:
``` 
gaby-{version-aplicativo}-{version-debian}.tgz 
```

Así por ejemplo, para la versión 
```
gaby-v0.9.5-11.4.tgz 
```

Significa que es la versión `v0.9.5` de la aplicación para la versión `11.4` de Debian

## Instalación
Una vez el paquete esté en el servidor, de descomprime:
```bash
tar xzvf gaby-v0.9.5-11.5.tgz
```

Se entra al directorio resultante
```bash
cd gaby/
```

Y se ejecuta el instalador
```bash
bash instalador.sh 
```


## Configuracion
El fichero en `/etc/gaby.yaml` debería ser descriptivo sobre los datos requeridos

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
