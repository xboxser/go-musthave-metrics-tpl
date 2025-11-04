# временный файл с данными метрики
echo '{"id":"test_counter","type":"counter","delta":1}' > counter_metric.json

# POST-запросы с метрикой
ab -t 4 -c 1 -p counter_metric.json -T "application/json" http://127.0.0.1:8080/update/
ab -t 4 -c 1 http://127.0.0.1:8080/
ab -t 4 -c 1 http://127.0.0.1:8080/ping
ab -t 4 -c 1 http://127.0.0.1:8080/value/counter/test_counter

ab -t 4 -c 1 -p counter_metric.json -T "application/json" http://127.0.0.1:8080/update/
ab -t 4 -c 1 http://127.0.0.1:8080/
ab -t 4 -c 1 http://127.0.0.1:8080/ping
ab -t 4 -c 1 http://127.0.0.1:8080/value/counter/test_counter

# удаляем временный файл
rm counter_metric.json