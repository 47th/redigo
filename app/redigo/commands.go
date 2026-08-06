package redigo

import (
	"strconv"
	"strings"
	"time"
)

type Value struct {
	Env       Envelope
	Expiry    time.Time
	ExpirySet bool
}

var db map[string]Value = make(map[string]Value)

func commandRouter(command Envelope) Envelope {

	if command.Size == 0 {
		return handleErr("ERR no command given", SimpleError)
	}

	op := command.Array[0]

	switch com := strings.ToUpper(op.String); com {
	case "PING":
		return ping()
	case "ECHO":
		return echo(command)
	case "SET":
		return set(command)
	case "GET":
		return get(command)
	default:
		var args string
		for i := range command.Size {
			args = args + "'" + command.Array[i].String + "' "
		}
		errStr := "ERR unknown command '" + op.String + "', with args beginning with: " + args
		return handleErr(errStr, SimpleError)
	}

}

func ping() Envelope {
	env := Envelope{
		OpCode: SimpleString,
		String: "PONG",
	}

	return env
}

func echo(req Envelope) Envelope {
	if req.Size != 2 {
		return handleErr("ERR wrong number of arguments for 'echo' command", SimpleError)
	}

	env := Envelope{
		OpCode: BulkString,
		Size:   req.Array[1].Size,
		String: req.Array[1].String,
	}

	return env
}

func handleErr(str string, errType Type) Envelope {
	env := Envelope{
		OpCode: errType,
		String: str,
		Size:   len(str),
	}

	return env
}

func set(env Envelope) Envelope {
	switch env.Size {
	case 3:
		key := env.Array[1].String
		value := Value{
			Env: env.Array[2],
		}
		db[key] = value

		return Envelope{
			OpCode: SimpleString,
			String: "OK",
			Size:   2,
		}
	case 5:
		switch option := env.Array[3].String; option {
		case "EX", "PX":

			// calculate expiry
			interval, err := strconv.Atoi(env.Array[4].String)
			if err != nil {
				return handleErr("ERR error in converting seconds to integer", SimpleError)
			}
			expiry := time.Now()
			if option == "EX" {
				expiry = expiry.Add(time.Duration(interval) * time.Second)
			} else if option == "PX" {
				expiry = expiry.Add(time.Duration(interval) * time.Millisecond)
			}

			key := env.Array[1].String
			value := Value{
				Env:       env.Array[2],
				ExpirySet: true,
				Expiry:    expiry,
			}

			db[key] = value

			return Envelope{
				OpCode: SimpleString,
				String: "OK",
				Size:   2,
			}
		}
	default:
		return handleErr("ERR wrong number of arguments for 'set' command", SimpleError)
	}

	return handleErr("ERR wrong number of arguments for 'set' command", SimpleError)
}

func get(env Envelope) Envelope {
	NIL := Envelope{
		OpCode: BulkString,
		Size:   -1,
	}

	if env.Size != 2 {
		return handleErr("ERR wrong number of arguments for 'get' command", SimpleError)
	}
	key := env.Array[1].String
	value, ok := db[key]

	if ok {
		if value.ExpirySet && time.Now().After(value.Expiry) {
			delete(db, key)
			return NIL
		}

		return Envelope{
			OpCode: BulkString,
			Size:   len(value.Env.String),
			String: value.Env.String,
		}
	}

	return NIL
}
