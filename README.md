
![asme logo](assets/logo/logo_t.png)

# asme @ itba website


## How to generate site
1. **download `go` tool** from [golang.org](https://golang.org/).
2. **Check if Go is installed.** run `go version` in command prompt to check if tool is installed. 
you should see (the version may be higher):
    ```
    >> go version go1.14.3 windows/amd64
    ```
3. **Set `GOPATH` environment variable to some directory**. 
    
    a.  **On windows:** type `env`
    in search -> "Edit the system environment variables" ->
    environment variables -> New... -> 
    Variable Name:`GOPATH`, 
    Variable Value:`C:\*your directory*`

4. install dependencies with the following commands
     ```
    go get -u github.com/asme-itba/asme-itba.github.io/src
     ```
5. edit [`generate.yaml`](src/generate.yaml)
    to add changes you wish to see in site
6. **Run generator**. Open command prompt in 
   /src/ folder and run the generator with 
    ```
    go run .
    ```
7. Enjoy new site. You can now push changes to github.

### Extra information
* [YAML specification](https://yaml.org/spec/1.2/spec.html)
* [How to open command prompt easily Windows 10](https://www.itechtics.com/open-command-window-folder/)
