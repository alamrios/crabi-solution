# Crabi solution
¡Hola!, esta es mi solución a la evaluación práctica de Crabi.

## Desarrollo
Para ejecutar el servicio usar el comando `make serve`. Esto levantará el docker-compose configurado con los siguientes contenedores:

| Nombre | Imagen | Puerto |
| ------ | ------ | ------ |
| crabi-pld | vligascrabi/crabi-pld-test:v1 | 3000 |
| crabi-solution-mongo | mongo:4.2.0 | 27217 |
| crabi-solution-dev | crabi-solution | 8080 |

## Pruebas unitarias
Para ejecutar las pruebas unitarias ejecutar el comando `make test`

## Detalles
Algunos datos importantes sobre el funcionamiento del proyecto:
- El archivo de [configuracion](/config/config.go) permite cargar las variables de entorno requeridas por los servicios externos para funcionar
- Para ejecutar los servicios en local (pensado solamente para desarrollo) se usan los archivos Docker dentro de [infra](/infra/deploy/local/)
- El servicio usa Mongo como base de datos para persistir a los usuarios
- Se utiliza cifrado SHA256 para almacenar las contraseñas

## API
El contenedor crabi-solution-dev expone los siguientes endpoints:
### `/api/v1/users/ [POST]`
Para crear usuarios, espera un json con el siguiente formato:
```json
{
  "first_name": string,
  "last_name": string,
  "email": string,
  "password": string
}
```
### `/api/v1/login/ [POST]`
Para autenticación de los usuarios, espera un json con el siguiente formato:
```json
{
  "email": string,
  "password": string
}
```
### `/api/v1/users/{email} [GET]`
Para consultar los datos de un usuario mediante su email, espera el *email* del usuario como parámetro en la solicitud.

## Postman
En el repositorio se incluye una colección de Postman para probar los endpoints expuestos por el servicio.