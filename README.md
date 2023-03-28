# Outbox pattern

### Proyect structure
- config: process configuration using .env (viper)
- cmd: main entrypoint
- docs: swagger files
- db: schemas for MySQL 
- internal
  - application: configure cateogry modules (controllers, handlers, jobs)
  - infrastructure: configure MySQL, Rabbit, Logs

### Desing system
Design outbox pattern in go
![outbox drawio (1)](https://user-images.githubusercontent.com/5313452/228265009-5da4f318-c49c-424a-94c5-dce8c46e64e1.png)
