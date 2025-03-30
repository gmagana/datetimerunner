# datetimerunner
A utility written in Go to delete files. Intended to delete the oldest X logs or backup files.

---
# What datetimerunner does

## Test Mode

---

# Usage

`datetimerunner <executable_name> <argument_template> [-t]`

## Arguments

| Argument | Comments| 
| --- | --- |
| `<executable_name>` | (REQUIRED) The full pathname of the exceutable you want to run. Paths with spaces will be automatically quoted. |
| `<argument_template>` | (REQUIRED) The argument template you want to pass to the executable, using the below placeholders |
| `-t` | If specified, nothing is executed and the parsed arguments are shown, as well as the command that would run with all the placeholders replaced |

## Working with Quoted Arguments

When you need to include spaces or special characters in your arguments, you should use quotes. Here are some guidelines for properly specifying quoted arguments:

### Command Prompt (CMD)

In Windows Command Prompt, use double quotes to enclose arguments with spaces:

```
datetimerunner.exe "c:\program files\app.exe" "argument with spaces"
```

For arguments that themselves need quotes, use quotes to escaped via `\"`:
```
datetimerunner.exe app.exe "argument with \"quoted\" text"
```

### PowerShell

In PowerShell, use double quotes to enclose arguments with spaces:

```
.\datetimerunner.exe "c:\program files\app.exe" "argument with spaces"
```

For arguments that themselves need quotes, use backtick (`) to escape quotes:
```
.\datetimerunner.exe app.exe "argument with `"quoted`" text"
```

### Batch Files

In batch files, double quotes can be used as in Command Prompt, but special characters like < and > need to be escaped with ^ when used with placeholders:

```
datetimerunner.exe app.exe "echo Date: ^<yyyy^>-^<mm^>-^<dd^>"
```

## Placeholders

| Placeholder | Comments |
| --- | --- |
| `<y>`       | Two-digit year (without leading zero) |
| `<yy>`      | Two-digit year (with leading zero) |
| `<yyyy>`    | Four-digit year |
| `<m>`       | Month (without leading zero) |
| `<mm>`      | Month (with leading zero) |
| `<mmm>`     | Month (three-letter abbreviation) |
| `<mmmm>`    | Full month name |
| `<d>`       | Day of month (without leading zero) |
| `<dd>`      | Day of month (with leading zero) |
| `<ddd>`     | Day of year (with leading zeroes) |
| `<h24>`     | Hour (without leading zero) - 24 hour scheme |
| `<hh24>`    | Hour (with leading zero) - 24 hour sceheme |
| `<h12>`     | Hour (without leading zero) - 12 hour scheme |
| `<hh12>`    | Hour (with leading zero) - 12 hour sceheme |
| `<ampm>`    | AM/PM label |
| `<i>`       | Minute (without leading zero) |
| `<mi>`      | Minute (with leading zero) |
| `<s>`       | Second (without leading zero) |
| `<ss>`      | Second (with leading zero) |
| `<dow>`     | Day of week (three-letter abbreviation) |
| `<weekday>` | Day of week (full name) |

## Examples

```
DateTimeRunner "c:\utils\pkzipc.exe" "-add c:\backup\full-<yyyy><mm><dd>-<hh><mi><ss>.zip c:\sourcedata\*.*"
```

```
DateTimeRunner "C:\WINNT\system32\xcopy.exe" "m:\app\log.txt m:\app\log-<d>.<m>.<yy>.txt"
```

```
DateTimeRunner "c:\utils\pkzipc.exe" "-add c:\backup\full-<yyyy><mm><dd>-<hh><mi><ss>.zip c:\sourcedata\*.* ""abc.txt"""
