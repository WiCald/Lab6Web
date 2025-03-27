# Lab6Web

API en sí, bienvenida:
![image](https://github.com/user-attachments/assets/43b3cc0e-568a-4798-abeb-554051b58623)
Lista de partidas en el API, sin partidas:
![image](https://github.com/user-attachments/assets/20eb9659-03bb-40b2-a636-c17e3ddf60e9)
Lista de partidas en el API, con la partida de Transformers vs One Piece
![image](https://github.com/user-attachments/assets/7233f185-9e34-438a-b4ee-4feb7819ffa6)
![image](https://github.com/user-attachments/assets/24e0cd6a-1259-467e-b638-d04f40a9fb16)


Correr Frontend...
Frontend> docker run --rm -it -p 3000:80 -v ${pwd}:/usr/share/nginx/html:ro nginx

Correr backend...
Backend> docker build -t laliga-backend .
Backend> docker run -p 8080:8080 laliga-backend
