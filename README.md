## Description
Terminal base redmine client. Built with [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) and using [redmine API](https://www.redmine.org/projects/redmine/wiki/rest_api).

## Installation
1. Go to root project's directory `cd regent/` 
2. Create `.env` file
3. Set two environment

```
#/regent/.env
SOURCE=https://redmine.source.com
USER_API_KEY=examplekey12345   
```

*Note: You can find your user api key in redmine->my account*

4. Run regent with `go run .` or build `go build` and run with `./regent`
## TODO
- [x] Take key api from config file
- [x] Get all time entries 
- [ ] Notify if no time entries yestarday
- [ ] Implement help element from bubble library
- [ ] Add spiner
- [x] Get issues only current user
- [x] Add pagination
- [ ] Functional to add and change issues
- [ ] Menu
- [ ] View port for viewing issue and another objects
