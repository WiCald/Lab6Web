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

En caso de que dé un error al correr el backend, volver a correr con esto para asegurarse que sea siempre la versión más nueva de código...
Backend> docker build --no-cache -t laliga-backend .
Backend> docker run -p 8080:8080 laliga-backend

Postman: https://wc-2887871.postman.co/workspace/W-C's-Workspace~a86b466f-7d68-4b7b-8c3b-81963ec4eb72/collection/43555267-e121ad30-e2df-4111-abbd-ea51aa3aab6d?action=share&creator=43555267

NOTA: Asegurarse de tener Docker Desktop abierto antes de intentar correr ambos
