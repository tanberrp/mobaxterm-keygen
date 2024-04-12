# mobaxterm-keygen

Mobaxterm-keygen is a build cil tool used to solve the upper limit of mobaxterm session.

## Use steps

### Download cli tool
+ You can download the binary compression package for the corresponding platform from [release page](https://github.com/tanberrp/mobaxterm-keygen/releases).

+ If you have golang environment, you can download with go install:
  ```shell
  go install github.com/tanberrp/mobaxterm-keygen/cmd/mobaxterm-keygen@latest
  ```
  
### Generate MobaXterm key
```shell
mobaxterm-keygen --username tanber --version 23.5 --mobaxterm-dir  ~/MobaXterm_Portable_V23.5
```
Flags:
+ --username: The user name to licensed, just specify one at will. 
+ --version: The version of MobaXterm licensed to.
+ --mobaxterm-dir: The dir of mobaxterm has installed.

### Reboot MobaXterm
Reboot MobaXterm, and you can save session without limit.

