# Tech choices

## General architecture

This will generally be an event-driven system, partly for fun and partly because
the idea of sending in scores lends itself well to the idea. We can take in the
raw data of scores and generate various views and other data points off of them,
and enable real-time viewing more easily in particular.

## Assumptions

Suppose 1000 users. Each user plays 1 round of golf per week, with 3 average
players per round. Each player puts in 1 score per hole. Each score may
contain a few fields, such as putts, GIR, sand, etc. to give some flavor to the
score.

As a ballpark:

1 hole = player ID, round ID, score, par, length, hole number, putts, sand, and
a few other fields. For a conservative estimate, this could be 100 bytes of
data. Multiply by 3 per user per week and 18 holes per round, this would give
us a total of `5,400,000 B` or around `527 KB` per week total, or about `75 KB`
per day.

## Database selection

Because this is a hobby project, we are optimizing for cost. Amazon DynamoDB
will give us 25 GB of free storage, which means our 75 KB per day of raw data
would give us about 1 million years of free storage. And then maybe we'll pay a
few extra cents on top of that, or whatever galactic currency is in use at that
time.
