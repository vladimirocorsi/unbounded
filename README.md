# Unbounded

Same as https://github.com/jonbodner/unbounded but using a ring buffer.
Has an issue with backing array never being shrinked when space is not needed anymore, e.g. queue gets emptied.
