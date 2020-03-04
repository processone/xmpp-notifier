# Xmpp-notifier
Github action to notify xmpp users when some events occur on a given repository.  

You can either notify a single user or send a message to a channel.

## List of parameters
To have more information on parameters this action can accept, please refer to the 
[action.yml file](https://github.com/processone/xmpp-notifier/blob/master/action.yml). 


## Main.yml
This file could be named as you wish, but has to be placed in the .github.workflows directory of your project.
This is an example for the main configuration that could be used to call the action :  
```yaml
on:
  # Specifies that we only want to trigger the following jobs on pushes and pull request creations for the master branch
  push:
    branches:
      - master
  pull_request:
    branches:
      - master
jobs:
  notif-script:
    runs-on: ubuntu-latest
    name: workflow that pushes repo news to xmpp server
    steps:
      - name: push_info_step
        id: push
        uses: processone/xmpp-notifier@master
        # Will only trigger when a push is made to the master branch
        if: github.event_name == 'push'
        with: # Set the secrets as inputs
          # jid expects the bot's bare jid (user@domain)
          jid: ${{ secrets.jid }}
          password: ${{ secrets.password }}
          server_host: ${{ secrets.server_host }}
          # Intended recipient of the notification such as a room or single user. Bare JID expected.
          recipient: ${{ secrets.recipient }}
          # Port is optional. Defaults to 5222
          server_port: ${{ secrets.server_port }}
          message: |
            ${{ github.actor }} pushed ${{ github.event.ref }} ${{ github.event.compare }} with message:
            ${{ join(github.event.commits.*.message) }}
          # Boolean to indicate if correspondent should be treated as a room (true) or a single user (false)
          recipient_is_room: true
      - name: pr_open_info_step
        id: pull_request_open
        uses: processone/xmpp-notifier@master
        # Will only get triggered when a pull request to master is created
        if: github.event_name == 'pull_request' && github.event.action == 'opened'
        with: # Set the secrets as inputs
          jid: ${{ secrets.jid }}
          password: ${{ secrets.password }}
          server_host: ${{ secrets.server_host }}
          recipient: ${{ secrets.recipient }}
          message: |
            ${{ github.actor }} opened a PR : ${{ github.event.pull_request.html_url }} with message :
            ${{ github.event.pull_request.title }}
          recipient_is_room: true
      - name: pr_edit_info_step
        id: pull_request_edit
        uses: processone/xmpp-notifier@master
        # Will only get triggered when a pull request to master is created
        if: github.event_name == 'pull_request' && github.event.action == 'edited'
        with: # Set the secrets as inputs
          jid: ${{ secrets.jid }}
          password: ${{ secrets.password }}
          server_host: ${{ secrets.server_host }}
          recipient: ${{ secrets.recipient }}
          message: |
            ${{ github.actor }} edited the following PR : ${{ github.event.pull_request.html_url }} with message :
            ${{ github.event.pull_request.title }}
          recipient_is_room: true
``` 

## action.yml  
This file must be placed at the project root, and should not be renamed (see github actions documentation).  
You should not modify it because the go program relies on it.  

## Dockerfile
The Dockerfile in this action is used to delpoy a docker container and run the go code that will notify users.  

## entrypoint.sh
Used as the entry point of the docker container. Meaning this is executed when the docker container is started.  
This script uses inputs from the github action.

## main.go
A small go program that will be compiled and ran in the docker container when the github action is executed.  
It uses the [native go-xmpp library](https://github.com/FluuxIO/go-xmpp).

## Example of configuration to trigger notifications on tests failure for a PR

```yaml
jobs:
  notif-script:
    runs-on: ubuntu-latest
    name: workflow that pushes test failures to xmpp server
    steps:
      - name: Run tests
        run: |
          go test ./... -v -race
      - name: Tests failed notif
        if: failure()  
        id: test_fail_notif
        uses: processone/xmpp-notifier@master
        with: # Set the secrets as inputs
          # Login expects the bot's bare jid (user@domain)
          jid: ${{ secrets.jid }}
          password: ${{ secrets.password }}
          server_host: ${{ secrets.server_host }}
          # The intended recipient of the notification such as a xmpp room or a single user. Bare JID is expected
          recipient: ${{ secrets.recipient }}
          # Port is optional. Defaults to 5222
          server_port: ${{ secrets.server_port }}
          message: |
            tests for the following PR have failed : ${{ github.event.pull_request.html_url }}
          # Boolean to indicate if correspondent should be treated as a room (true) or a single user
          recipient_is_room: true
```
