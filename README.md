# camera-adapter

Адаптер для подключения превращения внешних пультов к камерам onvif. Релаизована поддержка HID класса Keyboard c жестким ограничением нажатий 
Разработка велась под Raspberry PI3, однако система может быть собрана без привязки к архитектуре, в рамках одноплатного компьютера с убунтой или дебиан на борду.

В основе лежит Go, версии 1.19.

# Процесс разработки

Для начала работы
```bash
git clone https://github.com/Zeryoshka/camera-adapter.git
cd camera-adapter
go mod download
```

Сборка и запуск
```bash
go build cmd/main.go
sudo ./main -conf conf.yaml
```
*Запускаем с sudo, иначе hidapi не даст прав на открытие девайса (выявлено на RaspberryOS)

При добавлении новых пакетов выполнить: 
```bash
go mod tidy
```

# Структура проекта

`cmd/main.go` - входная точка проекта, в ней реализован основной процесс обработки событий
`camera/` - пакет, отвечающий за управление камерой, он скрывает в себе общий процесс выполнения команд, переключения между камерами и т.д. Содержит сущности: Camera, CameraManager, Command
`reader/` - пакет определяющий интерфейс для разработки новых ридеров и релализованные ридеры (сейчас Keyboard HID)
`confstore/` - пакет определяющий хранение конфига

# Сборка устройства
* Подготовить одноплатник с накаченным образом ОС
* Выполнить сборку и загрузку (лучше собирать на плате, чтобы выполнялось с учетом архитектуры процессора)
```bash
git clone https://github.com/Zeryoshka/camera-adapter.git
cd camera-adapter
go build cmd/main.go
cp main /usr/local/bin/cam-hub
cd ../
cd rm camera-adapter
``` 
* Подготовить службу и настроить запуск