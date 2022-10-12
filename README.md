# Jira tool workflow

- On open, check to see if there is a config file
  - If one exists, move forward
  - If one doesn't exist, advise the user to fill in the config file located at <dir>
- Check to see if we're able to make a connection with the config info
  - If not, inform the user and exit
  - If the connection is made, check to see if there is a default board in the config
    - If there is, use that board
    - If there isn't let the user select a board
      - Once the board is selected, ask them if they'd like to set that board as default
- Main view is laid out like the backlog view in jira. Epic filters in a bar on the left,
  current sprint issues in a box in the main content section, backlog issues below that.
  - Toggling epics toggles the filters on the issues. If no epics are toggles, everything is visible
  - Selecting an issue should bring up the issue view, which displays the child elements, description, etc
