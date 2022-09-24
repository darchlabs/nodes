FROM node

WORKDIR /home/node

RUN npm install -g ganache

CMD ["ganache", "--fork", "https://capable-quick-vineyard.discover.quiknode.pro/3f240555880826d7153fe818e34bf1aef2b75e50@15569350", "--host", "0.0.0.0", "--db", "/home/node/data"]
