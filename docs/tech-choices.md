# Tech choices

## Goals

We are optimizing for cost. This is a hobby project. At a higher scale and at
higher budgets, other choices are better.

## General architecture

This will generally be an event-driven system, partly for fun and partly because
the idea of sending in scores lends itself well to the idea. We can take in the
raw data of scores and generate various views and other data points off of them,
and enable real-time viewing more easily in particular.

## Cloud tech

Going with AWS because I already have an account, I'm familiar with it, and it's
the gold standard. The free tiers of things should be plenty for us to work
with as well.

## Assumptions

Suppose 1000 users. Each user plays 1 round of golf per week, with 3 average
players per round. Each player puts in 1 score per hole. Each score may
contain a few fields, such as putts, GIR, sand, etc. to give some flavor to the
score.

As a ballpark:

1 hole = player ID, round ID, score, par, length, hole number, putts, sand, and
a few other fields. For a conservative estimate, this could be 100 bytes of
data. Multiply by 3 per user per week and 18 holes per round, this would give
us a total of `5,400,000 B` or around `5.27 MB` per week total, or about `753 KB`
per day.

## Database selection

Because this is a hobby project, we are optimizing for cost. Amazon DynamoDB
will give us 25 GB of free storage, which means our 753 KB per day of raw data
would give us about 95 years of free storage.

It also has good integrations with Lambdas, which we'll make use of elsewhere.
These integrations are also low price or free, which is an additional bonus.
