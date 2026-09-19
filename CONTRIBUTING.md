// fix the contributing.md file later

### COMMIT RULES 

## structure : 
`<type>(<scope>): <summary>`

## types :
feat      - New feature
fix       - Bug fix
refactor  - Code restructuring
perf      - Performance improvement
test      - Tests
docs      - Documentation
chore     - Maintenance
style     - Formatting


## scope :
commands
resp
server
client
general


### Managing codecrafters and github branches

## 1. push to main and then codecrafters test

-> write code by switching to the main branch
-> commit and push from main to the backup remote
-> change branch to master branch
-> cherry pick the latest commit from the main branch using `git cherry-pick main/latest commit hash`
-> add, commit and push/ codecrafters test/ codecrafters submit on the master branch

## 2. changes made in main, dont wanna commit them but also wanna change the branches ? 

-> use `git stash`, pushes the code globally, you can do `git stash pop` anywhere
