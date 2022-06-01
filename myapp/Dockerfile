FROM node
WORKDIR /usr/src/app
COPY ./ ./
COPY ./src/utils/dummy.json ./
RUN npm install
EXPOSE 3000
CMD ["npm", "start"]