# opsbots
A Slackbot that connects Slack and Pagerduty together, allowing users to use slack slash-commands to invoke Clippy (the helper slack-bot).  

Each command will coordinate a function designed to carry out operational tasks that can/should be automated.  


Command: /incident

Create a PagerDuty Incident
Create a SlackChannel
Connect the SlackChannel ID and the PagerDuty ID
Invite identified users and promote to awareness Slack Channel
Maintain the Slack and PagerDuty communication/state until the channel or incident is resolved
