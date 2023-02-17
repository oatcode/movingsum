# movingsum

[![Docs](https://pkg.go.dev/badge/github.com/oatcode/movingsum)](https://pkg.go.dev/github.com/oatcode/movingsum)

movingsum is a set of utilities for calculating moving sum, the sum of a series of numbers bounded by entry count or time.
Moving average and moving rate can also be derived.

- MovingSum
  - calculates moving sum with a fixed queue size
- MovingSumByTime
  - calculates moving sum for a duration
  - Note that this uses an unbounded queue. For large data set, consider using AggregatedMovingSumByTime.
- AggregatedMovingSumByTime
  - calculates moving sum with entries aggregated into time slots
  - This has a bounded storage compare to MovingSumByTime

