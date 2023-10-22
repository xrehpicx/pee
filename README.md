# Pee

Project Environment Executor, Define your project workspace as code and config.

## Usecase
Creating tmux sessions with preconfigured panes layouts and commands. inspired from [tmuxinator](https://github.com/tmuxinator/tmuxinator)

### Get Started

![Screen Recording 2023-10-22 at 2 11 11 PM](https://github.com/xrehpicx/pee/assets/22765674/b355da63-1e2d-4833-8300-1bd879e2245f)

#### Installation
```bash
go install github.com/xrehpicx/pee@latest
```

#### Initialize a project
```bash
pee init myproject
```
Select a directory from the file picker by `<space>`, the directory that the init is run in will be the default starting directory for the file picker, you can go up and down directories by `<backspace>` and `<enter>`

#### Configure project
Config is similar to [tmuxinator](https://github.com/tmuxinator/tmuxinator)
```yml
name: ppec-ui
editor: nvim
root: /Users/raj.sharma/Documents/GitHub/ppec-ui
windows:
  - window_name: editor
    layout: 8070,202x58,0,0[202x46,0,0,89,202x11,0,47,92]
    panes:
      - shell_command:
          - nvim "+SessionManager load_current_dir_session"
      - shell_command:
          - echo 'npm run dev'
  - window_name: hosts
    layout: even-horizontal
    shell_command_before:
      - cd somewhere && activate env
    panes:
      - shell_command:
          - ssh stg-host1
      - shell_command:
          - ssh stg-host2
  - window_name: git
    panes:
      - shell_command:
          - lazygit
lastopened: 2023-10-22T14:03:54.071678+05:30
attach: true
```
editor here is used as the default editor for editing the config file.
`note: attach here does not tmux attach the new session, instead uses tmux switch-client for faster switching.`

#### Run a project
For this example we have a project called ppec with the above config
```bash
pee ppec
```
and this should open up ppec tmux session

#### Edit or Run from list
![Screen Recording 2023-10-22 at 2 12 41 PM](https://github.com/xrehpicx/pee/assets/22765674/f8bb6c1d-1a68-4194-8c4c-62ff4856cd2c)

You can also run
```bash
pee ls
```
and select the project from table and click `<enter>` to open/create the session or `<e>` to edit the config of the project 

### Roadmap
1. Supporting iTerm2
   As of now this supports only tmux windows and panes, would want to add support for iterm and other terminals if they have api's to do so.
2. Ability to save an opened session into a config or update a config
3. Parse cli args to config
4. Ability to save custom layouts as named layouts that can be used across multiple projects
5. Sync configs across devices

---
All contributions are welcome 

