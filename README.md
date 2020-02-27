# Xmpp-notifier
Github action to notify xmpp users when some events occur on a given repository.  

You can either notify a single user or send a message to a channel.

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
      - name: Checkout
        uses: actions/checkout@v2
        # Checkout the action repository 
        with:
            repository: processone/xmpp-notifier
      - name: push_info_step
        id: push
        uses: ./
        # Will only trigger when a push is made to the master branch
        if: github.event_name == 'push'
        with: # Set the secrets as inputs
          # Login expects the bot's bare jid (user@domain)
          login: ${{ secrets.bot_username }}
          pass: ${{ secrets.bot_password }}
          server_domain: ${{ secrets.server_rooms_domain }}
          # Correspondent is the intended recipient of the notification. 
          # If it is a single user, the bare Jid is expected (jid without resource)
          # If it is a chat room, only the name of it is expected, and "server_domain" will be used to complete the jid
          correspondant: ${{ secrets.room_correspondent }}
          # Port is optional. Defaults to 5222
          server_port: ${{ secrets.server_port }}
          message: |
            ${{ github.actor }} pushed ${{ github.event.ref }} ${{ github.event.compare }} with message:
            ${{ join(github.event.commits.*.message) }}
          # Boolean to indicate if correspondent should be treated as a room (true) or a single user 
          correspondent_is_room: true
      - name: pr_info_step
        id: pull_request
        uses: ./
        # Will only get triggered when a pull request to master is created
        if: github.event_name == 'pull_request'
        with: # Set the secrets as inputs
          login: ${{ secrets.bot_username }}
          pass: ${{ secrets.bot_password }}
          server_domain: ${{ secrets.server_rooms_domain }}
          correspondant: ${{ secrets.room_correspondent }}
          message: |
            ${{ github.actor }} opened a PR ${{ github.event.html_url }}
          correspondent_is_room: true
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
      # Checkout your project, then run tests in the following step.
      - name: Checkout
        uses: actions/checkout@v2
      - name: Run tests
        run: |
          go test ./... -v -race
      - name: Checkout action if tests failed
        if: failure() # If tests fail, let's pull the action in and start it 
        uses: actions/checkout@v2
        with:
          repository: processone/xmpp-notifier
      - name: Tests failed notif
        # Now that the action is here, start it.
        # Apparently github considers the "failure()" should return true if *any* of the previous steps fail
        # Even if the doc seem to say "only if the previous step fails" (https://help.github.com/en/actions/reference/contexts-and-expression-syntax-for-github-actions#job-status-check-functions)
        if: failure()  
        id: test_fail_notif
        uses: ./
        with: # Set the secrets as inputs
          # Login expects the bot's bare jid (user@domain)
          login: ${{ secrets.bot_username }}
          pass: ${{ secrets.bot_password }}
          server_domain: ${{ secrets.server_rooms_domain }}
          # Correspondent is the intended recipient of the notification.
          # If it is a single user, the bare Jid is expected (jid without resource)
          # If it is a chat room, only the name of it is expected, and "server_domain" will be used to complete the jid
          correspondant: ${{ secrets.room_correspondent }}
          # Port is optional. Defaults to 5222
          server_port: ${{ secrets.server_port }}
          message: |
            tests for the following PR have failed : ${{ github.event.pull_request.html_url }}
          # Boolean to indicate if correspondent should be treated as a room (true) or a single user
          correspondent_is_room: true
```