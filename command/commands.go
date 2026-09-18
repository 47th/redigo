package command

import (
	"bufio"
	"redigo/resp"
	"redigo/store"
	"strconv"
	"strings"
	"time"
)

// type Value struct {
// 	string    val
// 	Expiry    time.Time
// 	ExpirySet bool
// }

type Handler struct {
	store *store.Store
	// router map[string]func(resp.Envelope) string
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{
		store: store,
		// router: map[string]func(resp.Envelope) string{
		// 	"PING": ping,
		// 	"ECHO": echo,
		// 	"SET":  set,
		//	"GET":  get,
		// },
	}
}

func HandleRequest(h *Handler, reader *bufio.Reader) (string, error) {
	var res string
	var env resp.Envelope

	env, err := resp.ParseInput(reader)

	if err != nil {
		return res, err
	}

	res = commandRouter(env, h)
	return res, nil

}

func commandRouter(command resp.Envelope, h *Handler) string {

	if command.Size == 0 {
		return HandleErr("ERR no command given", resp.SimpleError)
	}

	op := command.Array[0]

	switch com := strings.ToUpper(op.String); com {
	case "PING":
		return ping()
	case "ECHO":
		return echo(command)
	case "SET":
		return set(command, h)
	case "GET":
		return get(command, h)
	case "INCR":
		return incr(command, h)
	default:
		var args string
		for i := range command.Size {
			args = args + "'" + command.Array[i].String + "' "
		}
		errStr := "ERR unknown command '" + op.String + "', with args beginning with: " + args
		return HandleErr(errStr, resp.SimpleError)
	}

}

func HandleErr(str string, errType resp.Type) string {
	env := resp.Envelope{
		OpCode: errType,
		String: str,
		Size:   len(str),
	}

	return resp.FormatMapper(env)

}

func ping() string {
	env := resp.Envelope{
		OpCode: resp.SimpleString,
		String: "PONG",
	}

	return resp.FormatMapper(env)
}

func echo(req resp.Envelope) string {
	if req.Size != 2 {
		return HandleErr("ERR wrong number of arguments for 'echo' command", resp.SimpleError)
	}

	env := resp.Envelope{
		OpCode: resp.BulkString,
		Size:   req.Array[1].Size,
		String: req.Array[1].String,
	}

	return resp.FormatMapper(env)
}

func set(env resp.Envelope, h *Handler) string {
	if env.Size < 3 || env.Size > 5 {
		return HandleErr("ERR wrong number of arguments for 'set' command", resp.SimpleError)
	}

	key := env.Array[1].String
	value := env.Array[2].String
	var expiry time.Time

	switch env.Size {
	case 3:
		expiry = time.Time{}
	case 5:
		switch option := strings.ToUpper(env.Array[3].String); option {
		case "EX", "PX":
			interval, err := strconv.Atoi(env.Array[4].String)
			if err != nil {
				return HandleErr("ERR error in converting seconds to integer", resp.SimpleError)
			}
			expiry = time.Now()
			if option == "EX" {
				expiry = expiry.Add(time.Duration(interval) * time.Second)
			} else {
				expiry = expiry.Add(time.Duration(interval) * time.Millisecond)
			}
		}
	default:
		return HandleErr("ERR wrong number of arguments for 'set' command", resp.SimpleError)
	}

	h.store.SET(key, value, expiry)
	resEnv := resp.Envelope{
		OpCode: resp.SimpleString,
		String: "OK",
		Size:   2,
	}

	return resp.FormatMapper(resEnv)

}

func get(env resp.Envelope, h *Handler) string {
	NIL := resp.Envelope{
		OpCode: resp.BulkString,
		Size:   -1,
	}

	if env.Size != 2 {
		return HandleErr("ERR wrong number of arguments for 'get' command", resp.SimpleError)
	}
	key := env.Array[1].String
	value, ok := h.store.GET(key)

	if ok {
		resEnv := resp.Envelope{
			OpCode: resp.BulkString,
			Size:   len(value),
			String: value,
		}
		return resp.FormatMapper(resEnv)
	}

	return resp.FormatMapper(NIL)
}

func incr(env resp.Envelope, h *Handler) string {
	if env.Size != 2 {
		return HandleErr("ERR wrong number of arguments for 'incr' command", resp.SimpleError)
	}

	key := env.Array[1].String
	value, ok := h.store.INCR(key)
	if ok {
		str := "(integer) " + value
		resEnv := resp.Envelope{
			OpCode: resp.SimpleString,
			Size:   len(str),
			String: str,
		}

		return resp.FormatMapper(resEnv)
	}

	return HandleErr("ERR value is not an integer or out of range", resp.SimpleError)

}
