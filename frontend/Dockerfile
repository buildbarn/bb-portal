FROM node:18-bullseye-slim

WORKDIR /frontend
COPY . /frontend

RUN npm install
#RUN npm run build

ENTRYPOINT ["npm", "run", "start"]
