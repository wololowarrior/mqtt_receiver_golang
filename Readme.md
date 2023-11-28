## IOT Project to receive speed from sensor

### Project Description

- Language: Golang
- Docker based deployment
- [mqtt](./mqtt) Contains the code that consumes the message received on the mqtt broker and updates it in redis
- [httpHandler](./httpHandler) contains the webserver that:
  1. Throws back a jwt token
  2. Using that token, we can `GET` the latest speed from redis.

### Api Flow

See the api flow ![here](https://raw.githubusercontent.com/wololowarrior/mqtt_receiver_golang/main/images/api_flow.png)

### Running the deployment
1. Using the docker compose file we'll bring up the topology. First time the docker will build the images.
2. First start redis and broker using : `docker-compose up redis emqx-broker -d`.
   1. Wait for 10-20 secs for emqx broker to initialize
   2. We could also start the mqtt handler directly, since i've used the `depends-on` setting, 
   but since I've not implemented a healthcheck flow we don't know when the mqtt broker has fully initialized.
3. Start the http server and mqtt handler using `docker-compose up mqtt-handler http-server -d`.
4. See for ref. 

    ![this](./images/dc_start_flow.png)  

### Running tests
1. Using mosquitto_pub to publish messages to the topic
2. Publish a speed update to the broker: 
    ```shell
     mosquitto_pub -h localhost -p 1883 -t "sensor/speed" -m "{\"speed\":10}"
    ```
3. Getting the JWT token:
    ```shell
    hgupta@Harshils-MacBook-Pro mqtt_receiver % curl --location 'localhost:4000/token' \                                                                                      
    --header 'Content-Type: application/json' \
    --data-raw '{
    "email": "h@g.com"
    }' | jq .token
    % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
    Dload  Upload   Total   Spent    Left  Speed
    100   189  100   165  100    24  15231   2215 --:--:-- --:--:-- --:--:-- 31500
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImhAZy5jb20iLCJleHAiOjE3MDExNjcxNjMsImlhdCI6MTcwMTE2Njg2M30.S9hZwqqzFSl5Z8QOVeAjA68Di_vd7fDqTkgnS7nxgsE"
    ```
4. Get the speed. Copy the token from above and insert it below:
    ```shell
    hgupta@Harshils-MacBook-Pro mqtt_receiver % curl --location 'localhost:4000/speed' \
    --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImhAZy5jb20iLCJleHAiOjE3MDExNjcxNjMsImlhdCI6MTcwMTE2Njg2M30.S9hZwqqzFSl5Z8QOVeAjA68Di_vd7fDqTkgnS7nxgsE'
    {"speed":"10"}

    ```
