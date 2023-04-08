docker build -t myapp .
docker run -p 8080:8080 -it -v --rm myapp