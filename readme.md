# Inverval Containment Sets implemened in GO

Interval containment sets (ICS) provide efficient solution for testing
containment in overlapping arithmetic intervals. This includes binary
classification for discrete and continuous values like integer of floating point
numbers, but can also be extended to arbitrary comparable and orderable values.

## Description

Assume you have a set of intervals, maybe overlapping, describing a certain
property of a value or attributing it to a class. To test if a value is
contained within the set, you can simply iterate over all the intervals
comparing it to lower and upper limits. This would be a straightforward approach
that works nicely if you have just a few intervals, or if the overall speed is
not important.

When performance and simplicity matters, however, this library offers a more
efficient solution there the intervals are linearised and flattened into a
structure that can then be used for efficient containment testing.

Some of the use cases: lookup tables for Unicode character properties, text
parsers, regexp optimization, etc.

For example:

```
A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
original intervals:

   └───────┘         └───────┘   └─┘
       └─────────┘             └─────────┘  
             └─┘         └───┘       └─┘
flattened intervals:

   └─────────────┘   └───────┘ └─────────┘
```            

Such flattened sets are represented as a sorted array of values where each
interval is formed by pairing adjacent elements. Values with even indices are
considered as inclusive lower boundaries, values with odd indices are considered
as exclusive upper boundaries. There are no duplicate elements in this array. If
the total number of elements is odd, then the last element is treaded as
open-ended interval. 

```
flattened intervals:
A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
   └─────────────┘   └───────┘ └─────────┘

in-memory ics representation:
elements: C J L P Q V
meaning:  [ ) [ ) [ )
```

Searching for containment on such a structure is extremely easy: perform a
binary search on the elements, effectively finding a lower bound index, then
check if the found index is an even or an odd number. Even indices then indicate
containment.

The sets are composited by starting with an empty set then inserting (merging)
the intervals one-by-one. The insertion routine then creates, extends, and
merges the elements in the set as required to keep it consistent with the
incoming intervals. This insertion/merging routine is quite efficient for high
performance runtime usage. 

For classification tasks that are static in nature, the resulting arrays of
elements then can be stored in memory, as a file, or even code-gened into the
source code. Later, then the containment lookup is required, use the array with
this library or as a direct input into a binary search algorithm followed by
even/odd check.

## Bonus Feature

The logic of the set can be inverted by inserting or removing a minimum value
(e.g. zero for unsigned integers) as a first element in the array:

```
values:    0 1 2 3 4 5 6 7 8 9 ... maxint
              └───┘     └───┘ 
elements:     [2,  4)   [7,  9)     

inverted:  0 1 2 3 4 5 6 7 8 9 ... maxint
          └───┘   └─────┘   └────────────┘ 

elements: [0,  2) [4,    7) [9  
```

## Unicode Intervals

The library features a couple of containment sets specializations for Unicode
codepoints and ASCII codeunits. For convenience, these sets provide insertion
and enumeration API that operates with fully closed `[a-z]`-style ranges instead
of half-open intervals.