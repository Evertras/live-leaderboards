# Data

Taking the advice of [this AWS blog
post](https://aws.amazon.com/blogs/database/single-table-vs-multi-table-design-in-amazon-dynamodb/),
we try to use a single table overall and composite sort keys to help
differentiate different groups of hierarchical data.

## Round start events

Round start events consist of the following fields as their composite primary key:

- Round ID (Partition key)
- `rnd_start` (Sort key)

It also contains:

- Start time (date)
- Course (string)

| Field | Type   | Description                        |
| ----- | ------ | ---------------------------------- |
| `p`   | list   | The list of player names, in order |
| `c`   | string | The course name                    |
| `s`   | string | The start time                     |

![Diagram](./diagrams/event_round_start.svg)

## Score events

Score events consist of the following fields as their composite primary key:

- Round ID (Partition key)
- `sc^[hole]^[player_index]` (Sort key)

It also contains:

| Field | Type   | Description                 |
| ----- | ------ | --------------------------- |
| `h`   | number | Hole number (usually 1-18)  |
| `pi`  | number | Player index (0 indexed)    |
| `pa`  | number | Par for the hole            |
| `sc`  | number | Player's score for the hole |

![Diagram](./diagrams/event_score.svg)
