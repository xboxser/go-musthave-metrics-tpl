# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## 17 инкремент
В рамках оптимизации дорабатывал метод service.GetAll.   
Его оптимизация заметна при просмотре alloc_space.   
Нагрузку задавал ghb при помощи ApacheBench `load_test/load_test.sh`.    

Проведение нагрузочного тестирования так же выявило проблему в работе проекта. Исправлена ошибка получения данных из метода GetAll при многопоточном режиме.

```
go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
File: main
Type: inuse_space
Time: Nov 4, 2025 at 12:30pm (+05)
Duration: 60.01s, Total samples = 512.17kB 
Showing nodes accounting for -2097.83kB, 409.59% of 512.17kB total
      flat  flat%   sum%        cum   cum%
 -902.59kB 176.23% 176.23% -2097.87kB 409.60%  compress/flate.NewWriter (inline)
 -650.62kB 127.03% 303.26% -1195.29kB 233.38%  compress/flate.(*compressor).init
 -544.67kB 106.34% 409.60%  -544.67kB 106.34%  compress/flate.newDeflateFast (inline)
  512.22kB   100% 309.59%   512.22kB   100%  runtime.malg
 -512.17kB   100% 409.59%  -512.17kB   100%  internal/profile.(*Profile).postDecode
         0     0% 409.59% -2097.87kB 409.60%  compress/gzip.(*Writer).Close
         0     0% 409.59% -2097.87kB 409.60%  compress/gzip.(*Writer).Write
         0     0% 409.59%  -512.17kB   100%  github.com/go-chi/chi/v5.(*Mux).Mount.func1
         0     0% 409.59% -2610.05kB 509.60%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 409.59%  -512.17kB   100%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 409.59%  -512.17kB   100%  internal/profile.Parse
         0     0% 409.59%  -512.17kB   100%  internal/profile.parseUncompressed
         0     0% 409.59% -2610.05kB 509.60%  metrics/internal/handler/middleware.GzipMiddleware.func1
         0     0% 409.59% -2610.05kB 509.60%  net/http.(*conn).serve
         0     0% 409.59% -2610.05kB 509.60%  net/http.HandlerFunc.ServeHTTP
         0     0% 409.59% -2610.05kB 509.60%  net/http.serverHandler.ServeHTTP
         0     0% 409.59%  -512.17kB   100%  net/http/pprof.Index
         0     0% 409.59%  -512.17kB   100%  net/http/pprof.collectProfile
         0     0% 409.59%  -512.17kB   100%  net/http/pprof.handler.ServeHTTP
         0     0% 409.59%  -512.17kB   100%  net/http/pprof.handler.serveDeltaProfile
         0     0% 409.59%   512.22kB   100%  runtime.newproc.func1
         0     0% 409.59%   512.22kB   100%  runtime.newproc1
         0     0% 409.59%   512.22kB   100%  runtime.systemstack
```

```
go tool pprof -alloc_space  -top -diff_base=profiles/base.pprof profiles/result.pprof 
File: main
Type: alloc_space
Time: Nov 4, 2025 at 12:30pm (+05)
Duration: 60.01s, Total samples = 2377094857B 
Showing nodes accounting for 109216266B, 4.59% of 2377094857B total
Dropped 58 nodes (cum <= 11885474B)
      flat  flat%   sum%        cum   cum%
-47734998B  2.01%  2.01% -11020657B  0.46%  metrics/internal/service.(*ServerService).GetAll
 36176688B  1.52%  0.49%  36176688B  1.52%  strconv.FormatFloat (inline)
 18367509B  0.77%  0.29%  18367509B  0.77%  text/template.addFuncs (inline)
 17826072B  0.75%  1.04%   2097124B 0.088%  reflect.Value.call
 16791765B  0.71%  1.74%  16791765B  0.71%  metrics/internal/model.(*MemStorage).GetAll
-16253224B  0.68%  1.06% -16778408B  0.71%  fmt.Sprintf
 14186902B   0.6%  1.66%  26770006B  1.13%  internal/fmtsort.Sort
 14156304B   0.6%  2.25%  13630236B  0.57%  metrics/internal/handler/middleware.(*loggingResponseWriter).Write
 13111200B  0.55%  2.80%  36178083B  1.52%  net/http.readRequest
```
## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**
