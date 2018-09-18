# exopulse kvstore package
Golang package kvstore provides support for local persistent key-value store.

[![CircleCI](https://circleci.com/gh/exopulse/kvstore.svg?style=svg)](https://circleci.com/gh/exopulse/kvstore)
[![Build Status](https://travis-ci.org/exopulse/kvstore.svg?branch=master)](https://travis-ci.org/exopulse/kvstore)
[![GitHub license](https://img.shields.io/github/license/exopulse/kvstore.svg)](https://github.com/exopulse/kvstore/blob/master/LICENSE)

# Overview

This package provides support for local persistent key-value store.
This implementation is a simple wrapper around boltdb. See https://github.com/boltdb/bolt

# Using kvstore package

## Installing package

Use go get to install the latest version of the library.

    $ go get github.com/exopulse/kvstore
 
Include kvstore in your application.
```go
import "github.com/exopulse/kvstore"
```

## Usage

```go
manager := kvstore.NewManager(&kvstore.Config{Filename:"/tmp/store.db"})

err := manager.Open()

defer manager.Close()

err := manager.Update(func(trx *kvstore.Trx) error {
	if err := trx.InitializeBucket("messages"); err != nil {
		return err
	}

	if err := trx.InitializeBucket("customers"); err != nil {
		return err
	}

    user, err := trx.Create("users", func(id kvstore.ID) interface{} {
        return newUser(id)
    })

	return nil
})

```

# About the project

## Contributors

* [exopulse](https://github.com/exopulse)

## License

Kvstore package is released under the MIT license. See
[LICENSE](https://github.com/exopulse/kvstore/blob/master/LICENSE)
