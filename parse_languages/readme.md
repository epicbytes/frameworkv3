# Парсер вызовов локализации

Собирает из кода, из папки с локализацией и матчит обратно в папку локализации. При добавлении в шаблон вызова локализации, данное приложение после удачной сборки шаблона сканирует все варианты вызова и матчит их с имеющимися на данный момент переводами

## Формат и название файлов переводов

### Название:

[lang].locale.yml

например `en.locale.yml` или `ru.locale.yml`

```yaml
en:
  general:
    actions: ""
    do_filtering: ""
    is_active: ""
    reset: ""
  layout:
    goto_main_page: ""
    loading: ""
    logout: ""
```

## Установка

Для глобальной установки:
```bash
go install github.com/epicbytes/frameworkv3/parser_languages@latest
```
Для установки при разработке:
```bash
go install ./parser_languages
````

## Подключение в Taskfile.yml гейтвея

```bash
./parse_languages `путь к шаблонам` `путь к папке с локализациями` `перечисление необходимых локалей`
```

```yaml
  collect_localizations:
    cmds:
      - parse_languages ./server/templates ./server/localizations en,ru

  build_templ:
    cmds:
      - templ fmt .
      - templ generate
      - task: collect_localizations
      - yarn build_tw
    sources:
      - ./server/templates/*.templ
      - ./server/templates/**/*.templ
```

## Вызов в коде шаблонов Templ

```html
<span>{ lib.T(ctx, "scope", "token", args...) }</span>
<span>{ lib.T(ctx, "layout","ordered","1знач","2знач","3333") }</span>
<span>{ lib.T(ctx, "layout","cnt", map[string]int{"Count":1 }) }</span>
```

где `scope` это раздел, страница, сущность, группа, `token` - токен перевода, `args` - дополнительные параметры в случае использования plural или map нотаций