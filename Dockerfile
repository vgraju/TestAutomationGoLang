FROM scratch

WORKDIR /app
ADD ./main /app/main

CMD [ "/app/main" ]
