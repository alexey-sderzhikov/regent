## Description
Terminal base redmine client. Build with [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) and using [redmine API](https://www.redmine.org/projects/redmine/wiki/rest_api).

## Installation
1. Go in root project's directory `cd regent/` 
2. Create `.env` file
3. Set two environment

```
#/regent/.env
SOURCE=https://redmine.source.com
USER_API_KEY=my_user_api_key   
```

*Note: You can find your user api key in redmaine->my account*

4. Run regent with `go run .` or build `go build` and run with `./regent`
## TODO
- [x] Take key api from config file
- [ ] Get all time entries 
- [ ] Notify if no time entries yestarday
- [ ] Implement help element from charm library
- [ ] Add spiner
- [x] Get issues only current user
- [x] Add pagination 
