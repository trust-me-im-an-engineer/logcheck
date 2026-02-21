# logcheck - линтер для проверки лог-записей

Тестовое задание для Selectel выполнено в соответствии с [техническими требованиями](https://tn-eoc.mck2.ru/c/M-UmAAQAAABeS1s1/sqQWBw/h9VzoB3JRdKhvY9j/?u=https%3A%2F%2Ffiles.selectel.ru%2Fdocs%2Fru%2Fbackend_golang_testovoe.pdf%3Futm_source%3Dmindbox_trig%26utm_medium%3Demail%26utm_campaign%3DHR_Career_Wave_confirmation_2026).

## Установка и использование

Линтер можно использовать как плагин для golangci-lint, либо как самостоятельное приложение.

### Плагин для golangci-lint

Чтобы добавить logcheck в качестве плагина для golangci-lint клонировать этот репозиторий не нужно. Нужно собрать golangci-lint с logcheck плагином и включить его.

#### Сборка golangci-lint
Для сборки нужны установленный golangci-lint актуальной версии (v2.10.1) и [файл конфигурации](.custom-gcl.yml) в текущей директории.

Собрать golangci-lint можно командой:
```bash
golangci-lint custom
```

#### Включение плагина

По умолчанию golangci-lint не включает добавленные модули, по этому нужно добавить logcheck в .golangci.yml

Если в проекте еще нет конфигурационного файла, добавьте [.golangci.yml](.golangci.yml).

Либо если в проекте уже есть .golangci.yml добавьте следующее в секцию linters:

```yml
  enable:
  - logcheck

  settings:
    custom:
      logcheck:
        type: "module"
        description: Linter for log messages
```

Запустить собранный линтер можно командой

```bash
./custom-gcl run
```

### Самостоятельное приложение

Так же logcheck можно установить как отдельный линтер:

```bash
go install github.com/trust-me-im-an-engineer/logcheck/cmd/logcheck@latest
```

Чтобы запустить проверку текущего проекта:

```bash
logcheck ./...
```

## Тестирование

Проект содержит unit-тесты для правил проверки и интеграционные тесты для анализатора.

Для запуска тестов:

```bash
go test ./...
```

## Бонусные задания

### Авто-исправление
Реализован SuggestedFixes для
автоматического исправления сообщений, начинающихся с заглавной буквы.

