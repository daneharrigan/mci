# mci: marvel comics informer

Get informed when new issues for a series you like comes to Marvel Unlimited

## about

I'm running this on Heroku and using the following addons:

- Mailgun
- Heroku PostgreSQL
- Heroku Scheduler

I have Heroku Scheduler running my `bin/monday` script on a daily basis. Because
I don't have the control to run only on Monday, my script will exit if
`date +%a` isn't "Mon".

## setup

```
git clone git@github.com:daneharrigan/mci.git
cd mci
bin/setup test
bin/setup development
```

## testing

```
bin/test
```
