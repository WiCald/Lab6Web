# Lab6Web

API en sí, bienvenida:
![image](https://github.com/user-attachments/assets/43b3cc0e-568a-4798-abeb-554051b58623)
Lista de partidas en el API:
![image](https://github.com/user-attachments/assets/20eb9659-03bb-40b2-a636-c17e3ddf60e9)

Correr Frontend...
Frontend> docker run --rm -it -p 3000:80 -v ${pwd}:/usr/share/nginx/html:ro nginx

Correr backend...
Backend> docker build -t laliga-backend .
Backend> docker run -p 8080:8080 laliga-backend
