## Overpunch

[![Build Status](https://travis-ci.org/barkimedes/go-deepcopy.svg?branch=master)](https://travis-ci.org/barkimedes/go-deepcopy) [![codecov](https://codecov.io/gh/barkimedes/go-deepcopy/branch/master/graph/badge.svg)](https://codecov.io/gh/barkimedes/go-deepcopy) [![](https://godoc.org/github.com/nathany/looper?status.svg)](https://godoc.org/github.com/barkimedes/go-deepcopy)

This package is a Golang implementation for creating deep copies of virtually any kind of Go type. 

This is a truly deep copy--every single value behind a pointer, every item in a slice or array, and every key and value in a map are all cloned so nothing is pointing to what it pointed to before.

To handle circular pointer references (e.g. a pointer to a struct with a pointer field that points back to the original struct), we keep track of a map of pointers that have already been visited. This serves two purposes. First, it keeps us from getting into any kind of infinite loop. Second, it ensures that the code will behave similarly to how it would have on the original struct -- if you expect two values to be pointing at the same value within the copied tree, then they'll both still point to the same thing.
